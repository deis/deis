"""
Classes to manage the presentation of Deis api models in the django
admin interface.
"""
# pylint: disable=R0903,R0904

from __future__ import unicode_literals

from django.contrib import admin

from api import models


class UuidAdmin(object):
    """Presents a UuidField, ensuring the actual UUID is read-only."""

    date_hierarchy = 'updated'

    def get_readonly_fields(self, _request, obj=None):
        """Override so once the UUID is set, it's read-only."""
        # pylint: disable=E1101
        if obj is not None:
            return self.readonly_fields + ('uuid',)
        else:
            return self.readonly_fields


# class SshKeysInline(admin.TabularInline):
#     model = models.SshKey.apps.through


# class SshKeyAdmin(admin.ModelAdmin):

#     inlines = [SshKeysInline,]
#     exclude = ('apps',)

# admin.site.register(models.SshKey, SshKeyAdmin)
#
#
# class InstanceAdmin(admin.ModelAdmin, UuidAdmin):
#     """Presents an Instance api model in the Django admin."""
#     pass
#
# admin.site.register(models.Instance, InstanceAdmin)
#
#
# class ProcessAdmin(admin.TabularInline, UuidAdmin):
#     """Presents a Process api model in the Django admin."""
#     model = models.Process
#
#
# class ProxyAdmin(admin.ModelAdmin, UuidAdmin):
#     """Presents a Proxy api model in the Django admin."""
#     pass
#
# admin.site.register(models.Proxy, ProxyAdmin)
#
#
# class RunAdmin(admin.ModelAdmin, UuidAdmin):
#     """Presents a Run api model in the Django admin."""
#     inlines = [ProcessAdmin]
#
# admin.site.register(models.Run, RunAdmin)


class BuildAdmin(admin.ModelAdmin, UuidAdmin):
    """Presents a Build api model in the Django admin."""
    pass

admin.site.register(models.Build, BuildAdmin)


class ConfigAdmin(admin.ModelAdmin, UuidAdmin):
    """Presents a Config api model in the Django admin."""

admin.site.register(models.Config, ConfigAdmin)


class ReleaseAdmin(admin.ModelAdmin, UuidAdmin):
    """Presents a Release api model in the Django admin."""
    pass

admin.site.register(models.Release, ReleaseAdmin)


class AccessAdmin(admin.ModelAdmin, UuidAdmin):
    """Presents an Access api model in the Django admin."""
    pass

admin.site.register(models.Access, AccessAdmin)


class EventAdmin(admin.ModelAdmin, UuidAdmin):
    """Presents an Event api model in the Django admin."""
    pass

admin.site.register(models.Event, EventAdmin)
