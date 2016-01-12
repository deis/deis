#!/usr/bin/env python
"""
Apply a "Security Group" to the members of an etcd cluster.

Usage: apply-firewall.py
"""
import re
import string
import argparse
import threading
from threading import Thread
import uuid

import paramiko
import requests
import sys
import yaml

from linodeapi import LinodeApiCommand
import linodeutils 

class FirewallCommand(LinodeApiCommand):


    def get_nodes_from_args(self):
        if not self.discovery_url:
            self.discovery_url = self.get_discovery_url_from_user_data()
        return self.get_nodes_from_discovery_url(self.discovery_url)
    
    def get_nodes_from_discovery_url(self, discovery_url):
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


    def get_discovery_url_from_user_data(self):
        name = 'linode-user-data.yaml'
        linodeutils.log_info('Loading discovery url from ' + name)
        try:
            user_data_file = linodeutils.get_file(name)
            user_data_yaml = yaml.safe_load(user_data_file)
            return user_data_yaml['coreos']['etcd2']['discovery']
        except:
            raise IOError('Could not load discovery url from ' + name)


    def validate_ip_address(self, ip):
        return True if re.match('([0-9]{1,3}\.){3}[0-9]{1,3}', ip) else False


    def get_firewall_contents(self, node_ips):
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
-A Firewall-INPUT -m conntrack --ctstate NEW -m multiport$multiport_private -p tcp --dports 22,2222,80,443$add_new_nodes -j ACCEPT
# Log and drop everything else
-A Firewall-INPUT -j REJECT
COMMIT
"""

        multiport_private = ' -s 192.168.0.0/16' if self.private else ''
        add_new_nodes = ',2379,2380' if self.adding_new_nodes else ''
        
        rules_template = string.Template(rules_template_text)
        return rules_template.substitute(node_ips=string.join(node_ips, ','), multiport_private=multiport_private, add_new_nodes=add_new_nodes)


    def apply_rules_to_all(self, host_ips, rules):
        pkey = self.detect_and_create_private_key()
        
        threads = []
        for ip in host_ips:
            t = Thread(target=self.apply_rules, args=(ip, rules, pkey))
            t.setDaemon(False)
            t.start()
            threads.append(t)
        for thread in threads:
            thread.join()


    def detect_and_create_private_key(self):
        private_key_text = self.private_key.read()
        self.private_key.seek(0)
        if '-----BEGIN RSA PRIVATE KEY-----' in private_key_text:
            return paramiko.RSAKey.from_private_key(self.private_key)
        elif '-----BEGIN DSA PRIVATE KEY-----' in private_key_text:
            return paramiko.DSSKey.from_private_key(self.private_key)
        else:
            raise ValueError('Invalid private key file ' + self.private_key.name)


    def apply_rules(self, host_ip, rules, private_key):
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

        linodeutils.log_success('Applied rule to ' + host_ip)

    def acquire_linode_ips(self):
        linodeutils.log_info('Getting info for Linodes from display group: ' + self.node_display_group)
        deis_grp = self.request('linode.list')
        deis_linodeids = [l.get('LINODEID','') for l in deis_grp if l.get('LPM_DISPLAYGROUP', '') == self.node_display_group]
        deis_grp_ips = self.request('linode.ip.list')
        self.deis_privateips = [ip.get('IPADDRESS','') for ip in deis_grp_ips if (ip.get('LINODEID','') in deis_linodeids) and (ip.get('ISPUBLIC', 1) == 0)]
        self.deis_publicips = [ip.get('IPADDRESS','') for ip in deis_grp_ips if (ip.get('LINODEID','') in deis_linodeids) and (ip.get('ISPUBLIC', 0) == 1)]

    def run(self):
        #NOTE: defaults to using display group, then manual input (via nodes/hosts), then discovery_url
        if self.node_display_group:
            self.acquire_linode_ips()
            nodes = self.deis_privateips
            hosts = self.deis_publicips
        else:
            nodes = self.nodes if self.nodes is not None else self.get_nodes_from_args()
            hosts = self.hosts if self.hosts is not None else nodes

        node_ips = []
        for ip in nodes:
            if self.validate_ip_address(ip):
                node_ips.append(ip)
            else:
                linodeutils.log_warning('Invalid IP will not be added to security group: ' + ip)

        if not len(node_ips) > 0:
            raise ValueError('No valid IP addresses in security group.')

        host_ips = []
        for ip in hosts:
            if self.validate_ip_address(ip):
                host_ips.append(ip)
            else:
                linodeutils.log_warning('Host has invalid IP address: ' + ip)

        if not len(host_ips) > 0:
            raise ValueError('No valid host addresses.')

        linodeutils.log_info('Generating iptables rules...')
        rules = self.get_firewall_contents(node_ips)
        linodeutils.log_success('Generated rules:')
        linodeutils.log_debug(rules)
        
        linodeutils.log_info('Applying rules...')
        self.apply_rules_to_all(host_ips, rules)
        linodeutils.log_success('Done!')


def main():
    linodeutils.init()

    parser = argparse.ArgumentParser(description='Apply a "Security Group" to a Deis cluster')
    parser.add_argument('--api-key', dest='linode_api_key', help='Linode API Key')
    parser.add_argument('--private-key', required=True, type=file, dest='private_key', help='Cluster SSH Private Key')
    parser.add_argument('--private', action='store_true', dest='private', help='Only allow access to the cluster from the private network')
    parser.add_argument('--adding-new-nodes', action='store_true', dest='adding_new_nodes', help='When adding new nodes to existing cluster, allows access to etcd')
    parser.add_argument('--discovery-url', dest='discovery_url', help='Etcd discovery url')
    parser.add_argument('--display-group', required=False, dest='node_display_group', help='Display group (used for Linode IP discovery).')
    parser.add_argument('--hosts', nargs='+', dest='hosts', help='The public IP addresses of the hosts to apply rules to (for ssh)')
    parser.add_argument('--nodes', nargs='+', dest='nodes', help='The private IP addresses of the hosts (for iptable setup)')
    parser.set_defaults(cmd=FirewallCommand)
    
    args = parser.parse_args()
    cmd = args.cmd(args)
    args.cmd(args).run()
    
if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        linodeutils.log_error(e.message)
        sys.exit(1)
