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
        """
        Setup a Cluster including router and log aggregator
        """
        return None

    def tearDown(self):
        """
        Tear down a cluster including router and log aggregator
        """
        return None

    # job api

    def create(self, name, image, command):
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
