from cStringIO import StringIO
import base64
import os
import random
import re
import subprocess


ROOT_DIR = os.path.join(os.getcwd(), 'coreos')
if not os.path.exists(ROOT_DIR):
    os.mkdir(ROOT_DIR)

MATCH = re.compile('(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z]+)?.(?P<c_num>[0-9]+)')

class FleetClient(object):

    def __init__(self, cluster_name, hosts, auth, domain, options):
        self.name = cluster_name
        self.hosts = hosts
        self.domain = domain
        self.options = options
        self.auth = auth
        self.auth_path = os.path.join(ROOT_DIR, 'ssh-{cluster_name}'.format(**locals()))
        with open(self.auth_path, 'w') as f:
            f.write(base64.b64decode(auth))
            os.chmod(self.auth_path, 0600)

        self.env = {
            'PATH': '/usr/local/bin:/usr/bin:/bin:{}'.format(
                os.path.abspath(os.path.join(__file__, '..'))),
            'FLEETW_KEY': self.auth_path,
            'FLEETW_HOST': random.choice(self.hosts.split(','))}

    # scheduler setup / teardown

    def setUp(self):
        """
        Setup a CoreOS cluster including router and log aggregator
        """
        #print 'Creating deis-router.1'
        #self._create_router('deis-router.1', 'deis/router', command='', port=80)
        #print 'Creating deis-logger.1'
        #self._create_logger('deis-logger.1', 'deis/logger', command='', port=514)
        return

    def _create_router(self, name, image, command, port):
        env = self.env.copy()
        self._create_container(name, image, command, ROUTER_TEMPLATE, env, port)
        self._create_announcer(name, image, command, ANNOUNCE_TEMPLATE, env, port)
        self._start_container(name, env)
        self._start_announcer(name, env)

    def _create_logger(self, name, image, command, port):
        env = self.env.copy()
        self._create_container(name, image, command, LOGGER_TEMPLATE, env, port)
        self._create_announcer(name, image, command, ANNOUNCE_TEMPLATE, env, port)
        self._start_container(name, env)
        self._start_announcer(name, env)

    def _destroy_router(self, name, env):
        env = self.env.copy()
        self._destroy_container(name, env)
        self._destroy_announcer(name, env)

    def _destroy_logger(self, name, env):
        self._destroy_container(name, env)
        self._destroy_announcer(name, env)

    def tearDown(self):
        """
        Tear down a CoreOS cluster including router and log aggregator
        """
        #env = self.env.copy()
        #print 'Destroying deis-router.1'
        #self._destroy_router('deis-router.1', env)
        #print 'Destroying deis-logger.1'
        #self._destroy_logger('deis-logger.1', env)
        return

    # job api

    def create(self, name, image, command='', template=None, port=5000):
        """
        Create a new job
        """
        print 'Creating {name}'.format(**locals())
        env = self.env.copy()
        self._create_container(name, image, command, template or CONTAINER_TEMPLATE, env, port)
        self._create_log(name, image, command, LOG_TEMPLATE, env, port)
        self._create_announcer(name, image, command, ANNOUNCE_TEMPLATE, env, port)

    # TODO: remove hardcoded ports

    def _create_container(self, name, image, command, template, env, port):
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        env.update({'FLEETW_UNIT': name + '.service'})
        env.update({'FLEETW_UNIT_DATA': base64.b64encode(template.format(**l))})
        return subprocess.check_call('fleetctl.sh submit {name}.service'.format(**l),
                                     shell=True, env=env)

    def _create_announcer(self, name, image, command, template, env, port):
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        env.update({'FLEETW_UNIT': name + '-announce' + '.service'})
        env.update({'FLEETW_UNIT_DATA': base64.b64encode(template.format(**l))})
        return subprocess.check_call('fleetctl.sh submit {name}-announce.service'.format(**l),  # noqa
                                     shell=True, env=env)

    def _create_log(self, name, image, command, template, env, port):
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        env.update({'FLEETW_UNIT': name + '-log' + '.service'})
        env.update({'FLEETW_UNIT_DATA': base64.b64encode(template.format(**l))})
        return subprocess.check_call('fleetctl.sh submit {name}-log.service'.format(**locals()),  # noqa
                                     shell=True, env=env)

    def start(self, name):
        """
        Start an idle job
        """
        print 'Starting {name}'.format(**locals())
        env = self.env.copy()
        self._start_log(name, env)
        self._start_container(name, env)
        self._start_announcer(name, env)

    def _start_container(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh start {name}.service'.format(**locals()),
            shell=True, env=env)

    def _start_announcer(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh start {name}-announce.service'.format(**locals()),
            shell=True, env=env)

    def _start_log(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh start {name}-log.service'.format(**locals()),
            shell=True, env=env)

    def stop(self, name):
        """
        Stop a running job
        """
        print 'Stopping {name}'.format(**locals())
        env = self.env.copy()
        self._stop_announcer(name, env)
        self._stop_container(name, env)
        self._stop_log(name, env)

    def _stop_container(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh stop {name}.service'.format(**locals()),
            shell=True, env=env)

    def _stop_announcer(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh stop {name}-announce.service'.format(**locals()),
            shell=True, env=env)

    def _stop_log(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh stop {name}-announce.service'.format(**locals()),
            shell=True, env=env)

    def destroy(self, name):
        """
        Destroy an existing job
        """
        print 'Destroying {name}'.format(**locals())
        env = self.env.copy()
        self._destroy_announcer(name, env)
        self._destroy_container(name, env)
        self._destroy_log(name, env)

    def _destroy_container(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh destroy {name}.service'.format(**locals()),
            shell=True, env=env)

    def _destroy_announcer(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh destroy {name}-announce.service'.format(**locals()),
            shell=True, env=env)

    def _destroy_log(self, name, env):
        return subprocess.check_call(
            'fleetctl.sh destroy {name}-log.service'.format(**locals()),
            shell=True, env=env)

    def run(self, name, image, command):
        """
        Run a one-off command
        """
        print 'Running {name}'.format(**locals())
        output = subprocess.PIPE
        p = subprocess.Popen('fleetrun.sh {command}'.format(**locals()), shell=True, env=self.env,
                             stdout=output, stderr=subprocess.STDOUT)
        rc = p.wait()
        return rc, p.stdout.read()

    def run(self, name, image, command):
        """
        Run a one-off command
        """
        print 'Running {name}'.format(**locals())
        output = subprocess.PIPE
        p = subprocess.Popen('fleetrun.sh {command}'.format(**locals()), shell=True, env=self.env,
                             stdout=output, stderr=subprocess.STDOUT)
        rc = p.wait()
        return rc, p.stdout.read()

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        return StringIO(), StringIO(), StringIO()

SchedulerClient = FleetClient


CONTAINER_TEMPLATE = """
[Unit]
Description={name}
After=docker.service
Requires=docker.service

[Service]
ExecStartPre=/usr/bin/docker pull {image}
ExecStart=-/usr/bin/docker run --name {name} -P -e PORT={port} -e ETCD=172.17.42.1:4001 {image} {command}
ExecStop=-/usr/bin/docker rm -f {name}
"""

ANNOUNCE_TEMPLATE = """
[Unit]
Description={name} announce
BindsTo={name}.service

[Service]
ExecStartPre=/bin/sh -c "until /usr/bin/docker port {name} {port} >/dev/null 2>&1; do sleep 2; done; port=$(docker port {name} {port} | cut -d ':' -f2); host=$(getent hosts deis | awk {{'print $1'}}); echo Waiting for $port/tcp...; until cat </dev/null>/dev/tcp/$host/$port; do sleep 1; done"
ExecStart=/bin/sh -c "port=$(docker port {name} {port} | cut -d ':' -f2); host=$(getent hosts deis | awk {{'print $1'}}); echo Connected to $host:$port/tcp, publishing to etcd...; while netstat -lnt | grep $port >/dev/null; do etcdctl set /deis/services/{app}/{name} $host:$port --ttl 60 >/dev/null; sleep 45; done"
ExecStop=/usr/bin/etcdctl rm --recursive /deis/services/{app}/{name}

[X-Fleet]
X-ConditionMachineOf={name}.service
"""

LOG_TEMPLATE = """
[Unit]
Description={name} log
BindsTo={name}.service

[Service]
ExecStartPre=/bin/sh -c "until /usr/bin/docker inspect {name} >/dev/null 2>&1; do sleep 1; done"
ExecStart=/bin/sh -c "/usr/bin/docker logs -f {name} 2>&1 | logger -p local0.info -t {app}[{c_type}.{c_num}] --tcp --server $(etcdctl get /deis/services/deis-logger/deis-logger.1 | cut -d ':' -f1) --port $(etcdctl get /deis/services/deis-logger/deis-logger.1 | cut -d ':' -f2)"

[X-Fleet]
X-ConditionMachineOf={name}.service
"""

ROUTER_TEMPLATE = """
[Unit]
Description={name} router
After=docker.service
Requires=docker.service

[Service]
ExecStartPre=/usr/bin/docker pull {image}
ExecStart=-/usr/bin/docker run --name {name} -p 80:80 -p 443:443 -e ETCD=172.17.42.1:4001 {image} {command}
ExecStop=-/usr/bin/docker rm -f {name}
TimeoutStartSec=10min
"""

LOGGER_TEMPLATE = """
[Unit]
Description={name} logger
After=docker.service
Requires=docker.service

[Service]
ExecStartPre=/usr/bin/docker pull {image}
ExecStart=-/usr/bin/docker run --name {name} -p 514:514 -e HOST=%H -e PORT=514 -e ETCD=172.17.42.1:4001 {image} {command}
ExecStop=-/usr/bin/docker rm -f {name}
"""
