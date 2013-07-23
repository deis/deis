"""
Deis API custom fields for representing data in Django forms.
"""

from __future__ import unicode_literals
from uuid import uuid4

from django import forms
from django.db import models
from json_field import JSONField
from yamlfield.fields import YAMLField


class UuidField(models.CharField):

    """A univerally unique ID field."""
    # pylint: disable=R0904

    description = __doc__

    def __init__(self, *args, **kwargs):
        kwargs.setdefault('auto_created', True)
        kwargs.setdefault('editable', False)
        kwargs.setdefault('max_length', 32)
        kwargs.setdefault('unique', True)
        super(UuidField, self).__init__(*args, **kwargs)

    def db_type(self, connection=None):
        """Return the database type for a UuidField."""
        db_type = None
        if connection and 'postgres' in connection.vendor:
            db_type = 'uuid'
        else:
            db_type = 'char({0})'.format(self.max_length)
        return db_type

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


class EnvVarsField(JSONField):

    """
    A text field that accepts a JSON object, coercing its keys to uppercase.
    """
    pass


class DataBagField(JSONField):
    """
    A text field that accepts a JSON object, used for storing Chef data bags.
    """
    pass


class ProcfileField(JSONField):
    """
    A text field that accepts a JSON object, used for Procfile data.
    """
    pass


class CredentialsField(JSONField):
    """
    A text field that accepts a JSON object, used for storing provider
    API Credentials.
    """
    pass


class ParamsField(JSONField):
    """
    A text field that accepts a JSON object, used for storing provider
    API Parameters.
    """

class CloudInitField(YAMLField):
    """
    A text field that accepts a YAML object, used for storing cloud-init
    boostrapping scripts.
    """
    pass


class NodeStatusField(JSONField):
    """
    A text field that accepts a YAML object, used for storing cloud-init
    boostrapping scripts.
    """
    pass


try:
    from south.modelsinspector import add_introspection_rules
    # Tell the South schema migration tool to handle a UuidField.
    add_introspection_rules([], [r'^api\.fields\.UuidField'])
    add_introspection_rules([], [r'^api\.fields\.EnvVarsField'])
    add_introspection_rules([], [r'^api\.fields\.DataBagField'])
    add_introspection_rules([], [r'^api\.fields\.ProcfileField'])
    add_introspection_rules([], [r'^api\.fields\.CredentialsField'])
    add_introspection_rules([], [r'^api\.fields\.ParamsField'])
    add_introspection_rules([], [r'^api\.fields\.CloudInitField'])
    add_introspection_rules([], [r'^api\.fields\.NodeStatusField'])
except ImportError:
    pass
