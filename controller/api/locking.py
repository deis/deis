import logging

from models import get_etcd_client

log = logging.getLogger(__name__)
LOCK_PATH_TEMPLATE = '/deis/controller/locks/{name}'


class AlreadyLocked(Exception):
    pass


def lock(name, ttl=3600):
    key = LOCK_PATH_TEMPLATE.format(name=name)
    value = 'locked'
    client = get_etcd_client()
    try:
        log.debug('Reading key: "{key}"'.format(key=key))
        client.read(key)
        raise AlreadyLocked
    except KeyError:
        log.debug('KeyError, try to create lock on "{key}"'.format(key=key))
        client.write(key, value, ttl=ttl)


def unlock(name):
    key = LOCK_PATH_TEMPLATE.format(name=name)
    try:
        log.debug('Removing key: "{key}"'.format(key=key))
        get_etcd_client().delete(key)
    except KeyError:
        log.warn('{key} already removed'.format(key=key))


