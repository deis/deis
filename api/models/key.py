from __future__ import unicode_literals

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible

from base import UuidAuditedModel


@python_2_unicode_compatible
class Key(UuidAuditedModel):
    """An SSH public key."""

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=128)
    public = models.TextField(unique=True)

    class Meta:
        app_label = 'api'
        verbose_name = 'SSH Key'
        unique_together = (('owner', 'id'))

    def __str__(self):
        return "{}...{}".format(self.public[:18], self.public[-31:])

    def save(self, *args, **kwargs):
        super(Key, self).save(*args, **kwargs)
        self.owner.publish()

    def delete(self, *args, **kwargs):
        super(Key, self).delete(*args, **kwargs)
        self.owner.publish()
