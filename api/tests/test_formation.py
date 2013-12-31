"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import os.path
import uuid

from django.test import TestCase
from django.test.utils import override_settings

from deis import settings


@override_settings(CELERY_ALWAYS_EAGER=True)
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
        body = {'id': 'auto_test-1', 'domain': 'localhost.localdomain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        self.assertIn('nodes', response.data)
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/api/formations/{formation_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        body = {'domain': 'new'}
        response = self.client.patch(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['domain'], 'new')
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_formation_delete(self):
        """
        Test that deleting a formation also deletes its apps.
        """
        url = '/api/formations'
        body = {'id': 'auto_test-1', 'domain': 'localhost.localdomain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        url = '/api/apps'
        body = {'formation': 'auto_test-1'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/api/formations/{formation_id}'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)

    def test_formation_cm(self):
        """
        Test that configuration management is updated on formation changes
        """
        url = '/api/formations'
        body = {'id': 'autotest-' + uuid.uuid4().hex[:4], 'domain': 'localhost.localdomain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        path = os.path.join(settings.TEMPDIR, 'formation-{}'.format(formation_id))
        with open(path) as f:
            data = json.loads(f.read())
        self.assertIn('id', data)
        self.assertEquals(data['id'], formation_id)
        url = '/api/formations/{formation_id}'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)
        self.assertFalse(os.path.exists(path))

    def test_formation_id(self):
        # Ensure that dashes and underscores are allowed in id, per Django's
        # definition of a SlugField.
        body = {'id': 'auto_test-1', 'domain': 'localhost.localdomain'}
        response = self.client.post('/api/formations', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue('id', response.data)
        body = {'id': response.data['id'], 'domain': 'localhost.localdomain'}
        response = self.client.post('/api/formations', json.dumps(body),
                                    content_type='application/json')
        self.assertContains(response, 'Formation with this Id already exists.', status_code=400)

    def test_formation_actions(self):
        url = '/api/formations'
        body = {'id': 'autotest', 'domain': 'localhost.localdomain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        # test calculate
        url = '/api/formations/{formation_id}/calculate'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertIn('nodes', response.data)
        # test converge
        url = '/api/formations/{formation_id}/converge'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        # test balance
        url = '/api/formations/{formation_id}/balance'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
