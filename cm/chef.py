
from __future__ import unicode_literals

import os
import re
import subprocess
import tempfile
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
    raise EnvironmentError('Failed to auto-configure Chef -- {}'.format(e))


def _get_client():
    """
    Return a new instance of a Chef API Client
    """
    return ChefAPI(CHEF_SERVER_URL, CHEF_CLIENT_NAME, CHEF_CLIENT_KEY)


def bootstrap_node(node):
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
        with open(pk_path, 'w') as f:
            f.write(node['ssh_private_key'])
        # build knife bootstrap command
        args = ['knife', 'bootstrap', node['fqdn']]
        args.extend(['--identity-file', pk_path])
        args.extend(['--node-name', node['id']])
        args.extend(['--sudo', '--ssh-user', node['ssh_username']])
        args.extend(['--ssh-port', str(node.get('ssh_port', 22))])
        args.extend(['--bootstrap-version', CHEF_CLIENT_VERSION])
        args.extend(['--no-host-key-verify'])
        args.extend(['--run-list', _construct_run_list(node)])
        print(' '.join(args))
        # TODO: figure out why home isn't being set correctly for knife exec
        env = os.environ.copy()
        env['HOME'] = '/opt/deis'
        # execute knife bootstrap
        p = subprocess.Popen(args, env=env, stderr=subprocess.PIPE)
        rc = p.wait()
        if rc != 0:
            print(p.stderr.read())
            raise RuntimeError('Node Bootstrap Error')
    # remove private key from fileystem
    finally:
        pass  # os.remove(pk_path)


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
    Purge the Node & Client records from Chef Server
    """
    client = _get_client()
    client.delete_node(node['id'])
    client.delete_client(node['id'])


def converge_controller():
    try:
        return subprocess.check_output(['sudo', 'chef-client'])
    except subprocess.CalledProcessError as e:
        print(e)
        print(e.output)
        raise e


def converge_node(node):
    ssh = connect_ssh(node['ssh_username'],
                      node['fqdn'], 22,
                      node['ssh_private_key'])
    output, rc = exec_ssh(ssh, 'sudo chef-client')
    if rc != 0:
        e = RuntimeError('Node converge error')
        e.output = output
        raise e
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


def publish_user(user, data):
    _publish('deis-users', user['username'], data)


def publish_app(app, data):
    _publish('deis-apps', app['id'], data)


def purge_app(app):
    _purge('deis-apps', app['id'])


def publish_formation(formation, data):
    _publish('deis-formations', formation['id'], data)


def purge_formation(formation):
    _purge('deis-formations', formation['id'])


def _publish(data_bag, item_name, item_value):
    client = _get_client()
    body, status = client.update_databag_item(data_bag, item_name, item_value)
    if status != 200:
        body, status = client.create_databag_item(data_bag, item_name, item_value)
        if status != 201:
            raise RuntimeError('Could not publish {item_name}: {body}'.format(**locals()))
    return body, status


def _purge(databag_name, item_name):
    client = _get_client()
    body, status = client.delete_databag_item(databag_name, item_name)
    if status == 200 or status == 404:
        return body, status
    raise RuntimeError('Could not purge {item_name}: {body}'.format(**locals()))
