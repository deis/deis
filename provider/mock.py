"""
Deis mock cloud provider implementation for testing.
"""

from __future__ import unicode_literals


def seed_flavors():
    """
    Seed the database with default flavors for the mock provider.

    :rtype: list of dicts containing flavor data
    """
    flavors = []
    for r in ('east', 'west'):
        flavors.append({'id': 'mock-{}'.format(r),
                        'provider': 'mock',
                        'params': '{}'})
    return flavors


def build_layer(layer):
    """
    Build a layer.

    This function is a no-op for the mock provider.

    :param layer: a dict containing formation, id, params, and creds info
    """
    return


def destroy_layer(layer):
    """
    Destroy a layer.

    This function is a no-op for the mock provider.

    :param layer: a dict containing formation, id, params, and creds info
    """
    return


def build_node(node):
    """
    Build a node.

    :param node: a dict containing formation, layer, params, and creds info.
    :rtype: a tuple of (provider_id, fully_qualified_domain_name, metadata)
    """
    provider_id = 'i-1234567'
    fqdn = node.get('fqdn') or 'localhost.localdomain.local'
    metadata = {'state': 'running'}
    return provider_id, fqdn, metadata


def destroy_node(node):
    """
    Destroy a node.

    This function is a no-op for the mock provider.

    :param node: a dict containing a node's provider_id, params, and creds
    """
    return
