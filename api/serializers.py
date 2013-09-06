"""
Classes to serialize the RESTful representation of Deis API models.
"""
# pylint: disable=R0903,W0232

from __future__ import unicode_literals

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


class KeySerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Key` model."""

    owner = serializers.Field(source='owner.username')

    class Meta:
        """Metadata options for a KeySerializer."""
        model = models.Key
        read_only_fields = ('created', 'updated')


class ProviderSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Provider` model."""

    owner = serializers.Field(source='owner.username')

    class Meta:
        """Metadata options for a ProviderSerializer."""
        model = models.Provider
        read_only_fields = ('created', 'updated')


class FlavorSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Flavor` model."""

    owner = serializers.Field(source='owner.username')
    provider = OwnerSlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`FlavorSerializer`."""
        model = models.Flavor
        read_only_fields = ('created', 'updated')


class ConfigSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Config` model."""

    owner = serializers.Field(source='owner.username')
    app = OwnerSlugRelatedField(slug_field='id')
    values = serializers.ModelField(
        model_field=models.Config()._meta.get_field('values'), required=False)

    class Meta:
        """Metadata options for a :class:`ConfigSerializer`."""
        model = models.Config
        read_only_fields = ('uuid', 'created', 'updated')


class BuildSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Build` model."""

    owner = serializers.Field(source='owner.username')
    app = OwnerSlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`BuildSerializer`."""
        model = models.Build
        read_only_fields = ('uuid', 'created', 'updated')


class ReleaseSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Release` model."""

    owner = serializers.Field(source='owner.username')
    app = OwnerSlugRelatedField(slug_field='id')
    config = serializers.SlugRelatedField(slug_field='uuid')
    build = serializers.SlugRelatedField(slug_field='uuid', required=False)

    class Meta:
        """Metadata options for a :class:`ReleaseSerializer`."""
        model = models.Release
        read_only_fields = ('uuid', 'created', 'updated')


class FormationSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Formation` model."""

    owner = serializers.Field(source='owner.username')

    class Meta:
        """Metadata options for a :class:`FormationSerializer`."""
        model = models.Formation
        read_only_fields = ('created', 'updated')


class LayerSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Layer` model."""

    owner = serializers.Field(source='owner.username')
    formation = OwnerSlugRelatedField(slug_field='id')
    flavor = OwnerSlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`LayerSerializer`."""
        model = models.Layer
        read_only_fields = ('created', 'updated')

    @property
    def data(self):
        """Custom data property that removes secure fields"""
        d = super(LayerSerializer, self).data
        for f in ('ssh_private_key',):
            if f in d:
                del d[f]
        return d


class NodeSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Node` model."""

    owner = serializers.Field(source='owner.username')
    formation = OwnerSlugRelatedField(slug_field='id')
    layer = OwnerSlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`NodeSerializer`."""
        model = models.Node
        read_only_fields = ('created', 'updated')


class AppSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.App` model."""

    owner = serializers.Field(source='owner.username')
    id = serializers.SlugField(default=utils.generate_app_name)
    formation = OwnerSlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`AppSerializer`."""
        model = models.App
        read_only_fields = ('created', 'updated')


class ContainerSerializer(serializers.ModelSerializer):
    """Serialize a :class:`~api.models.Container` model."""

    owner = serializers.Field(source='owner.username')
    formation = OwnerSlugRelatedField(slug_field='id')
    node = OwnerSlugRelatedField(slug_field='id')
    app = OwnerSlugRelatedField(slug_field='id')

    class Meta:
        """Metadata options for a :class:`ContainerSerializer`."""
        model = models.Container
        read_only_fields = ('created', 'updated')
