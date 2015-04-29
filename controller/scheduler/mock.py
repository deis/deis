import json
from cStringIO import StringIO
from .states import JobState, TransitionError


# HACK: MockSchedulerClient is not persistent across requests
jobs = {}


class MockSchedulerClient(object):

    def __init__(self, target, auth, options, pkey):
        self.target = target
        self.auth = auth
        self.options = options
        self.pkey = pkey

    # container api

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        return StringIO(), StringIO(), StringIO()

    def create(self, name, image, command, **kwargs):
        """
        Create a new container
        """
        job = jobs.get(name, {})
        job.update({'state': JobState.created})
        jobs[name] = job
        return

    def destroy(self, name):
        """
        Destroy a container
        """
        job = jobs.get(name, {})
        job.update({'state': JobState.destroyed})
        jobs[name] = job
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

    def start(self, name):
        """
        Start a container
        """
        if self.state(name) not in [JobState.created,
                                    JobState.up,
                                    JobState.down,
                                    JobState.crashed,
                                    JobState.error]:
            raise TransitionError(self.state(name),
                                  JobState.up,
                                  'the container must be stopped or up to start')
        job = jobs.get(name, {})
        job.update({'state': JobState.up})
        jobs[name] = job
        return

    def state(self, name):
        """
        Display the given job's running state
        """
        state = JobState.initialized
        job = jobs.get(name)
        if job:
            state = job.get('state')
        return state

    def stop(self, name):
        """
        Stop a container
        """
        job = jobs.get(name, {})
        if job.get('state') not in [JobState.up, JobState.crashed, JobState.error]:
            raise TransitionError(job.get('state'),
                                  JobState.up,
                                  'the container must be up to stop')
        job.update({'state': JobState.down})
        jobs[name] = job
        return

SchedulerClient = MockSchedulerClient
