import random

CREATE_ERROR_RATE = 0
DESTROY_ERROR_RATE = 0
START_ERROR_RATE = 0
STOP_ERROR_RATE = 0


class ChaosSchedulerClient(object):

    def __init__(self, name, hosts, auth, domain, options):
        self.name = name
        self.hosts = hosts
        self.auth = auth
        self.domain = domain
        self.options = options

    # scheduler setup / teardown

    def setUp(self):
        pass

    def tearDown(self):
        pass

    # job api

    def create(self, name, image, command, use_announcer, **kwargs):
        if random.random() < CREATE_ERROR_RATE:
            raise RuntimeError
        return True

    def start(self, name, use_announcer):
        """
        Start an idle job
        """
        if random.random() < START_ERROR_RATE:
            raise RuntimeError
        return True

    def stop(self, name, use_announcer):
        """
        Stop a running job
        """
        if random.random() < STOP_ERROR_RATE:
            raise RuntimeError
        return True

    def destroy(self, name, use_announcer):
        """
        Destroy an existing job
        """
        if random.random() < DESTROY_ERROR_RATE:
            raise RuntimeError
        return True

    def run(self, name, image, command):
        """
        Run a one-off command
        """
        if random.random() < CREATE_ERROR_RATE:
            raise RuntimeError
        return 0, ''

SchedulerClient = ChaosSchedulerClient
