# -*- coding: utf-8 -*-
"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import logging
import os.path
import requests
import tempfile

from django.contrib.auth.models import User
from django.test import TransactionTestCase
import etcd
import mock
from rest_framework.authtoken.models import Token
from simpleflock import SimpleFlock

import api.exceptions
from api.models import App, Config
from . import mock_status_ok


def mock_status_not_found(*args, **kwargs):
    resp = requests.Response()
    resp.status_code = 404
    resp._content_consumed = True
    return resp


def mock_request_connection_error(*args, **kwargs):
    raise requests.exceptions.ConnectionError("connection error")


class MockEtcdClient:

    def __init__(self, app):
        self.app = app

    def get(self, key, *args, **kwargs):
        node = {
            'key': '/deis/services/{}/{}_v2.web.1'.format(self.app, self.app),
            'value': '127.0.0.1:1234'
        }
        return etcd.EtcdResult(None, node)


@mock.patch('api.models.publish_release', lambda *args: None)
class ConfigTest(TransactionTestCase):

    """Tests setting and updating config values"""

    fixtures = ['tests.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.app = App.objects.all()[0]

    @mock.patch('requests.post', mock_status_ok)
    def test_config(self):
        """
        Test that config is auto-created for a new app and that
        config can be updated using a PATCH
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check to see that an initial/empty config was created
        url = "/v1/apps/{app_id}/config".format(**locals())
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('values', response.data)
        self.assertEqual(response.data['values'], {})
        config1 = response.data
        # set an initial config value
        body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        config2 = response.data
        self.assertNotEqual(config1['uuid'], config2['uuid'])
        self.assertIn('NEW_URL1', response.data['values'])
        # read the config
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        config3 = response.data
        self.assertEqual(config2, config3)
        self.assertIn('NEW_URL1', response.data['values'])
        # set an additional config value
        body = {'values': json.dumps({'NEW_URL2': 'http://localhost:8080/'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        config3 = response.data
        self.assertNotEqual(config2['uuid'], config3['uuid'])
        self.assertIn('NEW_URL1', response.data['values'])
        self.assertIn('NEW_URL2', response.data['values'])
        # read the config again
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        config4 = response.data
        self.assertEqual(config3, config4)
        self.assertIn('NEW_URL1', response.data['values'])
        self.assertIn('NEW_URL2', response.data['values'])
        # unset a config value
        body = {'values': json.dumps({'NEW_URL2': None})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        config5 = response.data
        self.assertNotEqual(config4['uuid'], config5['uuid'])
        self.assertNotIn('NEW_URL2', json.dumps(response.data['values']))
        # unset all config values
        body = {'values': json.dumps({'NEW_URL1': None})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertNotIn('NEW_URL1', json.dumps(response.data['values']))
        # disallow put/patch/delete
        response = self.client.put(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.patch(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        return config5

    @mock.patch('requests.post', mock_status_ok)
    def test_overlapping_config(self):
        """
        Test that config won't be created if a similar operation
        is in progress for that app.
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check to see that an initial/empty config was created
        url = "/v1/apps/{app_id}/config".format(**locals())
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('values', response.data)
        self.assertEqual(response.data['values'], {})
        config1 = response.data
        # create the lockfile as though a "deis config:set" were in progress
        lockfile = os.path.join(tempfile.gettempdir(), app_id + "-config")
        with SimpleFlock(lockfile):
            # set an initial config value
            body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
            response = self.client.post(url, json.dumps(body), content_type='application/json',
                                        HTTP_AUTHORIZATION='token {}'.format(self.token))
            self.assertEqual(response.status_code, 409)
            self.assertNotIn('values', response.data)
            self.assertNotIn('uuid', response.data)
        # read the config
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        config2 = response.data
        self.assertEqual(config1, config2)
        self.assertNotIn('NEW_URL1', response.data['values'])

    @mock.patch('requests.post', mock_status_ok)
    def test_response_data(self):
        """Test that the serialized response contains only relevant data."""
        body = {'id': 'test'}
        response = self.client.post('/v1/apps', json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        url = "/v1/apps/test/config"
        # set an initial config value
        body = {'values': json.dumps({'PORT': '5000'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        for key in response.data:
            self.assertIn(key, ['uuid', 'owner', 'created', 'updated', 'app', 'values', 'memory',
                                'cpu', 'tags'])
        expected = {
            'owner': self.user.username,
            'app': 'test',
            'values': {'PORT': '5000'},
            'memory': {},
            'cpu': {},
            'tags': {}
        }
        self.assertDictContainsSubset(expected, response.data)

    @mock.patch('requests.post', mock_status_ok)
    def test_response_data_types_converted(self):
        """Test that config data is converted into the correct type."""
        body = {'id': 'test'}
        response = self.client.post('/v1/apps', json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        url = "/v1/apps/test/config"

        body = {'values': json.dumps({'PORT': 5000}), 'cpu': json.dumps({'web': '1024'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        for key in response.data:
            self.assertIn(key, ['uuid', 'owner', 'created', 'updated', 'app', 'values', 'memory',
                                'cpu', 'tags'])
        expected = {
            'owner': self.user.username,
            'app': 'test',
            'values': {'PORT': '5000'},
            'memory': {},
            'cpu': {'web': 1024},
            'tags': {}
        }
        self.assertDictContainsSubset(expected, response.data)

        body = {'cpu': json.dumps({'web': 'this will fail'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        self.assertIn('CPU shares must be an integer', response.data['cpu'])

    @mock.patch('requests.post', mock_status_ok)
    def test_config_set_same_key(self):
        """
        Test that config sets on the same key function properly
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = "/v1/apps/{app_id}/config".format(**locals())
        # set an initial config value
        body = {'values': json.dumps({'PORT': '5000'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('PORT', response.data['values'])
        # reset same config value
        body = {'values': json.dumps({'PORT': '5001'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('PORT', response.data['values'])
        self.assertEqual(response.data['values']['PORT'], '5001')

    @mock.patch('requests.post', mock_status_ok)
    def test_config_set_unicode(self):
        """
        Test that config sets with unicode values are accepted.
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = "/v1/apps/{app_id}/config".format(**locals())
        # set an initial config value
        body = {'values': json.dumps({'POWERED_BY': 'Деис'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('POWERED_BY', response.data['values'])
        # reset same config value
        body = {'values': json.dumps({'POWERED_BY': 'Кроликов'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('POWERED_BY', response.data['values'])
        self.assertEqual(response.data['values']['POWERED_BY'], 'Кроликов')
        # set an integer to test unicode regression
        body = {'values': json.dumps({'INTEGER': 1})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('INTEGER', response.data['values'])
        self.assertEqual(response.data['values']['INTEGER'], '1')

    @mock.patch('requests.post', mock_status_ok)
    def test_config_str(self):
        """Test the text representation of a node."""
        config5 = self.test_config()
        config = Config.objects.get(uuid=config5['uuid'])
        self.assertEqual(str(config), "{}-{}".format(config5['app'], config5['uuid'][:7]))

    @mock.patch('requests.post', mock_status_ok)
    def test_valid_config_keys(self):
        """Test that valid config keys are accepted.
        """
        keys = ("FOO", "_foo", "f001", "FOO_BAR_BAZ_")
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = '/v1/apps/{app_id}/config'.format(**locals())
        for k in keys:
            body = {'values': json.dumps({k: "testvalue"})}
            resp = self.client.post(
                url, json.dumps(body), content_type='application/json',
                HTTP_AUTHORIZATION='token {}'.format(self.token))
            self.assertEqual(resp.status_code, 201)
            self.assertIn(k, resp.data['values'])

    @mock.patch('requests.post', mock_status_ok)
    def test_invalid_config_keys(self):
        """Test that invalid config keys are rejected.
        """
        keys = ("123", "../../foo", "FOO/", "FOO-BAR")
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = '/v1/apps/{app_id}/config'.format(**locals())
        for k in keys:
            body = {'values': json.dumps({k: "testvalue"})}
            resp = self.client.post(
                url, json.dumps(body), content_type='application/json',
                HTTP_AUTHORIZATION='token {}'.format(self.token))
            self.assertEqual(resp.status_code, 400)

    @mock.patch('requests.post', mock_status_ok)
    def test_admin_can_create_config_on_other_apps(self):
        """If a non-admin creates an app, an administrator should be able to set config
        values for that app.
        """
        user = User.objects.get(username='autotest2')
        token = Token.objects.get(user=user).key
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = "/v1/apps/{app_id}/config".format(**locals())
        # set an initial config value
        body = {'values': json.dumps({'PORT': '5000'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('PORT', response.data['values'])
        return response

    @mock.patch('requests.post', mock_status_ok)
    def test_limit_memory(self):
        """
        Test that limit is auto-created for a new app and that
        limits can be updated using a PATCH
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = '/v1/apps/{app_id}/config'.format(**locals())
        # check default limit
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('memory', response.data)
        self.assertEqual(response.data['memory'], {})
        # regression test for https://github.com/deis/deis/issues/1563
        self.assertNotIn('"', response.data['memory'])
        # set an initial limit
        mem = {'web': '1G'}
        body = {'memory': json.dumps(mem)}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        limit1 = response.data
        # check memory limits
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('memory', response.data)
        memory = response.data['memory']
        self.assertIn('web', memory)
        self.assertEqual(memory['web'], '1G')
        # set an additional value
        body = {'memory': json.dumps({'worker': '512M'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        limit2 = response.data
        self.assertNotEqual(limit1['uuid'], limit2['uuid'])
        memory = response.data['memory']
        self.assertIn('worker', memory)
        self.assertEqual(memory['worker'], '512M')
        self.assertIn('web', memory)
        self.assertEqual(memory['web'], '1G')
        # read the limit again
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        limit3 = response.data
        self.assertEqual(limit2, limit3)
        memory = response.data['memory']
        self.assertIn('worker', memory)
        self.assertEqual(memory['worker'], '512M')
        self.assertIn('web', memory)
        self.assertEqual(memory['web'], '1G')
        # regression test for https://github.com/deis/deis/issues/1613
        # ensure that config:set doesn't wipe out previous limits
        body = {'values': json.dumps({'NEW_URL2': 'http://localhost:8080/'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('NEW_URL2', response.data['values'])
        # read the limit again
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        memory = response.data['memory']
        self.assertIn('worker', memory)
        self.assertEqual(memory['worker'], '512M')
        self.assertIn('web', memory)
        self.assertEqual(memory['web'], '1G')
        # unset a value
        body = {'memory': json.dumps({'worker': None})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        limit4 = response.data
        self.assertNotEqual(limit3['uuid'], limit4['uuid'])
        self.assertNotIn('worker', json.dumps(response.data['memory']))
        # disallow put/patch/delete
        response = self.client.put(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.patch(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        return limit4

    @mock.patch('requests.post', mock_status_ok)
    def test_limit_cpu(self):
        """
        Test that CPU limits can be set
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = '/v1/apps/{app_id}/config'.format(**locals())
        # check default limit
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('cpu', response.data)
        self.assertEqual(response.data['cpu'], {})
        # regression test for https://github.com/deis/deis/issues/1563
        self.assertNotIn('"', response.data['cpu'])
        # set an initial limit
        body = {'cpu': json.dumps({'web': '1024'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        limit1 = response.data
        # check memory limits
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('cpu', response.data)
        cpu = response.data['cpu']
        self.assertIn('web', cpu)
        self.assertEqual(cpu['web'], 1024)
        # set an additional value
        body = {'cpu': json.dumps({'worker': '512'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        limit2 = response.data
        self.assertNotEqual(limit1['uuid'], limit2['uuid'])
        cpu = response.data['cpu']
        self.assertIn('worker', cpu)
        self.assertEqual(cpu['worker'], 512)
        self.assertIn('web', cpu)
        self.assertEqual(cpu['web'], 1024)
        # read the limit again
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        limit3 = response.data
        self.assertEqual(limit2, limit3)
        cpu = response.data['cpu']
        self.assertIn('worker', cpu)
        self.assertEqual(cpu['worker'], 512)
        self.assertIn('web', cpu)
        self.assertEqual(cpu['web'], 1024)
        # unset a value
        body = {'memory': json.dumps({'worker': None})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        limit4 = response.data
        self.assertNotEqual(limit3['uuid'], limit4['uuid'])
        self.assertNotIn('worker', json.dumps(response.data['memory']))
        # disallow put/patch/delete
        response = self.client.put(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.patch(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        return limit4

    @mock.patch('requests.post', mock_status_ok)
    def test_tags(self):
        """
        Test that tags can be set on an application
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = '/v1/apps/{app_id}/config'.format(**locals())
        # check default
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('tags', response.data)
        self.assertEqual(response.data['tags'], {})
        # set some tags
        body = {'tags': json.dumps({'environ': 'dev'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        tags1 = response.data
        # check tags again
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('tags', response.data)
        tags = response.data['tags']
        self.assertIn('environ', tags)
        self.assertEqual(tags['environ'], 'dev')
        # set an additional value
        body = {'tags': json.dumps({'rack': '1'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        tags2 = response.data
        self.assertNotEqual(tags1['uuid'], tags2['uuid'])
        tags = response.data['tags']
        self.assertIn('rack', tags)
        self.assertEqual(tags['rack'], '1')
        self.assertIn('environ', tags)
        self.assertEqual(tags['environ'], 'dev')
        # read the limit again
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        tags3 = response.data
        self.assertEqual(tags2, tags3)
        tags = response.data['tags']
        self.assertIn('rack', tags)
        self.assertEqual(tags['rack'], '1')
        self.assertIn('environ', tags)
        self.assertEqual(tags['environ'], 'dev')
        # unset a value
        body = {'tags': json.dumps({'rack': None})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        tags4 = response.data
        self.assertNotEqual(tags3['uuid'], tags4['uuid'])
        self.assertNotIn('rack', json.dumps(response.data['tags']))
        # set invalid values
        body = {'tags': json.dumps({'valid': 'in\nvalid'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        body = {'tags': json.dumps({'in.valid': 'valid'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        # disallow put/patch/delete
        response = self.client.put(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.patch(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)

    def test_config_owner_is_requesting_user(self):
        """
        Ensure that setting the config value is owned by the requesting user
        See https://github.com/deis/deis/issues/2650
        """
        response = self.test_admin_can_create_config_on_other_apps()
        self.assertEqual(response.data['owner'], self.user.username)

    def test_unauthorized_user_cannot_modify_config(self):
        """
        An unauthorized user should not be able to modify other config.

        Since an unauthorized user can't access the application, these
        requests should return a 403.
        """
        app_id = 'autotest'
        base_url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(base_url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        unauthorized_user = User.objects.get(username='autotest2')
        unauthorized_token = Token.objects.get(user=unauthorized_user).key
        url = '{}/{}/config'.format(base_url, app_id)
        body = {'values': {'FOO': 'bar'}}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(unauthorized_token))
        self.assertEqual(response.status_code, 403)

    def _test_app_healthcheck(self):
        # post a new build, expecting it to pass as usual
        url = "/v1/apps/{self.app}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # mock out the etcd client
        api.models._etcd_client = MockEtcdClient(self.app)
        # set an initial healthcheck url.
        url = "/v1/apps/{self.app}/config".format(**locals())
        body = {'values': json.dumps({'HEALTHCHECK_URL': '/'})}
        return self.client.post(url, json.dumps(body), content_type='application/json',
                                HTTP_AUTHORIZATION='token {}'.format(self.token))

    @mock.patch('requests.get', mock_status_ok)
    @mock.patch('time.sleep', lambda func: func)
    def test_app_healthcheck_good(self):
        """
        If a user deploys an app with a config value set for HEALTHCHECK_URL, the controller
        should check that the application responds with a 200 OK.
        """
        response = self._test_app_healthcheck()
        self.assertEqual(response.status_code, 201)
        self.assertEqual(self.app.release_set.latest().version, 3)

    @mock.patch('requests.get', mock_status_not_found)
    @mock.patch('api.models.get_etcd_client', lambda func: func)
    @mock.patch('time.sleep', lambda func: func)
    @mock.patch('api.models.logger')
    def test_app_healthcheck_bad(self, mock_logger):
        """
        If a user deploys an app with a config value set for HEALTHCHECK_URL, the controller
        should check that the application responds with a 200 OK. If it's down, the app should be
        rolled back.
        """
        response = self._test_app_healthcheck()
        self.assertEqual(response.status_code, 503)
        self.assertEqual(
            response.data,
            {'detail': 'aborting, app containers failed to respond to health check'})
        # check that only the build and initial release exist
        self.assertEqual(self.app.release_set.latest().version, 2)
        # assert that the reason why the containers failed was because
        # they failed the health check 4 times; we do this by looking
        # at logs-- there may be a better way
        exp_msg = "{}: app failed health check (got '404', expected: '200'); trying again in 0.0 \
seconds".format(self.app.id)
        exp_log_call = mock.call(logging.WARNING, exp_msg)
        log_calls = mock_logger.log.mock_calls
        self.assertEqual(log_calls.count(exp_log_call), 3)
        exp_msg = "{}: app failed health check (got '404', expected: '200')".format(self.app.id)
        exp_log_call = mock.call(logging.WARNING, exp_msg)
        self.assertEqual(log_calls.count(exp_log_call), 1)

    @mock.patch('requests.get', mock_status_not_found)
    @mock.patch('api.models.get_etcd_client', lambda func: func)
    @mock.patch('time.sleep')
    def test_app_backoff_interval(self, mock_time):
        """
        Ensure that when a healthcheck fails, a backoff strategy is used before trying again.
        """
        # post a new build, expecting it to pass as usual
        url = "/v1/apps/{self.app}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # mock out the etcd client
        api.models._etcd_client = MockEtcdClient(self.app)
        # set an initial healthcheck url.
        url = "/v1/apps/{self.app}/config".format(**locals())
        body = {'values': json.dumps({'HEALTHCHECK_URL': '/'})}
        return self.client.post(url, json.dumps(body), content_type='application/json',
                                HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(mock_time.call_count, 5)

    @mock.patch('requests.get', mock_status_ok)
    @mock.patch('time.sleep')
    def test_app_healthcheck_initial_delay(self, mock_time):
        """
        Ensure that when an initial delay is set, the request will sleep for x seconds, where
        x is the number of seconds in the initial timeout.
        """
        # post a new build, expecting it to pass as usual
        url = "/v1/apps/{self.app}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # mock out the etcd client
        api.models._etcd_client = MockEtcdClient(self.app)
        # set an initial healthcheck url.
        url = "/v1/apps/{self.app}/config".format(**locals())
        body = {'values': json.dumps({'HEALTHCHECK_URL': '/'})}
        return self.client.post(url, json.dumps(body), content_type='application/json',
                                HTTP_AUTHORIZATION='token {}'.format(self.token))
        # mock_time increments by one each time its called, so we should expect 2 calls to
        # mock_time; one for the call in the code, and one for this invocation.
        mock_time.assert_called_with(0)
        app = App.objects.all()[0]
        url = "/v1/apps/{app}/config".format(**locals())
        body = {'values': json.dumps({'HEALTHCHECK_INITIAL_DELAY': 10})}
        self.client.post(url, json.dumps(body), content_type='application/json',
                         HTTP_AUTHORIZATION='token {}'.format(self.token))
        mock_time.assert_called_with(10)

    @mock.patch('requests.get')
    @mock.patch('time.sleep', lambda func: func)
    def test_app_healthcheck_timeout(self, mock_request):
        """
        Ensure when a timeout value is set, the controller respects that value
        when making a request.
        """
        self._test_app_healthcheck()
        app = App.objects.all()[0]
        url = "/v1/apps/{app}/config".format(**locals())
        body = {'values': json.dumps({'HEALTHCHECK_TIMEOUT': 10})}
        self.client.post(url, json.dumps(body), content_type='application/json',
                         HTTP_AUTHORIZATION='token {}'.format(self.token))
        mock_request.assert_called_with('http://127.0.0.1:1234/', timeout=10)

    @mock.patch('requests.get', mock_request_connection_error)
    @mock.patch('time.sleep', lambda func: func)
    def test_app_healthcheck_connection_error(self):
        """
        If a user deploys an app with a config value set for HEALTHCHECK_URL but the app
        returns a connection error, the controller should continue checking until either the app
        responds or the app fails to respond within the timeout.

        NOTE (bacongobbler): the Docker userland proxy listens for connections and returns a
        ConnectionError, hence the unit test.
        """
        response = self._test_app_healthcheck()
        self.assertEqual(response.status_code, 503)
        self.assertEqual(
            response.data,
            {'detail': 'aborting, app containers failed to respond to health check'})
