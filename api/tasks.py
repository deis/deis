"""
Core Deis API functions that interact with providers and
configuration management.

This module orchestrates the real "heavy lifting" of Deis, and as such these
functions are decorated to run as asynchronous celery tasks.
"""

from __future__ import unicode_literals
import importlib

from celery import task

from deis import settings
from provider import import_provider_module
from .exceptions import BuildNodeError


# import user-defined config management module
CM = importlib.import_module(settings.CM_MODULE)


@task
def build_layer(layer):
    """
    Build a layer using its cloud provider.

    :param layer: a :class:`~api.models.Layer` to build
    """
    provider = import_provider_module(layer.flavor.provider.type)
    provider.build_layer(layer.flat())


@task
def destroy_layer(layer):
    """
    Destroy a layer.

    :param layer: a :class:`~api.models.Layer` to destroy
    """
    provider = import_provider_module(layer.flavor.provider.type)
    provider.destroy_layer(layer.flat())
    layer.delete()


@task
def build_node(node):
    """
    Build a node using its cloud provider.

    :param node: a :class:`~api.models.Node` to build
    """
    provider = import_provider_module(node.layer.flavor.provider.type)
    provider_id, fqdn, metadata = provider.build_node(node.flat())
    node.provider_id = provider_id
    node.fqdn = fqdn
    node.metadata = metadata
    node.save()
    try:
        CM.bootstrap_node(node.flat())
    except RuntimeError as err:
        raise BuildNodeError(str(err))


@task
def destroy_node(node):
    """
    Destroy a node.

    :param node: a :class:`~api.models.Node` to destroy
    """
    provider = import_provider_module(node.layer.flavor.provider.type)
    provider.destroy_node(node.flat())
    CM.purge_node(node.flat())
    node.delete()


@task
def converge_node(node):
    """
    Converge a node, aligning it with an intended configuration.

    :param node: a :class:`~api.models.Node` to converge
    """
    output, rc = CM.converge_node(node.flat())
    return output, rc


@task
def run_node(node, command):
    """
    Run a single shell command on a container on a node.

    Does not support interactive commands.

    :param node: a :class:`~api.models.Node` on which to run a command
    """
    output, rc = CM.run_node(node.flat(), command)
    if rc != 0 and 'failed to setup the container' in output:
        output = '\033[35mPlease run `git push deis master` first.\033[0m\n' + output
    return output, rc
