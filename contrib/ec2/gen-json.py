#!/usr/bin/env python
import json
import os

template = json.load(open("deis.template",'r'))

with open('../coreos/user-data','r') as f:
  lines = f.readlines()

template['Resources']['CoreOSServerLaunchConfig']['Properties']['UserData']['Fn::Base64']['Fn::Join'] = [ '', lines ]
template['Parameters']['ClusterSize']['Default'] = str(os.getenv('DEIS_NUM_INSTANCES', 3))

print json.dumps(template, indent=4, separators=(',', ': '))
