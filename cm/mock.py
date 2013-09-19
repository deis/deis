"""
Deis mock configuration management implementation for testing.
"""

from __future__ import unicode_literals

import json
import os

from deis import settings


def bootstrap_node(node):
    """
    Bootstrap configuration management tools onto a node.

    This is a no-op for the mock provider.

    :param node: a dict containing the node's fully-qualified domain name and SSH info
    """
    pass


def converge_node(node):
    """
    Converge a node.

    "Converge" means to change a node's configuration to match that defined by
    configuration management.

    This is a no-op for the mock provider.

    :param node: a dict containing the node's fully-qualified domain name and SSH info
    :returns: a tuple of the convergence command's (output, return_code)
    """
    return '', 0


def run_node(node, command):
    """
    Run a command on a node.

    This is a no-op for the mock provider.

    :param node: a dict containing the node's fully-qualified domain name and SSH info
    :param command: the command-line to execute on the node
    :returns: a tuple of the command's (output, return_code)
    """
    return '', 0


def purge_node(node):
    """
    Purge a node and its client from Chef configuration management.

    This is a no-op for the mock provider.

    :param node: a dict containing the id of a node to purge
    """
    pass


def publish_user(user, data):
    """
    Publish a user to configuration management.

    :param user: a dict containing the username
    :param data: data to store with the user
    """
    path = os.path.join(settings.TEMPDIR, 'user-{username}'.format(**user))
    with open(path, 'w') as f:
        f.write(json.dumps(data))


def purge_user(user):
    """
    Purge a user from configuration management.

    :param user: a dict containing the username
    """
    path = os.path.join(settings.TEMPDIR, 'user-{username}'.format(**user))
    os.remove(path)


def publish_app(app, data):
    """
    Publish an app to configuration management.

    :param app: a dict containing the id of the app
    :param data: data to store with the app
    """
    path = os.path.join(settings.TEMPDIR, 'app-{id}'.format(**app))
    with open(path, 'w') as f:
        f.write(json.dumps(data))


def purge_app(app):
    """
    Purge an app from configuration management.

    :param user: a dict containing the id of the app
    """
    path = os.path.join(settings.TEMPDIR, 'app-{id}'.format(**app))
    os.remove(path)


def publish_formation(formation, data):
    """
    Publish a formation to configuration management.

    :param formation: a dict containing the id of the formation
    :param data: data to store with the formation
    """
    path = os.path.join(settings.TEMPDIR, 'formation-{id}'.format(**formation))
    with open(path, 'w') as f:
        f.write(json.dumps(data))


def purge_formation(formation):
    """
    Purge a formation from configuration management.

    :param formation: a dict containing the id of the formation
    """
    path = os.path.join(settings.TEMPDIR, 'formation-{id}'.format(**formation))
    os.remove(path)


def converge_controller():
    """
    Converge this controller node.

    "Converge" means to change a node's configuration to match that defined by
    configuration management.

    This is a no-op for the mock provider.

    :returns: the output of the convergence command, in this case None
    """
    return None
