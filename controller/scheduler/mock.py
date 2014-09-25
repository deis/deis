from cStringIO import StringIO


class MockSchedulerClient(object):

    def __init__(self, name, hosts, auth, domain, options):
        self.name = name
        self.hosts = hosts
        self.auth = auth
        self.domain = domain
        self.options = options

    # scheduler setup / teardown

    def setUp(self):
        return None

    def tearDown(self):
        return None

    # container api

    def create(self, name, image, command, use_announcer, **kwargs):
        """
        Create a new container
        """
        return

    def start(self, name, use_announcer):
        """
        Start a container
        """
        return

    def stop(self, name, use_announcer):
        """
        Stop a container
        """
        return

    def destroy(self, name, use_announcer):
        """
        Destroy a container
        """
        return

    def run(self, name, image, command):
        """
        Run a one-off command
        """
        return 0, ''

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        return StringIO(), StringIO(), StringIO()

SchedulerClient = MockSchedulerClient
