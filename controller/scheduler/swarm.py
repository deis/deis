import logging
from docker import Client

from django.conf import settings



class SwarmClient(object):
    def __init__(self,target, auth, options, pkey):
        self.target =  settings.SWARM_HOST
        # single global connection
        self.registry = settings.REGISTRY_HOST+":"+settings.REGISTRY_PORT
        self.docker_cli = Client(base_url='tcp://'+self.target+':'+"2395",timeout=1200)

    def create(self, name, image, command='', template=None, **kwargs):
        """Create a container"""
        cimage=self.registry+"/"+image
        cname=name
        ccommand=command
        # self.docker_cli.pull(cimage, stream=False,insecure_registry=True)
        self.docker_cli.create_container(image=cimage,name=cname,command=ccommand)#,hostname=self._get_hostname(cname),ports=self._get_ports(cimage))
        self.docker_cli.start(cname, port_bindings=self._get_portbindings(cimage),publish_all_ports=True)

    def start(self, name):
        """
        Start a container
        """
        self.docker_cli.start(name)
        return

    def stop(self, name):
        """
        Stop a container
        """
        self.docker_cli.stop(name)
        return
    def destroy(self, name):
        """
        Destroy a container
        """
        self.docker_cli.stop(name)
        self.docker_cli.remove_container(name)
        return

    def run(self, name, image, entrypoint, command):
        """
        Run a one-off command
        """
        # dump input into a json object for testing purposes
        return 0, json.dumps({'name': name,
                              'image': image,
                              'entrypoint': entrypoint,
                              'command': command})

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        return StringIO(), StringIO(), StringIO()

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

    def _get_portbindings(self,image):
        dictports=self.docker_cli.inspect_image(image)["ContainerConfig"]["ExposedPorts"]
        for port,mapping in dictports.items():
            dictports[port]=None
        return dictports

    def _get_ports(self,image):
        ports=[]
        dictports=self.docker_cli.inspect_image(image)["ContainerConfig"]["ExposedPorts"]
        for port,mapping in dictports.items():
            ports.append(int(port.split('/')[0]))
        return ports

SchedulerClient = SwarmClient
