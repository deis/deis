"""
Long-running tasks for the Deis Controller API

This module orchestrates the real "heavy lifting" of Deis, and as such these
functions are decorated to run as asynchronous celery tasks.
"""

from __future__ import unicode_literals

import threading

from celery import task


@task
def create_cluster(cluster):
    cluster._scheduler.setUp()


@task
def destroy_cluster(cluster):
    for app in cluster.app_set.all():
        app.destroy()
    cluster._scheduler.tearDown()


@task
def deploy_release(app, release):
    containers = app.container_set.all()
    threads = []
    for c in containers:
        threads.append(threading.Thread(target=c.deploy, args=(release,)))
    [t.start() for t in threads]
    [t.join() for t in threads]


@task
def start_containers(containers):
    create_threads = []
    start_threads = []
    for c in containers:
        create_threads.append(threading.Thread(target=c.create))
        start_threads.append(threading.Thread(target=c.start))
    [t.start() for t in create_threads]
    [t.join() for t in create_threads]
    [t.start() for t in start_threads]
    [t.join() for t in start_threads]


@task
def stop_containers(containers):
    destroy_threads = []
    delete_threads = []
    for c in containers:
        destroy_threads.append(threading.Thread(target=c.destroy))
        delete_threads.append(threading.Thread(target=c.delete))
    [t.start() for t in destroy_threads]
    [t.join() for t in destroy_threads]
    [t.start() for t in delete_threads]
    [t.join() for t in delete_threads]


@task
def run_command(c, command):
    release = c.release
    version = release.version
    image = release.image
    try:
        # pull the image first
        rc, pull_output = c.run("docker pull {image}".format(**locals()))
        if rc != 0:
            raise EnvironmentError('Could not pull image: {pull_image}'.format(**locals()))
        # run the command
        docker_args = ' '.join(['--entrypoint=/bin/bash',
                                '-a', 'stdout', '-a', 'stderr', '--rm', image])
        escaped_command = command.replace("'", "'\\''")
        command = r"docker run {docker_args} -c \'{escaped_command}\'".format(**locals())
        return c.run(command)
    finally:
        c.delete()
