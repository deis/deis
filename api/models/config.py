from __future__ import unicode_literals
import logging

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from json_field.fields import JSONField

from base import UuidAuditedModel
from application import App

logger = logging.getLogger(__name__)


@python_2_unicode_compatible
class Config(UuidAuditedModel):
    """
    Set of configuration values applied as environment variables
    during runtime execution of the Application.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey(App)
    version = models.PositiveIntegerField()

    values = JSONField(default='{}', blank=True)

    class Meta:
        app_label = 'api'
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'version'),)

    def __str__(self):
        return "{0}-v{1}".format(self.app.id, self.version)
