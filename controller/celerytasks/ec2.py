
from __future__ import unicode_literals
import json
import time

from boto import ec2
from boto.exception import EC2ResponseError
from celery import task
import yaml

from . import util
from api.models import Node
from deis import settings
from celerytasks.chef import ChefAPI

EC2_IMAGE_MAP = {
    'ap-northeast-1': 'ami-51129850',
    'ap-southeast-1': 'ami-a02f66f2',
    'ap-southeast-2': 'ami-974ddead',
    'eu-west-1': 'ami-89b1a3fd',
    'sa-east-1': 'ami-5c7edb41',
    'us-east-1': 'ami-23d9a94a',
    'us-west-1': 'ami-c4072e81',
    # deis optimized ami, with 3.8 kernel and chef 11 deps pre-installed
    'us-west-2': 'ami-bf41d28f',  # 'ami-fb68f8cb',
}


@task(name='ec2.prepare_formation')
def prepare_formation(formation, creds, params):
    region = params.get('region', 'us-east-1')
    conn = create_ec2_connection(
        region, creds['access_key'], creds['secret_key'])
    # create a new sg and authorize all ports
    # use iptables on the host to firewall ports
    sg = conn.create_security_group(formation, 'Managed by Deis')
    sg.authorize(ip_protocol='tcp', from_port=1, to_port=65535,
                 cidr_ip='0.0.0.0/0')


@task(name='ec2.cleanup_formation')
def cleanup_formation(formation, creds, params):
    region = params.get('region', 'us-east-1')
    conn = create_ec2_connection(
        region, creds['access_key'], creds['secret_key'])
    sg_name = formation
    try:
        conn.delete_security_group(sg_name)
    except EC2ResponseError as e:
        if e.code != 'InvalidGroup.NotFound':
            raise e


@task(name='ec2.launch_node')
def launch_node(node_id, creds, params, init, ssh_username, ssh_private_key):
    region = params.get('region', 'us-east-1')
    conn = create_ec2_connection(
        region, creds['access_key'], creds['secret_key'])
    # find or create the security group for this formation
    sg_name = params['formation']
    sg = conn.get_all_security_groups(sg_name)[0]
    # add the security group to the list
    params.setdefault('security_groups', []).append(sg.name)
    # retrieve the ami for launching this node
    image_id = params.get(
        'image', getattr(settings, 'EC2_IMAGE_MAP', EC2_IMAGE_MAP)[region])
    images = conn.get_all_images([image_id])
    if len(images) != 1:
        raise LookupError('Could not find AMI: %s' % image_id)
    image = images[0]
    kwargs = prepare_run_kwargs(params, init)
    reservation = image.run(**kwargs)
    instances = reservation.instances
    boto = instances[0]
    # initial sleep
    time.sleep(10)
    boto.update()
    # try adding a tag
    boto.add_tag('Name', params['id'])
    # loop until running
    while(True):
        time.sleep(2)
        boto.update()
        if boto.state == 'running':
            break
    # update the node
    node = Node.objects.get(uuid=node_id)
    node.provider_id = boto.id
    node.fqdn = boto.public_dns_name
    node.metadata = format_metadata(boto)
    node.save()
    # loop until cloud-init is finished
    ssh = util.connect_ssh(ssh_username, boto.public_dns_name, 22,
                           ssh_private_key)
    initializing = True
    while initializing:
        time.sleep(10)
        initializing, _rc = util.exec_ssh(
            ssh, 'ps auxw | egrep "cloud-init" | grep -v egrep')
    # loop until node is registered with chef
    # if chef bootstrapping fails, the node will not complete registration
    if settings.CHEF_ENABLED:
        registered = False
        while not registered:
            # reinstatiate the client on each poll attempt
            # to avoid disconnect errors
            client = ChefAPI(settings.CHEF_SERVER_URL,
                             settings.CHEF_CLIENT_NAME,
                             settings.CHEF_CLIENT_KEY)
            resp, status = client.get_node(node.id)
            if status == 200:
                body = json.loads(resp)
                # wait until idletime is not null
                # meaning the node is registered
                if body.get('automatic', {}).get('idletime'):
                    break
            time.sleep(5)


@task(name='ec2.terminate_node')
def terminate_node(node_id, creds, params, provider_id):
    region = params.get('region', 'us-east-1')
    conn = create_ec2_connection(
        region, creds['access_key'], creds['secret_key'])
    conn.terminate_instances([provider_id])
    i = conn.get_all_instances([provider_id])[0].instances[0]
    while(True):
        time.sleep(2)
        i.update()
        if i.state == "terminated":
            break
    # pull the node from the database
    node = Node.objects.get(uuid=node_id)
    chef_id = node.id
    node.provider_id = None
    node.fqdn = None
    node.metadata = {}
    node.save()
    # delete the node itself from the database
    node.delete()
    # purge the node & client records from chef server
    client = ChefAPI(settings.CHEF_SERVER_URL,
                     settings.CHEF_CLIENT_NAME,
                     settings.CHEF_CLIENT_KEY)
    client.delete_node(chef_id)
    client.delete_client(chef_id)


@task(name='ec2.converge_node')
def converge_node(node_id, ssh_username, fqdn, ssh_private_key,
                  command='sudo chef-client'):
    ssh = util.connect_ssh(ssh_username, fqdn, 22, ssh_private_key)
    output, rc = util.exec_ssh(ssh, command)
    return output, rc


# utility functions

def create_ec2_connection(region, access_key, secret_key):
    return ec2.connect_to_region(region, aws_access_key_id=access_key,
                                 aws_secret_access_key=secret_key)


def prepare_run_kwargs(params, init):
    # start with sane defaults
    kwargs = {
        'min_count': 1, 'max_count': 1,
        'key_name': None,
        'user_data': None, 'addressing_type': None,
        'instance_type': 'm1.small', 'placement': None,
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
        'instance_type': params.get('size', 'm1.small'),
        'key_name': params.get('key_name', None),
        'security_groups': params['security_groups'],
        'placement': requested_zone,
        'kernel_id': params.get('kernel', None),
    }
    # update user_data
    cloud_config = '#cloud-config\n'+yaml.safe_dump(init)
    kwargs.update({'user_data': cloud_config})
    # params override defaults
    kwargs.update(param_kwargs)
    return kwargs


def format_metadata(boto):
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
