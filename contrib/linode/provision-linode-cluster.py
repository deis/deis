#!/usr/bin/env python
"""
Provision a Deis cluster on Linode

Usage: provision-linode-cluster.py
"""
import argparse
import random
import string
import threading
from threading import Thread
import sys

import paramiko
from linodeapi import LinodeApiCommand
import linodeutils

class ProvisionCommand(LinodeApiCommand):
    _created_linodes = []

    def run(self):
        # validate arguments
        self._check_num_nodes()
        self._check_plan_size()

        # create the linodes
        self._create_linodes()

        # print the results
        self._report_created()

    def _report_created(self):
        # set up the report data
        rows = []
        ips = []
        data_center = self._get_data_center().get('ABBR')
        plan = self._get_plan().get('RAM')

        for linode in self._created_linodes:
            rows.append((
                linode['hostname'],
                linode['public'],
                linode['private'],
                linode['gateway'],
                data_center,
                plan
            ))
            ips.append(linode['public'])

        firewall_command = './apply-firewall.py --private-key /path/to/key/deis --hosts ' + string.join(ips, ' ')
        header_msg = '  Successfully provisioned ' + str(self.num_nodes) + ' nodes!'
        footer_msg = '  Finish up your installation by securing your cluster with the following command:\n  ' + firewall_command + '\n'
        linodeutils.log_table(rows, header_msg, footer_msg)

    def _get_plan(self):
        if self._plan is None:
            plans = self.request('avail.linodeplans', params={'PlanID': self.node_plan})
            if len(plans) != 1:
                raise ValueError('The --plan specified is invalid. Use the `list-plans` subcommand to see valid ids.')
            self._plan = plans[0]
        return self._plan

    def _get_plan_id(self):
        return self._get_plan().get('PLANID')

    def _get_data_center(self):
        if self._data_center is None:
            data_centers = self.request('avail.datacenters')
            for data_center in data_centers:
                if data_center.get('DATACENTERID') == self.node_data_center:
                    self._data_center = data_center
            if self._data_center is None:
                raise ValueError('The --datacenter specified is invalid. Use the `list-data-centers` subcommand to see valid ids.')
        return self._data_center

    def _get_data_center_id(self):
        return self._get_data_center().get('DATACENTERID')

    def _check_plan_size(self):
        ram = self._get_plan().get('RAM')
        if ram < 4096:
            raise ValueError('Deis cluster members must have at least 4GB of memory. Please choose a plan with more memory.')

    def _check_num_nodes(self):
        if self.num_nodes < 1:
            raise ValueError('Must provision at least one node.')
        elif self.num_nodes < 3:
            linodeutils.log_warning('A Deis cluster must have 3 or more nodes, only continue if you adding to a current cluster.')
            linodeutils.log_warning('Continue? (y/n)')
            accept = None
            while True:
                if accept == 'y':
                    return
                elif accept == 'n':
                    raise StandardError('User canceled provisioning')
                else:
                    accept = self._get_user_input('--> ').strip().lower()

    def _get_user_input(self, prompt):
        if sys.version_info[0] < 3:
            return raw_input(prompt)
        else:
            return input(prompt)

    def _create_linodes(self):
        threads = []
        for i in range(0, self.num_nodes):
            t = Thread(target=self._create_linode,
                       args=(self._get_plan_id(), self._get_data_center_id(), self.node_name_prefix, self.node_display_group))
            t.setDaemon(False)
            t.start()

            threads.append(t)

        for thread in threads:
            thread.join()

    def _create_linode(self, plan_id, data_center_id, name_prefix, display_group):
        self.info('Creating the Linode...')

        # create the linode
        node_id = self.request('linode.create', params={
            'DatacenterID': data_center_id,
            'PlanID': plan_id
        }).get('LinodeID')

        # update the configuration
        self.request('linode.update', params={
            'LinodeID': node_id,
            'Label': name_prefix + str(node_id),
            'lpm_displayGroup': display_group,
            'Alert_cpu_enabled': False,
            'Alert_diskio_enabled': False,
            'Alert_bwin_enabled': False,
            'Alert_bwout_enabled': False,
            'Alert_bwquota_enabled': False
        })

        self.success('Linode ' + str(node_id) + ' created!')
        hostname = name_prefix + str(node_id)
        threading.current_thread().name = hostname

        # configure the networking
        network = self._configure_networking(node_id)
        network['hostname'] = hostname

        # generate a password for the provisioning disk
        password = ''.join(random.SystemRandom().choice(string.ascii_uppercase + string.digits) for _ in range(24))

        # configure the disks
        total_hd = self.request('linode.list', params={'LinodeID': node_id})[0]['TOTALHD']
        provision_disk_mb = 600
        coreos_disk_mb = total_hd - provision_disk_mb
        provision_disk_id = self._create_provisioning_disk(node_id, provision_disk_mb, password)
        coreos_disk_id = self._create_coreos_disk(node_id, coreos_disk_mb)

        # create the provision config
        provision_config_id = self._create_provision_profile(node_id, provision_disk_id, coreos_disk_id)

        # create the CoreOS config
        coreos_config_id = self._create_coreos_profile(node_id, coreos_disk_id)

        # install CoreOS
        self._install_coreos(node_id, provision_config_id, network, password)

        # boot in to coreos
        self.info('Booting into CoreOS configuration profile...')
        self.request('linode.reboot', params={'LinodeID': node_id, 'ConfigID': coreos_config_id})

        # append the linode to the created list
        self._created_linodes.append(network)

    def _configure_networking(self, node_id):
        self.info('Configuring network...')

        # add the private network
        self.request('linode.ip.addprivate', params={'LinodeID': node_id})

        # pull the network config
        ip_data = self.request('linode.ip.list', params={'LinodeID': node_id})

        network = {'public': None, 'private': None, 'gateway': None}

        for ip in ip_data:
            if ip.get('ISPUBLIC') == 1:
                network['public'] = ip.get('IPADDRESS')
                # the gateway is the public ip with the last octet set to 1
                split_ip = str(network['public']).split('.')
                split_ip[3] = '1'
                network['gateway'] = string.join(split_ip, '.')
            else:
                network['private'] = ip.get('IPADDRESS')

        if network.get('public') is None:
            raise RuntimeError('Public IP address could not be found.')

        if network.get('private') is None:
            raise RuntimeError('Private IP address could not be found.')

        self.success('Network configured!')
        self.success('    Public IP:  ' + str(network['public']))
        self.success('    Private IP: ' + str(network['private']))
        self.success('    Gateway:    ' + str(network['gateway']))

        return network

    def _create_provisioning_disk(self, node_id, size, root_password):
        self.info('Creating provisioning disk...')

        disk_id = self.request('linode.disk.createfromdistribution', params={
            'LinodeID': node_id,
            'Label': 'Provision',
            'DistributionID': 130,
            'Type': 'ext4',
            'Size': size,
            'rootPass': root_password
        }).get('DiskID')

        self.success('Created provisioning disk!')

        return disk_id

    def _create_coreos_disk(self, node_id, size):
        self.info('Creating CoreOS disk...')

        disk_id = self.request('linode.disk.create', params={
            'LinodeID': node_id,
            'Label': 'CoreOS',
            'Type': 'ext4',
            'Size': size
        }).get('DiskID')

        self.success('Created CoreOS disk!')

        return disk_id

    def _create_provision_profile(self, node_id, provision_disk_id, coreos_disk_id):
        self.info('Creating Provision configuration profile...')

        # create a disk the total hd size
        config_id = self.request('linode.config.create', params={
            'LinodeID': node_id,
            'KernelID': 138,
            'Label': 'Provision',
            'DiskList': str(provision_disk_id) + ',' + str(coreos_disk_id)
        }).get('ConfigID')

        self.success('Provision profile created!')

        return config_id

    def _create_coreos_profile(self, node_id, coreos_disk_id):
        self.info('Creating CoreOS configuration profile...')

        # create a disk the total hd size
        config_id = self.request('linode.config.create', params={
            'LinodeID': node_id,
            'KernelID': 213,
            'Label': 'CoreOS',
            'DiskList': str(coreos_disk_id)
        }).get('ConfigID')

        self.success('CoreOS profile created!')

        return config_id

    def _get_cloud_config(self):
        if self.cloud_config_text is None:
            self.cloud_config_text = self.cloud_config.read()
        return self.cloud_config_text

    def _install_coreos(self, node_id, provision_config_id, network, password):
        self.info('Installing CoreOS...')

        # boot in to the provision configuration
        self.info('Booting into Provision configuration profile...')
        self.request('linode.boot', params={'LinodeID': node_id, 'ConfigID': provision_config_id})

        # connect to the server via ssh
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

        while True:
            try:
                ssh.connect(str(network['public']), username='root', password=password, allow_agent=False, look_for_keys=False)
                break
            except:
                continue

        # copy the cloud config
        self.info('Pushing cloud config...')
        cloud_config_template = string.Template(self._get_cloud_config())
        cloud_config = cloud_config_template.safe_substitute(public_ipv4=network['public'], private_ipv4=network['private'], gateway=network['gateway'],
                                                             hostname=network['hostname'])

        sftp = ssh.open_sftp()
        sftp.open('cloud-config.yaml', 'w').write(cloud_config)

        self.info('Installing...')

        commands = [
            'wget https://raw.githubusercontent.com/coreos/init/master/bin/coreos-install -O $HOME/coreos-install',
            'chmod +x $HOME/coreos-install',
            '$HOME/coreos-install -d /dev/sdb -C ' + self.coreos_channel + ' -V ' + self.coreos_version + ' -c $HOME/cloud-config.yaml -t /dev/shm'
        ]

        for command in commands:
            stdin, stdout, stderr = ssh.exec_command(command)
            stdout.channel.recv_exit_status()
            print stdout.read()

        ssh.close()


class ListDataCentersCommand(LinodeApiCommand):
    def run(self):
        data = self.request('avail.datacenters')
        column_format = "{:<4} {:}"
        linodeutils.log_success(column_format.format(*('ID', 'LOCATION')))
        for data_center in data:
            row = (
                data_center.get('DATACENTERID'),
                data_center.get('LOCATION')
            )
            linodeutils.log_minor_success(column_format.format(*row))


class ListPlansCommand(LinodeApiCommand):
    def run(self):
        data = self.request('avail.linodeplans')
        column_format = "{:<4} {:<16} {:<8} {:<12} {:}"
        linodeutils.log_success(column_format.format(*('ID', 'LABEL', 'CORES', 'RAM', 'PRICE')))
        for plan in data:
            row = (
                plan.get('PLANID'),
                plan.get('LABEL'),
                plan.get('CORES'),
                str(plan.get('RAM')) + 'MB',
                '$' + str(plan.get('PRICE'))
            )
            linodeutils.log_minor_success(column_format.format(*row))


def main():
    linodeutils.init()

    parser = argparse.ArgumentParser(description='Provision Linode Deis Cluster')
    parser.add_argument('--api-key', required=True, dest='linode_api_key', help='Linode API Key')
    subparsers = parser.add_subparsers()

    provision_parser = subparsers.add_parser('provision', help="Provision the Deis cluster")
    provision_parser.add_argument('--num', required=False, default=3, type=int, dest='num_nodes', help='Number of nodes to provision')
    provision_parser.add_argument('--name-prefix', required=False, default='deis', dest='node_name_prefix', help='Node name prefix')
    provision_parser.add_argument('--display-group', required=False, default='deis', dest='node_display_group', help='Node display group')
    provision_parser.add_argument('--plan', required=False, default=4, type=int, dest='node_plan', help='Node plan id. Use list-plans to find the id.')
    provision_parser.add_argument('--datacenter', required=False, default=2, type=int, dest='node_data_center',
                                  help='Node data center id. Use list-data-centers to find the id.')
    provision_parser.add_argument('--cloud-config', required=False, default='linode-user-data.yaml', type=file, dest='cloud_config',
                                  help='CoreOS cloud config user-data file')
    provision_parser.add_argument('--coreos-version', required=False, default='899.13.0', dest='coreos_version',
                                  help='CoreOS version number to install')
    provision_parser.add_argument('--coreos-channel', required=False, default='stable', dest='coreos_channel',
                                  help='CoreOS channel to install from')
    provision_parser.set_defaults(cmd=ProvisionCommand)

    list_data_centers_parser = subparsers.add_parser('list-data-centers', help="Lists the available Linode data centers.")
    list_data_centers_parser.set_defaults(cmd=ListDataCentersCommand)

    list_plans_parser = subparsers.add_parser('list-plans', help="Lists the available Linode plans.")
    list_plans_parser.set_defaults(cmd=ListPlansCommand)

    args = parser.parse_args()
    args.cmd(args).run()


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        linodeutils.log_error(e.message)
        sys.exit(1)
