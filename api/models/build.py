from __future__ import unicode_literals
import logging

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from json_field.fields import JSONField

from api.models import User
from base import UuidAuditedModel

import application

logger = logging.getLogger(__name__)


@python_2_unicode_compatible
class Build(UuidAuditedModel):
    """
    Instance of a software build used by runtime nodes
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey(application.App)
    sha = models.CharField('SHA', max_length=255, blank=True)
    output = models.TextField(blank=True)

    image = models.CharField(max_length=256, default='deis/slugbuilder')

    procfile = JSONField(blank=True)
    dockerfile = models.TextField(blank=True)
    config = JSONField(blank=True)

    url = models.URLField('URL')
    size = models.IntegerField(blank=True, null=True)
    checksum = models.CharField(max_length=255, blank=True)

    class Meta:
        app_label = 'api'
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'uuid'),)

    def __str__(self):
        return "{0}-{1}".format(self.app.id, self.sha[:7])

    @classmethod
    def push(cls, push):
        """Process a push from a local Git server.

        Creates a new Build and returns the application's
        databag for processing by the git-receive hook
        """
        # SECURITY:
        # we assume the first part of the ssh key name
        # is the authenticated user because we trust gitosis
        username = push.pop('username').split('_')[0]
        # retrieve the user and app instances
        user = User.objects.get(username=username)
        app = application.App.objects.get(id=push.pop('app'))
        # merge the push with the required model instances
        push['owner'] = user
        push['app'] = app
        # create the build
        new_build = cls.objects.create(**push)
        # send a release signal
        release.release_signal.send(sender=user, build=new_build, app=app, user=user)
        # see if we need to scale an initial web container
        if len(app.formation.node_set.filter(layer__runtime=True)) > 0 and \
           len(app.container_set.filter(type='web')) < 1:
            # scale an initial web containers
            container.Container.objects.scale(app, {'web': 1})
        # publish and converge the application
        return app.converge()


import container
import release
