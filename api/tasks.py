
from __future__ import unicode_literals
import importlib

from celery import task

from deis import settings
from provider import import_provider_module

# import user-defined config management module
CM = importlib.import_module(settings.CM_MODULE)


@task
def build_layer(layer):
    provider = import_provider_module(layer.flavor.provider.type)
    provider.build_layer(layer.flat())


@task
def destroy_layer(layer):
    provider = import_provider_module(layer.flavor.provider.type)
    provider.destroy_layer(layer.flat())
    layer.delete()


@task
def build_node(node):
    provider = import_provider_module(node.layer.flavor.provider.type)
    provider_id, fqdn, metadata = provider.build_node(node.flat())
    node.provider_id = provider_id
    node.fqdn = fqdn
    node.metadata = metadata
    node.save()
    CM.bootstrap_node(node.flat())


@task
def destroy_node(node):
    provider = import_provider_module(node.layer.flavor.provider.type)
    provider.destroy_node(node.flat())
    CM.purge_node(node.flat())
    node.delete()


@task
def converge_node(node):
    output, rc = CM.converge_node(node.flat())
    return output, rc


@task
def run_node(node, command):
    output, rc = CM.run_node(node.flat(), command)
    if rc != 0 and 'failed to setup the container' in output:
        output = '\033[35mPlease run `git push deis master` first.\033[0m\n' + output
    return output, rc


@task
def converge_controller():
    CM.converge_controller()
    return None
