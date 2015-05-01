import re
import time
from docker import Client
from django.conf import settings
from .states import JobState

MATCH = re.compile(
    '(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z-_]+)?.(?P<c_num>[0-9]+)')


class SwarmClient(object):
    def __init__(self, target, auth, options, pkey):
        self.target = settings.SWARM_HOST
        # single global connection
        self.registry = settings.REGISTRY_HOST+":"+settings.REGISTRY_PORT
        self.docker_cli = Client(base_url='tcp://'+self.target+':'+"2395",
                                 timeout=1200, version='1.17')

    def create(self, name, image, command='', template=None, **kwargs):
        """Create a container"""
        cimage = self.registry+"/"+image
        cname = name
        ccommand = command
        affinity = "affinity:container!=~/"+(re.split("_v\d.", cname)[0]) + "*/"
        l = locals().copy()
        l.update(re.match(MATCH, name).groupdict())
        mem = kwargs.get('memory', {}).get(l['c_type'], None)
        if mem:
            mem = mem.lower()
            if mem[-2:-1].isalpha() and mem[-1].isalpha():
                mem = mem[:-1]
        cpu = kwargs.get('cpu', {}).get(l['c_type'], None)
        self.docker_cli.create_container(image=cimage, name=cname,
                                         command=ccommand.encode('utf-8'), mem_limit=mem,
                                         cpu_shares=cpu,
                                         environment=[affinity])
        self.docker_cli.stop(name)

    def start(self, name):
        """
        Start a container
        """
        self.docker_cli.start(name, publish_all_ports=True)

    def stop(self, name):
        """
        Stop a container
        """
        self.docker_cli.stop(name)

    def destroy(self, name):
        """
        Destroy a container
        """
        self.docker_cli.stop(name)
        self.docker_cli.remove_container(name)

    def run(self, name, image, entrypoint, command):
        """
        Run a one-off command
        """
        cimage = self.registry+"/"+image
        cname = name
        ccommand = command
        affinity = "affinity:container!=~/"+(re.split("_v\d.", cname)[0]) + "*/"
        self.docker_cli.create_container(image=cimage, name=cname,
                                         command=ccommand.encode('utf-8'),
                                         environment=[affinity],
                                         entrypoint=[entrypoint])
        time.sleep(2)
        self.docker_cli.start(cname)
        rc = 0
        while (True):
            if self._get_container_state(name) == JobState.created:
                break
            time.sleep(1)
        try:
            output = self.docker_cli.logs(name)
            return rc, output
        except Exception:
            rc = 1
            return rc, output

    def _get_container_state(self, name):
        try:
            if self.docker_cli.inspect_container(name)["State"]["Running"]:
                return JobState.up
            else:
                return JobState.created
        except Exception:
            return JobState.destroyed

    def state(self, name):
        try:
            for _ in range(30):
                return self._get_container_state(name)
                time.sleep(1)
            # FIXME (smothiki): should be able to send JobState.crashed
        except KeyError:
            return JobState.error
        except RuntimeError:
            return JobState.destroyed

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        raise NotImplementedError

    def _get_hostname(self, application_name):
        hostname = settings.UNIT_HOSTNAME
        if hostname == "default":
            return ''
        elif hostname == "application":
            # replace underscore with dots, since underscore is not valid in DNS hostnames
            dns_name = application_name.replace("_", ".")
            return dns_name
        elif hostname == "server":
            raise NotImplementedError
        else:
            raise RuntimeError('Unsupported hostname: ' + hostname)

    def _get_portbindings(self, image):
        dictports = self.docker_cli.inspect_image(image)["ContainerConfig"]["ExposedPorts"]
        for port, mapping in dictports.items():
            dictports[port] = None
        return dictports

    def _get_ports(self, image):
        ports = []
        dictports = self.docker_cli.inspect_image(image)["ContainerConfig"]["ExposedPorts"]
        for port, mapping in dictports.items():
            ports.append(int(port.split('/')[0]))
        return ports

SchedulerClient = SwarmClient
