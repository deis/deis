"""
https://github.com/coderanger/pychef

We want a simpler version for making API calls
"""
import json


class ChefAPI(object):

    def __init__(self, chef_url, client_name, client_key):
        self.server_url = chef_url
        self.client_name = client_name
        self.client_key = client_key

    def request(self, verb, path, body=''):
        assert verb in ('GET', 'DELETE', 'PUT', 'POST')
        assert path
        assert body

    def create_databag(self, name):
        body = json.dumps({'name': name, 'id': name})
        resp = self.request('POST', '/data', body)
        return resp

    def create_databag_item(self, name, item_name, item_value):
        item_dict = {'id': item_name}
        item_dict.update(item_value)
        body = json.dumps(item_dict)
        resp = self.request('POST', "/data/%s" % name, body)
        return resp

    def get_databag(self, bag_name):
        return self.request('GET', "/data/%s" % bag_name)

    def delete_databag(self, bag_name):
        return self.request('DELETE', "/data/%s" % bag_name)

    def update_databag_item(self, bag_name, item_name, item_value):
        body = json.dumps(item_value)
        return self.request('PUT', "/data/%s/%s" % (bag_name, item_name), body)

    def get_databag_item(self, bag_name, item_name):
        return self.request('GET', "/data/%s/%s" % (bag_name, item_name))

    def get_all_cookbooks(self):
        return self.request('GET', '/cookbooks')
