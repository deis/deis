from django.db import models

from api import fields


class AuditedModel(models.Model):
    """Add created and updated fields to a model."""

    created = models.DateTimeField(auto_now_add=True)
    updated = models.DateTimeField(auto_now=True)

    class Meta:
        """Mark :class:`AuditedModel` as abstract."""
        abstract = True


class UuidAuditedModel(AuditedModel):
    """Add a UUID primary key to an :class:`AuditedModel`."""

    uuid = fields.UuidField('UUID', primary_key=True)

    class Meta:
        """Mark :class:`UuidAuditedModel` as abstract."""
        abstract = True
