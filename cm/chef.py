"""
Deis configuration management implementation for Opscode Chef.
"""

from __future__ import unicode_literals

import os
import re
import subprocess
import tempfile
import time
import socket

from celery.canvas import group

from api.ssh import exec_ssh, connect_ssh
from cm.chef_api import ChefAPI


CHEF_CONFIG_PATH = '/etc/chef'
CHEF_INSTALL_TYPE = 'gems'
CHEF_RUBY_VERSION = '1.9.1'
CHEF_ENVIRONMENT = '_default'
CHEF_CLIENT_VERSION = '11.6.2'

# load chef config using CHEF_CONFIG_PATH
try:
    # parse controller's chef config for server_url and client_name
    _client_cfg_path = os.path.join(CHEF_CONFIG_PATH, 'client.rb')
    if not os.path.exists(_client_cfg_path):
        raise EnvironmentError('Could not find {}'.format(_client_cfg_path))
    with open(_client_cfg_path) as f:
        _data = f.read()
    # construct a dict from the ruby client.rb
    _d = {}
    for m in re.findall(r'''^([a-zA-Z0-9_]+)[ \t]+(.*)$''',
                        _data, re.MULTILINE):
        _d[m[0]] = m[1].strip("'").strip('"')
    # set global variables from client.rb
    CHEF_SERVER_URL = _d['chef_server_url']
    CHEF_NODE_NAME = _d.get('node_name', socket.gethostname())
    CHEF_CLIENT_NAME = _d.get('node_name', socket.gethostname())
    CHEF_VALIDATION_NAME = _d['validation_client_name']
    # read the client key
    _client_pem_path = os.path.join(CHEF_CONFIG_PATH, 'client.pem')
    CHEF_CLIENT_KEY = subprocess.check_output(
        ['/bin/cat', _client_pem_path]).strip('\n')
    # read the validation key
    _valid_pem_path = os.path.join(CHEF_CONFIG_PATH, 'validation.pem')
    CHEF_VALIDATION_KEY = subprocess.check_output(
        ['/bin/cat', _valid_pem_path]).strip('\n')
except Exception as err:
    msg = "Failed to auto-configure Chef -- {}".format(err)
    if os.environ.get('READTHEDOCS'):
        # Just print the error if Sphinx is running
        print(msg)
    else:
        raise EnvironmentError(msg)


def _get_client():
    """
    Return a new instance of a Chef API Client

    :rtype: a :class:`~cm.chef_api.ChefAPI` object
    """
    return ChefAPI(CHEF_SERVER_URL, CHEF_CLIENT_NAME, CHEF_CLIENT_KEY)


def bootstrap_node(node):
    """
    Bootstrap the Chef configuration management tools onto a node.

    :param node: a dict containing the node's fully-qualified domain name and SSH info
    :raises: RuntimeError
    """
    # block until we can connect over ssh
    ssh = connect_ssh(node['ssh_username'], node['fqdn'], node.get('ssh_port', 22),
                      node['ssh_private_key'], timeout=120)
    # block until ubuntu cloud-init is finished
    initializing = True
    while initializing:
        time.sleep(10)
        initializing, _rc = exec_ssh(ssh, 'ps auxw | egrep "cloud-init" | grep -v egrep')
    # write out private key and prepare to `knife bootstrap`
    try:
        _, pk_path = tempfile.mkstemp()
        _, output_path = tempfile.mkstemp()
        with open(pk_path, 'w') as f:
            f.write(node['ssh_private_key'])
        # build knife bootstrap command
        args = ['knife', 'bootstrap', node['fqdn']]
        args.extend(['--config', '/etc/chef/client.rb'])
        args.extend(['--identity-file', pk_path])
        args.extend(['--node-name', node['id']])
        args.extend(['--sudo', '--ssh-user', node['ssh_username']])
        args.extend(['--ssh-port', str(node.get('ssh_port', 22))])
        args.extend(['--bootstrap-version', CHEF_CLIENT_VERSION])
        args.extend(['--no-host-key-verify'])
        args.extend(['--run-list', _construct_run_list(node)])
        print(' '.join(args))
        # tee the command's output to a tempfile
        args.extend(['|', 'tee', output_path])
        # TODO: figure out why home isn't being set correctly for knife exec
        env = os.environ.copy()
        env['HOME'] = '/opt/deis'
        # execute knife bootstrap
        p = subprocess.Popen(' '.join(args), env=env, shell=True)
        rc = p.wait()
        # always print knife output
        with open(output_path) as f:
            output = f.read()
        print(output)
        # raise an exception if bootstrap failed
        if rc != 0 or 'incorrect password' in output:
            raise RuntimeError('Node Bootstrap Error:\n' + output)
    # remove temp files from filesystem
    finally:
        os.remove(pk_path)
        os.remove(output_path)


def _construct_run_list(node):
    config = node['config']
    # if run_list override specified, use it (assumes csv)
    run_list = config.get('run_list', [])
    # otherwise construct a run_list using proxy/runtime flags
    if not run_list:
        run_list = ['recipe[deis]']
        if node.get('runtime') is True:
            run_list.append('recipe[deis::runtime]')
        if node.get('proxy') is True:
            run_list.append('recipe[deis::proxy]')
    return ','.join(run_list)


def purge_node(node):
    """
    Purge a node and its client from Chef configuration management.

    :param node: a dict containing the id of a node to purge
    """
    client = _get_client()
    node_id = node['id']
    body, status = client.delete_node(node_id)
    if status not in [200, 404]:
        raise RuntimeError("Could not purge node {node_id}: {body}".format(**locals()))
    body, status = client.delete_client(node_id)
    if status not in [200, 404]:
        raise RuntimeError("Could not purge node client {node_id}: {body}".format(**locals()))


def converge_controller():
    """
    Converge this controller node.

    "Converge" means to change a node's configuration to match that defined by
    configuration management.

    :returns: the output of the convergence command, in this case `sudo chef-client`
    """
    #try:
    #    # we only need to run the gitosis recipe to update `git push` ACLs
    #    return subprocess.check_output(['sudo', 'chef-client', '-o', 'recipe[deis::gitosis]'])
    #except subprocess.CalledProcessError as err:
    #    print(err)
    #    print(err.output)
    #    raise err
    pass  # TODO: replace this with new key lookup


def converge_node(node):
    """
    Converge a node.

    "Converge" means to change a node's configuration to match that defined by
    configuration management.

    :param node: a dict containing the node's fully-qualified domain name and SSH info
    :returns: a tuple of the convergence command's (output, return_code)
    """
    ssh = connect_ssh(node['ssh_username'],
                      node['fqdn'], 22,
                      node['ssh_private_key'])
    output, rc = exec_ssh(ssh, 'sudo chef-client')
    print(output)
    if rc != 0:
        e = RuntimeError('Node converge error')
        e.output = output
        raise e
    return output, rc


def run_node(node, command):
    """
    Run a command on a node.

    :param node: a dict containing the node's fully-qualified domain name and SSH info
    :param command: the command-line to execute on the node
    :returns: a tuple of the command's (output, return_code)
    """
    ssh = connect_ssh(node['ssh_username'], node['fqdn'],
                      node['ssh_port'], node['ssh_private_key'])
    output, rc = exec_ssh(ssh, command, pty=True)
    return output, rc


def converge_formation(formation):
    """
    Converge all nodes in a formation.

    "Converge" means to change a node's configuration to match that defined by
    configuration management.

    :param formation: a :class:`~api.models.Formation` to converge
    :returns: the combined output of the nodes' convergence commands
    """
    nodes = formation.node_set.all()
    subtasks = []
    for n in nodes:
        subtask = converge_node.s(n.id,
                                  n.layer.flavor.ssh_username,
                                  n.fqdn,
                                  n.layer.flavor.ssh_private_key)
        subtasks.append(subtask)
    job = group(*subtasks)
    return job.apply_async().join()


def publish_user(user, data):
    """
    Publish a user to configuration management.

    :param user: a dict containing the username
    :param data: data to store with the user
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    _publish('deis-users', user['username'], data)


def purge_user(user):
    """
    Purge a user from configuration management.

    :param app: a dict containing the username of the user
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    _purge('deis-users', user['username'])


def publish_app(app, data):
    """
    Publish an app to configuration management.

    :param app: a dict containing the id of the app
    :param data: data to store with the app
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    _publish('deis-apps', app['id'], data)


def purge_app(app):
    """
    Purge an app from configuration management.

    :param app: a dict containing the id of the app
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    _purge('deis-apps', app['id'])


def publish_formation(formation, data):
    """
    Publish a formation to configuration management.

    :param formation: a dict containing the id of the formation
    :param data: data to store with the formation
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    _publish('deis-formations', formation['id'], data)


def purge_formation(formation):
    """
    Purge a formation from configuration management.

    :param formation: a dict containing the id of the formation
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    _purge('deis-formations', formation['id'])


def _publish(data_bag, item_name, item_value):
    """
    Publish a data bag item to the Chef server.

    :param data_bag: the name of a Chef data bag
    :param item_name: the name of the item to publish
    :param item_value: the value of the item to publish
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    client = _get_client()
    body, status = client.update_databag_item(data_bag, item_name, item_value)
    if status != 200:
        body, status = client.create_databag_item(data_bag, item_name, item_value)
        if status != 201:
            raise RuntimeError('Could not publish {item_name}: {body}'.format(**locals()))
    return body, status


def _purge(databag_name, item_name):
    """
    Purge a data bag item from the Chef server.

    :param databag_name: the name of a Chef data bag
    :param item_name: the name of the item to purge
    :returns: a tuple of (body, status) from the underlying HTTP response
    :raises: RuntimeError
    """
    client = _get_client()
    body, status = client.delete_databag_item(databag_name, item_name)
    if status in [200, 404]:
        return body, status
    raise RuntimeError('Could not purge {item_name}: {body}'.format(**locals()))
