from __future__ import unicode_literals

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from json_field.fields import JSONField

from api import tasks
from base import UuidAuditedModel
from formation import Formation
from flavor import Flavor


@python_2_unicode_compatible
class Layer(UuidAuditedModel):
    """
    Layer of nodes used by the formation

    All nodes in a layer share the same flavor and configuration.

    The layer stores SSH settings used to trigger node convergence,
    as well as other configuration used during node bootstrapping
    (e.g. Chef Run List, Chef Environment)
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)

    formation = models.ForeignKey(Formation)
    flavor = models.ForeignKey(Flavor)

    proxy = models.BooleanField(default=False)
    runtime = models.BooleanField(default=False)

    ssh_username = models.CharField(max_length=64, default='ubuntu')
    ssh_private_key = models.TextField()
    ssh_public_key = models.TextField()
    ssh_port = models.SmallIntegerField(default=22)

    # example: {'run_list': [deis::runtime'], 'environment': 'dev'}
    config = JSONField(default='{}', blank=True)

    class Meta:
        app_label = 'api'
        unique_together = (('formation', 'id'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'provider_type': self.flavor.provider.type,
                'creds': dict(self.flavor.provider.creds),
                'formation': self.formation.id,
                'flavor': self.flavor.id,
                'params': dict(self.flavor.params),
                'proxy': self.proxy,
                'runtime': self.runtime,
                'ssh_username': self.ssh_username,
                'ssh_private_key': self.ssh_private_key,
                'ssh_public_key': self.ssh_public_key,
                'ssh_port': self.ssh_port,
                'config': dict(self.config)}

    def build(self):
        return tasks.build_layer.delay(self).wait()

    def destroy(self):
        return tasks.destroy_layer.delay(self).wait()
