
from __future__ import unicode_literals

import json
import os

from deis import settings


def bootstrap_node(node):
    return '', 0


def converge_node(node):
    return '', 0


def purge_node(node):
    return


def publish_user(user, data):
    path = os.path.join(settings.TEMPDIR, 'user-{username}'.format(**user))
    with open(path, 'w') as f:
        f.write(json.dumps(data))


def purge_user(user):
    path = os.path.join(settings.TEMPDIR, 'user-{username}'.format(**user))
    os.remove(path)


def publish_app(app, data):
    path = os.path.join(settings.TEMPDIR, 'app-{id}'.format(**app))
    with open(path, 'w') as f:
        f.write(json.dumps(data))


def purge_app(app):
    path = os.path.join(settings.TEMPDIR, 'app-{id}'.format(**app))
    os.remove(path)


def publish_formation(formation, data):
    path = os.path.join(settings.TEMPDIR, 'formation-{id}'.format(**formation))
    with open(path, 'w') as f:
        f.write(json.dumps(data))


def purge_formation(formation):
    path = os.path.join(settings.TEMPDIR, 'formation-{id}'.format(**formation))
    os.remove(path)


def converge_controller():
    return None
