#!/usr/bin/env python
"""
Create CoreOS user-data by merging contrib/coreos/user-data and contrib/linode/linode-user-data-template

Usage: create-linode-user-data.py
"""
import base64
import sys
import argparse

import yaml
import requests
import linodeutils


def validate_public_key(key):
    try:
        type, key_string, comment = key.split()
        data = base64.decodestring(key_string)
        return data[4:11] == type
    except:
        return False


def generate_etcd_token():
    linodeutils.log_info('Generating new Etcd token...')
    data = requests.get('https://discovery.etcd.io/new').text
    token = data.replace('https://discovery.etcd.io/', '')
    linodeutils.log_success('Generated new token: ' + token)
    return token


def validate_etcd_token(token):
    try:
        int(token, 16)
        return True
    except:
        return False


def main():
    linodeutils.init()

    parser = argparse.ArgumentParser(description='Create Linode User Data')
    parser.add_argument('--public-key', action='append', required=True, type=file, dest='public_key_files', help='Authorized SSH Keys')
    parser.add_argument('--etcd-token', required=False, default=None, dest='etcd_token', help='Etcd Token')
    args = parser.parse_args()

    etcd_token = args.etcd_token
    if etcd_token is None:
        etcd_token = generate_etcd_token()
    else:
        if not validate_etcd_token(args.etcd_token):
            raise ValueError('Invalid Etcd Token. You can generate a new token at https://discovery.etcd.io/new.')

    public_keys = []
    for public_key_file in args.public_key_files:
        public_key = public_key_file.read()
        if validate_public_key(public_key):
            public_keys.append(public_key)
        else:
            linodeutils.log_warning('Invalid public key: ' + public_key_file.name)

    if not len(public_keys) > 0:
        raise ValueError('Must supply at least one valid public key')

    linode_user_data = linodeutils.get_file("linode-user-data.yaml", "w", True)
    linode_template = linodeutils.get_file("linode-user-data-template.yaml")
    coreos_template = linodeutils.get_file("../coreos/user-data.example")

    coreos_template_string = coreos_template.read()
    coreos_template_string = coreos_template_string.replace('#DISCOVERY_URL', 'https://discovery.etcd.io/' + str(etcd_token))

    configuration_linode_template = yaml.safe_load(linode_template)
    configuration_coreos_template = yaml.safe_load(coreos_template_string)

    configuration = linodeutils.combine_dicts(configuration_coreos_template, configuration_linode_template)
    configuration['ssh_authorized_keys'] = public_keys

    dump = yaml.dump(configuration, default_flow_style=False, default_style='|')

    with linode_user_data as outfile:
        outfile.write("#cloud-config\n\n" + dump)
        linodeutils.log_success('Wrote Linode user data to ' + linode_user_data.name)

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        linodeutils.log_error(e.message)
        sys.exit(1)
