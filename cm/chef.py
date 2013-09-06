
from __future__ import unicode_literals

import json
import os.path
import re
import subprocess
import time

from celery.canvas import group

from api.ssh import exec_ssh, connect_ssh
from cm.chef_api import ChefAPI


CHEF_CONFIG_PATH = '/etc/chef'
CHEF_INSTALL_TYPE = 'gems'
CHEF_RUBY_VERSION = '1.9.1'
CHEF_ENVIRONMENT = '_default'
CHEF_CLIENT_VERSION = '11.4.4'

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
    CHEF_NODE_NAME = _d['node_name']
    CHEF_CLIENT_NAME = _d['node_name']
    CHEF_VALIDATION_NAME = _d['validation_client_name']
    # read the client key
    _client_pem_path = os.path.join(CHEF_CONFIG_PATH, 'client.pem')
    CHEF_CLIENT_KEY = subprocess.check_output(
        ['sudo', '/bin/cat', _client_pem_path]).strip('\n')
    # read the validation key
    _valid_pem_path = os.path.join(CHEF_CONFIG_PATH, 'validation.pem')
    CHEF_VALIDATION_KEY = subprocess.check_output(
        ['sudo', '/bin/cat', _valid_pem_path]).strip('\n')
except Exception as e:
    print 'Error: failed to auto-configure Chef -- {}'.format(e)


def configure_node(node):
    config = node.layer.config.copy()
    # http://cloudinit.readthedocs.org/en/latest/topics/examples.html#install-and-run-chef-recipes
    chef = config['chef'] = {}
    chef['node_name'] = node.id
    # if run_list specified in layer, use it (assumes csv)
    run_list = node.layer.config.get('run_list', [])
    if run_list:
        chef['run_list'] = run_list.split(',')
    # otherwise construct a run_list using proxy/runtime flags
    else:
        run_list = ['recipe[deis]']
        if node.layer.runtime is True:
            run_list.append('recipe[deis::runtime]')
        if node.layer.proxy is True:
            run_list.append('recipe[deis::proxy]')
        chef['run_list'] = run_list
    attrs = node.layer.config.get('initial_attributes')
    if attrs:
        chef['initial_attributes'] = attrs
    # add global chef config
    chef['version'] = CHEF_CLIENT_VERSION
    chef['ruby_version'] = CHEF_RUBY_VERSION
    chef['server_url'] = CHEF_SERVER_URL
    chef['install_type'] = CHEF_INSTALL_TYPE
    chef['environment'] = CHEF_ENVIRONMENT
    chef['validation_name'] = CHEF_VALIDATION_NAME
    chef['validation_key'] = CHEF_VALIDATION_KEY
    return config


def bootstrap_node(node):
    # loop until node is registered with chef
    # if chef bootstrapping fails, the node will not complete registration
    registered = False
    while not registered:
        # reinstatiate the client on each poll attempt
        # to avoid disconnect errors
        client = ChefAPI(CHEF_SERVER_URL,
                         CHEF_CLIENT_NAME,
                         CHEF_CLIENT_KEY)
        resp, status = client.get_node(node.id)
        if status == 200:
            body = json.loads(resp)
            # wait until idletime is not null
            # meaning the node is registered
            if body.get('automatic', {}).get('idletime'):
                break
        time.sleep(5)
    return node


def destroy_node(node):
    """
    Purge the Node & Client records from Chef Server
    """
    client = ChefAPI(CHEF_SERVER_URL,
                     CHEF_CLIENT_NAME,
                     CHEF_CLIENT_KEY)
    client.delete_node(node.id)
    client.delete_client(node.id)
    return node


def update_user(user):
    client = ChefAPI(CHEF_SERVER_URL,
                     CHEF_CLIENT_NAME,
                     CHEF_CLIENT_KEY)
    # client.create_databag_item('deis-users', user.username, user.calculate())
    client.update_databag_item('deis-users', user.username, user.calculate())


def update_app(app):
    client = ChefAPI(CHEF_SERVER_URL,
                     CHEF_CLIENT_NAME,
                     CHEF_CLIENT_KEY)
    client.update_databag_item('deis-apps', app.id, app.calculate())


def update_formation(formation, client):
    client.update_databag_item('deis-formations', formation.id, formation.calculate())


def converge_controller():
    # NOTE: converging the controller can overwrite any in-place
    # changes to application code
    return subprocess.check_output(
        ['sudo', 'chef-client', '--override-runlist', 'recipe[deis::gitosis]'])


def converge_node(node):
    ssh = connect_ssh(node.layer.ssh_username,
                      node.fqdn, 22,
                      node.layer.ssh_private_key)
    output, rc = exec_ssh(ssh, 'sudo chef-client')
    return output, rc


def converge_formation(formation):
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
