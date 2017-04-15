"""
The database models for the **api** Django app.
"""

import importlib

from django.conf import settings
from django.contrib.auth.models import User
from django.db.models.signals import post_delete
from django.db.models.signals import post_save

# import user-defined configuration management module
CM = importlib.import_module(settings.CM_MODULE)


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


def _publish_user_to_cm(**kwargs):
    if kwargs.get('update_fields') == frozenset(['last_login']):
        return
    kwargs['instance'].publish()


def _purge_user_from_cm(**kwargs):
    kwargs['instance'].purge()


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


from formation import Formation  # noqa
from flavor import Flavor  # noqa
from key import Key  # noqa
from providers import Provider  # noqa
from push import Push  # noqa

from layer import Layer  # noqa

from node import Node  # noqa

from container import Container  # noqa

from application import App, log_event  # noqa

from build import Build  # noqa
from config import Config  # noqa
from release import Release, release_signal  # noqa


# Connect Django model signals
# Sync database updates with the configuration management backend
post_save.connect(_publish_to_cm, sender=App, dispatch_uid='api.models')
post_save.connect(_publish_to_cm, sender=Formation, dispatch_uid='api.models')
post_save.connect(_publish_user_to_cm, sender=User, dispatch_uid='api.models')
post_delete.connect(_purge_user_from_cm, sender=User, dispatch_uid='api.models')
# Log significant app-related events
post_save.connect(_log_build_created, sender=Build, dispatch_uid='api.models')
post_save.connect(_log_release_created, sender=Release, dispatch_uid='api.models')
post_save.connect(_log_config_updated, sender=Config, dispatch_uid='api.models')
