"""
Classes to serialize the RESTful representation of Deis API models.
"""

from __future__ import unicode_literals

import json
import re

from django.conf import settings
from django.contrib.auth.models import User
from django.utils import timezone
from rest_framework import serializers
from rest_framework.validators import UniqueTogetherValidator

from api import models


PROCTYPE_MATCH = re.compile(r'^(?P<type>[a-z]+)')
MEMLIMIT_MATCH = re.compile(r'^(?P<mem>[0-9]+(MB|KB|GB|[BKMG]))$', re.IGNORECASE)
CPUSHARE_MATCH = re.compile(r'^(?P<cpu>[0-9]+)$')
TAGKEY_MATCH = re.compile(r'^[A-Za-z]+$')
TAGVAL_MATCH = re.compile(r'^\w+$')
CONFIGKEY_MATCH = re.compile(r'^[a-z_]+[a-z0-9_]*$', re.IGNORECASE)


class JSONFieldSerializer(serializers.Field):
    """
    A Django REST framework serializer for JSON data.
    """

    def to_representation(self, obj):
        """Serialize the field's JSON data, for read operations."""
        return obj

    def to_internal_value(self, data):
        """Deserialize the field's JSON data, for write operations."""
        try:
            val = json.loads(data)
        except TypeError:
            val = data
        return val


class JSONIntFieldSerializer(JSONFieldSerializer):
    """
    A JSON serializer that coerces its data to integers.
    """

    def to_internal_value(self, data):
        """Deserialize the field's JSON integer data."""
        field = super(JSONIntFieldSerializer, self).to_internal_value(data)

        for k, v in field.viewitems():
            if v is not None:  # NoneType is used to unset a value
                try:
                    field[k] = int(v)
                except ValueError:
                    field[k] = v
                    # Do nothing, the validator will catch this later
        return field


class JSONStringFieldSerializer(JSONFieldSerializer):
    """
    A JSON serializer that coerces its data to strings.
    """

    def to_internal_value(self, data):
        """Deserialize the field's JSON string data."""
        field = super(JSONStringFieldSerializer, self).to_internal_value(data)

        for k, v in field.viewitems():
            if v is not None:  # NoneType is used to unset a value
                field[k] = unicode(v)

        return field


class ModelSerializer(serializers.ModelSerializer):

    uuid = serializers.ReadOnlyField()

    def get_validators(self):
        """
        Hack to remove DRF's UniqueTogetherValidator when it concerns the UUID.

        See https://github.com/deis/deis/pull/2898#discussion_r23105147
        """
        validators = super(ModelSerializer, self).get_validators()
        for v in validators:
            if isinstance(v, UniqueTogetherValidator) and 'uuid' in v.fields:
                validators.remove(v)
        return validators


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ['email', 'username', 'password', 'first_name', 'last_name', 'is_superuser',
                  'is_staff', 'groups', 'user_permissions', 'last_login', 'date_joined',
                  'is_active']
        read_only_fields = ['is_superuser', 'is_staff', 'groups',
                            'user_permissions', 'last_login', 'date_joined', 'is_active']
        extra_kwargs = {'password': {'write_only': True}}

    def create(self, validated_data):
        now = timezone.now()
        user = User(
            email=validated_data.get('email'),
            username=validated_data.get('username'),
            last_login=now,
            date_joined=now,
            is_active=True
        )
        if validated_data.get('first_name'):
            user.first_name = validated_data['first_name']
        if validated_data.get('last_name'):
            user.last_name = validated_data['last_name']
        user.set_password(validated_data['password'])
        # Make the first signup an admin / superuser
        if not User.objects.filter(is_superuser=True).exists():
            user.is_superuser = user.is_staff = True
        user.save()
        return user


class AdminUserSerializer(serializers.ModelSerializer):
    """Serialize admin status for a User model."""

    class Meta:
        model = User
        fields = ['username', 'is_superuser']
        read_only_fields = ['username']


class AppSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.App` model."""

    owner = serializers.ReadOnlyField(source='owner.username')
    structure = JSONFieldSerializer(required=False)
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a :class:`AppSerializer`."""
        model = models.App
        fields = ['uuid', 'id', 'owner', 'url', 'structure', 'created', 'updated']
        read_only_fields = ['uuid']


class BuildSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Build` model."""

    app = serializers.SlugRelatedField(slug_field='id', queryset=models.App.objects.all())
    owner = serializers.ReadOnlyField(source='owner.username')
    procfile = JSONFieldSerializer(required=False)
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a :class:`BuildSerializer`."""
        model = models.Build
        fields = ['owner', 'app', 'image', 'sha', 'procfile', 'dockerfile', 'created',
                  'updated', 'uuid']
        read_only_fields = ['uuid']


class ConfigSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Config` model."""

    app = serializers.SlugRelatedField(slug_field='id', queryset=models.App.objects.all())
    owner = serializers.ReadOnlyField(source='owner.username')
    values = JSONStringFieldSerializer(required=False)
    memory = JSONStringFieldSerializer(required=False)
    cpu = JSONIntFieldSerializer(required=False)
    tags = JSONStringFieldSerializer(required=False)
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a :class:`ConfigSerializer`."""
        model = models.Config

    def validate_values(self, value):
        for k, v in value.viewitems():
            if not re.match(CONFIGKEY_MATCH, k):
                raise serializers.ValidationError(
                    "Config keys must start with a letter or underscore and "
                    "only contain [A-z0-9_]")
        return value

    def validate_memory(self, value):
        for k, v in value.viewitems():
            if v is None:  # use NoneType to unset a value
                continue
            if not re.match(PROCTYPE_MATCH, k):
                raise serializers.ValidationError("Process types can only contain [a-z]")
            if not re.match(MEMLIMIT_MATCH, str(v)):
                raise serializers.ValidationError(
                    "Limit format: <number><unit>, where unit = B, K, M or G")
        return value

    def validate_cpu(self, value):
        for k, v in value.viewitems():
            if v is None:  # use NoneType to unset a value
                continue
            if not re.match(PROCTYPE_MATCH, k):
                raise serializers.ValidationError("Process types can only contain [a-z]")
            shares = re.match(CPUSHARE_MATCH, str(v))
            if not shares:
                raise serializers.ValidationError("CPU shares must be an integer")
            for v in shares.groupdict().viewvalues():
                try:
                    i = int(v)
                except ValueError:
                    raise serializers.ValidationError("CPU shares must be an integer")
                if i > 1024 or i < 0:
                    raise serializers.ValidationError("CPU shares must be between 0 and 1024")
        return value

    def validate_tags(self, value):
        for k, v in value.viewitems():
            if v is None:  # use NoneType to unset a value
                continue
            if not re.match(TAGKEY_MATCH, k):
                raise serializers.ValidationError("Tag keys can only contain [a-z]")
            if not re.match(TAGVAL_MATCH, str(v)):
                raise serializers.ValidationError("Invalid tag value")
        return value


class ReleaseSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Release` model."""

    app = serializers.SlugRelatedField(slug_field='id', queryset=models.App.objects.all())
    owner = serializers.ReadOnlyField(source='owner.username')
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a :class:`ReleaseSerializer`."""
        model = models.Release


class ContainerSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Container` model."""

    app = serializers.SlugRelatedField(slug_field='id', queryset=models.App.objects.all())
    owner = serializers.ReadOnlyField(source='owner.username')
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    release = serializers.SerializerMethodField()

    class Meta:
        """Metadata options for a :class:`ContainerSerializer`."""
        model = models.Container
        fields = ['owner', 'app', 'release', 'type', 'num', 'state', 'created', 'updated', 'uuid']

    def get_release(self, obj):
        return "v{}".format(obj.release.version)


class KeySerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Key` model."""

    owner = serializers.ReadOnlyField(source='owner.username')
    fingerprint = serializers.CharField(read_only=True)
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a KeySerializer."""
        model = models.Key


class DomainSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Domain` model."""

    app = serializers.SlugRelatedField(slug_field='id', queryset=models.App.objects.all())
    owner = serializers.ReadOnlyField(source='owner.username')
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a :class:`DomainSerializer`."""
        model = models.Domain
        fields = ['uuid', 'owner', 'created', 'updated', 'app', 'domain']

    def validate_domain(self, value):
        """
        Check that the hostname is valid
        """
        if len(value) > 255:
            raise serializers.ValidationError('Hostname must be 255 characters or less.')
        if value[-1:] == ".":
            value = value[:-1]  # strip exactly one dot from the right, if present
        labels = value.split('.')
        if 'xip.io' in value:
            return value
        if labels[0] == '*':
            raise serializers.ValidationError(
                'Adding a wildcard subdomain is currently not supported.')
        allowed = re.compile("^(?!-)[a-z0-9-]{1,63}(?<!-)$", re.IGNORECASE)
        for label in labels:
            match = allowed.match(label)
            if not match or '--' in label or label.isdigit() or \
               len(labels) == 1 and any(char.isdigit() for char in label):
                raise serializers.ValidationError('Hostname does not look valid.')
        if models.Domain.objects.filter(domain=value).exists():
            raise serializers.ValidationError(
                "The domain {} is already in use by another app".format(value))
        return value


class CertificateSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Cert` model."""

    owner = serializers.ReadOnlyField(source='owner.username')
    expires = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a DomainCertSerializer."""
        model = models.Certificate
        extra_kwargs = {'certificate': {'write_only': True},
                        'key': {'write_only': True},
                        'common_name': {'required': False}}
        read_only_fields = ['expires', 'created', 'updated']


class PushSerializer(ModelSerializer):
    """Serialize a :class:`~api.models.Push` model."""

    app = serializers.SlugRelatedField(slug_field='id', queryset=models.App.objects.all())
    owner = serializers.ReadOnlyField(source='owner.username')
    created = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)
    updated = serializers.DateTimeField(format=settings.DEIS_DATETIME_FORMAT, read_only=True)

    class Meta:
        """Metadata options for a :class:`PushSerializer`."""
        model = models.Push
        fields = ['uuid', 'owner', 'app', 'sha', 'fingerprint', 'receive_user', 'receive_repo',
                  'ssh_connection', 'ssh_original_command', 'created', 'updated']
