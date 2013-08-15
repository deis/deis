"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import os.path

from django.test import TestCase
from deis import settings


class FormationTest(TestCase):

    """Tests creation of different node formations"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        url = '/api/providers'
        creds = {'secret_key': 'x' * 64, 'access_key': 1 * 20}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {'id': 'autotest', 'provider': 'autotest',
                'params': json.dumps({'region': 'us-west-2', 'instance_size': 'm1.medium'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_formation(self):
        """
        Test that a user can create, read, update and delete a node formation
        """
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        self.assertIn('layers', response.data)
        self.assertIn('containers', response.data)
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/api/formations/{formation_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        body = {'id': 'new'}
        response = self.client.patch(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_formation_auto_id(self):
        body = {'id': 'autotest'}
        response = self.client.post('/api/formations', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['id'])
        return response

    def test_formation_errors(self):
        # test duplicate id
        body = {}
        response = self.client.post('/api/formations', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['id'])
        body = {'id': response.data['id']}
        response = self.client.post('/api/formations', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(json.loads(response.content), 'Formation with this Id already exists.')

    def test_formation_scale_errors(self):
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        # scaling containers without a runtime layer should throw an error
        url = '/api/formations/{formation_id}/scale/containers'.format(**locals())
        body = {'web': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(json.loads(response.content),
                         'Must create a "runtime" layer to host containers')
        # scaling containers without any runtime nodes should throw an error
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'run_list': 'recipe[deis::runtime]'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/scale/containers'.format(**locals())
        body = {'web': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(json.loads(response.content),
                         'Must scale runtime nodes > 0 to host containers')

    def test_formation_actions(self):
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        # test calculate
        url = '/api/formations/{formation_id}/calculate'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertIn('nodes', response.data)
        self.assertIn('containers', response.data)
        self.assertIn('proxy', response.data)
        self.assertIn('release', response.data)
        # test converge
        url = '/api/formations/{formation_id}/converge'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertIn('nodes', response.data)
        self.assertIn('containers', response.data)
        self.assertIn('proxy', response.data)
        self.assertIn('release', response.data)
        # test balance
        url = '/api/formations/{formation_id}/balance'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertIn('nodes', response.data)
        self.assertIn('containers', response.data)
        self.assertIn('proxy', response.data)
        self.assertIn('release', response.data)
        # test logs
        if not os.path.exists(settings.DEIS_LOG_DIR):
            os.mkdir(settings.DEIS_LOG_DIR)
        path = os.path.join(settings.DEIS_LOG_DIR, formation_id + '.log')
        with open(path, 'w') as f:
            f.write(FAKE_LOG_DATA)
        url = '/api/formations/{formation_id}/logs'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, FAKE_LOG_DATA)


FAKE_LOG_DATA = """
2013-08-15 12:41:25 [33454] [INFO] Starting gunicorn 17.5
2013-08-15 12:41:25 [33454] [INFO] Listening at: http://0.0.0.0:5000 (33454)
2013-08-15 12:41:25 [33454] [INFO] Using worker: sync
2013-08-15 12:41:25 [33457] [INFO] Booting worker with pid 33457
"""
