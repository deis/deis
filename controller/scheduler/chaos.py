import random

CREATE_ERROR_RATE = 0
DESTROY_ERROR_RATE = 0
START_ERROR_RATE = 0
STOP_ERROR_RATE = 0


class ChaosSchedulerClient(object):

    def __init__(self, target, auth, options, pkey):
        self.target = target
        self.auth = auth
        self.options = options
        self.pkey = pkey

    # job api

    def create(self, name, image, command, **kwargs):
        if random.random() < CREATE_ERROR_RATE:
            raise RuntimeError
        return True

    def start(self, name):
        """
        Start an idle job
        """
        if random.random() < START_ERROR_RATE:
            raise RuntimeError
        return True

    def stop(self, name):
        """
        Stop a running job
        """
        if random.random() < STOP_ERROR_RATE:
            raise RuntimeError
        return True

    def destroy(self, name):
        """
        Destroy an existing job
        """
        if random.random() < DESTROY_ERROR_RATE:
            raise RuntimeError
        return True

    def run(self, name, image, entrypoint, command):
        """
        Run a one-off command
        """
        if random.random() < CREATE_ERROR_RATE:
            raise RuntimeError('exit code 1')
        return 0, ''

SchedulerClient = ChaosSchedulerClient
