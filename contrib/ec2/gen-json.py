#!/usr/bin/env python
import json
import os
import yaml

CURR_DIR = os.path.dirname(os.path.realpath(__file__))

# Add EC2-specific units to the shared user-data
FORMAT_EPHEMERAL_VOLUME = '''
  [Unit]
  Description=Formats the ephemeral volume
  ConditionPathExists=!/etc/ephemeral-volume-formatted
  [Service]
  Type=oneshot
  RemainAfterExit=yes
  ExecStart=/usr/sbin/wipefs -f /dev/xvdb
  ExecStart=/usr/sbin/mkfs.ext4 /dev/xvdb
  ExecStart=/bin/touch /etc/ephemeral-volume-formatted
'''
MOUNT_EPHEMERAL_VOLUME = '''
  [Unit]
  Description=Formats and mounts the ephemeral drive
  Requires=format-ephemeral-volume.service
  After=format-ephemeral-volume.service
  [Mount]
  What=/dev/xvdb
  Where=/media/ephemeral
  Type=ext4
'''
PREPARE_ETCD_DATA_DIRECTORY = '''
  [Unit]
  Description=Prepares the etcd data directory
  Requires=media-ephemeral.mount
  After=media-ephemeral.mount
  Before=etcd.service
  [Service]
  Type=oneshot
  RemainAfterExit=yes
  ExecStart=/usr/bin/mkdir -p /media/ephemeral/etcd
  ExecStart=/usr/bin/chown -R etcd /media/ephemeral/etcd
'''
FORMAT_DOCKER_VOLUME = '''
  [Unit]
  Description=Formats the added EBS volume for Docker
  ConditionPathExists=!/etc/docker-volume-formatted
  [Service]
  Type=oneshot
  RemainAfterExit=yes
  ExecStart=/usr/sbin/wipefs -f /dev/xvdf
  ExecStart=/usr/sbin/mkfs.btrfs -f /dev/xvdf
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
  Type=btrfs
'''

data = yaml.load(file(os.path.join(CURR_DIR, '..', 'coreos', 'user-data'), 'r'))
data['coreos']['units'].append(dict({'name': 'format-ephemeral-volume.service', 'command': 'start', 'content': FORMAT_EPHEMERAL_VOLUME}))
data['coreos']['units'].append(dict({'name': 'media-ephemeral.mount', 'command': 'start', 'content': MOUNT_EPHEMERAL_VOLUME}))
data['coreos']['units'].append(dict({'name': 'prepare-etcd-data-directory.service', 'command': 'start', 'content': PREPARE_ETCD_DATA_DIRECTORY}))
data['coreos']['units'].append(dict({'name': 'format-docker-volume.service', 'command': 'start', 'content': FORMAT_DOCKER_VOLUME}))
data['coreos']['units'].append(dict({'name': 'var-lib-docker.mount', 'command': 'start', 'content': MOUNT_DOCKER_VOLUME}))

# configure etcd to use the ephemeral drive
data['coreos']['etcd']['data-dir'] = '/media/ephemeral/etcd'

header = ["#cloud-config", "---"]
dump = yaml.dump(data, default_flow_style=False)

template = json.load(open(os.path.join(CURR_DIR, 'deis.template.json'),'r'))

template['Resources']['CoreOSServerLaunchConfig']['Properties']['UserData']['Fn::Base64']['Fn::Join'] = [ "\n", header + dump.split("\n") ]
template['Parameters']['ClusterSize']['Default'] = str(os.getenv('DEIS_NUM_INSTANCES', 3))

VPC_ID = os.getenv('VPC_ID', None)
VPC_SUBNETS = os.getenv('VPC_SUBNETS', None)
VPC_ZONES = os.getenv('VPC_ZONES', None)
ELB_SCHEME = os.getenv('ELB_SCHEME', None)

if ELB_SCHEME:
  # Update ELB Scheme
  template['Resources']['DeisWebELB']['Properties']['Scheme'] = ELB_SCHEME

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
  template['Resources']['CoreOSServerAutoScale']['Properties']['VPCZoneIdentifier'] = VPC_SUBNETS.split(',')
  template['Resources']['DeisWebELB']['Properties']['Subnets'] = VPC_SUBNETS.split(',')

print json.dumps(template)
