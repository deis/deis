import cStringIO
import base64
import copy
import json
import httplib
import paramiko
import socket
import re
import time
from django.conf import settings


MATCH = re.compile(
    '(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z-_]+)?.(?P<c_num>[0-9]+)')
RETRIES = 3


class UHTTPConnection(httplib.HTTPConnection):
    """Subclass of Python library HTTPConnection that uses a Unix domain socket.
    """

    def __init__(self, path):
        httplib.HTTPConnection.__init__(self, 'localhost')
        self.path = path

    def connect(self):
        sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        sock.connect(self.path)
        self.sock = sock


class FleetHTTPClient(object):

    def __init__(self, target, auth, options, pkey):
        self.target = target
        self.auth = auth
        self.options = options
        self.pkey = pkey
        # single global connection
        self.conn = UHTTPConnection(self.target)

    # connection helpers

    def _put_unit(self, name, body):
        headers = {'Content-Type': 'application/json'}
        self.conn.request('PUT', '/v1-alpha/units/{name}.service'.format(**locals()),
                          headers=headers, body=json.dumps(body))
        resp = self.conn.getresponse()
        data = resp.read()
        if not 200 <= resp.status <= 299:
            errmsg = "Failed to create unit: {} {} - {}".format(
                resp.status, resp.reason, data)
            raise RuntimeError(errmsg)
        return data

    def _delete_unit(self, name):
        headers = {'Content-Type': 'application/json'}
        self.conn.request('DELETE', '/v1-alpha/units/{name}.service'.format(**locals()),
                          headers=headers)
        resp = self.conn.getresponse()
        data = resp.read()
        if resp.status not in (404, 204):
            errmsg = "Failed to delete unit: {} {} - {}".format(
                resp.status, resp.reason, data)
            raise RuntimeError(errmsg)
        return data

    def _get_state(self, name=None):
        headers = {'Content-Type': 'application/json'}
        url = '/v1-alpha/state'
        if name:
            url += '?unitName={name}.service'.format(**locals())
        self.conn.request('GET', url, headers=headers)
        resp = self.conn.getresponse()
        data = resp.read()
        if resp.status not in (200,):
            errmsg = "Failed to retrieve state: {} {} - {}".format(
                resp.status, resp.reason, data)
            raise RuntimeError(errmsg)
        return json.loads(data)

    def _get_machines(self):
        headers = {'Content-Type': 'application/json'}
        url = '/v1-alpha/machines'
        self.conn.request('GET', url, headers=headers)
        resp = self.conn.getresponse()
        data = resp.read()
        if resp.status not in (200,):
            errmsg = "Failed to retrieve machines: {} {} - {}".format(
                resp.status, resp.reason, data)
            raise RuntimeError(errmsg)
        return json.loads(data)

    # container api

    def create(self, name, image, command='', template=None, **kwargs):
        """Create a container"""
        self._create_container(name, image, command,
                               template or copy.deepcopy(CONTAINER_TEMPLATE), **kwargs)

    def _create_container(self, name, image, command, unit, **kwargs):
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        # prepare memory limit for the container type
        mem = kwargs.get('memory', {}).get(l['c_type'], None)
        if mem:
            l.update({'memory': '-m {}'.format(mem.lower())})
        else:
            l.update({'memory': ''})
        # prepare memory limit for the container type
        cpu = kwargs.get('cpu', {}).get(l['c_type'], None)
        if cpu:
            l.update({'cpu': '-c {}'.format(cpu)})
        else:
            l.update({'cpu': ''})
        # set unit hostname
        l.update({'hostname': self._get_hostname(name)})
        # should a special entrypoint be used
        entrypoint = kwargs.get('entrypoint')
        if entrypoint:
            l.update({'entrypoint': '{}'.format(entrypoint)})
        # encode command as utf-8
        if isinstance(l.get('command'), basestring):
            l['command'] = l['command'].encode('utf-8')
        # construct unit from template
        for f in unit:
            f['value'] = f['value'].format(**l)
        # prepare tags only if one was provided
        tags = kwargs.get('tags', {})
        if tags:
            tagset = ' '.join(['"{}={}"'.format(k, v) for k, v in tags.items()])
            unit.append({"section": "X-Fleet", "name": "MachineMetadata",
                         "value": tagset})
        # post unit to fleet and retry
        for attempt in range(RETRIES):
            try:
                self._put_unit(name, {"desiredState": "launched", "options": unit})
                break
            except:
                if attempt == (RETRIES - 1):  # account for 0 indexing
                    raise

    def _get_hostname(self, application_name):
        hostname = settings.UNIT_HOSTNAME
        if hostname == "default":
            return ''
        elif hostname == "application":
            # replace underscore with dots, since underscore is not valid in DNS hostnames
            dns_name = application_name.replace("_", ".")
            return '-h ' + dns_name
        elif hostname == "server":
            return '-h %H'
        else:
            raise RuntimeError('Unsupported hostname: ' + hostname)

    def start(self, name):
        """Start a container"""
        self._wait_for_container(name)

    def _wait_for_container(self, name):
        failures = 0
        # we bump to 20 minutes here to match the timeout on the router and in the app unit files
        for _ in range(1200):
            states = self._get_state(name)
            if states and len(states.get('states', [])) == 1:
                state = states.get('states')[0]
                subState = state.get('systemdSubState')
                if subState == 'running' or subState == 'exited':
                    break
                elif subState == 'failed':
                    # FIXME: fleet unit state reports failed when containers are fine
                    failures += 1
                    if failures == 10:
                        raise RuntimeError('container failed to start')
            time.sleep(1)
        else:
            raise RuntimeError('container timeout on start')

    def _wait_for_destroy(self, name):
        for _ in range(30):
            states = self._get_state(name)
            if not states:
                break
            time.sleep(1)
        else:
            raise RuntimeError('timeout on container destroy')

    def stop(self, name):
        """Stop a container"""
        raise NotImplementedError

    def destroy(self, name):
        """Destroy a container"""
        # call all destroy functions, ignoring any errors
        try:
            self._destroy_container(name)
        except:
            pass
        self._wait_for_destroy(name)

    def _destroy_container(self, name):
        for attempt in range(RETRIES):
            try:
                self._delete_unit(name)
                break
            except:
                if attempt == (RETRIES - 1):  # account for 0 indexing
                    raise

    def run(self, name, image, entrypoint, command):  # noqa
        """Run a one-off command"""
        self._create_container(name, image, command, copy.deepcopy(RUN_TEMPLATE),
                               entrypoint=entrypoint)

        # wait for the container to get scheduled
        for _ in range(30):
            states = self._get_state(name)
            if states and len(states.get('states', [])) == 1:
                state = states.get('states')[0]
                break
            time.sleep(1)
        else:
            raise RuntimeError('container did not report state')

        try:
            machineID = state.get('machineID')

            # find the machine
            machines = self._get_machines()
            if not machines:
                raise RuntimeError('no available hosts to run command')

            # find the machine's primaryIP
            primaryIP = None
            for m in machines.get('machines', []):
                if m['id'] == machineID:
                    primaryIP = m['primaryIP']
            if not primaryIP:
                raise RuntimeError('could not find host')

            # prepare ssh key
            file_obj = cStringIO.StringIO(base64.b64decode(self.pkey))
            pkey = paramiko.RSAKey(file_obj=file_obj)

            # grab output via docker logs over SSH
            ssh = paramiko.SSHClient()
            ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            ssh.connect(primaryIP, username="core", pkey=pkey)
            # share a transport
            tran = ssh.get_transport()

            def _do_ssh(cmd):
                chan = tran.open_session()
                # get a pty so stdout/stderr look right
                chan.get_pty()
                out = chan.makefile()
                chan.exec_command(cmd)
                rc, output = chan.recv_exit_status(), out.read()
                return rc, output

            # wait for container to launch
            # we loop indefinitely here, as we have no idea how long the docker pull will take
            while True:
                rc, _ = _do_ssh('docker inspect {name}'.format(**locals()))
                if rc == 0:
                    break
                time.sleep(1)
            else:
                raise RuntimeError('failed to create container')

            # wait for container to start
            for _ in range(2):
                _rc, _output = _do_ssh('docker inspect {name}'.format(**locals()))
                if _rc != 0:
                    raise RuntimeError('failed to inspect container')
                _container = json.loads(_output)
                started_at = _container[0]["State"]["StartedAt"]
                if not started_at.startswith('0001'):
                    break
                time.sleep(1)
            else:
                raise RuntimeError('container failed to start')

            # wait for container to complete
            for _ in range(1200):
                _rc, _output = _do_ssh('docker inspect {name}'.format(**locals()))
                if _rc != 0:
                    raise RuntimeError('failed to inspect container')
                _container = json.loads(_output)
                finished_at = _container[0]["State"]["FinishedAt"]
                if not finished_at.startswith('0001'):
                    break
                time.sleep(1)
            else:
                raise RuntimeError('container timed out')

            # gather container output
            _rc, output = _do_ssh('docker logs {name}'.format(**locals()))
            if _rc != 0:
                raise RuntimeError('could not attach to container')

            # determine container exit code
            _rc, _output = _do_ssh('docker inspect {name}'.format(**locals()))
            if _rc != 0:
                raise RuntimeError('could not determine exit code')
            container = json.loads(_output)
            rc = container[0]["State"]["ExitCode"]

        finally:
            # cleanup
            self._destroy_container(name)
            self._wait_for_destroy(name)

        # return rc and output
        return rc, output

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        raise NotImplementedError

SchedulerClient = FleetHTTPClient


CONTAINER_TEMPLATE = [
    {"section": "Unit", "name": "Description", "value": "{name}"},
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; docker pull $IMAGE"'''},  # noqa
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "docker inspect {name} >/dev/null 2>&1 && docker rm -f {name} || true"'''},  # noqa
    {"section": "Service", "name": "ExecStart", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; port=$(docker inspect -f '{{{{range $k, $v := .ContainerConfig.ExposedPorts }}}}{{{{$k}}}}{{{{end}}}}' $IMAGE | cut -d/ -f1) ; docker run --name {name} {memory} {cpu} {hostname} -P -e PORT=$port $IMAGE {command}"'''},  # noqa
    {"section": "Service", "name": "ExecStop", "value": '''/usr/bin/docker rm -f {name}'''},
    {"section": "Service", "name": "TimeoutStartSec", "value": "20m"},
    {"section": "Service", "name": "RestartSec", "value": "5"},
    {"section": "Service", "name": "Restart", "value": "on-failure"},
]


RUN_TEMPLATE = [
    {"section": "Unit", "name": "Description", "value": "{name} admin command"},
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; docker pull $IMAGE"'''},  # noqa
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "docker inspect {name} >/dev/null 2>&1 && docker rm -f {name} || true"'''},  # noqa
    {"section": "Service", "name": "ExecStart", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; docker run --name {name} --entrypoint={entrypoint} -a stdout -a stderr $IMAGE {command}"'''},  # noqa
    {"section": "Service", "name": "TimeoutStartSec", "value": "20m"},
]
