# -*- coding: utf-8 -*-

"""
Django admin app configuration for Deis API models.
"""

from __future__ import unicode_literals

from django.contrib import admin
from guardian.admin import GuardedModelAdmin

from .models import App
from .models import Build
from .models import Cluster
from .models import Config
from .models import Container
from .models import Key
from .models import Release


class AppAdmin(GuardedModelAdmin):
    """Set presentation options for :class:`~api.models.App` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner', 'cluster')
    list_filter = ('owner', 'cluster')
admin.site.register(App, AppAdmin)


class BuildAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Build` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('sha', 'owner', 'app')
    list_filter = ('owner', 'app')
admin.site.register(Build, BuildAdmin)


class ClusterAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Cluster` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner',)
    list_filter = ('owner',)
admin.site.register(Cluster, ClusterAdmin)


class ConfigAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Config` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('version', 'owner', 'app')
    list_filter = ('owner', 'app')
admin.site.register(Config, ConfigAdmin)


class ContainerAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Container` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('short_name', 'owner', 'cluster', 'app', 'state')
    list_filter = ('owner', 'cluster', 'app', 'state')
admin.site.register(Container, ContainerAdmin)


class KeyAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Key` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('id', 'owner', '__str__')
    list_filter = ('owner',)
admin.site.register(Key, KeyAdmin)


class ReleaseAdmin(admin.ModelAdmin):
    """Set presentation options for :class:`~api.models.Release` models
    in the Django admin.
    """
    date_hierarchy = 'created'
    list_display = ('owner', 'app', 'version')
    list_filter = ('owner', 'app')
admin.site.register(Release, ReleaseAdmin)
