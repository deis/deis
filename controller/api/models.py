# -*- coding: utf-8 -*-

"""
Data models for the Deis API.
"""

from __future__ import unicode_literals
import base64
from datetime import datetime
import etcd
import importlib
import logging
import os
import re
import subprocess
import time
from threading import Thread

from django.conf import settings
from django.contrib.auth import get_user_model
from django.core.exceptions import ValidationError, SuspiciousOperation
from django.db import models
from django.db.models import Count
from django.db.models import Max
from django.db.models.signals import post_delete, post_save
from django.dispatch import receiver
from django.utils.encoding import python_2_unicode_compatible
from docker.utils import utils as dockerutils
from json_field.fields import JSONField
from OpenSSL import crypto
import requests
from rest_framework.authtoken.models import Token

from api import fields, utils
from registry import publish_release
from utils import dict_diff, fingerprint


logger = logging.getLogger(__name__)


def close_db_connections(func, *args, **kwargs):
    """
    Decorator to explicitly close db connections during threaded execution

    Note this is necessary to work around:
    https://code.djangoproject.com/ticket/22420
    """
    def _close_db_connections(*args, **kwargs):
        ret = None
        try:
            ret = func(*args, **kwargs)
        finally:
            from django.db import connections
            for conn in connections.all():
                conn.close()
        return ret
    return _close_db_connections


def log_event(app, msg, level=logging.INFO):
    # controller needs to know which app this log comes from
    logger.log(level, "{}: {}".format(app.id, msg))
    app.log(msg)


def validate_base64(value):
    """Check that value contains only valid base64 characters."""
    try:
        base64.b64decode(value.split()[1])
    except Exception as e:
        raise ValidationError(e)


def validate_id_is_docker_compatible(value):
    """
    Check that the ID follows docker's image name constraints
    """
    match = re.match(r'^[a-z0-9-]+$', value)
    if not match:
        raise ValidationError("App IDs can only contain [a-z0-9-].")


def validate_app_structure(value):
    """Error if the dict values aren't ints >= 0."""
    try:
        if any(int(v) < 0 for v in value.viewvalues()):
            raise ValueError("Must be greater than or equal to zero")
    except ValueError, err:
        raise ValidationError(err)


def validate_reserved_names(value):
    """A value cannot use some reserved names."""
    if value in settings.DEIS_RESERVED_NAMES:
        raise ValidationError('{} is a reserved name.'.format(value))


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


def validate_certificate(value):
    try:
        crypto.load_certificate(crypto.FILETYPE_PEM, value)
    except crypto.Error as e:
        raise ValidationError('Could not load certificate: {}'.format(e))


class AuditedModel(models.Model):
    """Add created and updated fields to a model."""

    created = models.DateTimeField(auto_now_add=True)
    updated = models.DateTimeField(auto_now=True)

    class Meta:
        """Mark :class:`AuditedModel` as abstract."""
        abstract = True


def select_app_name():
    """Select a unique randomly generated app name"""
    name = utils.generate_app_name()

    while App.objects.filter(id=name).exists():
        name = utils.generate_app_name()

    return name


class UuidAuditedModel(AuditedModel):
    """Add a UUID primary key to an :class:`AuditedModel`."""

    uuid = fields.UuidField('UUID', primary_key=True)

    class Meta:
        """Mark :class:`UuidAuditedModel` as abstract."""
        abstract = True


@python_2_unicode_compatible
class App(UuidAuditedModel):
    """
    Application used to service requests on behalf of end-users
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64, unique=True, default=select_app_name,
                          validators=[validate_id_is_docker_compatible,
                                      validate_reserved_names])
    structure = JSONField(default={}, blank=True, validators=[validate_app_structure])

    class Meta:
        permissions = (('use_app', 'Can use app'),)

    @property
    def _scheduler(self):
        mod = importlib.import_module(settings.SCHEDULER_MODULE)
        return mod.SchedulerClient(settings.SCHEDULER_TARGET,
                                   settings.SCHEDULER_AUTH,
                                   settings.SCHEDULER_OPTIONS,
                                   settings.SSH_PRIVATE_KEY)

    def __str__(self):
        return self.id

    @property
    def url(self):
        return self.id + '.' + settings.DEIS_DOMAIN

    def log(self, message):
        """Logs a message to the application's log file.

        This is a workaround for how Django interacts with Python's logging module. Each app
        needs its own FileHandler instance so it can write to its own log file. That won't work in
        Django's case because logging is set up before you run the server and it disables all
        existing logging configurations.
        """
        with open(os.path.join(settings.DEIS_LOG_DIR, self.id + '.log'), 'a') as f:
            msg = "{} deis[api]: {}\n".format(time.strftime(settings.DEIS_DATETIME_FORMAT),
                                              message)
            f.write(msg.encode('utf-8'))

    def create(self, *args, **kwargs):
        """Create a new application with an initial config and release"""
        config = Config.objects.create(owner=self.owner, app=self)
        Release.objects.create(version=1, owner=self.owner, app=self, config=config, build=None)

    def delete(self, *args, **kwargs):
        """Delete this application including all containers"""
        try:
            # attempt to remove containers from the scheduler
            self._destroy_containers([c for c in self.container_set.exclude(type='run')])
        except RuntimeError:
            pass
        self._clean_app_logs()
        return super(App, self).delete(*args, **kwargs)

    def restart(self, **kwargs):
        to_restart = self.container_set.all()
        if kwargs.get('type'):
            to_restart = to_restart.filter(type=kwargs.get('type'))
        if kwargs.get('num'):
            to_restart = to_restart.filter(num=kwargs.get('num'))
        self._restart_containers(to_restart)
        return to_restart

    def _clean_app_logs(self):
        """Delete application logs stored by the logger component"""
        path = os.path.join(settings.DEIS_LOG_DIR, self.id + '.log')
        if os.path.exists(path):
            os.remove(path)

    def scale(self, user, structure):  # noqa
        """Scale containers up or down to match requested structure."""
        if self.release_set.latest().build is None:
            raise EnvironmentError('No build associated with this release')
        requested_structure = structure.copy()
        release = self.release_set.latest()
        # test for available process types
        available_process_types = release.build.procfile or {}
        for container_type in requested_structure:
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
        # iterate on a copy of the container_type keys
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
                if settings.SCHEDULER_MODULE == "scheduler.k8s" :
                    if len(containers) == 0:
                        t = Thread(target=c.scale)
                        t.start()
                        t.join()
                    c.delete()
                else:
                    to_remove.append(c)
                diff += 1
            if diff == 0 and len(containers) != 0:
                c = containers.pop()
                self._start_containers([c])
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
                if settings.SCHEDULER_MODULE == "scheduler.k8s" :
                    self._start_containers([to_add.pop()])
                else:
                    self._start_containers(to_add)
            if to_remove:
                self._destroy_containers(to_remove)
        # save new structure to the database
        vals = self.container_set.exclude(type='run').values(
            'type').annotate(Count('pk')).order_by()
        new_structure = structure.copy()
        new_structure.update({v['type']: v['pk__count'] for v in vals})
        self.structure = new_structure
        self.save()
        return changed

    def _start_containers(self, to_add):
        """Creates and starts containers via the scheduler"""
        if not to_add:
            return
        create_threads = [Thread(target=c.create) for c in to_add]
        start_threads = [Thread(target=c.start) for c in to_add]
        [t.start() for t in create_threads]
        [t.join() for t in create_threads]
        if settings.SCHEDULER_MODULE != "scheduler.k8s" and any(c.state != 'created' for c in to_add):
            err = 'aborting, failed to create some containers '+c.state
            log_event(self, err, logging.ERROR)
            self._destroy_containers(to_add)
            raise RuntimeError(err)
        [t.start() for t in start_threads]
        [t.join() for t in start_threads]
        if set([c.state for c in to_add]) != set(['up']):
            err = 'warning, some containers failed to start'+c.state
            log_event(self, err, logging.WARNING)

    def _restart_containers(self, to_restart):
        """Restarts containers via the scheduler"""
        if not to_restart:
            return
        stop_threads = [Thread(target=c.stop) for c in to_restart]
        start_threads = [Thread(target=c.start) for c in to_restart]
        [t.start() for t in stop_threads]
        [t.join() for t in stop_threads]
        if any(c.state != 'created' for c in to_restart):
            err = 'warning, some containers failed to stop'
            log_event(self, err, logging.WARNING)
        [t.start() for t in start_threads]
        [t.join() for t in start_threads]
        if any(c.state != 'up' for c in to_restart):
            err = 'warning, some containers failed to start'
            log_event(self, err, logging.WARNING)

    def _destroy_containers(self, to_destroy):
        """Destroys containers via the scheduler"""
        if not to_destroy:
            return
        if settings.SCHEDULER_MODULE == "scheduler.k8s" :
            destroy_threads = [Thread(target=to_destroy[0].destroy)]
        else:
            destroy_threads = [Thread(target=c.destroy) for c in to_destroy]
        [t.start() for t in destroy_threads]
        [t.join() for t in destroy_threads]
        if settings.SCHEDULER_MODULE == "scheduler.k8s" :
            [c.delete() for c in to_destroy]
        else :
            [c.delete() for c in to_destroy if c.state == 'destroyed']
        if any(c.state != 'destroyed' for c in to_destroy):
            err = 'aborting, failed to destroy some containers'+c.state
            log_event(self, err, logging.ERROR)
            raise RuntimeError(err)

    def deploy(self, user, release):
        """Deploy a new release to this application"""
        existing = self.container_set.exclude(type='run')
        new = []
        for e in existing:
            n = e.clone(release)
            n.save()
            new.append(n)

        if new and settings.SCHEDULER_MODULE == "scheduler.k8s" :
            self._start_containers([new.pop()])
        else :
            self._start_containers(new)

        # destroy old containers
        if existing:
            self._destroy_containers(existing)

        # perform default scaling if necessary
        if self.structure == {} and release.build is not None:
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

    def logs(self, log_lines):
        """Return aggregated log data for this application."""
        path = os.path.join(settings.DEIS_LOG_DIR, self.id + '.log')
        if not os.path.exists(path):
            raise EnvironmentError('Could not locate logs')
        data = subprocess.check_output(['tail', '-n', log_lines, path])
        return data

    def run(self, user, command):
        """Run a one-off command in an ephemeral app container."""
        # FIXME: remove the need for SSH private keys by using
        # a scheduler that supports one-off admin tasks natively
        if not settings.SSH_PRIVATE_KEY:
            raise EnvironmentError('Support for admin commands is not configured')
        if self.release_set.latest().build is None:
            raise EnvironmentError('No build associated with this release to run this command')
        # TODO: add support for interactive shell
        msg = "{} runs '{}'".format(user.username, command)
        log_event(self, msg)
        c_num = max([c.num for c in self.container_set.filter(type='run')] or [0]) + 1

        # create database record for run process
        c = Container.objects.create(owner=self.owner,
                                     app=self,
                                     release=self.release_set.latest(),
                                     type='run',
                                     num=c_num)
        image = c.release.image

        # check for backwards compatibility
        def _has_hostname(image):
            repo, tag = dockerutils.parse_repository_tag(image)
            return True if '/' in repo and '.' in repo.split('/')[0] else False

        if not _has_hostname(image):
            image = '{}:{}/{}'.format(settings.REGISTRY_HOST,
                                      settings.REGISTRY_PORT,
                                      image)
        # SECURITY: shell-escape user input
        escaped_command = command.replace("'", "'\\''")
        return c.run(escaped_command)


@python_2_unicode_compatible
class Container(UuidAuditedModel):
    """
    Docker container used to securely host an application process.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    release = models.ForeignKey('Release')
    type = models.CharField(max_length=128, blank=False)
    num = models.PositiveIntegerField()

    @property
    def _scheduler(self):
        return self.app._scheduler

    @property
    def state(self):
        return self._scheduler.state(self._job_id).name

    def short_name(self):
        return "{}.{}.{}".format(self.app.id, self.type, self.num)
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

    def _get_command(self):
        try:
            # if this is not procfile-based app, ensure they cannot break out
            # and run arbitrary commands on the host
            # FIXME: remove slugrunner's hardcoded entrypoint
            if self.release.build.dockerfile or not self.release.build.sha:
                return "bash -c '{}'".format(self.release.build.procfile[self.type])
            else:
                return 'start {}'.format(self.type)
        # if the key is not present or if a parent attribute is None
        except (KeyError, TypeError, AttributeError):
            # handle special case for Dockerfile deployments
            return '' if self.type == 'cmd' else 'start {}'.format(self.type)

    _command = property(_get_command)

    def clone(self, release):
        c = Container.objects.create(owner=self.owner,
                                     app=self.app,
                                     release=release,
                                     type=self.type,
                                     num=self.num)
        return c

    @close_db_connections
    def create(self):
        image = self.release.image
        kwargs = {'memory': self.release.config.memory,
                  'cpu': self.release.config.cpu,
                  'tags': self.release.config.tags,
                  'aname': self.app.id,
                  'num': self.num}
        job_id = self._job_id
        try:
            self._scheduler.create(
                name=job_id,
                image=image,
                command=self._command,
                **kwargs)
        except Exception as e:
            err = '{} (create): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise

    @close_db_connections
    def scale(self):
        image = self.release.image
        job_id = self._job_id
        try:
            self._scheduler.scale(
                name=job_id,
                image=image,
                num=0)
        except Exception as e:
            err = '{} (scale): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise

    @close_db_connections
    def start(self):
        job_id = self._job_id
        try:
            self._scheduler.start(job_id)
        except Exception as e:
            err = '{} (start): {}'.format(job_id, e)
            log_event(self.app, err, logging.WARNING)
            raise

    @close_db_connections
    def stop(self):
        job_id = self._job_id
        try:
            self._scheduler.stop(job_id)
        except Exception as e:
            err = '{} (stop): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise

    @close_db_connections
    def destroy(self):
        job_id = self._job_id
        try:
            self._scheduler.destroy(job_id)
        except Exception as e:
            err = '{} (destroy): {}'.format(job_id, e)
            log_event(self.app, err, logging.ERROR)
            raise

    def run(self, command):
        """Run a one-off command"""
        if self.release.build is None:
            raise EnvironmentError('No build associated with this release '
                                   'to run this command')
        image = self.release.image
        job_id = self._job_id
        entrypoint = '/bin/bash'
        # if this is a procfile-based app, switch the entrypoint to slugrunner's default
        # FIXME: remove slugrunner's hardcoded entrypoint
        if self.release.build.procfile and \
           self.release.build.sha and not \
           self.release.build.dockerfile:
            entrypoint = '/runner/init'
            command = "'{}'".format(command)
        else:
            command = "-c '{}'".format(command)
        try:
            rc, output = self._scheduler.run(job_id, image, entrypoint, command)
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

    def create(self, user, *args, **kwargs):
        latest_release = self.app.release_set.latest()
        source_version = 'latest'
        if self.sha:
            source_version = 'git-{}'.format(self.sha)
        new_release = latest_release.new(user,
                                         build=self,
                                         config=latest_release.config,
                                         source_version=source_version)
        try:
            self.app.deploy(user, new_release)
            return new_release
        except RuntimeError:
            new_release.delete()
            raise

    def save(self, **kwargs):
        try:
            previous_build = self.app.build_set.latest()
            to_destroy = []
            for proctype in previous_build.procfile:
                if proctype not in self.procfile:
                    for c in self.app.container_set.filter(type=proctype):
                        to_destroy.append(c)
            self.app._destroy_containers(to_destroy)
        except Build.DoesNotExist:
            pass
        return super(Build, self).save(**kwargs)

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

    def save(self, **kwargs):
        """merge the old config with the new"""
        try:
            previous_config = self.app.config_set.latest()
            for attr in ['cpu', 'memory', 'tags', 'values']:
                # Guard against migrations from older apps without fixes to
                # JSONField encoding.
                try:
                    data = getattr(previous_config, attr).copy()
                except AttributeError:
                    data = {}
                try:
                    new_data = getattr(self, attr).copy()
                except AttributeError:
                    new_data = {}
                data.update(new_data)
                # remove config keys if we provided a null value
                [data.pop(k) for k, v in new_data.viewitems() if v is None]
                setattr(self, attr, data)
        except Config.DoesNotExist:
            pass
        return super(Config, self).save(**kwargs)


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
    build = models.ForeignKey('Build', null=True)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'version'),)

    def __str__(self):
        return "{0}-v{1}".format(self.app.id, self.version)

    @property
    def image(self):
        return '{}:v{}'.format(self.app.id, str(self.version))

    def new(self, user, config, build, summary=None, source_version='latest'):
        """
        Create a new application release using the provided Build and Config
        on behalf of a user.

        Releases start at v1 and auto-increment.
        """
        # construct fully-qualified target image
        new_version = self.version + 1
        # create new release and auto-increment version
        release = Release.objects.create(
            owner=user, app=self.app, config=config,
            build=build, version=new_version, summary=summary)
        try:
            release.publish()
        except EnvironmentError as e:
            # If we cannot publish this app, just log and carry on
            log_event(self.app, e)
            pass
        return release

    def publish(self, source_version='latest'):
        if self.build is None:
            raise EnvironmentError('No build associated with this release to publish')
        source_tag = 'git-{}'.format(self.build.sha) if self.build.sha else source_version
        source_image = '{}:{}'.format(self.build.image, source_tag)
        # IOW, this image did not come from the builder
        # FIXME: remove check for mock registry module
        if not self.build.sha and 'mock' not in settings.REGISTRY_MODULE:
            # we assume that the image is not present on our registry,
            # so shell out a task to pull in the repository
            data = {
                'src': self.build.image
            }
            requests.post(
                '{}/v1/repositories/{}/tags'.format(settings.REGISTRY_URL,
                                                    self.app.id),
                data=data,
            )
            # update the source image to the repository we just imported
            source_image = self.app.id
            # if the image imported had a tag specified, use that tag as the source
            if ':' in self.build.image:
                if '/' not in self.build.image[self.build.image.rfind(':') + 1:]:
                    source_image += self.build.image[self.build.image.rfind(':'):]
        publish_release(source_image,
                        self.config.values,
                        self.image)

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

    def rollback(self, user, version):
        if version < 1:
            raise EnvironmentError('version cannot be below 0')
        summary = "{} rolled back to v{}".format(user, version)
        prev = self.app.release_set.get(version=version)
        new_release = self.new(
            user,
            build=prev.build,
            config=prev.config,
            summary=summary,
            source_version='v{}'.format(version))
        try:
            self.app.deploy(user, new_release)
            return new_release
        except RuntimeError:
            new_release.delete()
            raise

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
class Certificate(AuditedModel):
    """
    Public and private key pair used to secure application traffic at the router.
    """
    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    # there is no upper limit on the size of an x.509 certificate
    certificate = models.TextField(validators=[validate_certificate])
    key = models.TextField()
    # X.509 certificates allow any string of information as the common name.
    common_name = models.TextField(unique=True)
    expires = models.DateTimeField()

    def __str__(self):
        return self.common_name

    def _get_certificate(self):
        try:
            return crypto.load_certificate(crypto.FILETYPE_PEM, self.certificate)
        except crypto.Error as e:
            raise SuspiciousOperation(e)

    def save(self, *args, **kwargs):
        certificate = self._get_certificate()
        if not self.common_name:
            self.common_name = certificate.get_subject().CN
        if not self.expires:
            # convert openssl's expiry date format to Django's DateTimeField format
            self.expires = datetime.strptime(certificate.get_notAfter(), '%Y%m%d%H%M%SZ')
        return super(Certificate, self).save(*args, **kwargs)


@python_2_unicode_compatible
class Key(UuidAuditedModel):
    """An SSH public key."""

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=128)
    public = models.TextField(unique=True, validators=[validate_base64])
    fingerprint = models.CharField(max_length=128)

    class Meta:
        verbose_name = 'SSH Key'
        unique_together = (('owner', 'fingerprint'))

    def __str__(self):
        return "{}...{}".format(self.public[:18], self.public[-31:])

    def save(self, *args, **kwargs):
        self.fingerprint = fingerprint(self.public)
        return super(Key, self).save(*args, **kwargs)


# define update/delete callbacks for synchronizing
# models with the configuration management backend

def _log_build_created(**kwargs):
    if kwargs.get('created'):
        build = kwargs['instance']
        # log only to the controller; this event will be logged in the release summary
        logger.info("{}: build {} created".format(build.app, build))


def _log_release_created(**kwargs):
    if kwargs.get('created'):
        release = kwargs['instance']
        # log only to the controller; this event will be logged in the release summary
        logger.info("{}: release {} created".format(release.app, release))
        # append release lifecycle logs to the app
        release.app.log(release.summary)


def _log_config_updated(**kwargs):
    config = kwargs['instance']
    # log only to the controller; this event will be logged in the release summary
    logger.info("{}: config {} updated".format(config.app, config))


def _log_domain_added(**kwargs):
    domain = kwargs['instance']
    msg = "domain {} added".format(domain)
    log_event(domain.app, msg)


def _log_domain_removed(**kwargs):
    domain = kwargs['instance']
    msg = "domain {} removed".format(domain)
    log_event(domain.app, msg)


def _log_cert_added(**kwargs):
    cert = kwargs['instance']
    logger.info("cert {} added".format(cert))


def _log_cert_removed(**kwargs):
    cert = kwargs['instance']
    logger.info("cert {} removed".format(cert))


def _etcd_publish_key(**kwargs):
    key = kwargs['instance']
    _etcd_client.write('/deis/builder/users/{}/{}'.format(
        key.owner.username, fingerprint(key.public)), key.public)


def _etcd_purge_key(**kwargs):
    key = kwargs['instance']
    try:
        _etcd_client.delete('/deis/builder/users/{}/{}'.format(
            key.owner.username, fingerprint(key.public)))
    except KeyError:
        pass


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
    try:
        _etcd_client.delete('/deis/services/{}'.format(appname), dir=True, recursive=True)
    except KeyError:
        pass


def _etcd_publish_cert(**kwargs):
    cert = kwargs['instance']
    if kwargs['created']:
        _etcd_client.write('/deis/certs/{}/cert'.format(cert), cert.certificate)
        _etcd_client.write('/deis/certs/{}/key'.format(cert), cert.key)


def _etcd_purge_cert(**kwargs):
    cert = kwargs['instance']
    try:
        _etcd_client.delete('/deis/certs/{}'.format(cert),
                            prevExist=True, dir=True, recursive=True)
    except KeyError:
        pass


def _etcd_publish_domains(**kwargs):
    domain = kwargs['instance']
    if kwargs['created']:
        _etcd_client.write('/deis/domains/{}'.format(domain), domain.app)


def _etcd_purge_domains(**kwargs):
    domain = kwargs['instance']
    try:
        _etcd_client.delete('/deis/domains/{}'.format(domain),
                            prevExist=True, dir=True, recursive=True)
    except KeyError:
        pass


# Log significant app-related events
post_save.connect(_log_build_created, sender=Build, dispatch_uid='api.models.log')
post_save.connect(_log_release_created, sender=Release, dispatch_uid='api.models.log')
post_save.connect(_log_config_updated, sender=Config, dispatch_uid='api.models.log')
post_save.connect(_log_domain_added, sender=Domain, dispatch_uid='api.models.log')
post_save.connect(_log_cert_added, sender=Certificate, dispatch_uid='api.models.log')
post_delete.connect(_log_domain_removed, sender=Domain, dispatch_uid='api.models.log')
post_delete.connect(_log_cert_removed, sender=Certificate, dispatch_uid='api.models.log')


# automatically generate a new token on creation
@receiver(post_save, sender=get_user_model())
def create_auth_token(sender, instance=None, created=False, **kwargs):
    if created:
        Token.objects.create(user=instance)

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
    post_delete.connect(_etcd_purge_user, sender=get_user_model(), dispatch_uid='api.models')
    post_save.connect(_etcd_publish_domains, sender=Domain, dispatch_uid='api.models')
    post_delete.connect(_etcd_purge_domains, sender=Domain, dispatch_uid='api.models')
    post_save.connect(_etcd_create_app, sender=App, dispatch_uid='api.models')
    post_delete.connect(_etcd_purge_app, sender=App, dispatch_uid='api.models')
    post_save.connect(_etcd_publish_cert, sender=Certificate, dispatch_uid='api.models')
    post_delete.connect(_etcd_purge_cert, sender=Certificate, dispatch_uid='api.models')
