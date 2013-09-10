"""
Much of this has been copied from pyChef.
https://github.com/coderanger/pychef

We want a simpler version for making API calls
"""

import base64
import datetime
import hashlib
import httplib
import json
import re
import time
import urlparse

from chef_rsa import Key


def ruby_b64encode(value):
    """The Ruby function Base64.encode64 automatically breaks things up
    into 60-character chunks.
    """
    b64 = base64.b64encode(value)
    for i in xrange(0, len(b64), 60):
        yield b64[i:i + 60]


class UTC(datetime.tzinfo):
    """UTC timezone stub."""

    ZERO = datetime.timedelta(0)

    def utcoffset(self, dt):
        return self.ZERO

    def tzname(self, dt):
        return 'UTC'

    def dst(self, dt):
        return self.ZERO


utc = UTC()


def canonical_time(timestamp):
    if timestamp.tzinfo is not None:
        timestamp = timestamp.astimezone(utc).replace(tzinfo=None)
    return timestamp.replace(microsecond=0).isoformat() + 'Z'


canonical_path_regex = re.compile(r'/+')


def canonical_path(path):
    path = canonical_path_regex.sub('/', path)
    if len(path) > 1:
        path = path.rstrip('/')
    return path


def canonical_request(http_method, path, hashed_body, timestamp, user_id):
    # Canonicalize request parameters
    http_method = http_method.upper()
    path = canonical_path(path)
    if isinstance(timestamp, datetime.datetime):
        timestamp = canonical_time(timestamp)
    hashed_path = sha1_base64(path)
    return """\
Method:{}
Hashed Path:{}
X-Ops-Content-Hash:{}
X-Ops-Timestamp:{}
X-Ops-UserId:{}""".format(http_method, hashed_path, hashed_body, timestamp,
                          user_id)


def sha1_base64(value):
    return '\n'.join(ruby_b64encode(hashlib.sha1(value).digest()))


def create_authorization(blank_headers, verb, url, priv_key, user, body=''):
    headers = blank_headers.copy()
    rsa_key = Key(fp=priv_key)
    timestamp = canonical_time(datetime.datetime.utcnow())
    hashed_body = sha1_base64(body)

    canon = canonical_request(verb, url, hashed_body, timestamp, user)
    b64_priv = ruby_b64encode(rsa_key.private_encrypt(canon))

    for i, line in enumerate(b64_priv):
        headers['X-Ops-Authorization-' + str(i + 1)] = line

    headers['X-Ops-Timestamp'] = timestamp
    headers['X-Ops-Content-Hash'] = hashed_body
    headers['X-Ops-UserId'] = user
    return headers


class ChefAPI(object):

    headers = {
        'Accept': 'application/json',
        'X-Chef-Version': '11.0.4.x',
        'X-Ops-Sign': 'version=1.0',
        'Content-Type': 'application/json'
    }

    def __init__(self, server_url, client_name, client_key):
        self.server_url = server_url
        self.client_name = client_name
        self.client_key = client_key
        self.hostname = urlparse.urlsplit(self.server_url).netloc
        self.path = urlparse.urlsplit(self.server_url).path
        self.headers.update({'Host': self.hostname})
        self.conn = httplib.HTTPSConnection(self.hostname)
        self.conn.connect()

    def request(self, verb, path, body='', attempts=5, interval=5):
        url = self.path + path
        headers = create_authorization(
            self.headers, verb, url, self.client_key, self.client_name, body)
        # retry all chef api requests
        for _ in range(attempts):
            self.conn.request(verb, url, body=body, headers=headers)
            resp = self.conn.getresponse()
            if resp.status != 500:
                break
            time.sleep(interval)
        else:
            errmsg = 'Chef API requests failed: {}'.format(path)
            raise RuntimeError(errmsg)
        return resp.read(), resp.status

    def create_databag(self, name):
        body = json.dumps({'name': name, 'id': name})
        resp = self.request('POST', '/data', body)
        return resp

    def create_databag_item(self, name, item_name, item_value):
        item_dict = {'id': item_name}
        item_dict.update(item_value)
        body = json.dumps(item_dict)
        resp = self.request('POST', '/data/%s' % name, body)
        return resp

    def get_databag(self, bag_name):
        return self.request('GET', '/data/%s' % bag_name)

    def delete_databag(self, bag_name):
        return self.request('DELETE', '/data/%s' % bag_name)

    def delete_databag_item(self, bag_name, item_name):
        return self.request('DELETE', '/data/%s/%s' % (bag_name, item_name))

    def update_databag_item(self, bag_name, item_name, item_value):
        body = json.dumps(item_value)
        return self.request('PUT', '/data/%s/%s' % (bag_name, item_name), body)

    def get_databag_item(self, bag_name, item_name):
        return self.request('GET', '/data/%s/%s' % (bag_name, item_name))

    def get_all_cookbooks(self):
        return self.request('GET', '/cookbooks')

    def get_node(self, node_id):
        return self.request('GET', '/nodes/%s' % node_id)

    def delete_node(self, node_id):
        return self.request('DELETE', '/nodes/%s' % node_id)

    def delete_client(self, client_id):
        return self.request('DELETE', '/clients/%s' % client_id)

#     def create_cookbook(self, cookbook_name, cookbooks, priv_key, user, org):
#         checksums = {}
#         by_cb = {}
#         first = None
#         for c in cookbooks:
#             json_cb = json.dumps(c)
#             first = json_cb
#             hasher = hashlib.md5()
#             hasher.update(json_cb)
#             check = hasher.hexdigest()
#             checksums[check] = None
#             by_cb[c['name']] = check
#         body = json.dumps({'checksums': checksums})
#         sandbox = json.loads(self.request('POST', '/sandboxes'))
#         print 'Sandbox is ', sandbox
#         for k, v in sandbox['checksums'].items():
#             print 'URL ', v
#             if 'url' in v:
#                print 'Trigger it ', self.request(
#                    'PUT', v['url'][25:], json_cb, priv_key, user)
#
#        print 'Mark as uploaded ', self.request(
#            'PUT', sandbox['uri'][25:], '''{'is_completed':true}''', priv_key,
#            user)
#        print 'Mark as uploaded ', self.request(
#            'PUT', sandbox['uri'][25:], '''{'is_completed':true}''', priv_key,
#            user)
#        print 'Mark as uploaded ', self.request(
#            'PUT', sandbox['uri'][25:], '''{'is_completed':true}''', priv_key,
#            user)
#        print 'Mark as uploaded ', self.request(
#            'PUT', sandbox['uri'][25:], '''{'is_completed':true}''', priv_key,
#            user)
#
#         for c in cookbooks:
#             c['definitions'] = [{
#                 'name': 'unicorn_config.rb',
#                 'checksum': by_cb[c['name']],
#                 'path': 'definitions/unicorn_config.rb',
#                 'specificity': 'default'
#             }],
#             return self.request('PUT', '/organizations/%s/cookbooks/%s/1' %
#                                 (org, cookbook_name), body, priv_key, user)
#
# @task(name='chef.update_data_bag_item')
# def update_data_bag_item(conn_info, bag_name, item_name, item_value):
#     client = ChefAPI(conn_info['server_url'],
#                      conn_info['client_name'],
#                      conn_info['client_key'],
#                      conn_info['organization'])
#     client.update_databag_item(bag_name, item_name, item_value)
