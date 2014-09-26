# -*- coding: utf-8 -*-

"""
Data models for the Deis API.
"""

from __future__ import unicode_literals
import etcd
import importlib
import logging
import os
import re
import subprocess
import time
import threading

from django.conf import settings
from django.contrib.auth.models import User
from django.core.exceptions import ValidationError
from django.db import models
from django.db.models import Max
from django.db.models.signals import post_delete
from django.db.models.signals import post_save
from django.utils.encoding import python_2_unicode_compatible
from django_fsm import FSMField, transition
from django_fsm.signals import post_transition
from docker.utils import utils
from json_field.fields import JSONField
import requests

from api import fields
from registry import publish_release
from utils import dict_diff, fingerprint


logger = logging.getLogger(__name__)


def log_event(app, msg, level=logging.INFO):
    msg = "{}: {}".format(app.id, msg)
    logger.log(level, msg)  # django logger
    app.log(msg)            # local filesystem


def validate_app_structure(value):
    """Error if the dict values aren't ints >= 0."""
    try:
        for k, v in value.iteritems():
            if int(v) < 0:
                raise ValueError("Must be greater than or equal to zero")
    except ValueError, err:
        raise ValidationError(err)


def validate_comma_separated(value):
    """Error if the value doesn't look like a list of hostnames or IP addresses
    separated by commas.
    """
    if not re.search(r'^[a-zA-Z0-9-,\.]+$', value):
        raise ValidationError(
            "{} should be a comma-separated list".format(value))


def validate_domain(value):
    """Error if the domain contains unexpected characters."""
    if not re.search(r'^[a-zA-Z0-9-\.]+$', value):
        raise ValidationError('"{}" contains unexpected characters'.format(value))


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
                     ('chaos', 'Chaos Cluster'))

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=128, unique=True)
    type = models.CharField(max_length=16, choices=CLUSTER_TYPES, default='coreos')

    domain = models.CharField(max_length=128, validators=[validate_domain])
    hosts = models.CharField(max_length=256, validators=[validate_comma_separated])
    auth = models.TextField()
    options = JSONField(default={}, blank=True)

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
        return self._scheduler.setUp()

    def destroy(self):
        """
        Destroy a cluster's router and log aggregator
        """
        return self._scheduler.tearDown()


@python_2_unicode_compatible
class App(UuidAuditedModel):
    """
    Application used to service requests on behalf of end-users
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64, unique=True)
    cluster = models.ForeignKey('Cluster')
    structure = JSONField(default={}, blank=True, validators=[validate_app_structure])

    class Meta:
        permissions = (('use_app', 'Can use app'),)

    def __str__(self):
        return self.id

    @property
    def url(self):
        return self.id + '.' + self.cluster.domain

    def log(self, message):
        """Logs a message to the application's log file.

        This is a workaround for how Django interacts with Python's logging module. Each app
        needs its own FileHandler instance so it can write to its own log file. That won't work in
        Django's case because logging is set up before you run the server and it disables all
        existing logging configurations.
        """
        with open(os.path.join(settings.DEIS_LOG_DIR, self.id + '.log'), 'a') as f:
            msg = "{} deis[api]: {}\n".format(time.strftime('%Y-%m-%d %H:%M:%S'), message)
            f.write(msg.encode('utf-8'))

    def create(self, *args, **kwargs):
        """Create a new application with an initial release"""
        config = Config.objects.create(owner=self.owner, app=self)
        build = Build.objects.create(owner=self.owner, app=self, image=settings.DEFAULT_BUILD)
        Release.objects.create(version=1, owner=self.owner, app=self, config=config, build=build)

    def delete(self, *args, **kwargs):
        """Delete this application including all containers"""
        for c in self.container_set.all():
            c.destroy()
        self._clean_app_logs()
        return super(App, self).delete(*args, **kwargs)

    def _clean_app_logs(self):
        """Delete application logs stored by the logger component"""
        path = os.path.join(settings.DEIS_LOG_DIR, self.id + '.log')
        if os.path.exists(path):
            os.remove(path)

    def scale(self, user, structure):  # noqa
        """Scale containers up or down to match requested structure."""
        requested_structure = structure.copy()
        release = self.release_set.latest()
        # test for available process types
        available_process_types = release.build.procfile or {}
        for container_type in requested_structure.keys():
            if container_type == 'cmd':
                continue  # allow docker cmd types in case we don't have the image source
            if container_type not in available_process_types:
                raise EnvironmentError(
                    'Container type {} does not exist in application'.format(container_type))
        msg = '{} scaled containers '.format(user.username) + ' '.join(
            "{}={}".format(k, v) for k, v in requested_structure.items())
        log_event(self, msg)
        # iterate and scale by container type (web, worker, etc)
        changed = False
        to_add, to_remove = [], []
        for container_type in requested_structure.keys():
            containers = list(self.container_set.filter(type=container_type).order_by('created'))
            # increment new container nums off the most recent container
            results = self.container_set.filter(type=container_type).aggregate(Max('num'))
            container_num = (results.get('num__max') or 0) + 1
            requested = requested_structure.pop(container_type)
            diff = requested - len(containers)
            if diff == 0:
                continue
            changed = True
            while diff < 0:
                c = containers.pop()
                to_remove.append(c)
                diff += 1
            while diff > 0:
                # create a database record
                c = Container.objects.create(owner=self.owner,
                                             app=self,
                                             release=release,
                                             type=container_type,
                                             num=container_num)
                to_add.append(c)
                container_num += 1
                diff -= 1
        if changed:
            if to_add:
                self._start_containers(to_add)
            if to_remove:
                self._destroy_containers(to_remove)
        # save new structure to the database
        self.structure = structure
        self.save()
        return changed

    def _start_containers(self, to_add):
        """Creates and starts containers via the scheduler"""
        create_threads = []
        start_threads = []
        for c in to_add:
            create_threads.append(threading.Thread(target=c.create))
            start_threads.append(threading.Thread(target=c.start))
        [t.start() for t in create_threads]
        [t.join() for t in create_threads]
        if set([c.state for c in to_add]) != set([Container.CREATED]):
            err = 'aborting, failed to create some containers'
            log_event(self, err, logging.ERROR)
            raise RuntimeError(err)
        [t.start() for t in start_threads]
        [t.join() for t in start_threads]
        if set([c.state for c in to_add]) != set([Container.UP]):
            err = 'warning, some containers failed to start'
            log_event(self, err, logging.WARNING)

    def _destroy_containers(self, to_destroy):
        """Destroys containers via the scheduler"""
        destroy_threads = []
        for c in to_destroy:
            destroy_threads.append(threading.Thread(target=c.destroy))
        [t.start() for t in destroy_threads]
        [t.join() for t in destroy_threads]
        [c.delete() for c in to_destroy if c.state == Container.DESTROYED]
        if set([c.state for c in to_destroy]) != set([Container.DESTROYED]):
            err = 'aborting, failed to destroy some containers'
            log_event(self, err, logging.ERROR)
            raise RuntimeError(err)

    def deploy(self, user, release, initial=False):
        """Deploy a new release to this application"""
        existing = self.container_set.all()
        new = []
        for e in existing:
            n = e.clone(release)
            n.save()
            new.append(n)

        # create new containers
        threads = []
        for c in new:
            threads.append(threading.Thread(target=c.create))
        [t.start() for t in threads]
        [t.join() for t in threads]

        # check for containers that failed to create
        if len(new) > 0 and set([c.state for c in new]) != set([Container.CREATED]):
            err = 'aborting, failed to create some containers'
            log_event(self, err, logging.ERROR)
            self._destroy_containers(new)
            raise RuntimeError(err)

        # start new containers
        threads = []
        for c in new:
            threads.append(threading.Thread(target=c.start))
        [t.start() for t in threads]
        [t.join() for t in threads]

        # check for containers that didn't come up correctly
        if len(new) > 0 and set([c.state for c in new]) != set([Container.UP]):
            # report the deploy error
            err = 'warning, some containers failed to start'
            log_event(self, err, logging.WARNING)

        # destroy old containers
        if existing:
            self._destroy_containers(existing)

        # perform default scaling if necessary
        if initial:
            self._default_scale(user, release)

    def _default_scale(self, user, release):
        """Scale to default structure based on release type"""
        # if there is no SHA, assume a docker image is being promoted
        if not release.build.sha:
            structure = {'cmd': 1}

        # if a dockerfile exists without a procfile, assume docker workflow
        elif release.build.dockerfile and not release.build.procfile:
            structure = {'cmd': 1}

        # if a procfile exists without a web entry, assume docker workflow
        elif release.build.procfile and 'web' not in release.build.procfile:
            structure = {'cmd': 1}

        # default to heroku workflow
        else:
            structure = {'web': 1}

        self.scale(user, structure)

    def logs(self):
        """Return aggregated log data for this application."""
        path = os.path.join(settings.DEIS_LOG_DIR, self.id + '.log')
        if not os.path.exists(path):
            raise EnvironmentError('Could not locate logs')
        data = subprocess.check_output(['tail', '-n', str(settings.LOG_LINES), path])
        return data

    def run(self, user, command):
        """Run a one-off command in an ephemeral app container."""
        # TODO: add support for interactive shell
        msg = "{} runs '{}'".format(user.username, command)
        log_event(self, msg)
        c_num = max([c.num for c in self.container_set.filter(type='admin')] or [0]) + 1
        try:
            # create database record for admin process
            c = Container.objects.create(owner=self.owner,
                                         app=self,
                                         release=self.release_set.latest(),
                                         type='admin',
                                         num=c_num)
            image = c.release.image + ':v' + str(c.release.version)

            # check for backwards compatibility
            def _has_hostname(image):
                repo, tag = utils.parse_repository_tag(image)
                return True if '/' in repo and '.' in repo.split('/')[0] else False

            if not _has_hostname(image):
                image = '{}:{}/{}'.format(settings.REGISTRY_HOST,
                                          settings.REGISTRY_PORT,
                                          image)
            # SECURITY: shell-escape user input
            escaped_command = command.replace("'", "'\\''")
            return c.run(escaped_command)
        # always cleanup admin containers
        finally:
            c.delete()


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
    CRASHED = 'crashed'
    ERROR = 'error'
    STATE_CHOICES = (
        (INITIALIZED, 'initialized'),
        (CREATED, 'created'),
        (UP, 'up'),
        (DOWN, 'down'),
        (DESTROYED, 'destroyed'),
        (CRASHED, 'crashed'),
        (ERROR, 'error'),
    )

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    release = models.ForeignKey('Release')
    type = models.CharField(max_length=128, blank=False)
    num = models.PositiveIntegerField()
    state = FSMField(default=INITIALIZED, choices=STATE_CHOICES,
                     protected=True, propagate=False)

    def short_name(self):
        return "{}.{}.{}".format(self.release.app.id, self.type, self.num)
    short_name.short_description = 'Name'

    def __str__(self):
        return self.short_name()

    class Meta:
        get_latest_by = '-created'
        ordering = ['created']

    def _get_job_id(self):
        app = self.app.id
        release = self.release
        version = "v{}".format(release.version)
        num = self.num
        job_id = "{app}_{version}.{self.type}.{num}".format(**locals())
        return job_id

    _job_id = property(_get_job_id)

    def _get_scheduler(self):
        return self.app.cluster._scheduler

    _scheduler = property(_get_scheduler)

    def _get_command(self):
        # handle special case for Dockerfile deployments
        if self.type == 'cmd':
            return ''
        else:
            return 'start {}'.format(self.type)

    _command = property(_get_command)

    def _command_announceable(self):
        return self._command.lower() in ['start web', '']

    def clone(self, release):
        c = Container.objects.create(owner=self.owner,
                                     app=self.app,
                                     release=release,
                                     type=self.type,
                                     num=self.num)
        return c

    @transition(field=state, source=INITIALIZED, target=CREATED, on_error=ERROR)
    def create(self):
        image = self.release.image + ':v' + str(self.release.version)
        kwargs = {'memory': self.release.config.memory,
                  'cpu': self.release.config.cpu,
                  'tags': self.release.config.tags}
        job_id = self._job_id
        try:
            self._scheduler.create(
                name=job_id,
                image=image,
                command=self._command,
                use_announcer=self._command_announceable(), **kwargs)
        except Exception as e:
            err = '{} (create): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise

    @transition(field=state, source=[CREATED, UP, DOWN], target=UP, on_error=CRASHED)
    def start(self):
        job_id = self._job_id
        try:
            self._scheduler.start(job_id, self._command_announceable())
        except Exception as e:
            err = '{} (start): {}'.format(job_id, e)
            log_event(self.app, err, logging.WARNING)
            raise

    @transition(field=state, source=UP, target=DOWN, on_error=ERROR)
    def stop(self):
        job_id = self._job_id
        try:
            self._scheduler.stop(job_id, self._command_announceable())
        except Exception as e:
            err = '{} (stop): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise

    @transition(field=state, source='*', target=DESTROYED, on_error=ERROR)
    def destroy(self):
        job_id = self._job_id
        try:
            self._scheduler.destroy(job_id, self._command_announceable())
        except Exception as e:
            err = '{} (destroy): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise

    def run(self, command):
        """Run a one-off command"""
        image = self.release.image + ':v' + str(self.release.version)
        job_id = self._job_id
        try:
            rc, output = self._scheduler.run(job_id, image, command)
            return rc, output
        except Exception as e:
            err = '{} (run): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise


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
    procfile = JSONField(default={}, blank=True)
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
    values = JSONField(default={}, blank=True)
    memory = JSONField(default={}, blank=True)
    cpu = JSONField(default={}, blank=True)
    tags = JSONField(default={}, blank=True)

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
    image = models.CharField(max_length=256, default=settings.DEFAULT_BUILD)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'version'),)

    def __str__(self):
        return "{0}-v{1}".format(self.app.id, self.version)

    def new(self, user, config=None, build=None, summary=None, source_version='latest'):
        """
        Create a new application release using the provided Build and Config
        on behalf of a user.

        Releases start at v1 and auto-increment.
        """
        if not config:
            config = self.config
        if not build:
            build = self.build
        # always create a release off the latest image
        source_image = '{}:{}'.format(build.image, source_version)
        # construct fully-qualified target image
        new_version = self.version + 1
        tag = 'v{}'.format(new_version)
        release_image = '{}:{}'.format(self.app.id, tag)
        target_image = '{}'.format(self.app.id)
        # create new release and auto-increment version
        release = Release.objects.create(
            owner=user, app=self.app, config=config,
            build=build, version=new_version, image=target_image, summary=summary)
        # IOW, this image did not come from the builder
        # FIXME: remove check for mock registry module
        if not build.sha and 'mock' not in settings.REGISTRY_MODULE:
            # we assume that the image is not present on our registry,
            # so shell out a task to pull in the repository
            data = {
                'src': build.image
            }
            requests.post(
                '{}/v1/repositories/{}/tags'.format(settings.REGISTRY_URL,
                                                    self.app.id),
                data=data,
            )
            # update the source image to the repository we just imported
            source_image = self.app.id
            # if the image imported had a tag specified, use that tag as the source
            if ':' in build.image:
                if '/' not in build.image[build.image.rfind(':') + 1:]:
                    source_image += build.image[build.image.rfind(':'):]

        publish_release(source_image,
                        config.values,
                        release_image,)
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

    def save(self, *args, **kwargs):  # noqa
        if not self.summary:
            self.summary = ''
            prev_release = self.previous()
            # compare this build to the previous build
            old_build = prev_release.build if prev_release else None
            old_config = prev_release.config if prev_release else None
            # if the build changed, log it and who pushed it
            if self.version == 1:
                self.summary += "{} created initial release".format(self.app.owner)
            elif self.build != old_build:
                if self.build.sha:
                    self.summary += "{} deployed {}".format(self.build.owner, self.build.sha[:7])
                else:
                    self.summary += "{} deployed {}".format(self.build.owner, self.build.image)
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
                # if the limits changed (memory or cpu), log the dict diff
                changes = []
                old_mem = old_config.memory if old_config else {}
                diff = dict_diff(self.config.memory, old_mem)
                if diff.get('added') or diff.get('changed') or diff.get('deleted'):
                    changes.append('memory')
                old_cpu = old_config.cpu if old_config else {}
                diff = dict_diff(self.config.cpu, old_cpu)
                if diff.get('added') or diff.get('changed') or diff.get('deleted'):
                    changes.append('cpu')
                if changes:
                    changes = 'changed limits for '+', '.join(changes)
                    self.summary += "{} {}".format(self.config.owner, changes)
                # if the tags changed, log the dict diff
                changes = []
                old_tags = old_config.tags if old_config else {}
                diff = dict_diff(self.config.tags, old_tags)
                # try to be as succinct as possible
                added = ', '.join(k for k in diff.get('added', {}))
                added = 'added tag ' + added if added else ''
                changed = ', '.join(k for k in diff.get('changed', {}))
                changed = 'changed tag ' + changed if changed else ''
                deleted = ', '.join(k for k in diff.get('deleted', {}))
                deleted = 'deleted tag ' + deleted if deleted else ''
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
class Domain(AuditedModel):
    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    domain = models.TextField(blank=False, null=False, unique=True)

    def __str__(self):
        return self.domain


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
        log_event(build.app, "build {} created".format(build))


def _log_release_created(**kwargs):
    if kwargs.get('created'):
        release = kwargs['instance']
        log_event(release.app, "release {} created".format(release))
        # append release lifecycle logs to the app
        release.app.log(release.summary)


def _log_config_updated(**kwargs):
    config = kwargs['instance']
    log_event(config.app, "config {} updated".format(config))


def _log_domain_added(**kwargs):
    domain = kwargs['instance']
    msg = "domain {} added".format(domain)
    log_event(domain.app, msg)
    # adding a domain does not create a release, so we have to log here
    domain.app.log(msg)


def _log_domain_removed(**kwargs):
    domain = kwargs['instance']
    msg = "domain {} removed".format(domain)
    log_event(domain.app, msg)
    # adding a domain does not create a release, so we have to log here
    domain.app.log(msg)


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
    try:
        _etcd_client.delete(
            '/deis/builder/users/{}'.format(username), dir=True, recursive=True)
    except KeyError:
        # If _etcd_publish_key() wasn't called, there is no user dir to delete.
        pass


def _etcd_create_app(**kwargs):
    appname = kwargs['instance']
    if kwargs['created']:
        _etcd_client.write('/deis/services/{}'.format(appname), None, dir=True)


def _etcd_purge_app(**kwargs):
    appname = kwargs['instance']
    _etcd_client.delete('/deis/services/{}'.format(appname), dir=True, recursive=True)


def _etcd_publish_domains(**kwargs):
    app = kwargs['instance'].app
    app_domains = app.domain_set.all()
    if app_domains:
        _etcd_client.write('/deis/domains/{}'.format(app),
                           ' '.join(str(d.domain) for d in app_domains))


def _etcd_purge_domains(**kwargs):
    app = kwargs['instance'].app
    _etcd_client.delete('/deis/domains/{}'.format(app))


# Log significant app-related events
post_save.connect(_log_build_created, sender=Build, dispatch_uid='api.models.log')
post_save.connect(_log_release_created, sender=Release, dispatch_uid='api.models.log')
post_save.connect(_log_config_updated, sender=Config, dispatch_uid='api.models.log')
post_save.connect(_log_domain_added, sender=Domain, dispatch_uid='api.models.log')
post_delete.connect(_log_domain_removed, sender=Domain, dispatch_uid='api.models.log')


# save FSM transitions as they happen
def _save_transition(**kwargs):
    kwargs['instance'].save()
    # close database connections after transition
    # to avoid leaking connections inside threads
    from django.db import connection
    connection.close()

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
    post_save.connect(_etcd_publish_domains, sender=Domain, dispatch_uid='api.models')
    post_delete.connect(_etcd_purge_domains, sender=Domain, dispatch_uid='api.models')
    post_save.connect(_etcd_create_app, sender=App, dispatch_uid='api.models')
    post_delete.connect(_etcd_purge_app, sender=App, dispatch_uid='api.models')
