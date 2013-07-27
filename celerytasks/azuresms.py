
from __future__ import unicode_literals

from azure.servicemanagement import ServiceManagementService
from azure.servicemanagement import LinuxConfigurationSet, OSVirtualHardDisk
from celery import task
import yaml

from . import util


@task(name='azuresms.launch_node')
def launch_node(node_id, creds, params, init, ssh_username, ssh_private_key):
    # "pip install azure"
    sms = ServiceManagementService(
        subscription_id='69581868-8a08-4d98-a5b0-1d111c616fc3',
        cert_file='/Users/dgriffin/certs/iOSWAToolkit.pem')
    for i in sms.list_os_images():
        print 'I is ', i.name, ' -- ', i.label, ' -- ', i.location, ' -- ', i.media_link
    media_link = \
        'http://opdemandstorage.blob.core.windows.net/communityimages/' + \
        'b39f27a8b8c64d52b05eac6a62ebad85__Ubuntu_DAILY_BUILD-' + \
        'precise-12_04_2-LTS-amd64-server-20130702-en-us-30GB.vhd'
    config = LinuxConfigurationSet(user_name="ubuntu", user_password="opdemand")
    hard_disk = OSVirtualHardDisk(
        'b39f27a8b8c64d52b05eac6a62ebad85__Ubuntu_DAILY_BUILD-' +
        'precise-12_04_2-LTS-amd64-server-20130702-en-us-30GB',
        media_link, disk_label='opdemandservice')
    ret = sms.create_virtual_machine_deployment(
        'opdemandservice', 'deploy1', 'production', 'opdemandservice2',
        'opdemandservice3', config, hard_disk)
       # service_name, deployment_name, deployment_slot, label, role_name
       # system_config, os_virtual_hard_disk
    print 'Ret ', ret
    return sms


@task(name='azuresms.terminate_node')
def terminate_node(node_id, creds, params, provider_id):
    pass


@task(name='azuresms.converge_node')
def converge_node(node_id, ssh_username, fqdn, ssh_private_key,
                  command='sudo chef-client'):
    ssh = util.connect_ssh(ssh_username, fqdn, 22, ssh_private_key)
    output, rc = util.exec_ssh(ssh, command)
    return output, rc


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
    cloud_config = '#cloud-config\n' + yaml.safe_dump(init)
    kwargs.update({'user_data': cloud_config})
    # params override defaults
    kwargs.update(param_kwargs)
    return kwargs


def format_metadata(boto):
    return {
        'architecture': boto.architecture,
        'block_device_mapping': {
            k: v.volume_id for k, v in boto.block_device_mapping.items()},
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


if __name__ == "__main__":
    print 'Checking '
    l = launch_node(None, None, None, None, None, None)
