"""
Deis API custom fields for representing data in Django forms.
"""

from __future__ import unicode_literals
from uuid import uuid4

from django import forms
from django.db import models


class UuidField(models.CharField):
    """A univerally unique ID field."""

    description = __doc__

    def __init__(self, *args, **kwargs):
        kwargs.setdefault('auto_created', True)
        kwargs.setdefault('editable', False)
        kwargs.setdefault('max_length', 32)
        kwargs.setdefault('unique', True)
        super(UuidField, self).__init__(*args, **kwargs)

    def db_type(self, connection=None):
        """Return the database column type for a UuidField."""
        if connection and 'postgres' in connection.vendor:
            return 'uuid'
        else:
            return "char({})".format(self.max_length)

    def pre_save(self, model_instance, add):
        """Initialize an empty field with a new UUID before it is saved."""
        value = getattr(model_instance, self.get_attname(), None)
        if not value and add:
            uuid = str(uuid4())
            setattr(model_instance, self.get_attname(), uuid)
            return uuid
        else:
            return super(UuidField, self).pre_save(model_instance, add)

    def formfield(self, **kwargs):
        """Tell forms how to represent this UuidField."""
        kwargs.update({
            'form_class': forms.CharField,
            'max_length': self.max_length,
        })
        return super(UuidField, self).formfield(**kwargs)


try:
    from south.modelsinspector import add_introspection_rules
    # Tell the South schema migration tool to handle our custom fields.
    add_introspection_rules([], [r'^api\.fields\.UuidField'])
except ImportError:  # pragma: no cover
    pass
