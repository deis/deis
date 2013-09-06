
from __future__ import unicode_literals

import json
import time
import yaml

from boto import ec2
from boto.exception import EC2ResponseError

from api.ssh import connect_ssh, exec_ssh
from deis import settings


# Deis-optimized EC2 amis -- with 3.8 kernel, chef 11 deps,
# and large docker images (e.g. buildstep) pre-installed
IMAGE_MAP = {
    'ap-northeast-1': 'ami-6da8356c',
    'ap-southeast-1': 'ami-a66f24f4',
    'ap-southeast-2': 'ami-d5f66bef',
    'eu-west-1': 'ami-acbf5adb',
    'sa-east-1': 'ami-f9fd5ae4',
    'us-east-1': 'ami-69f3bc00',
    'us-west-1': 'ami-f0695cb5',
    'us-west-2': 'ami-ea1e82da',
}


def build_layer(layer):
    region = layer.flavor.params.get('region', 'us-east-1')
    conn = _create_ec2_connection(layer.flavor.provider.creds, region)
    # create a new sg and authorize all ports
    # use iptables on the host to firewall ports
    sg_name = "{}-{}".format(layer.formation.id, layer.id)
    sg = conn.create_security_group(sg_name, 'Created by Deis')
    # loop until the sg is *actually* there
    for i in xrange(10):
        try:
            sg.authorize(ip_protocol='tcp', from_port=1, to_port=65535,
                         cidr_ip='0.0.0.0/0')
            break
        except EC2ResponseError:
            if i < 10:
                time.sleep(1.5)
                continue
            else:
                raise
    return layer


def destroy_layer(layer):
    # there's an ec2 race condition on instances terminating
    # successfully but still holding a lock on the security group
    # let's take a nap
    time.sleep(5)
    region = layer.flavor.params.get('region', 'us-east-1')
    sg_name = "{}-{}".format(layer.formation.id, layer.id)
    conn = _create_ec2_connection(layer.flavor.provider.creds, region)
    try:
        conn.delete_security_group(sg_name)
    except EC2ResponseError as e:
        if e.code != 'InvalidGroup.NotFound':
            raise e
    return layer


def build_node(node, config):
    creds = node.layer.flavor.provider.creds.copy()
    params = node.layer.flavor.params.copy()
    # connect to ec2
    region = params.get('region', 'us-east-1')
    conn = _create_ec2_connection(creds, region)
    # use the layer name as the security group name
    sg_name = "{}-{}".format(node.formation.id, node.layer.id)
    sg = conn.get_all_security_groups(sg_name)[0]
    params.setdefault('security_groups', []).append(sg.name)
    # retrieve the ami for launching this node
    image_id = params.get(
        'image', getattr(settings, 'IMAGE_MAP', IMAGE_MAP)[region])
    images = conn.get_all_images([image_id])
    if len(images) != 1:
        raise LookupError('Could not find AMI: %s' % image_id)
    image = images[0]
    kwargs = _prepare_run_kwargs(params, config)
    reservation = image.run(**kwargs)
    instances = reservation.instances
    boto = instances[0]
    # sleep before tagging
    time.sleep(10)
    boto.update()
    boto.add_tag('Name', node.id)
    # loop until running
    while(True):
        time.sleep(2)
        boto.update()
        if boto.state == 'running':
            break
    # save node updates
    node.provider_id = boto.id
    node.fqdn = boto.public_dns_name
    node.metadata = _format_metadata(boto)
    node.save()
    # loop until cloud-init is finished
    ssh = connect_ssh(node.layer.ssh_username,
                      boto.public_dns_name, 22,
                      node.layer.ssh_private_key,
                      timeout=120)
    initializing = True
    while initializing:
        time.sleep(10)
        initializing, _rc = exec_ssh(
            ssh, 'ps auxw | egrep "cloud-init" | grep -v egrep')
    return node


def destroy_node(node):
    region = node.layer.flavor.params.get('region', 'us-east-1')
    conn = _create_ec2_connection(node.layer.flavor.provider.creds, region)
    if node.provider_id:
        conn.terminate_instances([node.provider_id])
        i = conn.get_all_instances([node.provider_id])[0].instances[0]
        while(True):
            time.sleep(2)
            i.update()
            if i.state == "terminated":
                break
    return node


def seed_flavors(user):
    """Seed the database with default Flavors for each EC2 region."""
    flavors = []
    for r in ('us-east-1', 'us-west-1', 'us-west-2', 'eu-west-1',
              'ap-northeast-1', 'ap-southeast-1', 'ap-southeast-2',
              'sa-east-1'):
        flavors.append({'id': 'ec2-{}'.format(r),
                        'provider': 'ec2',
                        'params': json.dumps({
                            'region': r,
                            'image': IMAGE_MAP[r],
                            'zone': 'any',
                            'size': 'm1.medium'})})
    return flavors


def _create_ec2_connection(creds, region):
    if not creds:
        raise EnvironmentError('No credentials provided')
    return ec2.connect_to_region(region,
                                 aws_access_key_id=creds['access_key'],
                                 aws_secret_access_key=creds['secret_key'])


def _prepare_run_kwargs(params, init):
    # start with sane defaults
    kwargs = {
        'min_count': 1, 'max_count': 1,
        'key_name': None,
        'user_data': None, 'addressing_type': None,
        'instance_type': None, 'placement': None,
        'kernel_id': None, 'ramdisk_id': None,
        'monitoring_enabled': False, 'subnet_id': None,
        'block_device_map': None,
    }
    # convert zone "any" to NoneType
    requested_zone = params.get('zone')
    if requested_zone and requested_zone.lower() == 'any':
        requested_zone = None
    # lookup kwargs from params
    param_kwargs = {
        'instance_type': params.get('size', 'm1.medium'),
        'key_name': params.get('key_name', None),
        'security_groups': params['security_groups'],
        'placement': requested_zone,
        'kernel_id': params.get('kernel', None),
    }
    # update user_data
    cloud_config = '#cloud-config\n' + yaml.safe_dump(init)
    kwargs.update({'user_data': cloud_config})
    # params override defaults
    kwargs.update(param_kwargs)
    return kwargs


def _format_metadata(boto):
    return {
        'architecture': boto.architecture,
        'block_device_mapping': {
            k: v.volume_id for k, v in boto.block_device_mapping.items()
        },
        'client_token': boto.client_token,
        'dns_name': boto.dns_name,
        'ebs_optimized': boto.ebs_optimized,
        'eventsSet': boto.eventsSet,
        'group_name': boto.group_name,
        'groups': [g.id for g in boto.groups],
        'hypervisor': boto.hypervisor,
        'id': boto.id,
        'image_id': boto.image_id,
        'instance_profile': boto.instance_profile,
        'instance_type': boto.instance_type,
        'interfaces': list(boto.interfaces),
        'ip_address': boto.ip_address,
        'kernel': boto.kernel,
        'key_name': boto.key_name,
        'launch_time': boto.launch_time,
        'monitored': boto.monitored,
        'monitoring_state': boto.monitoring_state,
        'persistent': boto.persistent,
        'placement': boto.placement,
        'placement_group': boto.placement_group,
        'placement_tenancy': boto.placement_tenancy,
        'previous_state': boto.previous_state,
        'private_dns_name': boto.private_dns_name,
        'private_ip_address': boto.private_ip_address,
        'public_dns_name': boto.public_dns_name,
        'ramdisk': boto.ramdisk,
        'region': boto.region.name,
        'root_device_name': boto.root_device_name,
        'root_device_type': boto.root_device_type,
        'spot_instance_request_id': boto.spot_instance_request_id,
        'state': boto.state,
        'state_code': boto.state_code,
        'state_reason': boto.state_reason,
        'subnet_id': boto.subnet_id,
        'tags': dict(boto.tags),
        'virtualization_type': boto.virtualization_type,
        'vpc_id': boto.vpc_id,
    }
