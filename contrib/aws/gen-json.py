#!/usr/bin/env python
import argparse
import json
import os
import urllib2
import yaml

parser = argparse.ArgumentParser()
parser.add_argument('--channel', help='the CoreOS channel to use', default='stable')
parser.add_argument('--version', help='the CoreOS version to use', default='current')
args = vars(parser.parse_args())

CURR_DIR = os.path.dirname(os.path.realpath(__file__))

# Add AWS-specific units to the shared user-data
FORMAT_DOCKER_VOLUME = '''
  [Unit]
  Description=Formats the added EBS volume for Docker
  ConditionPathExists=!/etc/docker-volume-formatted
  [Service]
  Type=oneshot
  RemainAfterExit=yes
  ExecStart=/usr/sbin/wipefs -f /dev/xvdf
  ExecStart=/usr/sbin/mkfs.ext4 -i 4096 -b 4096 /dev/xvdf
  ExecStart=/bin/touch /etc/docker-volume-formatted
'''
MOUNT_DOCKER_VOLUME = '''
  [Unit]
  Description=Mount Docker volume to /var/lib/docker
  Requires=format-docker-volume.service
  After=format-docker-volume.service
  Before=docker.service
  [Mount]
  What=/dev/xvdf
  Where=/var/lib/docker
  Type=ext4
'''
DOCKER_DROPIN = '''
  [Unit]
  Requires=var-lib-docker.mount
  After=var-lib-docker.mount
'''
FORMAT_ETCD_VOLUME = '''
  [Unit]
  Description=Formats the etcd device
  ConditionPathExists=!/etc/etcd-volume-formatted
  [Service]
  Type=oneshot
  RemainAfterExit=yes
  ExecStart=/usr/sbin/wipefs -f /dev/xvdg
  ExecStart=/usr/sbin/mkfs.ext4 -i 4096 -b 4096 /dev/xvdg
  ExecStart=/bin/touch /etc/etcd-volume-formatted
'''
MOUNT_ETCD_VOLUME = '''
  [Unit]
  Description=Mounts the etcd volume
  Requires=format-etcd-volume.service
  After=format-etcd-volume.service
  [Mount]
  What=/dev/xvdg
  Where=/media/etcd
  Type=ext4
'''
PREPARE_ETCD_DATA_DIRECTORY = '''
  [Unit]
  Description=Prepares the etcd data directory
  Requires=media-etcd.mount
  After=media-etcd.mount
  Before=etcd2.service
  [Service]
  Type=oneshot
  RemainAfterExit=yes
  ExecStart=/usr/bin/chown -R etcd:etcd /media/etcd
'''
ETCD_DROPIN = '''
  [Unit]
  Requires=prepare-etcd-data-directory.service
  After=prepare-etcd-data-directory.service
'''

def coreos_amis(channel, version):
    url = "http://{channel}.release.core-os.net/amd64-usr/{version}/coreos_production_ami_all.json".format(channel=channel, version=version)
    try:
        amis = json.load(urllib2.urlopen(url))
    except (IOError, ValueError):
        print "The URL {} is invalid.".format(url)
        raise

    return dict(map(lambda n: (n['name'], dict(PV=n['pv'], HVM=n['hvm'])), amis['amis']))

new_units = [
    dict({'name': 'format-docker-volume.service', 'command': 'start', 'content': FORMAT_DOCKER_VOLUME}),
    dict({'name': 'var-lib-docker.mount', 'command': 'start', 'content': MOUNT_DOCKER_VOLUME}),
    dict({'name': 'docker.service', 'drop-ins': [{'name': '90-after-docker-volume.conf', 'content': DOCKER_DROPIN}]}),
    dict({'name': 'format-etcd-volume.service', 'command': 'start', 'content': FORMAT_ETCD_VOLUME}),
    dict({'name': 'media-etcd.mount', 'command': 'start', 'content': MOUNT_ETCD_VOLUME}),
    dict({'name': 'prepare-etcd-data-directory.service', 'command': 'start', 'content': PREPARE_ETCD_DATA_DIRECTORY}),
    dict({'name': 'etcd2.service', 'drop-ins': [{'name': '90-after-etcd-volume.conf', 'content': ETCD_DROPIN}]})
]

with open(os.path.join(CURR_DIR, '..', 'coreos', 'user-data'), 'r') as f:
    data = yaml.safe_load(f)

# coreos-cloudinit will start the units in order, so we want these to be processed before etcd/fleet
# are started
data['coreos']['units'] = new_units + data['coreos']['units']

# Point to the right data directory
data['coreos']['etcd2']['data-dir'] = '/media/etcd'

header = ["#cloud-config", "---"]
dump = yaml.dump(data, default_flow_style=False)

template = json.load(open(os.path.join(CURR_DIR, 'deis.template.json'), 'r'))

template['Resources']['CoreOSServerLaunchConfig']['Properties']['UserData']['Fn::Base64']['Fn::Join'] = ["\n", header + dump.split("\n")]
template['Parameters']['ClusterSize']['Default'] = str(os.getenv('DEIS_NUM_INSTANCES', 3))
template['Mappings']['CoreOSAMIs'] = coreos_amis(args['channel'], args['version'])

VPC_ID = os.getenv('VPC_ID', None)
VPC_SUBNETS = os.getenv('VPC_SUBNETS', None)
VPC_PRIVATE_SUBNETS = os.getenv('VPC_PRIVATE_SUBNETS', VPC_SUBNETS)
VPC_ZONES = os.getenv('VPC_ZONES', None)

if VPC_ID and VPC_SUBNETS and VPC_ZONES and len(VPC_SUBNETS.split(',')) == len(VPC_ZONES.split(',')):
    # skip VPC, subnet, route, and internet gateway creation
    del template['Resources']['VPC']
    del template['Resources']['Subnet1']
    del template['Resources']['Subnet2']
    del template['Resources']['Subnet1RouteTableAssociation']
    del template['Resources']['Subnet2RouteTableAssociation']
    del template['Resources']['InternetGateway']
    del template['Resources']['GatewayToInternet']
    del template['Resources']['PublicRouteTable']
    del template['Resources']['PublicRoute']
    del template['Resources']['CoreOSServerLaunchConfig']['DependsOn']
    del template['Resources']['DeisWebELB']['DependsOn']

    # update VpcId fields
    template['Resources']['DeisWebELBSecurityGroup']['Properties']['VpcId'] = VPC_ID
    template['Resources']['VPCSecurityGroup']['Properties']['VpcId'] = VPC_ID

    # update subnets and zones
    template['Resources']['CoreOSServerAutoScale']['Properties']['AvailabilityZones'] = VPC_ZONES.split(',')
    template['Resources']['CoreOSServerAutoScale']['Properties']['VPCZoneIdentifier'] = VPC_PRIVATE_SUBNETS.split(',')
    template['Resources']['DeisWebELB']['Properties']['Subnets'] = VPC_SUBNETS.split(',')

print json.dumps(template)
