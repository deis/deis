"""
Classes to serialize the RESTful representation of Deis API models.
"""
# pylint: disable=R0903,W0232

from __future__ import unicode_literals

from rest_framework import serializers

from api import models, utils

from django.contrib.auth.models import User



class UserSerializer(serializers.ModelSerializer):

    """Serializes a User model."""

    class Meta:
        """Metadata options for a UserSerializer."""
        model = User
        read_only_fields = ('is_superuser', 'is_staff', 'groups', 'user_permissions',
                            'last_login', 'date_joined')

    @property
    def data(self):
        "Custom data property that removes secure user fields"
        d = super(UserSerializer, self).data
        for f in ('password',):
            if f in d:
                del d[f]
        return d


class KeySerializer(serializers.ModelSerializer):

    """Serializes a Key model."""

    owner = serializers.Field(source='owner.username')

    class Meta:
        """Metadata options for a KeySerializer."""
        model = models.Key
        read_only_fields = ('created', 'updated')


class ProviderSerializer(serializers.ModelSerializer):

    """Serializes a Provider model."""

    owner = serializers.Field(source='owner.username')

    class Meta:
        """Metadata options for a ProviderSerializer."""
        model = models.Provider
        read_only_fields = ('created', 'updated')


class FlavorSerializer(serializers.ModelSerializer):

    """Serializes a Flavor model."""

    owner = serializers.Field(source='owner.username')
    provider = serializers.SlugRelatedField(slug_field='id')
    
    class Meta:
        """Metadata options for a FlavorSerializer."""
        model = models.Flavor
        read_only_fields = ('created', 'updated')


class ConfigSerializer(serializers.ModelSerializer):

    """Serializes a Config model."""

    owner = serializers.Field(source='owner.username')
    formation = serializers.SlugRelatedField(slug_field='id')
    values = serializers.ModelField(model_field=models.Config()._meta.get_field('values'),
                                    required=False)
    
    class Meta:
        """Metadata options for a ConfigSerializer."""
        model = models.Config
        read_only_fields = ('uuid', 'created', 'updated')


class BuildSerializer(serializers.ModelSerializer):

    """Serializes a Build model."""

    owner = serializers.Field(source='owner.username')
    formation = serializers.SlugRelatedField(slug_field='id')
    
    class Meta:
        """Metadata options for a BuildSerializer."""
        model = models.Build
        read_only_fields = ('uuid', 'created', 'updated')


class ReleaseSerializer(serializers.ModelSerializer):

    """Serializes a Release model."""

    owner = serializers.Field(source='owner.username')
    formation = serializers.SlugRelatedField(slug_field='id')
    config = serializers.SlugRelatedField(slug_field='uuid')
    build = serializers.SlugRelatedField(slug_field='uuid', required=False)
    
    class Meta:
        """Metadata options for a ReleaseSerializer."""
        model = models.Release
        read_only_fields = ('uuid', 'created', 'updated')


class FormationSerializer(serializers.ModelSerializer):

    """Serializes a Formation model."""

    owner = serializers.Field(source='owner.username')
    id = serializers.SlugField(default=utils.generate_app_name)
    flavor = serializers.SlugRelatedField(slug_field='id')
    structure = serializers.ModelField(
        model_field=models.Formation()._meta.get_field('structure'), required=False)
                                    
    class Meta:
        """Metadata options for a FormationSerializer."""
        model = models.Formation
        read_only_fields = ('created', 'updated')


class NodeSerializer(serializers.ModelSerializer):

    """Serializes a Node model."""

    owner = serializers.Field(source='owner.username')
    formation = serializers.SlugRelatedField(slug_field='id')
    
    class Meta:
        """Metadata options for a NodeSerializer."""
        model = models.Node
        read_only_fields = ('created', 'updated')


class BackendSerializer(serializers.ModelSerializer):
 
    """Serializes a Backend model."""
    
    owner = serializers.Field(source='owner.username')
    formation = serializers.SlugRelatedField(slug_field='id')
    node = serializers.SlugRelatedField(slug_field='uuid')
    
    class Meta:
        """Metadata options for a BackendSerializer."""
        model = models.Backend
        read_only_fields = ('created', 'updated')


class ProxySerializer(serializers.ModelSerializer):
 
    """Serializes a Proxy model."""
    
    owner = serializers.Field(source='owner.username')
    formation = serializers.SlugRelatedField(slug_field='id')
    node = serializers.SlugRelatedField(slug_field='uuid')
    
    class Meta:
        """Metadata options for a ProxySerializer."""
        model = models.Proxy
        read_only_fields = ('created', 'updated')


class ContainerSerializer(serializers.ModelSerializer):

    """Serializes a Container model."""

    owner = serializers.Field(source='owner.username')
    formation = serializers.SlugRelatedField(slug_field='id')
    node = serializers.SlugRelatedField(slug_field='uuid')
    
    class Meta:
        """Metadata options for a ContainerSerializer."""
        model = models.Container
        read_only_fields = ('created', 'updated')
