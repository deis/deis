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
    try:
        [t.start() for t in threads]
        [t.join() for t in threads]
    except Exception:
        for c in containers:
            c.state = 'error'
            c.save()
        raise


@task
def start_containers(containers):
    create_threads = []
    start_threads = []
    for c in containers:
        create_threads.append(threading.Thread(target=c.create))
        start_threads.append(threading.Thread(target=c.start))
    try:
        [t.start() for t in create_threads]
        [t.join() for t in create_threads]
        [t.start() for t in start_threads]
        [t.join() for t in start_threads]
    except Exception:
        for c in containers:
            c.state = 'error'
            c.save()
            raise


@task
def stop_containers(containers):
    destroy_threads = []
    delete_threads = []
    for c in containers:
        destroy_threads.append(threading.Thread(target=c.destroy))
        delete_threads.append(threading.Thread(target=c.delete))
    try:
        [t.start() for t in destroy_threads]
        [t.join() for t in destroy_threads]
        [t.start() for t in delete_threads]
        [t.join() for t in delete_threads]
    except Exception:
        for c in containers:
            c.state = 'error'
            c.save()
            raise
