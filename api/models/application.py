from __future__ import unicode_literals
import logging
import os
import subprocess

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible
from guardian.shortcuts import get_users_with_perms
from json_field.fields import JSONField

from api.models import CM
from base import UuidAuditedModel
from formation import Formation

logger = logging.getLogger(__name__)


@python_2_unicode_compatible
class App(UuidAuditedModel):
    """
    Application used to service requests on behalf of end-users
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64, unique=True)
    formation = models.ForeignKey(Formation)
    containers = JSONField(default='{}', blank=True)

    class Meta:
        app_label = 'api'
        permissions = (('use_app', 'Can use app'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'formation': self.formation.id,
                'containers': dict(self.containers)}

    def build(self):
        initial_config = config.Config.objects.create(
            version=1, owner=self.owner, app=self, values={})
        initial_build = build.Build.objects.create(owner=self.owner, app=self)
        release.Release.objects.create(
            version=1, owner=self.owner, app=self, config=initial_config, build=initial_build)
        self.formation.publish()

    def destroy(self):
        CM.purge_app(self.flat())
        self.delete()
        self.formation.publish()

    def publish(self):
        """Publish the application to configuration management"""
        data = self.calculate()
        CM.publish_app(self.flat(), data)
        return data

    def converge(self):
        databag = self.publish()
        self.formation.converge()
        return databag

    def calculate(self):
        """Return a representation for configuration management"""
        d = {}
        d['id'] = self.id
        d['release'] = {}
        releases = self.release_set.all().order_by('-created')
        if releases:
            release = releases[0]
            d['release']['version'] = release.version
            d['release']['config'] = release.config.values
            d['release']['build'] = {'image': release.build.image}
            if release.build.url:
                d['release']['build']['url'] = release.build.url
                d['release']['build']['procfile'] = release.build.procfile
        d['containers'] = {}
        containers = self.container_set.all()
        if containers:
            for c in containers:
                d['containers'].setdefault(c.type, {})[str(c.num)] = c.status
        d['domains'] = []
        if self.formation.domain:
            d['domains'].append('{}.{}'.format(self.id, self.formation.domain))
        else:
            for n in self.formation.node_set.filter(layer__proxy=True):
                d['domains'].append(n.fqdn)
        # add proper sharing and access controls
        d['users'] = {self.owner.username: 'owner'}
        for u in (get_users_with_perms(self)):
            d['users'][u.username] = 'user'
        return d

    def logs(self):
        """Return aggregated log data for this application."""
        path = os.path.join(settings.DEIS_LOG_DIR, self.id + '.log')
        if not os.path.exists(path):
            raise EnvironmentError('Could not locate logs')
        data = subprocess.check_output(['tail', '-n', str(settings.LOG_LINES), path])
        return data

    def run(self, command):
        """Run a one-off command in an ephemeral app container."""
        # TODO: add support for interactive shell
        nodes = self.formation.node_set.filter(layer__runtime=True).order_by('?')
        if not nodes:
            raise EnvironmentError('No nodes available to run command')
        app_id, node = self.id, nodes[0]
        release = self.release_set.order_by('-created')[0]
        # prepare ssh command
        version = release.version
        docker_args = ' '.join(
            ['-a', 'stdout', '-a', 'stderr', '-rm',
             '-v', '/opt/deis/runtime/slugs/{app_id}-v{version}:/app'.format(**locals()),
             'deis/slugrunner'])
        env_args = ' '.join(["-e '{k}={v}'".format(**locals())
                             for k, v in release.config.values.items()])
        log_event(self, "deis run '{}'".format(command))
        command = "sudo docker run {env_args} {docker_args} {command}".format(**locals())
        return node.run(command)


def log_event(app, msg, level=logging.INFO):
    msg = "{}: {}".format(app.id, msg)
    logger.log(level, msg)


import build
import config
import release
