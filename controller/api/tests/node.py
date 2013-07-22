"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase


class NodeTest(TestCase):

    """Tests creation of nodes"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        url = '/api/providers'
        creds = {'secret_key': 'x'*64, 'access_key': 'A'*20}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {'id': 'autotest', 'provider': 'autotest', 'ssh_username': 'ubuntu',
                'params': json.dumps({'region': 'us-west-2', 'instance_size': 'm1.medium'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        
    def test_node(self):
        """
        Test that a user can create, read, update and delete a node
        """
        url = '/api/formations'
        body = {'id': 'autotest', 'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        # should start with zero
        url = '/api/formations/{formation_id}/backends'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'backends': 8, 'proxies': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 10)
        # scale down
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'backends': 2, 'proxies': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        # scale down to 0
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'backends': 0, 'proxies': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)