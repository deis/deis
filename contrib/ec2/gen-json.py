#!/usr/bin/env python
import json
import os
import yaml

CURR_DIR = os.path.dirname(os.path.realpath(__file__))

# Add EC2-specific units to the shared user-data
FORMAT_EPHEMERAL = '''
  [Unit]
  Description=Formats the ephemeral drive
  ConditionPathExists=!/etc/docker-volume-formatted
  [Service]
  Type=oneshot
  RemainAfterExit=yes
  ExecStart=/usr/sbin/wipefs -f /dev/xvdf
  ExecStart=/usr/sbin/mkfs.btrfs -f /dev/xvdf
  ExecStart=/bin/touch /etc/docker-volume-formatted
'''
DOCKER_MOUNT = '''
  [Unit]
  Description=Mount ephemeral to /var/lib/docker
  Requires=format-ephemeral.service
  After=format-ephemeral.service
  Before=docker.service
  [Mount]
  What=/dev/xvdf
  Where=/var/lib/docker
  Type=btrfs
'''

data = yaml.load(file(os.path.join(CURR_DIR, '..', 'coreos', 'user-data'), 'r'))
data['coreos']['units'].append(dict({'name': 'format-ephemeral.service', 'command': 'start', 'content': FORMAT_EPHEMERAL}))
data['coreos']['units'].append(dict({'name': 'var-lib-docker.mount', 'command': 'start', 'content': DOCKER_MOUNT}))

header = ["#cloud-config", "---"]
dump = yaml.dump(data, default_flow_style=False)

template = json.load(open(os.path.join(CURR_DIR, 'deis.template.json'),'r'))

template['Resources']['CoreOSServerLaunchConfig']['Properties']['UserData']['Fn::Base64']['Fn::Join'] = [ "\n", header + dump.split("\n") ]
template['Parameters']['ClusterSize']['Default'] = str(os.getenv('DEIS_NUM_INSTANCES', 3))

VPC_ID = os.getenv('VPC_ID', None)
VPC_SUBNETS = os.getenv('VPC_SUBNETS', None)
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
  template['Resources']['CoreOSServerAutoScale']['Properties']['VPCZoneIdentifier'] = VPC_SUBNETS.split(',')
  template['Resources']['DeisWebELB']['Properties']['Subnets'] = VPC_SUBNETS.split(',')

print json.dumps(template)
