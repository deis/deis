
from __future__ import unicode_literals

from api.models import Node
from celery import task


@task(name='mock.build_layer')
def build_layer(layer, creds, params):
    # create security group any other infrastructure
    return


@task(name='mock.destroy_layer')
def destroy_layer(layer, creds, params):
    # delete security group and any other infrastructure
    return


@task(name='mock.launch_node')
def launch_node(node_id, creds, params, init, ssh_username, ssh_private_key):
    node = Node.objects.get(uuid=node_id)
    node.provider_id = 'i-1234567'
    node.metadata = {'state': 'running'}
    node.fqdn = 'localhost.localdomain.local'
    node.save()


@task(name='mock.terminate_node')
def terminate_node(node_id, creds, params, provider_id):
    node = Node.objects.get(uuid=node_id)
    node.metadata = {'state': 'terminated'}
    node.save()
    # delete the node itself from the database
    node.delete()


@task(name='mock.converge_node')
def converge_node(node_id, ssh_username, fqdn, ssh_private_key):
    output = ""
    rc = 0
    return output, rc


@task(name='mock.run_node')
def run_node(node_id, ssh_username, fqdn, ssh_private_key, docker_args, command):
    output = """\
total 80
drwxr-xr-x  11 matt  staff   374 Aug 14 10:57 .
drwxr-xr-x  34 matt  staff  1156 Aug 19 12:15 ..
drwxr-xr-x  14 matt  staff   476 Aug 20 09:48 .git
-rw-r--r--   1 matt  staff     5 Aug 14 10:57 .gitignore
-rw-r--r--   1 matt  staff    11 Aug 14 10:57 .ruby-version
-rw-r--r--   1 matt  staff    67 Aug 14 10:57 Gemfile
-rw-r--r--   1 matt  staff   277 Aug 14 10:57 Gemfile.lock
-rw-r--r--   1 matt  staff   553 Aug 14 10:57 LICENSE
-rw-r--r--   1 matt  staff    37 Aug 14 10:57 Procfile
-rw-r--r--   1 matt  staff  9165 Aug 14 10:57 README.md
-rw-r--r--   1 matt  staff   127 Aug 14 10:57 web.rb
"""
    rc = 0
    return output, rc
