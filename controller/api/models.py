# -*- coding: utf-8 -*-

"""
Data models for the Deis API.
"""

from __future__ import unicode_literals
import etcd
import importlib
import logging
import os
import subprocess

from celery.canvas import group
from django.conf import settings
from django.contrib.auth.models import User
from django.db import models
from django.db.models import Max
from django.db.models.signals import post_delete
from django.db.models.signals import post_save
from django.utils.encoding import python_2_unicode_compatible
from django_fsm import FSMField, transition
from django_fsm.signals import post_transition
from json_field.fields import JSONField

from api import fields, tasks
from registry import publish_release
from utils import dict_diff, fingerprint


logger = logging.getLogger(__name__)


def log_event(app, msg, level=logging.INFO):
    msg = "{}: {}".format(app.id, msg)
    logger.log(level, msg)


class AuditedModel(models.Model):
    """Add created and updated fields to a model."""

    created = models.DateTimeField(auto_now_add=True)
    updated = models.DateTimeField(auto_now=True)

    class Meta:
        """Mark :class:`AuditedModel` as abstract."""
        abstract = True


class UuidAuditedModel(AuditedModel):
    """Add a UUID primary key to an :class:`AuditedModel`."""

    uuid = fields.UuidField('UUID', primary_key=True)

    class Meta:
        """Mark :class:`UuidAuditedModel` as abstract."""
        abstract = True


@python_2_unicode_compatible
class Cluster(UuidAuditedModel):
    """
    Cluster used to run jobs
    """

    CLUSTER_TYPES = (('mock', 'Mock Cluster'),
                     ('coreos', 'CoreOS Cluster'),
                     ('faulty', 'Faulty Cluster'))

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=128, unique=True)
    type = models.CharField(max_length=16, choices=CLUSTER_TYPES, default='coreos')

    domain = models.CharField(max_length=128)
    hosts = models.CharField(max_length=256)
    auth = models.TextField()
    options = JSONField(default='{}', blank=True)

    def __str__(self):
        return self.id

    def _get_scheduler(self, *args, **kwargs):
        module_name = 'scheduler.' + self.type
        mod = importlib.import_module(module_name)
        return mod.SchedulerClient(self.id, self.hosts, self.auth,
                                   self.domain, self.options)

    _scheduler = property(_get_scheduler)

    def create(self):
        """
        Initialize a cluster's router and log aggregator
        """
        return tasks.create_cluster.delay(self).get()

    def destroy(self):
        """
        Destroy a cluster's router and log aggregator
        """
        return tasks.destroy_cluster.delay(self).get()


@python_2_unicode_compatible
class App(UuidAuditedModel):
    """
    Application used to service requests on behalf of end-users
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64, unique=True)
    cluster = models.ForeignKey('Cluster')
    structure = JSONField(default='{}', blank=True)

    class Meta:
        permissions = (('use_app', 'Can use app'),)

    def __str__(self):
        return self.id

    def create(self, *args, **kwargs):
        config = Config.objects.create(owner=self.owner, app=self, values={})
        build = Build.objects.create(owner=self.owner, app=self, image=settings.DEFAULT_BUILD)
        Release.objects.create(version=1, owner=self.owner, app=self, config=config, build=build)

    def delete(self, *args, **kwargs):
        for c in self.container_set.all():
            c.destroy()
        return super(App, self).delete(*args, **kwargs)

    def deploy(self, release):
        tasks.deploy_release.delay(self, release).get()

    def destroy(self, *args, **kwargs):
        return self.delete(*args, **kwargs)

    def scale(self, **kwargs):  # noqa
        """Scale containers up or down to match requested."""
        requested_containers = self.structure.copy()
        release = self.release_set.latest()
        # test for available process types
        available_process_types = release.build.procfile or {}
        for container_type in requested_containers.keys():
            if container_type == 'cmd':
                continue  # allow docker cmd types in case we don't have the image source
            if not container_type in available_process_types:
                raise EnvironmentError(
                    'Container type {} does not exist in application'.format(container_type))
        msg = 'Containers scaled ' + ' '.join(
            "{}={}".format(k, v) for k, v in requested_containers.items())
        # iterate and scale by container type (web, worker, etc)
        changed = False
        to_add, to_remove = [], []
        for container_type in requested_containers.keys():
            containers = list(self.container_set.filter(type=container_type).order_by('created'))
            # increment new container nums off the most recent container
            results = self.container_set.filter(type=container_type).aggregate(Max('num'))
            container_num = (results.get('num__max') or 0) + 1
            requested = requested_containers.pop(container_type)
            diff = requested - len(containers)
            if diff == 0:
                continue
            changed = True
            while diff < 0:
                c = containers.pop()
                to_remove.append(c)
                diff += 1
            while diff > 0:
                c = Container.objects.create(owner=self.owner,
                                             app=self,
                                             release=release,
                                             type=container_type,
                                             num=container_num)
                to_add.append(c)
                container_num += 1
                diff -= 1
        if changed:
            subtasks = []
            if to_add:
                subtasks.append(tasks.start_containers.s(to_add))
            if to_remove:
                subtasks.append(tasks.stop_containers.s(to_remove))
            group(*subtasks).apply_async().join()
            log_event(self, msg)
        return changed

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
        log_event(self, "deis run '{}'".format(command))
        c_num = max([c.num for c in self.container_set.filter(type='admin')] or [0]) + 1
        c = Container.objects.create(owner=self.owner,
                                     app=self,
                                     release=self.release_set.latest(),
                                     type='admin',
                                     num=c_num)
        rc, output = tasks.run_command.delay(c, command).get()
        return rc, output


@python_2_unicode_compatible
class Container(UuidAuditedModel):
    """
    Docker container used to securely host an application process.
    """
    INITIALIZED = 'initialized'
    CREATED = 'created'
    UP = 'up'
    DOWN = 'down'
    DESTROYED = 'destroyed'
    STATE_CHOICES = (
        (INITIALIZED, 'initialized'),
        (CREATED, 'created'),
        (UP, 'up'),
        (DOWN, 'down'),
        (DESTROYED, 'destroyed')
    )

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    release = models.ForeignKey('Release')
    type = models.CharField(max_length=128, blank=True)
    num = models.PositiveIntegerField()
    state = FSMField(default=INITIALIZED, choices=STATE_CHOICES, protected=True)

    def short_name(self):
        if self.type:
            return "{}.{}.{}".format(self.release.app.id, self.type, self.num)
        return "{}.{}".format(self.release.app.id, self.num)
    short_name.short_description = 'Name'

    def __str__(self):
        return self.short_name()

    class Meta:
        get_latest_by = '-created'
        ordering = ['created']
        unique_together = (('type', 'num'),)

    def _get_job_id(self):
        app = self.app.id
        release = self.release
        version = "v{}".format(release.version)
        num = self.num
        c_type = self.type
        if not c_type:
            job_id = "{app}_{version}.{num}".format(**locals())
        else:
            job_id = "{app}_{version}.{c_type}.{num}".format(**locals())
        return job_id

    _job_id = property(_get_job_id)

    def _get_scheduler(self):
        return self.app.cluster._scheduler

    _scheduler = property(_get_scheduler)

    def _get_command(self):
        c_type = self.type
        if c_type:
            # handle special case for Dockerfile deployments
            if c_type == 'cmd':
                return ''
            else:
                return 'start {c_type}'
        else:
            return ''

    _command = property(_get_command)

    @transition(field=state, source=INITIALIZED, target=CREATED)
    def create(self):
        image = self.release.image
        c_type = self.type
        self._scheduler.create(self._job_id, image, self._command.format(**locals()))

    @transition(field=state,
                source=[CREATED, UP, DOWN],
                target=UP, crashed=DOWN)
    def start(self):
        self._scheduler.start(self._job_id)

    @transition(field=state,
                source=[INITIALIZED, CREATED, UP, DOWN],
                target=UP,
                crashed=DOWN)
    def deploy(self, release):
        old_job_id = self._job_id
        # update release
        self.release = release
        self.save()
        # deploy new container
        new_job_id = self._job_id
        image = self.release.image
        c_type = self.type
        self._scheduler.create(new_job_id, image, self._command.format(**locals()))
        self._scheduler.start(new_job_id)
        # destroy old container
        self._scheduler.destroy(old_job_id)

    @transition(field=state, source=UP, target=DOWN)
    def stop(self):
        self._scheduler.stop(self._job_id)

    @transition(field=state,
                source=[INITIALIZED, CREATED, UP, DOWN],
                target=DESTROYED)
    def destroy(self):
        # TODO: add check for active connections before killing
        self._scheduler.destroy(self._job_id)

    @transition(field=state,
                source=[INITIALIZED, CREATED, DESTROYED],
                target=DESTROYED)
    def run(self, command):
        """Run a one-off command"""
        rc, output = self._scheduler.run(self._job_id, self.release.image, command)
        return rc, output


@python_2_unicode_compatible
class Push(UuidAuditedModel):
    """
    Instance of a push used to trigger an application build
    """
    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    sha = models.CharField(max_length=40)

    fingerprint = models.CharField(max_length=255)
    receive_user = models.CharField(max_length=255)
    receive_repo = models.CharField(max_length=255)

    ssh_connection = models.CharField(max_length=255)
    ssh_original_command = models.CharField(max_length=255)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'uuid'),)

    def __str__(self):
        return "{0}-{1}".format(self.app.id, self.sha[:7])


@python_2_unicode_compatible
class Build(UuidAuditedModel):
    """
    Instance of a software build used by runtime nodes
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    image = models.CharField(max_length=256)

    # optional fields populated by builder
    sha = models.CharField(max_length=40, blank=True)
    procfile = JSONField(default='{}', blank=True)
    dockerfile = models.TextField(blank=True)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'uuid'),)

    def __str__(self):
        return "{0}-{1}".format(self.app.id, self.uuid[:7])


@python_2_unicode_compatible
class Config(UuidAuditedModel):
    """
    Set of configuration values applied as environment variables
    during runtime execution of the Application.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    values = JSONField(default='{}', blank=True)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'uuid'),)

    def __str__(self):
        return "{}-{}".format(self.app.id, self.uuid[:7])


@python_2_unicode_compatible
class Release(UuidAuditedModel):
    """
    Software release deployed by the application platform

    Releases contain a :class:`Build` and a :class:`Config`.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    version = models.PositiveIntegerField()
    summary = models.TextField(blank=True, null=True)

    config = models.ForeignKey('Config')
    build = models.ForeignKey('Build')
    # NOTE: image contains combined build + config, ready to run
    image = models.CharField(max_length=256)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'version'),)

    def __str__(self):
        return "{0}-v{1}".format(self.app.id, self.version)

    def new(self, user, config=None, build=None, summary=None):
        """
        Create a new application release using the provided Build and Config
        on behalf of a user.

        Releases start at v1 and auto-increment.
        """
        if not config:
            config = self.config
        if not build:
            build = self.build
        # prepare release tag
        new_version = self.version + 1
        tag = 'v{}'.format(new_version)
        image = build.image + ':{tag}'.format(**locals())
        # create new release and auto-increment version
        release = Release.objects.create(
            owner=user, app=self.app, config=config,
            build=build, version=new_version, image=image, summary=summary)
        # publish release to registry as new docker image
        repository_path = self.app.id
        publish_release(repository_path, config.values, tag)
        return release

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
            if self.build != old_build:
                if self.build.sha:
                    self.summary += "{} deployed {}".format(self.build.owner, self.build.sha[:7])
                else:
                    self.summary += "{} deployed {}".format(self.build.owner, self.build.image)
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


@python_2_unicode_compatible
class Key(UuidAuditedModel):
    """An SSH public key."""

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=128)
    public = models.TextField(unique=True)

    class Meta:
        verbose_name = 'SSH Key'
        unique_together = (('owner', 'id'))

    def __str__(self):
        return "{}...{}".format(self.public[:18], self.public[-31:])


# define update/delete callbacks for synchronizing
# models with the configuration management backend


def _log_build_created(**kwargs):
    if kwargs.get('created'):
        build = kwargs['instance']
        log_event(build.app, "Build {} created".format(build))


def _log_release_created(**kwargs):
    if kwargs.get('created'):
        release = kwargs['instance']
        log_event(release.app, "Release {} created".format(release))


def _log_config_updated(**kwargs):
    config = kwargs['instance']
    log_event(config.app, "Config {} updated".format(config))


def _etcd_publish_key(**kwargs):
    key = kwargs['instance']
    _etcd_client.write('/deis/builder/users/{}/{}'.format(
        key.owner.username, fingerprint(key.public)), key.public)


def _etcd_purge_key(**kwargs):
    key = kwargs['instance']
    _etcd_client.delete('/deis/builder/users/{}/{}'.format(
        key.owner.username, fingerprint(key.public)))


def _etcd_purge_user(**kwargs):
    username = kwargs['instance'].username
    _etcd_client.delete('/deis/builder/users/{}'.format(username), dir=True, recursive=True)


# Log significant app-related events
post_save.connect(_log_build_created, sender=Build, dispatch_uid='api.models')
post_save.connect(_log_release_created, sender=Release, dispatch_uid='api.models')
post_save.connect(_log_config_updated, sender=Config, dispatch_uid='api.models')


# save FSM transitions as they happen
def _save_transition(**kwargs):
    kwargs['instance'].save()

post_transition.connect(_save_transition)

# wire up etcd publishing if we can connect
try:
    _etcd_client = etcd.Client(host=settings.ETCD_HOST, port=int(settings.ETCD_PORT))
    _etcd_client.get('/deis')
except etcd.EtcdException:
    logger.log(logging.WARNING, 'Cannot synchronize with etcd cluster')
    _etcd_client = None

if _etcd_client:
    post_save.connect(_etcd_publish_key, sender=Key, dispatch_uid='api.models')
    post_delete.connect(_etcd_purge_key, sender=Key, dispatch_uid='api.models')
    post_delete.connect(_etcd_purge_user, sender=User, dispatch_uid='api.models')
