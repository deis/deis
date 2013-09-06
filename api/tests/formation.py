"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase


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
        self.assertIn('nodes', response.data)
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

    def test_formation_id(self):
        body = {'id': 'autotest'}
        response = self.client.post('/api/formations', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue('id', response.data)
        body = {'id': response.data['id']}
        response = self.client.post('/api/formations', json.dumps(body),
                                    content_type='application/json')
        self.assertContains(response, 'Formation with this Id already exists.', status_code=400)

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
        # test converge
        url = '/api/formations/{formation_id}/converge'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        # test balance
        url = '/api/formations/{formation_id}/balance'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
