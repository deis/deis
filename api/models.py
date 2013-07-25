#!/usr/bin/python
# -*- coding: utf-8 -*-

"""
Data models for the Deis API.
"""
# pylint: disable=R0903,W0232

from __future__ import unicode_literals
import importlib
import json

from django.conf import settings
from django.db import models
from django.db.models.signals import post_save
from django.dispatch import receiver
from django.dispatch.dispatcher import Signal
from django.utils.encoding import python_2_unicode_compatible
from rest_framework.authtoken.models import Token

from api import fields
from celerytasks import chef, controller
from celery.canvas import group
import yaml
import os.path


# define custom signals
scale_signal = Signal(providing_args=['formation', 'user'])
release_signal = Signal(providing_args=['formation', 'user'])

def import_tasks(provider_type):
    "Return Celery tasks for a given provider type"
    try:
        tasks = importlib.import_module('celerytasks.'+provider_type)
    except ImportError as e:
        raise e
    return tasks


class AuditedModel(models.Model):

    """
    Adds created and update fields to a model.
    """

    created = models.DateTimeField(auto_now_add=True)
    updated = models.DateTimeField(auto_now=True)

    class Meta:
        """
        Metadata options for AuditedModel, marking this class as abstract.
        """
        abstract = True


class UuidAuditedModel(AuditedModel):

    """
    Adds a UUID primary key to an audited model.
    """
    uuid = fields.UuidField('UUID', primary_key=True)

    class Meta:
        """
        Metadata options for UuidAuditedModel, marking this class as abstract.
        """
        abstract = True


@python_2_unicode_compatible
class Key(UuidAuditedModel):

    """An SSH public key."""

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.CharField(max_length=128)
    public = models.TextField()

    class Meta:
        verbose_name = 'SSH Key'
        unique_together = (('owner', 'id'))

    def __str__(self):
        return '{0} : {1}'.format(self.owner.username, self.id)


@python_2_unicode_compatible
class ProviderManager(models.Manager):

    def seed(self, user, **kwargs):
        providers = (('ec2', 'ec2'),)
        for p_id, p_type in providers:
            self.create(owner=user, id=p_id, type=p_type, creds='{}')


@python_2_unicode_compatible
class Provider(UuidAuditedModel):

    """
    Cloud provider information for a user. Available as `user.provider_set`.
    """
    objects = ProviderManager()

    PROVIDERS = (
        ('ec2', 'Amazon Elastic Compute Cloud (EC2)'),
        ('mock', 'Mock Reference Provider'),
    )

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    type = models.SlugField(max_length=16, choices=PROVIDERS)
    creds = fields.CredentialsField()

    class Meta:
        unique_together = (('owner', 'id'),)


@python_2_unicode_compatible
class FlavorManager(models.Manager):

    def load_cloud_config_base(self):
        # load cloud-config-base yaml_
        _cloud_config_path = os.path.abspath(
                os.path.join(__file__, '..', 'files', 'cloud-config-base.yml'))
        with open(_cloud_config_path) as f:
            _data = f.read()
        return yaml.safe_load(_data)

    def seed(self, user, **kwargs):
        # TODO: add optimized AMIs to default flavors
        flavors = (
            {'id': 'ec2-us-east-1',
             'provider': 'ec2',
             'params': json.dumps({'region': 'us-east-1'})},
            {'id': 'ec2-us-west-1',
             'provider': 'ec2',
             'params': json.dumps({'region': 'us-west-1'})},
            {'id': 'ec2-us-west-2',
             'provider': 'ec2',
             'params': json.dumps({'region': 'us-west-2'})},
            {'id': 'ec2-eu-west-1',
             'provider': 'ec2',
             'params': json.dumps({'region': 'eu-west-1'})},
            {'id': 'ec2-ap-northeast-1',
             'provider': 'ec2',
             'params': json.dumps({'region': 'ap-northeast-1'})},
            {'id': 'ec2-ap-southeast-1',
             'provider': 'ec2',
             'params': json.dumps({'region': 'ap-southeast-1'})},
            {'id': 'ec2-ap-southeast-2',
             'provider': 'ec2',
             'params': json.dumps({'region': 'ap-southeast-2'})},
            {'id': 'ec2-sa-east-1',
             'provider': 'ec2',
             'params': json.dumps({'region': 'sa-east-1'})},
        )
        cloud_config = self.load_cloud_config_base()
        for flavor in flavors:
            provider = flavor.pop('provider')
            flavor['provider'] = Provider.objects.get(owner=user, id=provider)
            flavor['init'] = cloud_config
            self.create(owner=user, **flavor)


@python_2_unicode_compatible
class Flavor(UuidAuditedModel):

    """
    Virtual machine flavors available as `user.flavor_set`.
    """
    objects = FlavorManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    provider = models.ForeignKey('Provider')
    params = fields.ParamsField()
    init = fields.CloudInitField()

    class Meta:
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return self.id


class ScalingError(Exception):
    pass


@python_2_unicode_compatible
class FormationManager(models.Manager):

    def publish(self, **kwargs):
        # build data bag
        formations = self.all()
        databag = {
            'id': 'gitosis',
            'ssh_keys': {},
            'admins': [],
            'formations': {}
        }
        # add all ssh keys on the system
        for key in Key.objects.all():
            key_id = '{0}_{1}'.format(key.owner.username, key.id)
            databag['ssh_keys'][key_id] = key.public
        # TODO: add sharing-based key lookup, for now just owner's keys
        for formation in formations:
            keys = databag['formations'][formation.id] = []
            owner_keys = [ '{0}_{1}'.format(k.owner.username, k.id) for k in formation.owner.key_set.all() ]
            keys.extend(owner_keys)
        # call a celery task to update gitosis
        if settings.CHEF_ENABLED:
            controller.update_gitosis.delay(databag).wait()  # @UndefinedVariable

    def next_container_node(self, formation, container_type):
        count = []
        layer = formation.layer_set.get(id='runtime')
        runtime_nodes = list(Node.objects.filter(formation=formation, layer=layer).order_by('created'))
        container_map = { n: [] for n in runtime_nodes }
        containers = list(Container.objects.filter(formation=formation, type=container_type).order_by('created'))
        for c in containers:
            container_map[c.node].append(c)
        for n in container_map.keys():
            # (2, node3), (2, node2), (3, node1)
            count.append((len(container_map[n]), n))
        if not count:
            raise ScalingError('No nodes available for containers')
        count.sort()
        return count[0][1]


@python_2_unicode_compatible
class Formation(UuidAuditedModel):

    """
    Formation of machine instances, list of nodes available
    as `formation.nodes`
    """
    objects = FormationManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    flavor = models.ForeignKey('Flavor')
    layers = fields.JSONField(default='{}', blank=True)
    containers = fields.JSONField(default='{}', blank=True)

    environment = models.CharField(max_length=64, default='_default')
    
    ssh_username = models.CharField(max_length=64, default='ubuntu')
    ssh_private_key = models.TextField()
    ssh_public_key = models.TextField()
    
    class Meta:
        unique_together = (('owner', 'id'),)

    def scale_layers(self, **kwargs):
        layers = self.layers.copy()
        funcs = []
        for layer_id, requested in layers.items():
            layer = self.layer_set.get(id=layer_id)
            nodes = list(layer.node_set.all().order_by('created'))
            diff = requested - len(nodes)
            if diff == 0:
                return
            while diff < 0:
                node = nodes.pop(0)
                funcs.append(node.terminate)
                diff = requested - len(nodes)
            while diff > 0: 
                node = Node.objects.new(self, layer)
                nodes.append(node)
                funcs.append(node.launch)
                diff = requested - len(nodes)
        # http://docs.celeryproject.org/en/latest/userguide/canvas.html#groups
        job = [func() for func in funcs ]
        # balance containers
        containers_balanced = self._balance_containers()
        # launch/terminate nodes in parallel
        if job:
            group(*job).apply_async().join()
        # once nodes are in place, recalculate the formation and update the data bag
        databag = self.calculate()
        # force-converge nodes if there were changes
        if job or containers_balanced:
            self.converge(databag)
        return databag

    def scale_containers(self, **kwargs):
        requested_containers = self.containers.copy()
        runtime_layers = self.layer_set.filter(id='runtime')
        if len(runtime_layers) < 1:
            raise ScalingError('Must create a "runtime" layer to host containers')
        runtime_nodes = runtime_layers[0].node_set.all()
        if len(runtime_nodes) < 1:
            raise ScalingError('Must scale runtime nodes > 0 to host containers')
        # increment new container nums off the most recent container
        all_containers = self.container_set.all().order_by('-created')
        container_num = 1 if not all_containers else all_containers[0].num + 1
        # iterate and scale by container type (web, worker, etc)
        change = False
        for container_type in requested_containers.keys():
            containers = list(self.container_set.filter(type=container_type).order_by('created'))
            requested = requested_containers.pop(container_type)
            diff = requested - len(containers)
            change = 0
            if diff == 0:
                return change
            while diff < 0:
                c = containers.pop(0)
                c.delete()
                diff = requested - len(containers)
                change -= 1
            while diff > 0:
                node = Formation.objects.next_container_node(self, container_type)
                c = Container.objects.create(owner=self.owner,
                                             formation=self,
                                             type=container_type,
                                             num=container_num,
                                             node=node)
                containers.append(c)
                container_num += 1
                diff = requested - len(containers)
                change +=1
        # once nodes are in place, recalculate the formation and update the data bag
        databag = self.calculate()
        if change is True:
            self.converge(databag)
        return databag

    def balance(self, **kwargs):
        containers_balanced = self._balance_containers()
        databag = self.calculate()
        if containers_balanced:
            self.converge(databag)
        return databag
        
    def _balance_containers(self, **kwargs):
        runtime_nodes = self.node_set.filter(layer__id='runtime').order_by('created')
        if len(runtime_nodes) < 2:
            return # there's nothing to balance with 1 runtime node
        all_containers = Container.objects.filter(formation=self).order_by('-created')
        # get the next container number (e.g. web.19)
        container_num = 1 if not all_containers else all_containers[0].num + 1
        changed = False
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
                # create a container on the most under-utilized node
                Container.objects.create(owner=self.owner,
                                         formation=self,
                                         type=container_type,
                                         num=container_num,
                                         node=n_under)
                container_num +=1
                # delete the oldest container from the most over-utilized node
                c = n_over.container_set.filter(type=container_type).order_by('created')[0]
                c.delete()
                # update the n_map accordingly
                for n in (n_over, n_under):
                    ct = len(n.container_set.filter(type=container_type))
                    n_map.setdefault(ct, []).append(n)
                changed = True
        return changed
        
    def __str__(self):
        return self.id

    def prepare_provider(self, *args, **kwargs):
        tasks = import_tasks(self.flavor.provider.type)
        args = (self.id, self.flavor.provider.creds.copy(),
                self.flavor.params.copy())
        return tasks.prepare_formation.subtask(args)
    
    def calculate(self):
        "Return a Chef data bag item for this formation"
        release = self.release_set.all().order_by('-created')[0]
        d = {}
        d['id'] = self.id
        d['release'] = {}
        d['release']['version'] = release.version
        d['release']['config'] = release.config.values
        d['release']['image'] = release.image
        d['release']['build'] = {}
        if release.build:
            d['release']['build']['url'] = release.build.url
            d['release']['build']['procfile'] = release.build.procfile
        # calculate proxy
        d['proxy'] = {}
        d['proxy']['algorithm'] =  'round_robin'
        d['proxy']['port'] = 80
        d['proxy']['backends'] = []
        # calculate container formation
        d['containers'] = {}
        for c in self.container_set.all().order_by('created'):
            # all container types get an exposed port starting at 5001
            port = 5000 + c.num
            d['containers'].setdefault(c.type, {})
            d['containers'][c.type].update({ c.num: "{0}:{1}".format(c.node.id, port) })
            # only proxy to 'web' containers
            if c.type == 'web':
                d['proxy']['backends'].append("{0}:{1}".format(c.node.fqdn, port))
        # add all the participating nodes
        d['nodes'] = {}
        for n in self.node_set.all():
            d['nodes'].setdefault(n.layer.id, {})[n.id] = n.fqdn
        # call a celery task to update the formation data bag
        if settings.CHEF_ENABLED:
            controller.update_formation.delay(self.id, d).wait()  # @UndefinedVariable
        return d

    def converge(self, databag):
        # call a celery task to update the formation data bag
        if settings.CHEF_ENABLED:
            controller.update_formation.delay(self.id, databag).wait()  # @UndefinedVariable
            nodes = [ node for node in self.node_set.all() ]
            job = group(*[ n.converge() for n in nodes ])
            _results = job.apply_async().join()
        return databag

    def destroy(self):
        tasks = import_tasks(self.flavor.provider.type)
        subtasks = []
        # call a celery task to update the formation data bag
        if settings.CHEF_ENABLED:
            subtasks.extend([controller.destroy_formation.s(self.id)])  # @UndefinedVariable
        # create subtasks to terminate all nodes in parallel
        subtasks.extend([ node.terminate() for node in self.node_set.all() ])
        job = group(*subtasks)
        job.apply_async().join() # block for termination
        # purge other hosting provider infrastructure
        tasks.cleanup_formation.delay(self.id,
            self.flavor.provider.creds.copy(),
            self.flavor.params.copy()).wait()


@python_2_unicode_compatible
class Layer(UuidAuditedModel):
    
    """
    Layer of nodes used by the formation
    
    All nodes in a layer share the same flavor and configuration
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    id = models.SlugField(max_length=64)
    
    formation = models.ForeignKey('Formation')
    flavor = models.ForeignKey('Flavor')
    nodes = models.PositiveSmallIntegerField(default=0)

    run_list = models.CharField(max_length=512)
    initial_attributes = fields.JSONField(default='{}', blank=True)

    def __str__(self):
        return self.id


@python_2_unicode_compatible
class NodeManager(models.Manager):

    def new(self, formation, layer):
        existing_nodes = self.filter(formation=formation, layer=layer).order_by('-created')
        if existing_nodes:
            next_num = existing_nodes[0].num + 1
        else:
            next_num = 1
        node = self.create(owner=formation.owner,
                           formation=formation,
                           layer=layer,
                           num=next_num,
                           id="{0}-{1}-{2}".format(formation.id, layer.id, next_num))
        return node


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
    num = models.PositiveIntegerField(default=1)
    
    # synchronized with node after creation
    provider_id = models.SlugField(max_length=64, blank=True, null=True)
    fqdn = models.CharField(max_length=256, blank=True, null=True)
    status = fields.NodeStatusField(blank=True, null=True)

    class Meta:
        unique_together = (('owner', 'id'),)

    def __str__(self):
        return self.id

    def launch(self, *args, **kwargs):
        tasks = import_tasks(self.formation.flavor.provider.type)
        args = self._prepare_launch_args()
        return tasks.launch_node.subtask(args)

    def _prepare_launch_args(self):
        creds = self.formation.flavor.provider.creds.copy()
        params = self.formation.flavor.params.copy()
        params['formation'] = self.formation.id
        params['id'] = self.id
        init = self.formation.flavor.init.copy()
        if settings.CHEF_ENABLED:
            chef = init['chef'] = {}
            chef['ruby_version'] = settings.CHEF_RUBY_VERSION
            chef['server_url'] = settings.CHEF_SERVER_URL
            chef['install_type'] = settings.CHEF_INSTALL_TYPE
            chef['environment'] = settings.CHEF_ENVIRONMENT
            chef['validation_name'] = settings.CHEF_VALIDATION_NAME
            chef['validation_key'] = settings.CHEF_VALIDATION_KEY
            chef['node_name'] = self.id
            if self.layer.run_list:
                chef['run_list'] = self.layer.run_list
            if self.layer.initial_attributes:
                chef['initial_attributes'] = self.layer.initial_attributes
        # add the formation's ssh pubkey
        init.setdefault('ssh_authorized_keys', []).append(
                                self.formation.ssh_public_key)
        # add all of the owner's SSH keys
        init['ssh_authorized_keys'].extend([k.public for k in self.formation.owner.key_set.all() ])
        ssh_username = self.formation.ssh_username
        ssh_private_key = self.formation.ssh_private_key
        args = (self.uuid, creds, params, init, ssh_username, ssh_private_key)
        return args

    def converge(self, *args, **kwargs):
        tasks = import_tasks(self.formation.flavor.provider.type)
        args = self._prepare_converge_args()
        # TODO: figure out how to store task return values in model
        return tasks.converge_node.subtask(args)

    def _prepare_converge_args(self):
        ssh_username = self.formation.ssh_username
        fqdn = self.fqdn
        ssh_private_key = self.formation.ssh_private_key
        args = (self.uuid, ssh_username, fqdn, ssh_private_key)
        return args

    def terminate(self, *args, **kwargs):
        tasks = import_tasks(self.formation.flavor.provider.type)
        args = self._prepare_terminate_args()
        # TODO: figure out how to store task return values in model
        return tasks.terminate_node.subtask(args)

    def _prepare_terminate_args(self):
        creds = self.formation.flavor.provider.creds.copy()
        params = self.formation.flavor.params.copy()
        args = (self.uuid, creds, params, self.provider_id)
        return args


@python_2_unicode_compatible
class Container(UuidAuditedModel):

    """
    Docker container used to securely host an application process
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey('Formation')
    node = models.ForeignKey('Node')
    type = models.CharField(max_length=128, default='web')
    num = models.PositiveIntegerField()
    # synchronized with container after creation
    id = models.CharField(max_length=128, blank=True)
    port = models.IntegerField(blank=True, null=True)
    metadata = fields.JSONField(blank=True)

    def __str__(self):
        if self.id:
            return self.id
        return '{0} {1}.{2}'.format(self.formation.id, self.type, self.num)

    class Meta:
        unique_together = (('formation', 'type', 'num'),)


@python_2_unicode_compatible
class Config(UuidAuditedModel):

    """
    Set of configuration values applied as environment variables
    during runtime execution of the Application.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey('Formation')
    version = models.PositiveIntegerField(default=1)
    
    values = fields.EnvVarsField(default='{}', blank=True)

    class Meta:
        get_latest_by = 'created'
        ordering = ['-created']

    def __str__(self):
        return '{0}-v{1}'.format(self.formation.id, self.version)


@python_2_unicode_compatible
class Build(UuidAuditedModel):

    """
    The software build process and creation of executable binaries and assets.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey('Formation')
    version = models.PositiveIntegerField(default=1)
    
    sha = models.CharField('SHA', max_length=255, blank=True)
    output = models.TextField(blank=True)
    # metadata
    procfile = fields.ProcfileField(blank=True)
    dockerfile = models.TextField(blank=True)
    config = fields.EnvVarsField(blank=True)
    # slug info, TODO: replace default URL with something more user friendly
    url = models.URLField('URL', default='https://s3.amazonaws.com/gabrtv-slugs/nodejs.tar.gz')
    size = models.IntegerField(blank=True, null=True)
    checksum = models.CharField(max_length=255, blank=True)

    class Meta:
        ordering = ('-created',)

    def __str__(self):
        return '{0}-v{1}'.format(self.formation.id, self.version)


@python_2_unicode_compatible
class Release(UuidAuditedModel):

    """
    The deployment of a Build to Instances and the restarting of Processes.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey('Formation')
    version = models.PositiveIntegerField(default=1)
    
    config = models.ForeignKey('Config')
    image = models.CharField(max_length=256, default='ubuntu')
    # build only required for heroku-style apps
    build = models.ForeignKey('Build', blank=True, null=True)

    class Meta:
        ordering = ('-created',)

    def __str__(self):
        return '{0}-v{1}'.format(self.formation.id, self.version)
    
    def rollback(self):
        # create a rollback log entry
        # call run
        raise NotImplementedError


@receiver(release_signal)
def new_release(sender, **kwargs):
    formation, user = kwargs['formation'], kwargs['user']
    last_release = Release.objects.filter(
                formation=formation).order_by('-created')[0]
    image = kwargs.get('image', last_release.image)
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
                owner=user, formation=formation, values=new_values)
    # create new release and auto-increment version
    new_version = last_release.version + 1
    release = Release.objects.create(owner=user, formation=formation, 
        image=image, config=config, build=build, version=new_version)
    return release


@python_2_unicode_compatible
class Access(UuidAuditedModel):

    """
    An access control list (ACL) entry specifying what role a user has for
    an app.

    A user is considered always to have "admin" access to his or her own
    apps whether or not a specific Access entry exists.
    """

    VIEWER = 'viewer'
    USER = 'user'
    ADMIN = 'admin'
    ACCESS_ROLES = (
        (VIEWER, 'Viewer'),
        (USER, 'User'),
        (ADMIN, 'Administrator'),
    )

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey('Formation')
    role = models.CharField(max_length=6, choices=ACCESS_ROLES, default=USER)

    class Meta:
        """Metadata options for an Access model."""
        verbose_name_plural = 'accesses'

    def __str__(self):
        return '{0}: {1} is {2}'.format(self.app, self.user, self.role)


@python_2_unicode_compatible
class Event(UuidAuditedModel):

    """
    A change in an Application's state worth persisting so it can be
    searched for in the future.
    """

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey('Formation')

    def __str__(self):
        # TODO: what's a useful string representation of this object?
        return '{0} event ?'.format(self.app)


@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_auth_token(sender, instance=None, created=False, **kwargs):
    """Adds an auth Token to each newly created user."""
    if created:
        # pylint: disable=E1101
        Token.objects.create(user=instance)

