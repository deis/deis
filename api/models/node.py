from __future__ import unicode_literals

from celery.canvas import group
from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from json_field.fields import JSONField

from api import tasks

from base import UuidAuditedModel
from formation import Formation
from layer import Layer


class NodeManager(models.Manager):

    def new(self, formation, layer, fqdn=None):
        existing_nodes = self.filter(formation=formation, layer=layer).order_by('-created')
        if existing_nodes:
            next_num = existing_nodes[0].num + 1
        else:
            next_num = 1
        node = self.create(owner=formation.owner,
                           formation=formation,
                           layer=layer,
                           num=next_num,
                           id="{0}-{1}-{2}".format(formation.id, layer.id, next_num),
                           fqdn=fqdn)
        return node

    def scale(self, formation, structure, **kwargs):
        """Scale layers up or down to match requested structure."""
        funcs = []
        changed = False
        for layer_id, requested in structure.items():
            layer = formation.layer_set.get(id=layer_id)
            nodes = list(layer.node_set.all().order_by('created'))
            diff = requested - len(nodes)
            if diff == 0:
                continue
            while diff < 0:
                node = nodes.pop(0)
                funcs.append(tasks.destroy_node.si(node))
                diff = requested - len(nodes)
                changed = True
            while diff > 0:
                node = self.new(formation, layer)
                nodes.append(node)
                funcs.append(tasks.build_node.si(node))
                diff = requested - len(nodes)
                changed = True
        # launch/terminate nodes in parallel
        if funcs:
            group(*funcs).apply_async().join()
        # always scale and balance every application
        if nodes:
            for app in formation.app_set.all():
                container.Container.objects.scale(app, app.containers)
                container.Container.objects.balance(formation)
        # save new structure now that scaling was successful
        formation.nodes.update(structure)
        formation.save()
        # force-converge nodes if there were new nodes or container rebalancing
        if changed:
            return formation.converge()
        return formation.calculate()

    def next_runtime_node(self, formation, container_type, reverse=False):
        count = []
        layers = formation.layer_set.filter(runtime=True)
        runtime_nodes = []
        for l in layers:
            runtime_nodes.extend(Node.objects.filter(
                formation=formation, layer=l).order_by('created'))
        container_map = {n: [] for n in runtime_nodes}
        containers = list(container.Container.objects.filter(
            formation=formation, type=container_type).order_by('created'))
        for c in containers:
            container_map[c.node].append(c)
        for n in container_map.keys():
            # (2, node3), (2, node2), (3, node1)
            count.append((len(container_map[n]), n))
        if not count:
            raise EnvironmentError('No nodes available for containers')
        count.sort()
        # reverse means order by greatest # of containers, otherwise fewest
        if reverse:
            count.reverse()
        return count[0][1]

    def next_runtime_port(self, formation):
        containers = container.Container.objects.filter(formation=formation).order_by('-port')
        if not containers:
            return 10001
        return containers[0].port + 1


@python_2_unicode_compatible
class Node(UuidAuditedModel):
    """
    Node used to host containers

    List of nodes available as `formation.nodes`
    """

    objects = NodeManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=64)
    formation = models.ForeignKey(Formation)
    layer = models.ForeignKey(Layer)
    num = models.PositiveIntegerField()

    # TODO: add celery beat tasks for monitoring node health
    status = models.CharField(max_length=64, default='up')

    provider_id = models.SlugField(max_length=64, blank=True, null=True)
    fqdn = models.CharField(max_length=256, blank=True, null=True)
    status = JSONField(blank=True, null=True)

    class Meta:
        app_label = 'api'
        unique_together = (('formation', 'id'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'provider_type': self.layer.flavor.provider.type,
                'formation': self.formation.id,
                'layer': self.layer.id,
                'creds': dict(self.layer.flavor.provider.creds),
                'params': dict(self.layer.flavor.params),
                'runtime': self.layer.runtime,
                'proxy': self.layer.proxy,
                'ssh_username': self.layer.ssh_username,
                'ssh_public_key': self.layer.ssh_public_key,
                'ssh_private_key': self.layer.ssh_private_key,
                'ssh_port': self.layer.ssh_port,
                'config': dict(self.layer.config),
                'provider_id': self.provider_id,
                'fqdn': self.fqdn}

    def build(self):
        return tasks.build_node.delay(self).wait()

    def destroy(self):
        return tasks.destroy_node.delay(self).wait()

    def converge(self):
        return tasks.converge_node.delay(self).wait()

    def run(self, command, **kwargs):
        return tasks.run_node.delay(self, command).wait()


import container
