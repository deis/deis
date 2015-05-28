import re
import time
from django.conf import settings
from marathon import MarathonClient
from marathon.models import MarathonApp
from .states import JobState
from docker import Client
from .fleet import FleetHTTPClient

# turn down standard marathon logging

MATCH = re.compile(
    '(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z-_]+)?.(?P<c_num>[0-9]+)')
RETRIES = 3


class MarathonHTTPClient(object):

    def __init__(self, target, auth, options, pkey):
        self.target = settings.MARATHON_HOST
        self.auth = auth
        self.options = options
        self.pkey = pkey
        self.registry = settings.REGISTRY_HOST + ':' + settings.REGISTRY_PORT
        self.client = MarathonClient('http://'+self.target+':8180')
        self.fleet = FleetHTTPClient('/var/run/fleet.sock', auth, options, pkey)

    # helpers
    def _app_id(self, name):
        return name.replace('_', '.')

    # container api
    def create(self, name, image, command='', **kwargs):
        """Create a container"""
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
        cpu = kwargs.get('cpu', {}).get(l['c_type'])
        self.client.create_app(app_id,
                               MarathonApp(cmd="docker run --name "+name+" -P "+image+" "+command,
                                           mem=m, cpus=c))
        self.client.scale_app(app_id, 0, force=True)
        for _ in xrange(30):
            if self.client.get_app(self._app_id(name)).tasks_running == 0:
                return
            time.sleep(1)

    def start(self, name):
        """Start a container"""
        self.client.scale_app(self._app_id(name), 1, force=True)
        for _ in xrange(30):
            if self.client.get_app(self._app_id(name)).tasks_running == 1:
                return
            time.sleep(1)
        raise RuntimeError("App Not Started")

    def stop(self, name):
        """Stop a container"""
        raise NotImplementedError

    def destroy(self, name):
        """Destroy a container"""
        try:
            host = self.client.get_app(self._app_id(name)).tasks[0].host
            self.client.delete_app(self._app_id(name), force=True)
            self._delete_container(host, name)
        except:
            self.client.delete_app(self._app_id(name), force=True)

    def _delete_container(self, host, name):
        docker_cli = Client("tcp://{}:2375".format(host), timeout=1200, version='1.17')
        try:
            if docker_cli.inspect_container(name)['State']:
                docker_cli.remove_container(name, force=True)
        except:
            pass

    def run(self, name, image, entrypoint, command):  # noqa
        """Run a one-off command"""
        return self.fleet.run(name, image, entrypoint, command)

    def state(self, name):
        try:
            for _ in xrange(30):
                if self.client.get_app(self._app_id(name)).tasks_running == 1:
                    return JobState.up
                elif self.client.get_app(self._app_id(name)).tasks_running == 0:
                    return JobState.created
                time.sleep(1)
        except:
            return JobState.destroyed

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        raise NotImplementedError

SchedulerClient = MarathonHTTPClient
