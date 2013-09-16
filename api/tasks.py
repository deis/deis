
from __future__ import unicode_literals
import importlib

from celery import task
from celery.canvas import group

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
    return output, rc


@task
def build_formation(formation):
    return


@task
def destroy_formation(formation):
    app_tasks = [destroy_app.si(a) for a in formation.app_set.all()]
    node_tasks = [destroy_node.si(n) for n in formation.node_set.all()]
    layer_tasks = [destroy_layer.si(l) for l in formation.layer_set.all()]
    group(app_tasks + node_tasks).apply_async().join()
    group(layer_tasks).apply_async().join()
    CM.purge_formation(formation.flat())
    formation.delete()


@task
def converge_formation(formation):
    nodes = formation.node_set.all()
    subtasks = []
    for n in nodes:
        subtask = converge_node.si(n)
        subtasks.append(subtask)
    group(*subtasks).apply_async().join()


@task
def build_app(app):
    return


@task
def destroy_app(app):
    CM.purge_app(app.flat())
    app.delete()
    app.formation.publish()


@task
def converge_controller():
    CM.converge_controller()
    return None
