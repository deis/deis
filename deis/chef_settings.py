import os.path
import re
import subprocess

CHEF_CONFIG_PATH = '/etc/chef'
CHEF_INSTALL_TYPE = 'gems'
CHEF_RUBY_VERSION = '1.9.1'
CHEF_ENVIRONMENT = '_default'


# try to load chef config using CHEF_CONFIG_PATH
try:
    CHEF_ENABLED = True
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
    print 'Warning: failed to auto-configure Chef -- {}'.format(e)
    CHEF_ENABLED = False
