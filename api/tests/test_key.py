"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import os.path

from django.test import TestCase
from django.test.utils import override_settings

from deis import settings


@override_settings(CELERY_ALWAYS_EAGER=True)
class KeyTest(TestCase):

    """Tests cloud provider credentials"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))

    def test_key(self):
        """
        Test that a user can add, remove and manage their SSH public keys
        """
        url = '/api/keys'
        body = {'id': 'mykey@box.local', 'public': 'ssh-rsa XXX'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        key_id = response.data['id']
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/api/keys/{key_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(body['id'], response.data['id'])
        self.assertEqual(body['public'], response.data['public'])
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_key_cm(self):
        """
        Test that creating and deleting a key updates configuration management
        """
        url = '/api/keys'
        body = {'id': 'mykey@box.local', 'public': 'ssh-rsa XXX'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        key_id = response.data['id']
        path = os.path.join(settings.TEMPDIR, 'user-autotest')
        with open(path) as f:
            data = json.loads(f.read())
        self.assertIn('id', data)
        self.assertEquals(data['id'], 'autotest')
        self.assertIn(body['id'], data['ssh_keys'])
        self.assertEqual(body['public'], data['ssh_keys'][body['id']])
        url = '/api/keys/{key_id}'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)
        with open(path) as f:
            data = json.loads(f.read())
        self.assertNotIn(body['id'], data['ssh_keys'])

    def test_key_duplicate(self):
        """
        Test that a user cannot add a duplicate key
        """
        url = '/api/keys'
        body = {'id': 'mykey@box.local', 'public': 'ssh-rsa XXX'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)
