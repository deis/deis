"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import mock
import requests

from django.test import TransactionTestCase
from django.test.utils import override_settings

from api.models import Limit


def mock_import_repository_task(*args, **kwargs):
    resp = requests.Response()
    resp.status_code = 200
    resp._content_consumed = True
    return resp


@override_settings(CELERY_ALWAYS_EAGER=True)
class LimitTest(TransactionTestCase):

    """Tests setting and updating resource limits"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        body = {'id': 'autotest', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': {}}
        response = self.client.post('/api/clusters', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)

    @mock.patch('requests.post', mock_import_repository_task)
    def test_limit_memory(self):
        """
        Test that limit is auto-created for a new app and that
        limits can be updated using a PATCH
        """
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = '/api/apps/{app_id}/limits'.format(**locals())
        # check default limit
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertIn('memory', response.data)
        self.assertEqual(json.loads(response.data['memory']), json.dumps({}))
        # set an initial limit
        mem = {'web': '1G'}
        body = {'memory': json.dumps(mem)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertIn('x-deis-release', response._headers)
        limit1 = response.data
        # check memory limits
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertIn('memory', response.data)
        memory = json.loads(response.data['memory'])
        self.assertIn('web', memory)
        self.assertEqual(memory['web'], '1G')
        # set an additional value
        body = {'memory': json.dumps({'worker': '512M'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        limit2 = response.data
        self.assertNotEqual(limit1['uuid'], limit2['uuid'])
        memory = json.loads(response.data['memory'])
        self.assertIn('worker', memory)
        self.assertEqual(memory['worker'], '512M')
        self.assertIn('web', memory)
        self.assertEqual(memory['web'], '1G')
        # read the limit again
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        limit3 = response.data
        self.assertEqual(limit2, limit3)
        memory = json.loads(response.data['memory'])
        self.assertIn('worker', memory)
        self.assertEqual(memory['worker'], '512M')
        self.assertIn('web', memory)
        self.assertEqual(memory['web'], '1G')
        # unset a value
        body = {'memory': json.dumps({'worker': None})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        limit4 = response.data
        self.assertNotEqual(limit3['uuid'], limit4['uuid'])
        self.assertNotIn('worker', json.dumps(response.data['memory']))
        # disallow put/patch/delete
        self.assertEqual(self.client.put(url).status_code, 405)
        self.assertEqual(self.client.patch(url).status_code, 405)
        self.assertEqual(self.client.delete(url).status_code, 405)
        return limit4

    @mock.patch('requests.post', mock_import_repository_task)
    def test_limit_cpu(self):
        """
        Test that CPU limits can be set
        """
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = '/api/apps/{app_id}/limits'.format(**locals())
        # check default limit
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertIn('cpu', response.data)
        self.assertEqual(json.loads(response.data['memory']), json.dumps({}))
        # set an initial limit
        body = {'cpu': json.dumps({'web': '1024'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertIn('x-deis-release', response._headers)
        limit1 = response.data
        # check memory limits
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertIn('cpu', response.data)
        cpu = json.loads(response.data['cpu'])
        self.assertIn('web', cpu)
        self.assertEqual(cpu['web'], '1024')
        # set an additional value
        body = {'cpu': json.dumps({'worker': '512'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        limit2 = response.data
        self.assertNotEqual(limit1['uuid'], limit2['uuid'])
        cpu = json.loads(response.data['cpu'])
        self.assertIn('worker', cpu)
        self.assertEqual(cpu['worker'], '512')
        self.assertIn('web', cpu)
        self.assertEqual(cpu['web'], '1024')
        # read the limit again
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        limit3 = response.data
        self.assertEqual(limit2, limit3)
        cpu = json.loads(response.data['cpu'])
        self.assertIn('worker', cpu)
        self.assertEqual(cpu['worker'], '512')
        self.assertIn('web', cpu)
        self.assertEqual(cpu['web'], '1024')
        # unset a value
        body = {'memory': json.dumps({'worker': None})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        limit4 = response.data
        self.assertNotEqual(limit3['uuid'], limit4['uuid'])
        self.assertNotIn('worker', json.dumps(response.data['memory']))
        # disallow put/patch/delete
        self.assertEqual(self.client.put(url).status_code, 405)
        self.assertEqual(self.client.patch(url).status_code, 405)
        self.assertEqual(self.client.delete(url).status_code, 405)
        return limit4

    @mock.patch('requests.post', mock_import_repository_task)
    def test_limit_str(self):
        """Test the text representation of a Limit."""
        limit = self.test_limit_memory()
        l = Limit.objects.get(uuid=limit['uuid'])
        self.assertEqual(str(l), "{}-{}".format(limit['app'], limit['uuid'][:7]))
