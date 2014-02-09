"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

from django.test import TestCase
from django.test.utils import override_settings

import json


class ServicesTest(TestCase):

    """Tests creation of services"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(self.client.login(username='autotest', password='password'))

    def test_services(self):
        """
        Test that a user can list and enable services
        """
        # Listing providers
        url = '/api/services'
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        mock_data = {u'enabled': False, 'description': u'Mock service for testing purposes'}
        self.assertEqual(response.data['mock'], mock_data)

        # Enabling a provider
        url = '/api/services/mock'
        body = {'enabled': True}
        response = self.client.put(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/services'
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.data['mock']['enabled'], True)

        # Disabling a provider
        url = '/api/services/mock'
        body = {'enabled': False}
        response = self.client.put(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/services'
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.data['mock']['enabled'], False)

        # Enabling a non-existent provider
        url = '/api/services/halflife3'
        body = {'enabled': True}
        response = self.client.put(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 404)


@override_settings(CELERY_ALWAYS_EAGER=True)
class AddonsTest(TestCase):

    """Tests creation of services"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(self.client.login(username='autotest', password='password'))

        # Setup to allow creation of apps
        url = '/api/providers'
        creds = {'secret_key': 'x' * 64, 'access_key': 1 * 20}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {'id': 'autotest', 'provider': 'autotest',
                'params': json.dumps({'region': 'us-west-2'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        response = self.client.post('/api/formations', json.dumps(
            {'id': 'autotest', 'domain': 'localhost.localdomain'}),
            content_type='application/json')
        self.assertEqual(response.status_code, 201)

        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.app_id = response.data['id']

        # Enable the 'mock' service provider
        url = '/api/services/mock'
        body = {'enabled': True}
        response = self.client.put(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_services(self):
        """
        Test that a user can add, remove and retrieve details about service instances
        """

        # Listing addons requires a call to /api/services so is covered by Service
        # tests.

        # Adding an addon
        url = '/api/apps/{}/addons/mock'.format(self.app_id)
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 201)
        addon = response.data
        self.assertIn("mock://asfdfsdfg:hfghgdfg@localhost:1337/", addon['uri'])
        self.assertEqual(addon['docs'], "/services/mock/docs")
        self.assertEqual(addon['dashboard'], "/services/mock/dashboard")

        # Removing an addon
        url = '/api/apps/{}/addons/mock'.format(self.app_id)
        response = self.client.delete(url, content_type='application/json')
        self.assertEqual(response.status_code, 204)

        # Adding a nonexistent service
        url = '/api/apps/{}/addons/halflife3'.format(self.app_id)
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 404)

        # Adding a disabled service
        url = '/api/services/mock'
        body = {'enabled': False}
        response = self.client.put(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/apps/{}/addons/mock'.format(self.app_id)
        response = self.client.delete(url, content_type='application/json')
        self.assertEqual(response.status_code, 404)
