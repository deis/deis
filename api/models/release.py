from __future__ import unicode_literals

from django.conf import settings
from django.db import models
from django.dispatch import receiver
from django.dispatch.dispatcher import Signal
from django.utils.encoding import python_2_unicode_compatible

from api.utils import dict_diff
from base import UuidAuditedModel
import application
from build import Build
from config import Config


# define custom signals
release_signal = Signal(providing_args=['user', 'app'])


@python_2_unicode_compatible
class Release(UuidAuditedModel):
    """
    Software release deployed by the application platform

    Releases contain a :class:`Build` and a :class:`Config`.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey(application.App)
    version = models.PositiveIntegerField()
    summary = models.TextField(blank=True, null=True)

    config = models.ForeignKey(Config)
    build = models.ForeignKey(Build, blank=True, null=True)

    class Meta:
        app_label = 'api'
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'version'),)

    def __str__(self):
        return "{0}-v{1}".format(self.app.id, self.version)

    def previous(self):
        """
        Return the previous Release to this one.

        :return: the previous :class:`Release`, or None
        """
        releases = self.app.release_set
        if self.pk:
            releases = releases.exclude(pk=self.pk)
        try:
            # Get the Release previous to this one
            prev_release = releases.latest()
        except Release.DoesNotExist:
            prev_release = None
        return prev_release

    def save(self, *args, **kwargs):
        if not self.summary:
            self.summary = ''
            prev_release = self.previous()
            # compare this build to the previous build
            old_build = prev_release.build if prev_release else None
            # if the build changed, log it and who pushed it
            if self.build != old_build and self.build.sha:
                self.summary += "{} deployed {}".format(self.build.owner, self.build.sha[:7])
            # compare this config to the previous config
            old_config = prev_release.config if prev_release else None
            # if the config data changed, log the dict diff
            if self.config != old_config:
                dict1 = self.config.values
                dict2 = old_config.values if old_config else {}
                diff = dict_diff(dict1, dict2)
                # try to be as succinct as possible
                added = ', '.join(k for k in diff.get('added', {}))
                added = 'added ' + added if added else ''
                changed = ', '.join(k for k in diff.get('changed', {}))
                changed = 'changed ' + changed if changed else ''
                deleted = ', '.join(k for k in diff.get('deleted', {}))
                deleted = 'deleted ' + deleted if deleted else ''
                changes = ', '.join(i for i in (added, changed, deleted) if i)
                if changes:
                    if self.summary:
                        self.summary += ' and '
                    self.summary += "{} {}".format(self.config.owner, changes)
                if not self.summary:
                    if self.version == 1:
                        self.summary = "{} created the initial release".format(self.owner)
                    else:
                        self.summary = "{} changed nothing".format(self.owner)
        super(Release, self).save(*args, **kwargs)


@receiver(release_signal)
def new_release(sender, **kwargs):
    """
    Catch a release_signal and create a new release
    using the latest Build and Config for an application.

    Releases start at v1 and auto-increment.
    """
    user, app, = kwargs['user'], kwargs['app']
    last_release = app.release_set.latest()
    config = kwargs.get('config', last_release.config)
    build = kwargs.get('build', last_release.build)
    # overwrite config with build.config if the keys don't exist
    if build and build.config:
        new_values = {}
        for k, v in build.config.items():
            if not k in config.values:
                new_values[k] = v
        if new_values:
            # update with current config
            new_values.update(config.values)
            config = Config.objects.create(
                version=config.version + 1, owner=user,
                app=app, values=new_values)
    # create new release and auto-increment version
    new_version = last_release.version + 1
    release = Release.objects.create(
        owner=user, app=app, config=config,
        build=build, version=new_version)
    return release
