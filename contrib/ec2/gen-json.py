#!/usr/bin/env python
import json
import os

template = json.load(open("deis.template",'r'))

with open('../coreos/user-data','r') as f:
  lines = f.readlines()

template['Resources']['CoreOSServerLaunchConfig']['Properties']['UserData']['Fn::Base64']['Fn::Join'] = [ '', lines ]
template['Parameters']['ClusterSize']['Default'] = str(os.getenv('DEIS_NUM_INSTANCES', 3))

if os.getenv("VPC_ID", None) and os.getenv("VPC_SUBNETS", None):
  for resource in template['Resources'].keys():
    resource_type = template['Resources'][resource]['Type']
    if resource_type == 'AWS::EC2::SecurityGroup':
      template['Resources'][resource]['Properties']['VpcId'] = os.getenv("VPC_ID")
    elif resource_type == 'AWS::EC2::SecurityGroupIngress':
      template['Resources'][resource]['Properties']['GroupId'] = template['Resources'][resource]['Properties']['GroupName']
      del template['Resources'][resource]['Properties']['GroupName']
      template['Resources'][resource]['Properties']['SourceSecurityGroupId'] = {
        'Ref': template['Resources'][resource]['Properties']['SourceSecurityGroupId']['Fn::GetAtt'][0]
      }
    elif resource_type == 'AWS::AutoScaling::LaunchConfiguration':
      template['Resources'][resource]['Properties']['AssociatePublicIpAddress'] = False
    elif resource_type == 'AWS::AutoScaling::AutoScalingGroup':
      template['Resources'][resource]['Properties']['VPCZoneIdentifier'] = os.getenv('VPC_SUBNETS').split(',')


print json.dumps(template)
