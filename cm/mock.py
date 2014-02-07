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
    if 'error' in node.get('fqdn'):
        raise RuntimeError('Node Bootstrap Error:\nmock testing')


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
    if command.endswith('ls -al'):
        output, rc = """\
total 48
drwxr-xr-x  8 deis deis 4096 Dec 21 10:08 .
drwxr-xr-x  5 root root 4096 Dec 21 10:00 ..
-rw-------  1 deis deis  164 Dec 22 12:39 .bash_history
-rw-r--r--  1 deis deis  220 Mar 28  2013 .bash_logout
-rw-r--r--  1 deis deis 3486 Mar 28  2013 .bashrc
drwxr-xr-x  5 deis deis 4096 Dec 21 10:05 build
drwx------  2 deis deis 4096 Dec 21 10:08 .cache
drwx------  2 deis deis 4096 Dec 21 10:00 .chef
drwxr-xr-x 15 deis deis 4096 Dec 21 10:09 controller
-rw-r--r--  1 root root    0 Dec 21 10:00 prevent-apt-update
-rw-r--r--  1 deis deis  675 Mar 28  2013 .profile
drwxr-xr-x  2 deis deis 4096 Dec 21 10:07 .ssh
""", 0
    else:
        output, rc = 'failed to setup the container', 1
    return output, rc


def purge_node(node):
    """
    Purge a node and its client from Chef configuration management.

    This is a no-op for the mock provider.

    :param node: a dict containing the id of a node to purge
    """
    pass


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
