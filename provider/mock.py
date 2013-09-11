
from __future__ import unicode_literals


def seed_flavors():
    """Seed the database with default Flavors for the mock provider"""
    flavors = []
    for r in ('east', 'west'):
        flavors.append({'id': 'mock-{}'.format(r),
                        'provider': 'mock',
                        'params': '{}'})
    return flavors


def build_layer(layer):
    return


def destroy_layer(layer):
    return


def build_node(node):
    provider_id = 'i-1234567'
    fqdn = 'localhost.localdomain.local'
    metadata = {'state': 'running'}
    return provider_id, fqdn, metadata


def destroy_node(node):
    return
