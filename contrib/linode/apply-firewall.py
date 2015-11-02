#!/usr/bin/env python
"""
Apply a "Security Group" to the members of an etcd cluster.

Usage: apply-firewall.py
"""
import os
import re
import string
import argparse
from threading import Thread
import uuid

import colorama
from colorama import Fore, Style
import paramiko
import requests
import sys
import yaml


def get_nodes_from_args(args):
    if args.discovery_url is not None:
        return get_nodes_from_discovery_url(args.discovery_url)

    return get_nodes_from_discovery_url(get_discovery_url_from_user_data())


def get_nodes_from_discovery_url(discovery_url):
    try:
        nodes = []
        json = requests.get(discovery_url).json()
        discovery_nodes = json['node']['nodes']
        for node in discovery_nodes:
            value = node['value']
            ip = re.search('([0-9]{1,3}\.){3}[0-9]{1,3}', value).group(0)
            nodes.append(ip)
        return nodes
    except:
        raise IOError('Could not load nodes from discovery url ' + discovery_url)


def get_discovery_url_from_user_data():
    name = 'linode-user-data.yaml'
    log_info('Loading discovery url from ' + name)
    try:
        current_dir = os.path.dirname(__file__)
        user_data_file = file(os.path.abspath(os.path.join(current_dir, name)), 'r')
        user_data_yaml = yaml.safe_load(user_data_file)
        return user_data_yaml['coreos']['etcd2']['discovery']
    except:
        raise IOError('Could not load discovery url from ' + name)


def validate_ip_address(ip):
    return True if re.match('([0-9]{1,3}\.){3}[0-9]{1,3}', ip) else False


def get_firewall_contents(node_ips, private=False):
    rules_template_text = """*filter
:INPUT DROP [0:0]
:FORWARD DROP [0:0]
:OUTPUT ACCEPT [0:0]
:DOCKER - [0:0]
:Firewall-INPUT - [0:0]
-A INPUT -j Firewall-INPUT
-A FORWARD -j Firewall-INPUT
-A Firewall-INPUT -i lo -j ACCEPT
-A Firewall-INPUT -p icmp --icmp-type echo-reply -j ACCEPT
-A Firewall-INPUT -p icmp --icmp-type destination-unreachable -j ACCEPT
-A Firewall-INPUT -p icmp --icmp-type time-exceeded -j ACCEPT
# Ping
-A Firewall-INPUT -p icmp --icmp-type echo-request -j ACCEPT
# Accept any established connections
-A Firewall-INPUT -m conntrack --ctstate  ESTABLISHED,RELATED -j ACCEPT
# Enable the traffic between the nodes of the cluster
-A Firewall-INPUT -s $node_ips -j ACCEPT
# Allow connections from docker container
-A Firewall-INPUT -i docker0 -j ACCEPT
# Accept ssh, http, https and git
-A Firewall-INPUT -m conntrack --ctstate NEW -m multiport$multiport_private -p tcp --dports 22,2222,80,443 -j ACCEPT
# Log and drop everything else
-A Firewall-INPUT -j REJECT
COMMIT
"""

    multiport_private = ' -s 192.168.0.0/16' if private else ''

    rules_template = string.Template(rules_template_text)
    return rules_template.substitute(node_ips=string.join(node_ips, ','), multiport_private=multiport_private)


def apply_rules_to_all(host_ips, rules, private_key):
    pkey = detect_and_create_private_key(private_key)

    threads = []
    for ip in host_ips:
        t = Thread(target=apply_rules, args=(ip, rules, pkey))
        t.setDaemon(False)
        t.start()
        threads.append(t)
    for thread in threads:
        thread.join()


def detect_and_create_private_key(private_key):
    private_key_text = private_key.read()
    private_key.seek(0)
    if '-----BEGIN RSA PRIVATE KEY-----' in private_key_text:
        return paramiko.RSAKey.from_private_key(private_key)
    elif '-----BEGIN DSA PRIVATE KEY-----' in private_key_text:
        return paramiko.DSSKey.from_private_key(private_key)
    else:
        raise ValueError('Invalid private key file ' + private_key.name)


def apply_rules(host_ip, rules, private_key):
    # connect to the server via ssh
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect(host_ip, username='core', allow_agent=False, look_for_keys=False, pkey=private_key)

    # copy the rules to the temp directory
    temp_file = '/tmp/' + str(uuid.uuid4())

    ssh.open_sftp()
    sftp = ssh.open_sftp()
    sftp.open(temp_file, 'w').write(rules)

    # move the rules in to place and enable and run the iptables-restore.service
    commands = [
        'sudo mv ' + temp_file + ' /var/lib/iptables/rules-save',
        'sudo chown root:root /var/lib/iptables/rules-save',
        'sudo systemctl enable iptables-restore.service',
        'sudo systemctl start iptables-restore.service'
    ]

    for command in commands:
        stdin, stdout, stderr = ssh.exec_command(command)
        stdout.channel.recv_exit_status()

    ssh.close()

    log_success('Applied rule to ' + host_ip)


def main():
    colorama.init()

    parser = argparse.ArgumentParser(description='Apply a "Security Group" to a Deis cluster')
    parser.add_argument('--private-key', required=True, type=file, dest='private_key', help='Cluster SSH Private Key')
    parser.add_argument('--private', action='store_true', dest='private', help='Only allow access to the cluster from the private network')
    parser.add_argument('--discovery-url', dest='discovery_url', help='Etcd discovery url')
    parser.add_argument('--hosts', nargs='+', dest='hosts', help='The IP addresses of the hosts to apply rules to')
    args = parser.parse_args()

    nodes = get_nodes_from_args(args)
    hosts = args.hosts if args.hosts is not None else nodes

    node_ips = []
    for ip in nodes:
        if validate_ip_address(ip):
            node_ips.append(ip)
        else:
            log_warning('Invalid IP will not be added to security group: ' + ip)

    if not len(node_ips) > 0:
        raise ValueError('No valid IP addresses in security group.')

    host_ips = []
    for ip in hosts:
        if validate_ip_address(ip):
            host_ips.append(ip)
        else:
            log_warning('Host has invalid IP address: ' + ip)

    if not len(host_ips) > 0:
        raise ValueError('No valid host addresses.')

    log_info('Generating iptables rules...')
    rules = get_firewall_contents(node_ips, args.private)
    log_success('Generated rules:')
    log_debug(rules)

    log_info('Applying rules...')
    apply_rules_to_all(host_ips, rules, args.private_key)
    log_success('Done!')


def log_debug(message):
    print(Style.DIM + Fore.MAGENTA + message + Fore.RESET + Style.RESET_ALL)


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
