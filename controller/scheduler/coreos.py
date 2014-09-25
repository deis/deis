import cStringIO
import base64
import copy
import functools
import json
import httplib
import paramiko
import socket
import re
import time


MATCH = re.compile(
    '(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z]+)?.(?P<c_num>[0-9]+)')


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

    def __init__(self, cluster_name, hosts, auth, domain, options):
        self.name = cluster_name
        self.hosts = hosts
        self.auth = auth
        self.domain = domain
        self.options = options
        # single global connection
        self.conn = UHTTPConnection('/var/run/fleet.sock')

    # scheduler setup / teardown

    def setUp(self):
        pass

    def tearDown(self):
        pass

    # connection helpers

    def _put_unit(self, name, body):
        headers = {'Content-Type': 'application/json'}
        self.conn.request('PUT', '/v1-alpha/units/{name}.service'.format(**locals()),
                          headers=headers, body=json.dumps(body))
        resp = self.conn.getresponse()
        data = resp.read()
        if resp.status != 204:
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

    def create(self, name, image, command='', template=None, use_announcer=True, **kwargs):
        """Create a container"""
        self._create_container(name, image, command,
                               template or copy.deepcopy(CONTAINER_TEMPLATE), **kwargs)
        self._create_log(name, image, command, copy.deepcopy(LOG_TEMPLATE))

        if use_announcer:
            self._create_announcer(name, image, command, copy.deepcopy(ANNOUNCE_TEMPLATE))

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
        # construct unit from template
        for f in unit:
            f['value'] = f['value'].format(**l)
        # prepare tags only if one was provided
        tags = kwargs.get('tags', {})
        if tags:
            tagset = ' '.join(['"{}={}"'.format(k, v) for k, v in tags.items()])
            unit.append({"section": "X-Fleet", "name": "MachineMetadata",
                         "value": tagset})
        # post unit to fleet
        self._put_unit(name, {"desiredState": "launched", "options": unit})

    def _create_log(self, name, image, command, unit):
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        # construct unit from template
        for f in unit:
            f['value'] = f['value'].format(**l)
        # post unit to fleet
        self._put_unit(name+'-log', {"desiredState": "launched", "options": unit})

    def _create_announcer(self, name, image, command, unit):
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        # construct unit from template
        for f in unit:
            f['value'] = f['value'].format(**l)
        # post unit to fleet
        self._put_unit(name+'-announce', {"desiredState": "launched", "options": unit})

    def start(self, name, use_announcer=True):
        """Start a container"""
        self._wait_for_container(name)

        if use_announcer:
            self._wait_for_announcer(name)

    def _wait_for_container(self, name):
        # we bump to 20 minutes here to match the timeout on the router and in the app unit files
        for _ in range(1200):
            states = self._get_state(name)
            if states and len(states.get('states', [])) == 1:
                state = states.get('states')[0]
                subState = state.get('systemdSubState')
                if subState == 'running' or subState == 'exited':
                    break
                elif subState == 'failed':
                    raise RuntimeError('container failed to start')
            time.sleep(1)
        else:
            raise RuntimeError('container failed to start')

    def _wait_for_announcer(self, name):
        # wait a bit for the announcer to come up, otherwise we may have hit
        # https://github.com/docker/docker/issues/8022
        for _ in range(30):
            states = self._get_state(name)
            if states and len(states.get('states', [])) == 1:
                state = states.get('states')[0]
                subState = state.get('systemdSubState')
                if subState == 'running':
                    # wait for the router to be reconfigured
                    time.sleep(10)
                    break
                elif subState == 'failed':
                    raise RuntimeError('announcer failed to start')
            time.sleep(1)
        else:
            raise RuntimeError('announcer timeout on start')

    def _wait_for_destroy(self, name):
        for _ in range(30):
            states = self._get_state(name)
            if not states:
                break
            time.sleep(1)
        else:
            raise RuntimeError('timeout on container destroy')

    def stop(self, name, use_announcer=True):
        """Stop a container"""
        raise NotImplementedError

    def destroy(self, name, use_announcer=True):
        """Destroy a container"""
        funcs = []
        if use_announcer:
            funcs.append(functools.partial(self._destroy_announcer, name))
        funcs.append(functools.partial(self._destroy_container, name))
        funcs.append(functools.partial(self._destroy_log, name))
        # call all destroy functions, ignoring any errors
        for f in funcs:
            try:
                f()
            except:
                pass
        self._wait_for_destroy(name)

    def _destroy_container(self, name):
        return self._delete_unit(name)

    def _destroy_announcer(self, name):
        return self._delete_unit(name+'-announce')

    def _destroy_log(self, name):
        return self._delete_unit(name+'-log')

    def run(self, name, image, command):
        """Run a one-off command"""
        self._create_container(name, image, command, copy.deepcopy(RUN_TEMPLATE))

        # wait for the container to return something
        for _ in range(1200):
            states = self._get_state(name)
            if states and len(states.get('states', [])) == 1:
                state = states.get('states')[0]
                subState = state.get('systemdSubState')
                if subState == 'exited' or subState == 'failed' or subState == 'dead':
                    break
            time.sleep(1)
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
        file_obj = cStringIO.StringIO(base64.b64decode(self.auth))
        pkey = paramiko.RSAKey(file_obj=file_obj)

        # grab output via docker logs over SSH
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(primaryIP, username="core", pkey=pkey)

        # get a pty so stdout/stderr look right
        tran = ssh.get_transport()
        chan = tran.open_session()
        chan.get_pty()
        out = chan.makefile()

        # exec the command to gather container output
        chan.exec_command('docker logs {name}'.format(**locals()))
        rc, output = chan.recv_exit_status(), out.read()
        if rc != 0:
            raise RuntimeError('could not attach to container')

        # use another channel to inspect the container
        chan = tran.open_session()
        chan.get_pty()
        out = chan.makefile()
        chan.exec_command('docker inspect {name}'.format(**locals()))
        rc, inspect_output = chan.recv_exit_status(), out.read()
        if rc != 0:
            raise RuntimeError('could not determine exit code')
        container = json.loads(inspect_output)
        rc = container[0]["State"]["ExitCode"]

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
    {"section": "Service", "name": "ExecStart", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; port=$(docker inspect -f '{{{{range $k, $v := .ContainerConfig.ExposedPorts }}}}{{{{$k}}}}{{{{end}}}}' $IMAGE | cut -d/ -f1) ; docker run --name {name} {memory} {cpu} -P -e PORT=$port $IMAGE {command}"'''},  # noqa
    {"section": "Service", "name": "ExecStop", "value": '''/usr/bin/docker rm -f {name}'''},
    {"section": "Service", "name": "TimeoutStartSec", "value": "20m"},
]


LOG_TEMPLATE = [
    {"section": "Unit", "name": "Description", "value": "{name} log"},
    {"section": "Unit", "name": "BindsTo", "value": "{name}.service"},
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "until docker inspect {name} >/dev/null 2>&1; do sleep 1; done"'''},  # noqa
    {"section": "Service", "name": "ExecStart", "value": '''/bin/sh -c "docker logs -f {name} 2>&1 | logger -p local0.info -t {app}[{c_type}.{c_num}] --udp --server $(etcdctl get /deis/logs/host) --port $(etcdctl get /deis/logs/port)"'''},  # noqa
    {"section": "Service", "name": "TimeoutStartSec", "value": "20m"},
    {"section": "X-Fleet", "name": "MachineOf", "value": "{name}.service"},
]


ANNOUNCE_TEMPLATE = [
    {"section": "Unit", "name": "Description", "value": "{name} announce"},
    {"section": "Unit", "name": "BindsTo", "value": "{name}.service"},
    {"section": "Service", "name": "EnvironmentFile", "value": "/etc/environment"},
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "until docker inspect -f '{{{{range $i, $e := .NetworkSettings.Ports }}}}{{{{$p := index $e 0}}}}{{{{$p.HostPort}}}}{{{{end}}}}' {name} >/dev/null 2>&1; do sleep 2; done; port=$(docker inspect -f '{{{{range $i, $e := .NetworkSettings.Ports }}}}{{{{$p := index $e 0}}}}{{{{$p.HostPort}}}}{{{{end}}}}' {name}); if [[ -z $port ]]; then echo We have no port...; exit 1; fi; echo Waiting for $port/tcp...; until netstat -lnt | grep :$port >/dev/null; do sleep 1; done"'''},  # noqa
    {"section": "Service", "name": "ExecStart", "value": '''/bin/sh -c "port=$(docker inspect -f '{{{{range $i, $e := .NetworkSettings.Ports }}}}{{{{$p := index $e 0}}}}{{{{$p.HostPort}}}}{{{{end}}}}' {name}); echo Connected to $COREOS_PRIVATE_IPV4:$port/tcp, publishing to etcd...; while netstat -lnt | grep :$port >/dev/null; do etcdctl set /deis/services/{app}/{name} $COREOS_PRIVATE_IPV4:$port --ttl 60 >/dev/null; sleep 45; done"'''},  # noqa
    {"section": "Service", "name": "ExecStop", "value": "/usr/bin/etcdctl rm --recursive /deis/services/{app}/{name}"},  # noqa
    {"section": "Service", "name": "TimeoutStartSec", "value": "20m"},
    {"section": "X-Fleet", "name": "MachineOf", "value": "{name}.service"},
]


RUN_TEMPLATE = [
    {"section": "Unit", "name": "Description", "value": "{name} admin command"},
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; docker pull $IMAGE"'''},  # noqa
    {"section": "Service", "name": "ExecStartPre", "value": '''/bin/sh -c "docker inspect {name} >/dev/null 2>&1 && docker rm -f {name} || true"'''},  # noqa
    {"section": "Service", "name": "ExecStart", "value": '''/bin/sh -c "IMAGE=$(etcdctl get /deis/registry/host 2>&1):$(etcdctl get /deis/registry/port 2>&1)/{image}; docker run --name {name} --entrypoint=/bin/bash -a stdout -a stderr $IMAGE -c '{command}'"'''},  # noqa
    {"section": "Service", "name": "TimeoutStartSec", "value": "20m"},
]
