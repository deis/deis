"""
Deis cloud provider implementation for local vagrant setups.
"""

from __future__ import unicode_literals

from api.ssh import exec_ssh, connect_ssh

import json
import string
import subprocess
import uuid

from api.models import Layer

# Collect details for connecting to the host machine
HOST_NODES_DIR = open('/home/vagrant/.host_nodes_dir').read().strip()
PKEY = open('/home/vagrant/.ssh/id_rsa').read()


def seed_flavors():
    """Seed the database with default flavors for each Rackspace region.

    :rtype: list of dicts containing flavor data
    """
    flavors = []
    for m in ['512', '1024', '2048']:
        flavors.append({
            'id': 'vagrant-{}'.format(m),
            'provider': 'vagrant',
            'params': json.dumps({
                'memory': m
            })
        })
    return flavors


def build_layer(layer):
    """
    Build a layer.

    :param layer: a dict containing formation, id, params, and creds info
    """

    # This can also be done with `deis layers:update` now.
    layer_ = Layer.objects.get(id=layer['id'])
    layer_.ssh_username = 'vagrant'
    layer_.save()


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

    # Can't use the vagrant UUID because it's not booted yet
    uid = str(uuid.uuid1())

    # Create a new Vagrantfile from a template
    node['params'].setdefault('memory', '512')
    template = open('/opt/deis/controller/contrib/vagrant/nodes_vagrantfile_template.rb')
    raw = string.Template(template.read())
    result = raw.substitute({ 'id':uid, 'memory':node['params']['memory'] })

    # Make a folder for the VM with its own Vagrantfile. Vagrant will then create a .vagrant folder
    # there too when it first gets booted.
    node_dir = HOST_NODES_DIR + '/' + uid
    mkdir = 'mkdir ' + node_dir
    cp_tpl = 'echo "' + result.replace('"', '\\"') + '" > ' + node_dir + '/Vagrantfile'
    _host_ssh(commands = [mkdir, cp_tpl], creds = node['creds'])

    # Boot the VM
    _run_vagrant_command(uid, args = ['up'], creds = node['creds'])

    # Copy the layer's public SSH key to the VM so that the Controller can access it.
    _run_vagrant_command(
        uid,
        args = [
            'ssh',
            '-c',
            '"echo \\"' + node['ssh_public_key'] + '\\" >> /home/vagrant/.ssh/authorized_keys"'
        ],
        creds = node['creds'],
    )

    provider_id = uid
    fqdn = provider_id + '.local' # hostname is broadcast via avahi-daemon
    metadata = {
        'id': uid,
        'fqdn': fqdn,
        'flavor': node['params']['memory']
    }
    return provider_id, fqdn, metadata


def destroy_node(node):
    """
    Destroy a node.

    :param node: a dict containing a node's provider_id, params, and creds
    """

    # This is useful if node creation failed. So that there's a record in the DB, but it has no
    # ID associated with it.
    if node['provider_id'] == None:
        return

    # Shut the VM down and destroy it
    _run_vagrant_command(node['provider_id'], args = ['destroy', '--force'], creds = node['creds'])
    node_dir = HOST_NODES_DIR + '/' + node['provider_id']

    # Sanity check before `rm -rf`
    if 'contrib/vagrant' not in node_dir:
        raise RuntimeError("Aborted node destruction: attempting to 'rm -rf' unexpected directory")

    # Completely remove the folder that contained the VM
    rm_vagrantfile = 'rm ' + node_dir + '/Vagrantfile'
    rm_node_dir = 'rm -rf ' + node_dir
    _host_ssh(commands = [rm_vagrantfile, rm_node_dir], creds = node['creds'])


def _run_vagrant_command(node_id, args = [], creds = {}):
    """
    args: A tuple of arguments to a vagrant command line.
    e.g. ['up', 'my_vm_name', '--no-provision']
    """

    cd = 'cd ' + HOST_NODES_DIR + '/' + node_id
    command = ['vagrant'] + [arg for arg in args if arg is not None]
    return _host_ssh(commands = [cd, ' '.join(command)], creds = creds)


def _host_ssh(creds = {}, commands = []):
    """
    Connect to the host machine. Namely the user's local machine.
    """
    if creds == {}:
        raise RuntimeError("No credentials provided to _host_ssh()")
    command = ' && '.join(commands)

    # First check if we can access the host machine. It's likely that their
    # IP address changes every time they request a DHCP lease.
    # TODO: Find a way of passing this error onto the CLI client.
    try:
        subprocess.check_call([
            'nc', '-z', '-w2', creds['ip'], '22'
        ], stderr=subprocess.PIPE)
    except subprocess.CalledProcessError:
        raise RuntimeError("Couldn't ping port 22 at host with IP " + creds['ip'])

    ssh = connect_ssh(creds['user'], creds['ip'], 22, PKEY, timeout=120)
    result, status = exec_ssh(ssh, command)
    if status > 0:
        raise RuntimeError(
            'SSH to Vagrant host error: ' + result +
            'Command: ' + command)
    return result
