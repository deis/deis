import random

from .mock import MockSchedulerClient, jobs
from .states import JobState


CREATE_ERROR_RATE = 0
DESTROY_ERROR_RATE = 0
START_ERROR_RATE = 0
STOP_ERROR_RATE = 0


class ChaosSchedulerClient(MockSchedulerClient):

    def create(self, name, image, command, **kwargs):
        """Create a new container."""
        if random.random() < CREATE_ERROR_RATE:
            jobs.setdefault(name, {})['state'] = JobState.error
        else:
            super(ChaosSchedulerClient, self).create(name, image, command, **kwargs)

    def destroy(self, name):
        """Destroy a container."""
        if random.random() < DESTROY_ERROR_RATE:
            jobs.setdefault(name, {})['state'] = JobState.error
        else:
            super(ChaosSchedulerClient, self).destroy(name)

    def run(self, name, image, entrypoint, command):
        """Run a one-off command."""
        if random.random() < CREATE_ERROR_RATE:
            raise RuntimeError('exit code 1')
        else:
            super(ChaosSchedulerClient, self).run(name, image, entrypoint, command)

    def start(self, name):
        """Start a container."""
        if random.random() < START_ERROR_RATE:
            jobs.setdefault(name, {})['state'] = JobState.crashed
        else:
            super(ChaosSchedulerClient, self).start(name)

    def stop(self, name):
        """Stop a container."""
        if random.random() < STOP_ERROR_RATE:
            jobs.setdefault(name, {})['state'] = JobState.crashed
        else:
            super(ChaosSchedulerClient, self).stop(name)

SchedulerClient = ChaosSchedulerClient
