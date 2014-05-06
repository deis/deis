"""
Classes to serialize the RESTful representation of Deis API models.
"""

from __future__ import unicode_literals

import re

from django.contrib.auth.models import User
from rest_framework import serializers

from api import models
from api import utils


class OwnerSlugRelatedField(serializers.SlugRelatedField):
    """Filter queries by owner as well as slug_field."""

    def from_native(self, data):
        """Fetch model object from its 'native' representation.
        TODO: request.user is not going to work in a team environment...
        """
        self.queryset = self.queryset.filter(owner=self.context['request'].user)
        return serializers.SlugRelatedField.from_native(self, data)


class UserSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.User` model."""

    class Meta:
        """Metadata options for a UserSerializer."""
        model = User
        read_only_fields = ('is_superuser', 'is_staff', 'groups',
                            'user_permissions', 'last_login', 'date_joined')

    @property
    def data(self):
        """Custom data property that removes secure user fields"""
        d = super(UserSerializer, self).data
        for f in ('password',):
            if f in d:
                del d[f]
        return d


class AdminUserSerializer(serializers.ModelSerializer):
    """Serialize admin status for a :class:`~api.models.User` model."""

    class Meta:
        model = User
        fields = ('username', 'is_superuser')
        read_only_fields = ('username',)


class ClusterSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Cluster` model."""

    owner = serializers.Field(source='owner.username')

    class Meta:
        """Metadata options for a :class:`ClusterSerializer`."""
        model = models.Cluster
        read_only_fields = ('created', 'updated')


class PushSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Push` model."""

    owner = serializers.Field(source='owner.username')
    app = serializers.SlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`PushSerializer`."""
        model = models.Push
        read_only_fields = ('uuid', 'created', 'updated')


class BuildSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Build` model."""

    owner = serializers.Field(source='owner.username')
    app = serializers.SlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`BuildSerializer`."""
        model = models.Build
        read_only_fields = ('uuid', 'created', 'updated')


class ConfigSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Config` model."""

    owner = serializers.Field(source='owner.username')
    app = serializers.SlugRelatedField(slug_field='id')
    values = serializers.ModelField(
        model_field=models.Config()._meta.get_field('values'), required=False)

    class Meta:
        """Metadata options for a :class:`ConfigSerializer`."""
        model = models.Config
        read_only_fields = ('uuid', 'created', 'updated')


class ReleaseSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Release` model."""

    owner = serializers.Field(source='owner.username')
    app = serializers.SlugRelatedField(slug_field='id')
    config = serializers.SlugRelatedField(slug_field='uuid')
    build = serializers.SlugRelatedField(slug_field='uuid')

    class Meta:
        """Metadata options for a :class:`ReleaseSerializer`."""
        model = models.Release
        read_only_fields = ('uuid', 'created', 'updated')


class AppSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.App` model."""

    owner = serializers.Field(source='owner.username')
    id = serializers.SlugField(default=utils.generate_app_name)
    cluster = serializers.SlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`AppSerializer`."""
        model = models.App
        read_only_fields = ('created', 'updated')

    def validate_id(self, attrs, source):
        """
        Check that the ID is all lowercase
        """
        value = attrs[source]
        match = re.match(r'^[a-z0-9-]+$', value)
        if not match:
            raise serializers.ValidationError("App IDs can only contain [a-z0-9-]")
        return attrs


class ContainerSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Container` model."""

    owner = serializers.Field(source='owner.username')
    app = OwnerSlugRelatedField(slug_field='id')
    release = serializers.SlugRelatedField(slug_field='uuid')

    class Meta:
        """Metadata options for a :class:`ContainerSerializer`."""
        model = models.Container
        read_only_fields = ('created', 'updated')

    def transform_release(self, obj, value):
        return "v{}".format(obj.release.version)


class KeySerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Key` model."""

    owner = serializers.Field(source='owner.username')

    class Meta:
        """Metadata options for a KeySerializer."""
        model = models.Key
        read_only_fields = ('created', 'updated')
