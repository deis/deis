
from __future__ import unicode_literals


def seed_flavors(username):
    """Seed the database with default Flavors for the mock provider"""
    flavors = []
    for r in ('east', 'west'):
        flavors.append({'id': 'mock-{}'.format(r),
                        'provider': 'mock',
                        'params': '{}'})
    return flavors


def build_layer(layer):
    return layer


def destroy_layer(layer):
    return layer


def build_node(node, config):
    node.provider_id = 'i-1234567'
    node.metadata = {'state': 'running'}
    node.fqdn = 'localhost.localdomain.local'
    node.save()
    return node


def destroy_node(node):
    return node
