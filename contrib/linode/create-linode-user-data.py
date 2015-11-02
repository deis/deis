#!/usr/bin/env python
"""
Create CoreOS user-data by merging contrib/coreos/user-data and contrib/linode/linode-user-data-template

Usage: create-linode-user-data.py
"""
import base64
import sys
import os
import collections
import argparse

import yaml
import colorama
from colorama import Fore, Style
import requests


def combine_dicts(orig_dict, new_dict):
    for key, val in new_dict.iteritems():
        if isinstance(val, collections.Mapping):
            tmp = combine_dicts(orig_dict.get(key, {}), val)
            orig_dict[key] = tmp
        elif isinstance(val, list):
            orig_dict[key] = (orig_dict.get(key, []) + val)
        else:
            orig_dict[key] = new_dict[key]
    return orig_dict


def get_file(name, mode="r", abspath=False):
    current_dir = os.path.dirname(__file__)

    if abspath:
        return file(os.path.abspath(os.path.join(current_dir, name)), mode)
    else:
        return file(os.path.join(current_dir, name), mode)


def main():
    colorama.init()

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
            log_warning('Invalid public key: ' + public_key_file.name)

    if not len(public_keys) > 0:
        raise ValueError('Must supply at least one valid public key')

    linode_user_data = get_file("linode-user-data.yaml", "w", True)
    linode_template = get_file("linode-user-data-template.yaml")
    coreos_template = get_file("../coreos/user-data.example")

    coreos_template_string = coreos_template.read()
    coreos_template_string = coreos_template_string.replace('#DISCOVERY_URL', 'https://discovery.etcd.io/' + str(etcd_token))

    configuration_linode_template = yaml.safe_load(linode_template)
    configuration_coreos_template = yaml.safe_load(coreos_template_string)

    configuration = combine_dicts(configuration_coreos_template, configuration_linode_template)
    configuration['ssh_authorized_keys'] = public_keys

    dump = yaml.dump(configuration, default_flow_style=False, default_style='|')

    with linode_user_data as outfile:
        outfile.write("#cloud-config\n\n" + dump)
        log_success('Wrote Linode user data to ' + linode_user_data.name)


def validate_public_key(key):
    try:
        type, key_string, comment = key.split()
        data = base64.decodestring(key_string)
        return data[4:11] == type
    except:
        return False


def generate_etcd_token():
    log_info('Generating new Etcd token...')
    data = requests.get('https://discovery.etcd.io/new').text
    token = data.replace('https://discovery.etcd.io/', '')
    log_success('Generated new token: ' + token)
    return token


def validate_etcd_token(token):
    try:
        int(token, 16)
        return True
    except:
        return False


def log_info(message):
    print(Fore.CYAN + message + Fore.RESET)


def log_warning(message):
    print(Fore.YELLOW + message + Fore.RESET)


def log_success(message):
    print(Style.BRIGHT + Fore.GREEN + message + Fore.RESET + Style.RESET_ALL)


def log_error(message):
    print(Style.BRIGHT + Fore.RED + message + Fore.RESET + Style.RESET_ALL)


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        log_error(e.message)
        sys.exit(1)
