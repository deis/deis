# -*- coding: utf-8 -*-

"""
Data models for the Deis API.
"""

from __future__ import unicode_literals
import importlib
import logging
import os
import subprocess

from celery.canvas import group
from django.conf import settings
from django.contrib.auth.models import User
from django.db import models
from django.db.models.signals import post_delete
from django.db.models.signals import post_save
from django.dispatch import receiver
from django.dispatch.dispatcher import Signal
from django.utils.encoding import python_2_unicode_compatible
import etcd
from guardian.shortcuts import get_users_with_perms
from json_field.fields import JSONField  # @UnusedImport

from api import fields, tasks
from provider import import_provider_module
from utils import dict_diff, fingerprint


logger = logging.getLogger(__name__)

# import user-defined configuration management module
CM = importlib.import_module(settings.CM_MODULE)


# define custom signals
release_signal = Signal(providing_args=['user', 'app'])


# base models

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


# deis core models

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

    def save(self, *args, **kwargs):
        super(Key, self).save(*args, **kwargs)
        self.owner.publish()

    def delete(self, *args, **kwargs):
        super(Key, self).delete(*args, **kwargs)
        self.owner.publish()


class ProviderManager(models.Manager):
    """Manage database interactions for :class:`Provider`."""

    def seed(self, user, **kwargs):
        """
        Seeds the database with Providers for clouds supported by Deis.
        """
        providers = [(p, p) for p in settings.PROVIDER_MODULES]
        for p_id, p_type in providers:
            self.create(owner=user, id=p_id, type=p_type, creds='{}')


@python_2_unicode_compatible
class Provider(UuidAuditedModel):
    """Cloud provider settings for a user.

    Available as `user.provider_set`.
    """

    objects = ProviderManager()

    PROVIDERS = (
        ('ec2', 'Amazon Elastic Compute Cloud (EC2)'),
        ('mock', 'Mock Reference Provider'),
        ('rackspace', 'Rackspace Open Cloud'),
        ('static', 'Static Node'),
        ('digitalocean', 'Digital Ocean'),
        ('vagrant', 'Local Vagrant VMs'),
    )

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    type = models.SlugField(max_length=16, choices=PROVIDERS)
    creds = JSONField(blank=True)

    class Meta:
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return "{}-{}".format(self.id, self.get_type_display())


class FlavorManager(models.Manager):
    """Manage database interactions for :class:`Flavor`\s."""

    def seed(self, user, **kwargs):
        """Seed the database with default Flavors for each cloud region."""
        for provider_type in settings.PROVIDER_MODULES:
            provider = import_provider_module(provider_type)
            flavors = provider.seed_flavors()
            p = Provider.objects.get(owner=user, id=provider_type)
            for flavor in flavors:
                flavor['provider'] = p
                Flavor.objects.create(owner=user, **flavor)


@python_2_unicode_compatible
class Flavor(UuidAuditedModel):
    """
    Virtual machine flavors associated with a Provider

    Params is a JSON field including unstructured data
    for provider API calls, like region, zone, and size.
    """
    objects = FlavorManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    provider = models.ForeignKey('Provider')
    params = JSONField(blank=True)

    class Meta:
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return self.id


@python_2_unicode_compatible
class Formation(UuidAuditedModel):
    """
    Formation of nodes used to host applications
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64, unique=True)
    domain = models.CharField(max_length=128, blank=True, null=True)
    nodes = JSONField(default='{}', blank=True)

    class Meta:
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'domain': self.domain,
                'nodes': self.nodes}

    def build(self):
        return

    def destroy(self, *args, **kwargs):
        for app in self.app_set.all():
            app.destroy()
        node_tasks = [tasks.destroy_node.si(n) for n in self.node_set.all()]
        layer_tasks = [tasks.destroy_layer.si(l) for l in self.layer_set.all()]
        group(node_tasks).apply_async().join()
        group(layer_tasks).apply_async().join()
        CM.purge_formation(self.flat())
        self.delete()
        tasks.converge_controller.apply_async().wait()

    def publish(self):
        data = self.calculate()
        CM.publish_formation(self.flat(), data)
        return data

    def converge(self, controller=False, **kwargs):
        databag = self.publish()
        nodes = self.node_set.all()
        subtasks = []
        for n in nodes:
            subtask = tasks.converge_node.si(n)
            subtasks.append(subtask)
        if controller is True:
            subtasks.append(tasks.converge_controller.si())
        group(*subtasks).apply_async().join()
        return databag

    def calculate(self):
        """Return a representation of this formation for config management"""
        d = {}
        d['id'] = self.id
        d['domain'] = self.domain
        d['nodes'] = {}
        proxies = []
        for n in self.node_set.all():
            d['nodes'][n.id] = {'fqdn': n.fqdn,
                                'runtime': n.layer.runtime,
                                'proxy': n.layer.proxy}
            if n.layer.proxy is True:
                proxies.append(n.fqdn)
        d['apps'] = {}
        for a in self.app_set.all():
            d['apps'][a.id] = a.calculate()
            d['apps'][a.id]['proxy'] = {}
            d['apps'][a.id]['proxy']['nodes'] = proxies
            d['apps'][a.id]['proxy']['algorithm'] = 'round_robin'
            d['apps'][a.id]['proxy']['port'] = 80
            d['apps'][a.id]['proxy']['backends'] = []
            d['apps'][a.id]['containers'] = containers = {}
            for c in a.container_set.all().order_by('created'):
                containers.setdefault(c.type, {})
                containers[c.type].update(
                    {c.num: "{0}:{1}".format(c.node.id, c.port)})
                if c.type == 'web':
                    d['apps'][a.id]['proxy']['backends'].append(
                        "{0}:{1}".format(c.node.fqdn, c.port))
        return d


@python_2_unicode_compatible
class Layer(UuidAuditedModel):
    """
    Layer of nodes used by the formation

    All nodes in a layer share the same flavor and configuration.

    The layer stores SSH settings used to trigger node convergence,
    as well as other configuration used during node bootstrapping
    (e.g. Chef Run List, Chef Environment)
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)

    formation = models.ForeignKey('Formation')
    flavor = models.ForeignKey('Flavor')

    proxy = models.BooleanField(default=False)
    runtime = models.BooleanField(default=False)

    ssh_username = models.CharField(max_length=64, default='ubuntu')
    ssh_private_key = models.TextField()
    ssh_public_key = models.TextField()
    ssh_port = models.SmallIntegerField(default=22)

    # example: {'run_list': [deis::runtime'], 'environment': 'dev'}
    config = JSONField(default='{}', blank=True)

    class Meta:
        unique_together = (('formation', 'id'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'provider_type': self.flavor.provider.type,
                'creds': dict(self.flavor.provider.creds),
                'formation': self.formation.id,
                'flavor': self.flavor.id,
                'params': dict(self.flavor.params),
                'proxy': self.proxy,
                'runtime': self.runtime,
                'ssh_username': self.ssh_username,
                'ssh_private_key': self.ssh_private_key,
                'ssh_public_key': self.ssh_public_key,
                'ssh_port': self.ssh_port,
                'config': dict(self.config)}

    def build(self):
        return tasks.build_layer.delay(self).wait()

    def destroy(self):
        return tasks.destroy_layer.delay(self).wait()


class NodeManager(models.Manager):

    def new(self, formation, layer, fqdn=None):
        existing_nodes = self.filter(formation=formation, layer=layer).order_by('-created')
        if existing_nodes:
            next_num = existing_nodes[0].num + 1
        else:
            next_num = 1
        node = self.create(owner=formation.owner,
                           formation=formation,
                           layer=layer,
                           num=next_num,
                           id="{0}-{1}-{2}".format(formation.id, layer.id, next_num),
                           fqdn=fqdn)
        return node

    def scale(self, formation, structure, **kwargs):
        """Scale layers up or down to match requested structure."""
        funcs = []
        changed = False
        for layer_id, requested in structure.items():
            layer = formation.layer_set.get(id=layer_id)
            nodes = list(layer.node_set.all().order_by('created'))
            diff = requested - len(nodes)
            if diff == 0:
                continue
            while diff < 0:
                node = nodes.pop(0)
                funcs.append(tasks.destroy_node.si(node))
                diff = requested - len(nodes)
                changed = True
            while diff > 0:
                node = self.new(formation, layer)
                nodes.append(node)
                funcs.append(tasks.build_node.si(node))
                diff = requested - len(nodes)
                changed = True
        # launch/terminate nodes in parallel
        if funcs:
            group(*funcs).apply_async().join()
        # always scale and balance every application
        if nodes:
            for app in formation.app_set.all():
                Container.objects.scale(app, app.containers)
                Container.objects.balance(formation)
        # save new structure now that scaling was successful
        formation.nodes.update(structure)
        formation.save()
        # force-converge nodes if there were new nodes or container rebalancing
        if changed:
            return formation.converge()
        return formation.calculate()

    def next_runtime_node(self, formation, container_type, reverse=False):
        count = []
        layers = formation.layer_set.filter(runtime=True)
        runtime_nodes = []
        for l in layers:
            runtime_nodes.extend(Node.objects.filter(
                formation=formation, layer=l).order_by('created'))
        container_map = {n: [] for n in runtime_nodes}
        containers = list(Container.objects.filter(
            formation=formation, type=container_type).order_by('created'))
        for c in containers:
            container_map[c.node].append(c)
        for n in container_map.keys():
            # (2, node3), (2, node2), (3, node1)
            count.append((len(container_map[n]), n))
        if not count:
            raise EnvironmentError('No nodes available for containers')
        count.sort()
        # reverse means order by greatest # of containers, otherwise fewest
        if reverse:
            count.reverse()
        return count[0][1]

    def next_runtime_port(self, formation):
        containers = Container.objects.filter(formation=formation).order_by('-port')
        if not containers:
            return 10001
        return containers[0].port + 1


@python_2_unicode_compatible
class Node(UuidAuditedModel):
    """
    Node used to host containers

    List of nodes available as `formation.nodes`
    """

    objects = NodeManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=64)
    formation = models.ForeignKey('Formation')
    layer = models.ForeignKey('Layer')
    num = models.PositiveIntegerField()

    # TODO: add celery beat tasks for monitoring node health
    status = models.CharField(max_length=64, default='up')

    provider_id = models.SlugField(max_length=64, blank=True, null=True)
    fqdn = models.CharField(max_length=256, blank=True, null=True)
    status = JSONField(blank=True, null=True)

    class Meta:
        unique_together = (('formation', 'id'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'provider_type': self.layer.flavor.provider.type,
                'formation': self.formation.id,
                'layer': self.layer.id,
                'creds': dict(self.layer.flavor.provider.creds),
                'params': dict(self.layer.flavor.params),
                'runtime': self.layer.runtime,
                'proxy': self.layer.proxy,
                'ssh_username': self.layer.ssh_username,
                'ssh_public_key': self.layer.ssh_public_key,
                'ssh_private_key': self.layer.ssh_private_key,
                'ssh_port': self.layer.ssh_port,
                'config': dict(self.layer.config),
                'provider_id': self.provider_id,
                'fqdn': self.fqdn}

    def build(self):
        return tasks.build_node.delay(self).wait()

    def destroy(self):
        return tasks.destroy_node.delay(self).wait()

    def converge(self):
        return tasks.converge_node.delay(self).wait()

    def run(self, command, **kwargs):
        return tasks.run_node.delay(self, command).wait()


def log_event(app, msg, level=logging.INFO):
    msg = "{}: {}".format(app.id, msg)
    logger.log(level, msg)


@python_2_unicode_compatible
class App(UuidAuditedModel):
    """
    Application used to service requests on behalf of end-users
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64, unique=True)
    formation = models.ForeignKey('Formation')
    containers = JSONField(default='{}', blank=True)

    class Meta:
        permissions = (('use_app', 'Can use app'),)

    def __str__(self):
        return self.id

    def flat(self):
        return {'id': self.id,
                'formation': self.formation.id,
                'containers': dict(self.containers)}

    def build(self):
        config = Config.objects.create(
            version=1, owner=self.owner, app=self, values={})
        build = Build.objects.create(owner=self.owner, app=self)
        Release.objects.create(
            version=1, owner=self.owner, app=self, config=config, build=build)
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


class ContainerManager(models.Manager):

    def scale(self, app, structure, **kwargs):
        """Scale containers up or down to match requested."""
        requested_containers = structure.copy()
        formation = app.formation
        # increment new container nums off the most recent container
        all_containers = app.container_set.all().order_by('-created')
        container_num = 1 if not all_containers else all_containers[0].num + 1
        msg = 'Containers scaled ' + ' '.join(
            "{}={}".format(k, v) for k, v in requested_containers.items())
        # iterate and scale by container type (web, worker, etc)
        changed = False
        for container_type in requested_containers.keys():
            containers = list(app.container_set.filter(type=container_type).order_by('created'))
            requested = requested_containers.pop(container_type)
            diff = requested - len(containers)
            if diff == 0:
                continue
            changed = True
            while diff < 0:
                # get the next node with the most containers
                node = Node.objects.next_runtime_node(
                    formation, container_type, reverse=True)
                # delete a container attached to that node
                for c in containers:
                    if node == c.node:
                        containers.remove(c)
                        c.delete()
                        diff += 1
                        break
            while diff > 0:
                # get the next node with the fewest containers
                node = Node.objects.next_runtime_node(formation, container_type)
                port = Node.objects.next_runtime_port(formation)
                c = Container.objects.create(owner=app.owner,
                                             formation=formation,
                                             node=node,
                                             app=app,
                                             type=container_type,
                                             num=container_num,
                                             port=port)
                containers.append(c)
                container_num += 1
                diff -= 1
        log_event(app, msg)
        return changed

    def balance(self, formation, **kwargs):
        runtime_nodes = formation.node_set.filter(layer__runtime=True).order_by('created')
        all_containers = self.filter(formation=formation).order_by('-created')
        # get the next container number (e.g. web.19)
        container_num = 1 if not all_containers else all_containers[0].num + 1
        changed = False
        app = None
        # iterate by unique container type
        for container_type in set([c.type for c in all_containers]):
            # map node container counts => { 2: [b3, b4], 3: [ b1, b2 ] }
            n_map = {}
            for node in runtime_nodes:
                ct = len(node.container_set.filter(type=container_type))
                n_map.setdefault(ct, []).append(node)
            # loop until diff between min and max is 1 or 0
            while max(n_map.keys()) - min(n_map.keys()) > 1:
                # get the most over-utilized node
                n_max = max(n_map.keys())
                n_over = n_map[n_max].pop(0)
                if len(n_map[n_max]) == 0:
                    del n_map[n_max]
                # get the most under-utilized node
                n_min = min(n_map.keys())
                n_under = n_map[n_min].pop(0)
                if len(n_map[n_min]) == 0:
                    del n_map[n_min]
                # delete the oldest container from the most over-utilized node
                c = n_over.container_set.filter(type=container_type).order_by('created')[0]
                app = c.app  # pull ref to app for recreating the container
                c.delete()
                # create a container on the most under-utilized node
                self.create(owner=formation.owner,
                            formation=formation,
                            app=app,
                            type=container_type,
                            num=container_num,
                            node=n_under,
                            port=Node.objects.next_runtime_port(formation))
                container_num += 1
                # update the n_map accordingly
                for n in (n_over, n_under):
                    ct = len(n.container_set.filter(type=container_type))
                    n_map.setdefault(ct, []).append(n)
                changed = True
        if app:
            log_event(app, 'Containers balanced')
        return changed


@python_2_unicode_compatible
class Container(UuidAuditedModel):
    """
    Docker container used to securely host an application process.
    """

    objects = ContainerManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey('Formation')
    node = models.ForeignKey('Node')
    app = models.ForeignKey('App')
    type = models.CharField(max_length=128)
    num = models.PositiveIntegerField()
    port = models.PositiveIntegerField()

    # TODO: add celery beat tasks for monitoring node health
    status = models.CharField(max_length=64, default='up')

    def short_name(self):
        return "{}.{}".format(self.type, self.num)
    short_name.short_description = 'Name'

    def __str__(self):
        return "{0} {1}".format(self.formation.id, self.short_name())

    class Meta:
        get_latest_by = '-created'
        ordering = ['created']
        unique_together = (('app', 'type', 'num'),
                           ('formation', 'port'))


@python_2_unicode_compatible
class Config(UuidAuditedModel):
    """
    Set of configuration values applied as environment variables
    during runtime execution of the Application.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    app = models.ForeignKey('App')
    version = models.PositiveIntegerField()

    values = JSONField(default='{}', blank=True)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'version'),)

    def __str__(self):
        return "{0}-v{1}".format(self.app.id, self.version)


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
        get_latest_by = 'created'
        ordering = ['-created']
        unique_together = (('app', 'uuid'),)

    def __str__(self):
        return "{0}-{1}".format(self.app.id, self.sha[:7])


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
    build = models.ForeignKey('Build', blank=True, null=True)

    class Meta:
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


def _user_flat(self):
    return {'username': self.username}


def _user_calculate(self):
    data = {'id': self.username, 'ssh_keys': {}}
    for k in self.key_set.all():
        data['ssh_keys'][k.id] = k.public
    return data


def _user_publish(self):
    CM.publish_user(self.flat(), self.calculate())


def _user_purge(self):
    CM.purge_user(self.flat())


# attach to built-in django user
User.flat = _user_flat
User.calculate = _user_calculate
User.publish = _user_publish
User.purge = _user_purge


# define update/delete callbacks for synchronizing
# models with the configuration management backend

def _publish_to_cm(**kwargs):
    kwargs['instance'].publish()


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


# Connect Django model signals
# Sync database updates with the configuration management backend
post_save.connect(_publish_to_cm, sender=App, dispatch_uid='api.models')
post_save.connect(_publish_to_cm, sender=Formation, dispatch_uid='api.models')
# Log significant app-related events
post_save.connect(_log_build_created, sender=Build, dispatch_uid='api.models')
post_save.connect(_log_release_created, sender=Release, dispatch_uid='api.models')
post_save.connect(_log_config_updated, sender=Config, dispatch_uid='api.models')

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
