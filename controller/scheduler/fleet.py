import cStringIO
import base64
import copy
import httplib
import json
import paramiko
import socket
import re
import time

from django.conf import settings

from .states import JobState


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

    def _request_unit(self, method, name, body=None):
        headers = {'Content-Type': 'application/json'}
        self.conn.request(method, '/v1-alpha/units/{name}.service'.format(**locals()),
                                  headers=headers, body=json.dumps(body))
        return self.conn.getresponse()

    def _get_unit(self, name):
        for attempt in xrange(RETRIES):
            try:
                resp = self._request_unit('GET', name)
                data = resp.read()
                if not 200 <= resp.status <= 299:
                    errmsg = "Failed to retrieve unit: {} {} - {}".format(
                        resp.status, resp.reason, data)
                    raise RuntimeError(errmsg)
                return data
            except:
                if attempt >= (RETRIES - 1):
                    raise

    def _put_unit(self, name, body):
        for attempt in xrange(RETRIES):
            try:
                resp = self._request_unit('PUT', name, body)
                data = resp.read()
                if not 200 <= resp.status <= 299:
                    errmsg = "Failed to create unit: {} {} - {}".format(
                        resp.status, resp.reason, data)
                    raise RuntimeError(errmsg)
                return data
            except:
                if attempt >= (RETRIES - 1):
                    raise

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
            tagset = ' '.join(['"{}={}"'.format(k, v) for k, v in tags.viewitems()])
            unit.append({"section": "X-Fleet", "name": "MachineMetadata",
                         "value": tagset})
        # post unit to fleet
        self._put_unit(name, {"desiredState": "loaded", "options": unit})

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
        self._put_unit(name, {'desiredState': 'launched'})
        self._wait_for_container_running(name)

    def _wait_for_container_state(self, name):
        # wait for container to get scheduled
        for _ in xrange(30):
            states = self._get_state(name)
            if states and len(states.get('states', [])) == 1:
                return states.get('states')[0]
            time.sleep(1)
        else:
            raise RuntimeError('container timeout while retrieving state')

    def _wait_for_container_running(self, name):
        # we bump to 20 minutes here to match the timeout on the router and in the app unit files
        try:
            self._wait_for_job_state(name, JobState.up)
        except RuntimeError:
            raise RuntimeError('container failed to start')

    def _wait_for_job_state(self, name, state):
        # we bump to 20 minutes here to match the timeout on the router and in the app unit files
        for _ in xrange(1200):
            if self.state(name) == state:
                return
            time.sleep(1)
        else:
            raise RuntimeError('timeout waiting for job state: {}'.format(state))

    def _wait_for_destroy(self, name):
        for _ in xrange(30):
            if not self._get_state(name):
                break
            time.sleep(1)
        else:
            raise RuntimeError('timeout on container destroy')

    def stop(self, name):
        """Stop a container"""
        self._put_unit(name, {"desiredState": "loaded"})
        self._wait_for_job_state(name, JobState.created)

    def destroy(self, name):
        """Destroy a container"""
        # call all destroy functions, ignoring any errors
        try:
            self._destroy_container(name)
        except:
            pass
        self._wait_for_destroy(name)

    def _destroy_container(self, name):
        for attempt in xrange(RETRIES):
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
        # launch the container
        self._put_unit(name, {'desiredState': 'launched'})
        # wait for the container to get scheduled
        state = self._wait_for_container_state(name)

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
                with tran.open_session() as chan:
                    chan.exec_command(cmd)
                    while not chan.exit_status_ready():
                        time.sleep(1)
                    out = chan.makefile()
                    output = out.read()
                    rc = chan.recv_exit_status()
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
            for _ in xrange(2):
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
            for _ in xrange(1200):
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

    def state(self, name):
        systemdActiveStateMap = {
            "active": "up",
            "reloading": "down",
            "inactive": "created",
            "failed": "crashed",
            "activating": "down",
            "deactivating": "down",
        }
        try:
            # NOTE (bacongobbler): this call to ._get_unit() acts as a pre-emptive check to
            # determine if the job no longer exists (will raise a RuntimeError on 404)
            self._get_unit(name)
            state = self._wait_for_container_state(name)
            activeState = state['systemdActiveState']
            # FIXME (bacongobbler): when fleet loads a job, sometimes it'll automatically start and
            # stop the container, which in our case will return as 'failed', even though
            # the container is perfectly fine.
            if activeState == 'failed':
                if state['systemdLoadState'] == 'loaded':
                    return JobState.created
            return getattr(JobState, systemdActiveStateMap[activeState])
        except KeyError:
            # failed retrieving a proper response from the fleet API
            return JobState.error
        except RuntimeError:
            # failed to retrieve a response from the fleet API,
            # which means it does not exist
            return JobState.destroyed

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
    {"section": "Service", "name": "ExecStart", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; docker run --name {name} --rm {memory} {cpu} {hostname} $IMAGE {command}"'''},  # noqa
    {"section": "Service", "name": "ExecStop", "value": '''/usr/bin/docker stop {name}'''},
    {"section": "Service", "name": "ExecStop", "value": '''/usr/bin/docker rm -f {name}'''},
    {"section": "Service", "name": "TimeoutStartSec", "value": "20m"},
    {"section": "Service", "name": "TimeoutStopSec", "value": "10"},
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
