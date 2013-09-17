#!/usr/bin/python
# -*- coding: utf-8 -*-

"""
Django admin app configuration for Deis API models.
"""

from __future__ import unicode_literals

from django.contrib import admin

from .models import Build
from .models import Config
from .models import Container
from .models import Flavor
from .models import Formation
from .models import Key
from .models import Layer
from .models import Node
from .models import Provider
from .models import Release


class BuildAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Build` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('sha', 'owner', 'app')
    list_filter = ('owner', 'app')
admin.site.register(Build, BuildAdmin)


class ConfigAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Config` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('version', 'owner', 'app')
    list_filter = ('owner', 'app')
admin.site.register(Config, ConfigAdmin)


class ReleaseAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Release` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('owner', 'app', 'version')
    list_filter = ('owner', 'app')
admin.site.register(Release, ReleaseAdmin)


class ContainerAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Container` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('short_name', 'owner', 'formation', 'app', 'status')
    list_filter = ('owner', 'formation', 'app', 'status')
admin.site.register(Container, ContainerAdmin)


class FlavorAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Flavor` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner', 'provider')
    list_filter = ('owner', 'provider')
admin.site.register(Flavor, FlavorAdmin)


class FormationAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Formation` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner')
    list_filter = ('owner',)
admin.site.register(Formation, FormationAdmin)


class KeyAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Key` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner', '__str__')
    list_filter = ('owner',)
admin.site.register(Key, KeyAdmin)


class LayerAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Layer` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner', 'formation', 'flavor', 'proxy', 'runtime', 'config')
    list_filter = ('owner', 'formation', 'flavor')
admin.site.register(Layer, LayerAdmin)


class NodeAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Node` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner', 'formation', 'fqdn')
    list_filter = ('owner', 'formation')
admin.site.register(Node, NodeAdmin)


class ProviderAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Provider` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner', 'type')
    list_filter = ('owner', 'type')
admin.site.register(Provider, ProviderAdmin)
