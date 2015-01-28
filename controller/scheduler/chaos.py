import random
from .mock import MockSchedulerClient, jobs
from .states import JobState


CREATE_ERROR_RATE = 0
DESTROY_ERROR_RATE = 0
START_ERROR_RATE = 0
STOP_ERROR_RATE = 0


class ChaosSchedulerClient(MockSchedulerClient):

    def create(self, name, image, command, **kwargs):
        if random.random() < CREATE_ERROR_RATE:
            job = jobs.get(name, {})
            job.update({'state': JobState.error})
            jobs[name] = job
            return
        return super(ChaosSchedulerClient, self).create(name, image, command, **kwargs)

    def destroy(self, name):
        """
        Destroy an existing job
        """
        if random.random() < DESTROY_ERROR_RATE:
            job = jobs.get(name, {})
            job.update({'state': JobState.error})
            jobs[name] = job
            return
        return super(ChaosSchedulerClient, self).destroy(name)

    def run(self, name, image, entrypoint, command):
        """
        Run a one-off command
        """
        if random.random() < CREATE_ERROR_RATE:
            raise RuntimeError('exit code 1')
        return super(ChaosSchedulerClient, self).run(name, image, entrypoint, command)

    def start(self, name):
        """
        Start an idle job
        """
        if random.random() < START_ERROR_RATE:
            job = jobs.get(name, {})
            job.update({'state': JobState.crashed})
            jobs[name] = job
            return
        return super(ChaosSchedulerClient, self).start(name)

    def stop(self, name):
        """
        Stop a running job
        """
        if random.random() < STOP_ERROR_RATE:
            job = jobs.get(name, {})
            job.update({'state': JobState.error})
            jobs[name] = job
            return
        return super(ChaosSchedulerClient, self).stop(name)

SchedulerClient = ChaosSchedulerClient
