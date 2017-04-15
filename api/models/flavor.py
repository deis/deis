from __future__ import unicode_literals
import logging

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from json_field.fields import JSONField

import provider

from base import UuidAuditedModel
import providers

logger = logging.getLogger(__name__)


class FlavorManager(models.Manager):
    """Manage database interactions for :class:`Flavor`\s."""

    def seed(self, user, **kwargs):
        """Seed the database with default Flavors for each cloud region."""
        for provider_type in settings.PROVIDER_MODULES:
            provider_module = provider.import_provider_module(provider_type)
            flavors = provider_module.seed_flavors()
            p = providers.Provider.objects.get(owner=user, id=provider_type)
            for flavor in flavors:
                flavor['provider'] = p
                Flavor.objects.create(owner=user, **flavor)


@python_2_unicode_compatible
class Flavor(UuidAuditedModel):
    """
    Virtual machine flavors associated with a Provider

    Params is a JSON field including unstructured data
    for provider API calls, like region, zone, and size.
    """
    objects = FlavorManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    provider = models.ForeignKey(providers.Provider)
    params = JSONField(blank=True)

    class Meta:
        app_label = 'api'
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return self.id
