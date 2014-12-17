
import json
from cStringIO import StringIO


class MockSchedulerClient(object):

    def __init__(self, target, auth, options, pkey):
        self.target = target
        self.auth = auth
        self.options = options
        self.pkey = pkey

    # container api

    def create(self, name, image, command, **kwargs):
        """
        Create a new container
        """
        return

    def start(self, name):
        """
        Start a container
        """
        return

    def stop(self, name):
        """
        Stop a container
        """
        return

    def destroy(self, name):
        """
        Destroy a container
        """
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

SchedulerClient = MockSchedulerClient
