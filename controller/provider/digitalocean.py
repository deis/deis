"""
Deis cloud provider implementation for Digital Ocean.
"""

from __future__ import unicode_literals

import json
import time

from dop import client
from api.models import Layer


REGION_LIST = {
    'new-york-1': 'New York 1',
    'amsterdam-1': 'Amsterdam 1',
    'san-francisco-1': 'San Francisco 1',
    'new-york-2': 'New York 2',
    'amsterdam-2': 'Amsterdam 2',
    'singapore-1': 'Singapore 1',
}


def seed_flavors():
    """Seed the database with default Flavors for digital ocean"""
    flavors = []
    for r in REGION_LIST.keys():
        flavors.append({
            'id': 'digitalocean-{}'.format(r),
            'provider': 'digitalocean',
            'params': json.dumps({
                'region': REGION_LIST[r],
                'size': '4GB',
                'image': 'deis-node-image'
            })
        })
    return flavors


def build_layer(layer):
    conn = _create_digitalocean_connection(layer['creds'])
    # create a new SSH key
    name = "deis-{formation}-{id}".format(**layer)
    # import a new keypair to the Digital Ocean Control Panel
    params = {'name': name, 'ssh_pub_key': layer['ssh_public_key']}
    conn.request('/ssh_keys/new/', method='GET', params=params)
    # Digital Ocean images only have the root user created by default
    l = Layer.objects.get(id=layer['id'], formation__id=layer['formation'])
    l.ssh_username = 'root'
    l.save()


def destroy_layer(layer):
    conn = _create_digitalocean_connection(layer['creds'])
    name = "deis-{formation}-{id}".format(**layer)
    # retrieve and delete the SSH key created
    for k in conn.ssh_keys():
        if k.name == name:
            conn.destroy_ssh_key(k.id)


def build_node(node):
    conn = _create_digitalocean_connection(node['creds'])
    # get digital ocean config for a new node
    kwargs = _get_droplet_kwargs(node, conn)
    # spawn the droplet
    droplet = conn.create_droplet(**kwargs)
    # wait until the droplet has become active
    while conn.show_droplet(droplet.id).status != 'active':
        time.sleep(1)
    provider_id = 'digitalocean'
    # the data in the old droplet is not there, so
    # we just create a new Droplet object
    fqdn = str(conn.show_droplet(droplet.id).ip_address)
    metadata = _get_droplet_metadata(conn.show_droplet(droplet.id))
    return provider_id, fqdn, metadata


def destroy_node(node):
    params, creds = node['params'], node['creds']
    conn = _create_digitalocean_connection(creds)
    name = node['id']
    for i in conn.show_active_droplets():
        if i.name == name:
            conn.destroy_droplet(i.id)
    return


def _get_droplet_kwargs(node, conn):
    params = node['params']
    return {
        'name': node['id'],
        'size_id': _get_id(conn.sizes(), params.get('size', '4GB')),
        'image_id': _get_id(conn.images(my_images=True), params.get('image', 'deis-node-image')),
        'region_id': _get_id(conn.regions(), params.get('region', 'San Francisco 1')),
        'ssh_key_ids': [str(_get_id(conn.ssh_keys(),
                        "deis-{formation}-{layer}".format(**node)))],
        'virtio': True,
        'private_networking': True
    }


def _get_droplet_metadata(droplet):
    return {
        'id': droplet.id,
        'name': droplet.name,
        'size_id': droplet.size_id,
        'image_id': droplet.image_id,
        'region_id': droplet.region_id,
        'event_id': droplet.event_id,
        'backups_active': droplet.backups_active,
        'status': droplet.status,
        'ip_address': droplet.ip_address
    }


def _create_digitalocean_connection(creds):
    if not creds:
        raise EnvironmentError('No credentials provided')
    return client.Client(creds['client_id'], creds['api_key'])


def _get_id(lst, name):
    for i in lst:
        if i.name == name:
            return i.id
    return None


def _get_name(lst, name):
    for i in lst:
        if i.name == name:
            return i.name
    return None
