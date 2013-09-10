
from __future__ import unicode_literals

import json
import os

from deis import settings


def configure_node(node):
    config = {}
    config['config'] = config
    config['node_id'] = node.id
    config['layer_id'] = node.layer.id
    return config


def bootstrap_node(node):
    return node


def converge_node(node):
    return '', 0


def destroy_node(node):
    return node


def destroy_formation(formation):
    return formation


def converge_controller():
    return None


def publish_user(username, data):
    path = os.path.join(settings.TEMPDIR, 'user-{}'.format(username))
    with open(path, 'w') as f:
        f.write(json.dumps(data))
    return username


def publish_app(app_id, data):
    path = os.path.join(settings.TEMPDIR, 'app-{}'.format(app_id))
    with open(path, 'w') as f:
        f.write(json.dumps(data))
    return app_id


def purge_app(app_id):
    path = os.path.join(settings.TEMPDIR, 'app-{}'.format(app_id))
    os.remove(path)
    return app_id


def publish_formation(formation_id, data):
    path = os.path.join(settings.TEMPDIR, 'formation-{}'.format(formation_id))
    with open(path, 'w') as f:
        f.write(json.dumps(data))
    return formation_id


def purge_formation(formation_id):
    path = os.path.join(settings.TEMPDIR, 'formation-{}'.format(formation_id))
    os.remove(path)
    return formation_id
