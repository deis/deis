
from __future__ import unicode_literals
import importlib

from celery import task
from celery.canvas import group
from django.contrib.auth.models import User

from api.models import Formation, Node, App, FlavorManager, Provider, Flavor, \
    Layer
from api.ssh import connect_ssh, exec_ssh
from deis import settings
from provider import import_provider_module

# now that we've defined models that may be imported by celery tasks
# import user-defined config management module
CM = importlib.import_module(settings.CM_MODULE)


@task
def seed_flavors(username):
    user = User.objects.get(username=username)
    cloud_config = FlavorManager.load_cloud_config_base()
    for provider_type in ('mock', 'ec2'):
        provider = import_provider_module(provider_type)
        flavors = provider.seed_flavors(username)
        p = Provider.objects.get(owner=user, id=provider_type)
        for flavor in flavors:
            flavor['provider'] = p
            flavor['init'] = cloud_config
            Flavor.objects.create(owner=user, **flavor)


@task
def build_layer(layer_id):
    layer = Layer.objects.get(id=layer_id)
    provider = import_provider_module(layer.flavor.provider.type)
    provider.build_layer(layer)
    return layer_id


@task
def destroy_layer(layer_id):
    layer = Layer.objects.get(id=layer_id)
    provider = import_provider_module(layer.flavor.provider.type)
    group([destroy_node.s(n.id) for n in layer.node_set.all()]).apply_async().join()
    provider.destroy_layer(layer)
    layer.delete()
    return layer_id


@task
def build_node(node_id):
    node = Node.objects.get(id=node_id)
    provider = import_provider_module(node.layer.flavor.provider.type)
    config = CM.configure_node(node)
    # add ssh keys for layer and owner
    config.setdefault(
        'ssh_authorized_keys', []).append(node.layer.ssh_public_key)
    config['ssh_authorized_keys'].extend(
        [k.public for k in node.formation.owner.key_set.all()])
    # build the node
    provider.build_node(node, config)
    # use CM to bootstrap the node
    CM.bootstrap_node(node)
    return node_id


@task
def destroy_node(node_id):
    node = Node.objects.get(id=node_id)
    provider = import_provider_module(node.layer.flavor.provider.type)
    CM.destroy_node(node)
    provider.destroy_node(node)
    node.delete()
    return node.id


@task
def converge_node(node_id):
    node = Node.objects.get(id=node_id)
    output, rc = CM.converge_node(node)
    return output, rc


@task
def run_node(app_id, command):
    app = App.objects.get(id=app_id)
    node = app.node_set.order_by('?')[0]
    release = app.release_set.order_by('-created')[0]
    # prepare ssh command
    version = release.version
    docker_args = ' '.join(
        ['-v',
         '/opt/deis/runtime/slugs/{app_id}-{version}/app:/app'.format(**locals()),
         release.image])
    base_cmd = "export HOME=/app; cd /app && for profile in " \
               "`find /app/.profile.d/*.sh -type f`; do . $profile; done"
    command = "/bin/sh -c '{base_cmd} && {command}'".format(**locals())
    command = "sudo docker run {docker_args} {command}".format(**locals())
    # connect and exec
    ssh = connect_ssh(node.layer.ssh_username,
                      node.fqdn, 22,
                      node.layer.ssh_private_key)
    output, rc = exec_ssh(ssh, command, pty=True)
    return output, rc


@task
def converge_formation(formation_id):
    formation = Formation.objects.get(id=formation_id)
    nodes = formation.node_set.all()
    subtasks = []
    for n in nodes:
        subtask = converge_node.s(n.id)
        subtasks.append(subtask)
    job = group(*subtasks)
    return job.apply_async().join()


@task
def destroy_formation(formation_id):
    formation = Formation.objects.get(id=formation_id)
    group([destroy_node.s(n.id) for n in formation.node_set.all()]).apply_async().join()
    group([destroy_layer.s(l.id) for l in formation.layer_set.all()]).apply_async().join()
    formation.delete()
    return formation_id
