class FaultyClient(object):
    """A faulty scheduler that will always fail"""

    def __init__(self, cluster_name, hosts, auth, domain, options):
        pass

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def create(self, name, image, command='', template=None, port=5000):
        raise Exception()

    def start(self, name):
        raise Exception()

    def stop(self, name):
        raise Exception()

    def destroy(self, name):
        raise Exception()

    def run(self, name, image, command):
        raise Exception()

    def attach(self, name):
        raise Exception()

SchedulerClient = FaultyClient
