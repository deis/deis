"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase
from django.test.utils import override_settings


@override_settings(CELERY_ALWAYS_EAGER=True)
class DomaainTest(TestCase):

    """Tests creation of domains"""

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
        formation_id = 'autotest'
        response = self.client.post('/api/formations', json.dumps({'id': formation_id,
                                                                   'domain': 'test.domain.org'}),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)

        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.app_id = response.data['id']  # noqa

    def test_manage_domain(self):
        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        body = {'domain': 'test-domain.example.com'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)

        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        response = self.client.get(url, content_type='application/json')
        result = response.data['results'][0]
        self.assertEqual('test-domain.example.com', result['domain'])

        url = '/api/domains/{hostname}'.format(hostname='test-domain.example.com')
        response = self.client.delete(url, content_type='application/json')
        self.assertEqual(response.status_code, 204)

        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(0, response.data['count'])

    def test_manage_domain_invalid_app(self):
        url = '/api/apps/{app_id}/domains'.format(app_id="this-app-does-not-exis")
        body = {'domain': 'test-domain.example.com'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 404)

        url = '/api/apps/{app_id}/domains'.format(app_id='this-app-does-not-exis')
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 404)

    def test_manage_domain_no_perms_on_app(self):
        self.client.logout()
        self.assertTrue(
            self.client.login(username='autotest2', password='password'))
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        new_app_id = response.data['id']  # noqa

        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        body = {'domain': 'test-domain2.example.com'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)

    def test_manage_domain_invalid_domain(self):
        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        body = {'domain': 'this_is_an.invalid.domain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)

        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        body = {'domain': 'this-is-an.invalid.a'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)

        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        body = {'domain': 'domain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)

    def test_manage_domain_wildcard(self):
        # Wildcards are not allowed for now.
        url = '/api/apps/{app_id}/domains'.format(app_id=self.app_id)
        body = {'domain': '*.deis.example.com'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)
