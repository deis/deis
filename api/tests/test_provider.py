"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase

from api.models import Provider


class ProviderTest(TestCase):

    """Tests cloud provider credentials"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))

    def test_provider(self):
        """
        Test that a user can create a config record containing
        environment variables
        """
        url = '/api/providers'
        creds = {'secret_key': 'x' * 64, 'access_key': 1 * 20}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        provider_id = response.data['id']
        response = self.client.get('/api/providers')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/api/providers/{provider_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        new_creds = {'access_key': 'new', 'secret_key': 'new'}
        body = {'type': 'mock', 'creds': json.dumps(new_creds)}
        response = self.client.patch(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['creds'], json.dumps(new_creds))
        self.assertEqual(response.data['type'], 'mock')
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_provider_str(self):
        """Test the text representation of a provider."""
        url = '/api/providers'
        body = {'id': 'autotest', 'type': 'ec2', 'creds': json.dumps({})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        provider = Provider.objects.get(uuid=response.data['uuid'])
        self.assertEqual(str(provider), 'autotest-Amazon Elastic Compute Cloud (EC2)')
