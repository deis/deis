
from __future__ import unicode_literals
import subprocess

from celery import task

from cm.chef_api import ChefAPI
from deis import settings


@task(name='controller.update_gitosis')
def update_gitosis(databag_item_value):
    # update the data bag
    client = ChefAPI(settings.CHEF_SERVER_URL,
                     settings.CHEF_CLIENT_NAME,
                     settings.CHEF_CLIENT_KEY)
    client.update_databag_item('deis-build', 'gitosis', databag_item_value)
    # call a chef update
    subprocess.check_call(
        ['sudo', 'chef-client', '--override-runlist', 'recipe[deis::gitosis]'])


@task(name='controller.update_formation')
def update_formation(formation_id, databag_item_value):
    # update the data bag
    client = ChefAPI(settings.CHEF_SERVER_URL,
                     settings.CHEF_CLIENT_NAME,
                     settings.CHEF_CLIENT_KEY)
    # TODO: move this logic into the chef API
    resp, code = client.update_databag_item(
        'deis-formations', formation_id, databag_item_value)
    if code == 200:
        return resp, code
    elif code == 404:
        resp, code = client.create_databag_item(
            'deis-formations', formation_id, databag_item_value)
        if code != 201:
            msg = 'Failed to create data bag: {code} => {resp}'
            raise RuntimeError(msg.format(**locals()))
    else:
        msg = 'Failed to update data bag: {code} => {resp}'
        raise RuntimeError(msg.format(**locals()))


@task(name='controller.update_application')
def update_application(app_id, databag_item_value):
    # update the data bag
    client = ChefAPI(settings.CHEF_SERVER_URL,
                     settings.CHEF_CLIENT_NAME,
                     settings.CHEF_CLIENT_KEY)
    # TODO: move this logic into the chef API
    resp, code = client.update_databag_item(
        'deis-apps', app_id, databag_item_value)
    if code == 200:
        return resp, code
    elif code == 404:
        resp, code = client.create_databag_item(
            'deis-apps', app_id, databag_item_value)
        if code != 201:
            msg = 'Failed to create data bag: {code} => {resp}'
            raise RuntimeError(msg.format(**locals()))
    else:
        msg = 'Failed to update data bag: {code} => {resp}'
        raise RuntimeError(msg.format(**locals()))


@task(name='controller.destroy_formation')
def destroy_formation(formation_id):
    client = ChefAPI(settings.CHEF_SERVER_URL,
                     settings.CHEF_CLIENT_NAME,
                     settings.CHEF_CLIENT_KEY)
    _resp, _code = client.delete_databag_item('deis-formations', formation_id)
    subprocess.check_call(['sudo', 'chef-client', '--override-runlist', 'recipe[deis::gitosis]'])
