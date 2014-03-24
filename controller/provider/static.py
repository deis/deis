"""
Deis cloud provider implementation for static nodes created outside Deis.
"""


def seed_flavors():
    """Seed the database with default static flavor.

    :rtype: list of dicts containing flavor data
    """
    return [{
        'id': 'static',
        'provider': 'static',
        'params': {},
    }]


def build_layer(layer):
    """
    Build a layer.

    :param layer: a dict containing formation, id, params, and creds info
    """
    pass


def destroy_layer(layer):
    """
    Destroy a layer.

    :param layer: a dict containing formation, id, params, and creds info
    """
    pass


def build_node(node):
    """
    Build a node.

    :param node: a dict containing formation, layer, params, and creds info.
    :rtype: a tuple of (provider_id, fully_qualified_domain_name, metadata)
    """
    return ('static', node['fqdn'], {})


def destroy_node(node):
    """
    Destroy a node.

    :param node: a dict containing a node's provider_id, params, and creds
    """
    pass
