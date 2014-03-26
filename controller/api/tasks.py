"""
Long-running tasks for the Deis Controller API

This module orchestrates the real "heavy lifting" of Deis, and as such these
functions are decorated to run as asynchronous celery tasks.
"""

from __future__ import unicode_literals

from celery import task


@task
def deploy_release(app, release):
    containers = app.container_set.all()
    # TODO: parallelize
    for c in containers:
        try:
            c.deploy(release)
        except Exception:
            c.state = 'error'
            c.save()
            raise


@task
def start_containers(containers):
    # TODO: parallelize
    for c in containers:
        try:
            c.create()
            c.start()
        except Exception:
            c.state = 'error'
            c.save()
            raise


@task
def stop_containers(containers):
    # TODO: parallelize
    for c in containers:
        try:
            c.destroy()
            c.delete()
        except Exception:
            c.state = 'error'
            c.save()
            raise
