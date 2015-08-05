import re
import time

from django.conf import settings
from docker import Client
from marathon import MarathonClient
from marathon.models import MarathonApp

from . import AbstractSchedulerClient
from .fleet import FleetHTTPClient
from .states import JobState

# turn down standard marathon logging

MATCH = re.compile(
    '(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z-_]+)?.(?P<c_num>[0-9]+)')
RETRIES = 3
POLL_ATTEMPTS = 30
POLL_WAIT = 100


class MarathonHTTPClient(AbstractSchedulerClient):

    def __init__(self, target, auth, options, pkey):
        super(MarathonHTTPClient, self).__init__(target, auth, options, pkey)
        self.target = settings.MARATHON_HOST
        self.registry = settings.REGISTRY_HOST + ':' + settings.REGISTRY_PORT
        self.client = MarathonClient('http://'+self.target+':8180')
        self.fleet = FleetHTTPClient('/var/run/fleet.sock', auth, options, pkey)

    # helpers
    def _app_id(self, name):
        return name.replace('_', '.')

    # container api
    def create(self, name, image, command='', **kwargs):
        """Create a new container"""
        app_id = self._app_id(name)
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        image = self.registry + '/' + image
        mems = kwargs.get('memory', {}).get(l['c_type'])
        m = 0
        if mems:
            mems = mems.lower()
            if mems[-2:-1].isalpha() and mems[-1].isalpha():
                mems = mems[:-1]
            m = int(mems[:-1])
        c = 0.5
        cpu = kwargs.get('cpu', {}).get(l['c_type'])
        if cpu:
            c = cpu
        cmd = "docker run --name {name} -P {image} {command}".format(**locals())
        self.client.create_app(app_id, MarathonApp(cmd=cmd, mem=m, cpus=c, instances=0))
        for _ in xrange(POLL_ATTEMPTS):
            if self.client.get_app(self._app_id(name)).tasks_running == 0:
                return
            time.sleep(1)

    def start(self, name):
        """Start a container."""
        self.client.scale_app(self._app_id(name), 1, force=True)
        for _ in xrange(POLL_ATTEMPTS):
            if self.client.get_app(self._app_id(name)).tasks_running == 1:
                break
            time.sleep(1)
        host = self.client.get_app(self._app_id(name)).tasks[0].host
        self._waitforcontainer(host, name)

    def destroy(self, name):
        """Destroy a container."""
        try:
            host = self.client.get_app(self._app_id(name)).tasks[0].host
            self.client.delete_app(self._app_id(name), force=True)
            self._delete_container(host, name)
        except:
            self.client.delete_app(self._app_id(name), force=True)

    def _get_container_state(self, host, name):
        docker_cli = Client("tcp://{}:2375".format(host), timeout=1200, version='1.17')
        try:
            if docker_cli.inspect_container(name)['State']['Running']:
                return JobState.up
        except:
            return JobState.destroyed

    def _waitforcontainer(self, host, name):
        for _ in xrange(POLL_WAIT):
            if self._get_container_state(host, name) == JobState.up:
                return
            time.sleep(1)
        raise RuntimeError("App container Not Started")

    def _delete_container(self, host, name):
        docker_cli = Client("tcp://{}:2375".format(host), timeout=1200, version='1.17')
        if docker_cli.inspect_container(name)['State']:
            docker_cli.remove_container(name, force=True)

    def run(self, name, image, entrypoint, command):  # noqa
        """Run a one-off command."""
        return self.fleet.run(name, image, entrypoint, command)

    def state(self, name):
        """Display the given job's running state."""
        try:
            for _ in xrange(POLL_ATTEMPTS):
                if self.client.get_app(self._app_id(name)).tasks_running == 1:
                    return JobState.up
                elif self.client.get_app(self._app_id(name)).tasks_running == 0:
                    return JobState.created
                time.sleep(1)
        except:
            return JobState.destroyed

SchedulerClient = MarathonHTTPClient
