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
  ExecStart=/usr/sbin/mkfs.ext4 -i 4096 -b 4096 /dev/xvdb
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
  ExecStart=/usr/bin/chown -R etcd:etcd /media/ephemeral/etcd
'''
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

new_units = [
  dict({'name': 'format-ephemeral-volume.service', 'command': 'start', 'content': FORMAT_EPHEMERAL_VOLUME}),
  dict({'name': 'media-ephemeral.mount', 'command': 'start', 'content': MOUNT_EPHEMERAL_VOLUME}),
  dict({'name': 'prepare-etcd-data-directory.service', 'command': 'start', 'content': PREPARE_ETCD_DATA_DIRECTORY}),
  dict({'name': 'format-docker-volume.service', 'command': 'start', 'content': FORMAT_DOCKER_VOLUME}),
  dict({'name': 'var-lib-docker.mount', 'command': 'start', 'content': MOUNT_DOCKER_VOLUME})
]

data = yaml.load(file(os.path.join(CURR_DIR, '..', 'coreos', 'user-data'), 'r'))

# coreos-cloudinit will start the units in order, so we want these to be processed before etcd/fleet
# are started
data['coreos']['units'] = new_units + data['coreos']['units']

# configure etcd to use the ephemeral drive
data['coreos']['etcd']['data-dir'] = '/media/ephemeral/etcd'

header = "#cloud-config\n---"
dump = yaml.dump(data, default_flow_style=False)
print header
print dump

