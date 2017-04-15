from __future__ import unicode_literals

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from json_field.fields import JSONField

from base import UuidAuditedModel


class ProviderManager(models.Manager):
    """Manage database interactions for :class:`Provider`."""

    def seed(self, user, **kwargs):
        """
        Seeds the database with Providers for clouds supported by Deis.
        """
        providers = [(p, p) for p in settings.PROVIDER_MODULES]
        for p_id, p_type in providers:
            self.create(owner=user, id=p_id, type=p_type, creds='{}')


@python_2_unicode_compatible
class Provider(UuidAuditedModel):
    """Cloud provider settings for a user.

    Available as `user.provider_set`.
    """

    objects = ProviderManager()

    PROVIDERS = (
        ('ec2', 'Amazon Elastic Compute Cloud (EC2)'),
        ('mock', 'Mock Reference Provider'),
        ('rackspace', 'Rackspace Open Cloud'),
        ('static', 'Static Node'),
        ('digitalocean', 'Digital Ocean'),
        ('vagrant', 'Local Vagrant VMs'),
    )

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    type = models.SlugField(max_length=16, choices=PROVIDERS)
    creds = JSONField(blank=True)

    class Meta:
        app_label = 'api'
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return "{}-{}".format(self.id, self.get_type_display())
