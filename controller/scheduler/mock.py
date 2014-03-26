from cStringIO import StringIO


class MockSchedulerClient(object):

    def __init__(self, name, hosts, auth):
        self.name = name
        self.hosts = hosts
        self.auth = auth

    def create(self, name, image, command, resources={}, constraints={}):
        """
        Create a new job
        """
        return {'state': 'inactive'}

    def start(self, name):
        """
        Start an idle job
        """
        return {'state': 'active'}

    def stop(self, name):
        """
        Stop a running job
        """
        return {'state': 'inactive'}

    def destroy(self, name):
        """
        Destroy an existing job
        """
        return {'state': 'inactive'}

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        return StringIO(), StringIO(), StringIO()

SchedulerClient = MockSchedulerClient
