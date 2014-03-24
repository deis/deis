"""
Deis cloud provider implementation for Amazon EC2.
"""

from __future__ import unicode_literals

import json
import time

from boto import ec2
from boto.exception import EC2ResponseError

# from api.ssh import connect_ssh, exec_ssh
from deis import settings


# Deis-optimized EC2 amis -- with 3.8 kernel, chef 11 deps,
# and large docker images (e.g. buildstep) pre-installed
IMAGE_MAP = {
    'ap-northeast-1': 'ami-ae85ec9e',
    'ap-southeast-1': 'ami-904919c2',
    'ap-southeast-2': 'ami-a9db4393',
    'eu-west-1': 'ami-01eb1576',
    'sa-east-1': 'ami-d3cc6ece',
    'us-east-1': 'ami-51382c38',
    'us-west-1': 'ami-ec0d33a9',
    'us-west-2': 'ami-a085ec90',
}


def seed_flavors():
    """Seed the database with default flavors for each EC2 region.

    :rtype: list of dicts containing flavor data
    """
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


def build_layer(layer):
    """
    Build a layer.

    :param layer: a dict containing formation, id, params, and creds info
    """
    region = layer['params'].get('region', 'us-east-1')
    conn = _create_ec2_connection(layer['creds'], region)
    # create a new sg and authorize all ports
    # use iptables on the host to firewall ports
    name = "{formation}-{id}".format(**layer)
    sg = conn.create_security_group(name, 'Created by Deis')
    # import a new keypair using the layer key material
    conn.import_key_pair(name, layer['ssh_public_key'])
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
                raise RuntimeError('Failed to authorize security group')


def destroy_layer(layer):
    """
    Destroy a layer.

    :param layer: a dict containing formation, id, params, and creds info
    """
    region = layer['params'].get('region', 'us-east-1')
    name = "{formation}-{id}".format(**layer)
    conn = _create_ec2_connection(layer['creds'], region)
    conn.delete_key_pair(name)
    # there's an ec2 race condition on instances terminating
    # successfully but still holding a lock on the security group
    for i in range(5):
        # let's take a nap
        time.sleep(i ** 1.25)  # 1, 2.4, 3.9, 5.6, 7.4
        try:
            conn.delete_security_group(name)
            return
        except EC2ResponseError as err:
            if err.code == 'InvalidGroup.NotFound':
                return
            elif err.code in ('InvalidGroup.InUse',
                              'DependencyViolation') and i < 4:
                continue  # retry
            else:
                raise


def build_node(node):
    """
    Build a node.

    :param node: a dict containing formation, layer, params, and creds info.
    :rtype: a tuple of (provider_id, fully_qualified_domain_name, metadata)
    """
    params, creds = node['params'], node['creds']
    region = params.setdefault('region', 'us-east-1')
    conn = _create_ec2_connection(creds, region)
    name = "{formation}-{layer}".format(**node)
    params['key_name'] = name
    sg = conn.get_all_security_groups(name)[0]
    params.setdefault('security_groups', []).append(sg.name)
    image_id = params.get(
        'image', getattr(settings, 'IMAGE_MAP', IMAGE_MAP)[region])
    images = conn.get_all_images([image_id])
    if len(images) != 1:
        raise LookupError('Could not find AMI: %s' % image_id)
    image = images[0]
    kwargs = _prepare_run_kwargs(params)
    reservation = image.run(**kwargs)
    instances = reservation.instances
    boto = instances[0]
    # sleep before tagging
    time.sleep(10)
    boto.update()
    boto.add_tag('Name', node['id'])
    # loop until running
    while(True):
        time.sleep(2)
        boto.update()
        if boto.state == 'running':
            break
    # prepare return values
    provider_id = boto.id
    fqdn = boto.public_dns_name
    metadata = _format_metadata(boto)
    return provider_id, fqdn, metadata


def destroy_node(node):
    """
    Destroy a node.

    :param node: a dict containing a node's provider_id, params, and creds
    """
    provider_id = node['provider_id']
    region = node['params'].get('region', 'us-east-1')
    conn = _create_ec2_connection(node['creds'], region)
    if provider_id:
        try:
            conn.terminate_instances([provider_id])
            i = conn.get_all_instances([provider_id])[0].instances[0]
            while(True):
                time.sleep(2)
                i.update()
                if i.state == "terminated":
                    break
        except EC2ResponseError as e:
            if e.code not in ('InvalidInstanceID.NotFound',):
                raise


def _create_ec2_connection(creds, region):
    """
    Connect to an EC2 region with the given credentials.

    :param creds: a dict containing an EC2 access_key and secret_key
    :region: the name of an EC2 region, such as "us-west-2"
    :rtype: a connected :class:`~boto.ec2.connection.EC2Connection`
    :raises EnvironmentError: if no credentials are provided
    """
    if not creds:
        raise EnvironmentError('No credentials provided')
    return ec2.connect_to_region(region,
                                 aws_access_key_id=creds['access_key'],
                                 aws_secret_access_key=creds['secret_key'])


def _prepare_run_kwargs(params):
    # start with sane defaults
    kwargs = {
        'min_count': 1, 'max_count': 1,
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
        'security_groups': params['security_groups'],
        'placement': requested_zone,
        'key_name': params['key_name'],
        'kernel_id': params.get('kernel', None),
    }
    # add user_data if provided in params
    user_data = params.get('user_data')
    if user_data:
        kwargs.update({'user_data': user_data})
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
