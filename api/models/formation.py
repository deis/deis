from __future__ import unicode_literals

from celery.canvas import group
from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from json_field.fields import JSONField

from api import tasks

from api.models import CM
from base import UuidAuditedModel


@python_2_unicode_compatible
class Formation(UuidAuditedModel):
    """
    Formation of nodes used to host applications
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64, unique=True)
    domain = models.CharField(max_length=128, blank=True, null=True)
    nodes = JSONField(default='{}', blank=True)

    class Meta:
        app_label = 'api'
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'domain': self.domain,
                'nodes': self.nodes}

    def build(self):
        return

    def destroy(self, *args, **kwargs):
        for app in self.app_set.all():
            app.destroy()
        node_tasks = [tasks.destroy_node.si(n) for n in self.node_set.all()]
        layer_tasks = [tasks.destroy_layer.si(l) for l in self.layer_set.all()]
        group(node_tasks).apply_async().join()
        group(layer_tasks).apply_async().join()
        CM.purge_formation(self.flat())
        self.delete()
        tasks.converge_controller.apply_async().wait()

    def publish(self):
        data = self.calculate()
        CM.publish_formation(self.flat(), data)
        return data

    def converge(self, controller=False, **kwargs):
        databag = self.publish()
        nodes = self.node_set.all()
        subtasks = []
        for n in nodes:
            subtask = tasks.converge_node.si(n)
            subtasks.append(subtask)
        if controller is True:
            subtasks.append(tasks.converge_controller.si())
        group(*subtasks).apply_async().join()
        return databag

    def calculate(self):
        """Return a representation of this formation for config management"""
        d = {}
        d['id'] = self.id
        d['domain'] = self.domain
        d['nodes'] = {}
        proxies = []
        for n in self.node_set.all():
            d['nodes'][n.id] = {'fqdn': n.fqdn,
                                'runtime': n.layer.runtime,
                                'proxy': n.layer.proxy}
            if n.layer.proxy is True:
                proxies.append(n.fqdn)
        d['apps'] = {}
        for a in self.app_set.all():
            d['apps'][a.id] = a.calculate()
            d['apps'][a.id]['proxy'] = {}
            d['apps'][a.id]['proxy']['nodes'] = proxies
            d['apps'][a.id]['proxy']['algorithm'] = 'round_robin'
            d['apps'][a.id]['proxy']['port'] = 80
            d['apps'][a.id]['proxy']['backends'] = []
            d['apps'][a.id]['containers'] = containers = {}
            for c in a.container_set.all().order_by('created'):
                containers.setdefault(c.type, {})
                containers[c.type].update(
                    {c.num: "{0}:{1}".format(c.node.id, c.port)})
                if c.type == 'web':
                    d['apps'][a.id]['proxy']['backends'].append(
                        "{0}:{1}".format(c.node.fqdn, c.port))
        return d
