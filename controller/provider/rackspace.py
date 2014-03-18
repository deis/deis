"""
Deis cloud provider implementation for Rackspace Open Cloud.
"""

from __future__ import unicode_literals
import json
import logging

import novaclient
import pyrax

# TODO: this seems like a bad import, but can't update the layer
# correctly otherwise.
from api.models import Layer


logger = logging.getLogger(__name__)

## TODO: add syd and hkg when they gain performance flavors in early 2014
RACKSPACE_REGIONS = {
    'dfw': 'Dallas',
    'iad': 'Northern Virginia',
    'ord': 'Chicago',
    'lon': 'London',
    #'syd': 'Sydney',
    #'hkg': 'Hong Kong',
}

RACKSPACE_DEFAULT_REGION = 'dfw'


def seed_flavors():
    """Seed the database with default flavors for each Rackspace region.

    :rtype: list of dicts containing flavor data
    """
    flavors = []
    for r in RACKSPACE_REGIONS.keys():
        flavors.append({
            'id': 'rackspace-{}'.format(r),
            'provider': 'rackspace',
            'params': json.dumps({
                'region': r,
                # 'image': ?
                'flavor': 'performance1-2',   # 2GB Performance Instance, 2 VCPUs, 2048 / 80GB
            })
        })
    return flavors


def build_layer(layer):
    """
    Build a layer.

    :param layer: a dict containing formation, id, params, and creds info
    """
    region = layer['params'].setdefault('region', RACKSPACE_DEFAULT_REGION)
    creds = layer['creds']
    conn = _create_rackspace_connection(creds, region)
    name = "deis-{formation}-{id}".format(**layer)
    # import a new keypair using the layer key material
    conn.keypairs.create(name, layer['ssh_public_key'])
    # Rackspace images only have the root user created by default
    layer_ = Layer.objects.get(id=layer['id'], formation__id=layer['formation'])
    layer_.ssh_username = 'root'
    layer_.save()


def destroy_layer(layer):
    """
    Destroy a layer.

    :param layer: a dict containing formation, id, params, and creds info
    """
    region = layer['params'].setdefault('region', RACKSPACE_DEFAULT_REGION)
    creds = layer['creds']
    conn = _create_rackspace_connection(creds, region)
    name = "deis-{formation}-{id}".format(**layer)
    # delete the keypair we created in build_layer
    try:
        conn.keypairs.delete(name)
    except (novaclient.exceptions.NotFound, pyrax.exceptions.NotFound) as err:
        logger.warning("rackspace.destroy_layer: {}".format(err))


def build_node(node):
    """
    Build a node.

    :param node: a dict containing formation, layer, params, and creds info.
    :rtype: a tuple of (provider_id, fully_qualified_domain_name, metadata)
    """
    params, creds = node['params'], node['creds']
    region = params.setdefault('region', RACKSPACE_DEFAULT_REGION)
    conn = _create_rackspace_connection(creds, region)
    name = 'deis-{formation}-{layer}-{id}'.format(**node)
    params['key_name'] = 'deis-{formation}-{layer}'.format(**node)
    tags = {'Name': name}
    # look up the saved image / snapshot by name 'deis-node-image', until we
    # can create a public image at Rackspace--at which point we call list_images().
    try:
        image = next(i for i in conn.list_snapshots() if i.name == 'deis-node-image')
    except StopIteration:
        msg = """\
Can't find saved image "deis-node-image" in region {}. Please follow the
instructions in prepare-rackspace-image.sh before scaling Deis nodes.
""".format(region)
        raise EnvironmentError(msg)
    srv = conn.servers.create(
        name, image.id, params['flavor'], meta=tags, key_name=params['key_name'])
    server = pyrax.utils.wait_for_build(srv)
    provider_id = server.id
    # TODO: is this the right way to fetch the fqdn?
    public_addrs = server.addresses.get('public', [])
    ipv4_addrs = [entry['addr'] for entry in public_addrs if entry['version'] == 4]
    fqdn = ipv4_addrs[0] if ipv4_addrs else None
    metadata = _format_metadata(server)
    return provider_id, fqdn, metadata


def destroy_node(node):
    """
    Destroy a node.

    :param node: a dict containing a node's id, params, and creds
    """
    params, creds = node['params'], node['creds']
    region = params.setdefault('region', 'dfw')
    conn = _create_rackspace_connection(creds, region)
    name = 'deis-{formation}-{layer}-{id}'.format(**node)
    try:
        server = conn.servers.get(name)
        server.delete()
        pyrax.utils.wait_until(server, 'status', ['DELETED', 'ERROR'])
    except (novaclient.exceptions.NotFound, pyrax.exceptions.NotFound) as err:
        logger.warning("rackspace.destroy_node: {}".format(err))


def _create_rackspace_connection(creds, region):
    """
    Connect to a Rackspace region with the given credentials.

    :param creds: a dict containing a Rackspace username and api_key
    :region: the name of a Rackspace region, such as "ord"
    :rtype: a connected :class:`~novaclient.v1_1.client.Client`
    :raises EnvironmentError: if no credentials are provided
    """
    if not creds:
        raise EnvironmentError('No credentials provided')
    pyrax.set_setting('identity_type', creds['identity_type'])
    pyrax.set_credentials(creds['username'], creds['api_key'])
    return pyrax.connect_to_cloudservers(region.upper())


def _format_metadata(server):
    return {
        'addresses': server.addresses,
        'created': server.created,
        'flavor': server.flavor,
        'hostId': server.hostId,
        'human_id': server.human_id,
        'id': server.id,
        'image': server.image,
        'key_name': server.key_name,
        'is_loaded': server.is_loaded(),
        'metadata': server.metadata,
        'name': server.name,
        'networks': server.networks,
        'progress': server.progress,
        'status': server.status,
        'tenant_id': server.tenant_id,
        'updated': server.updated,
        'user_id': server.user_id,
    }
