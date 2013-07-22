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
        creds = {'secret_key': 'x'*64, 'access_key': 1*20}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {'id': 'autotest', 'provider': 'autotest', 'ssh_username': 'ubuntu',
                'params': json.dumps({'region': 'us-west-2', 'instance_size': 'm1.medium'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        
    def test_formation(self):
        """
        Test that a user can create, read, update and delete a node formation
        """
        url = '/api/formations'
        body = {'id': 'autotest', 'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        self.assertIn('flavor', response.data)
        self.assertIn('image', response.data)
        self.assertIn('structure', response.data)
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
        body = {'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post('/api/formations', json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['id'])
        return response

    def test_formation_errors(self):
        # test duplicate id
        body = {'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post('/api/formations', json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['id'])
        body = {'id': response.data['id'], 'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post('/api/formations', json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(json.loads(response.content), 'Formation with this Id already exists.')

    def test_formation_scale_errors(self):
        url = '/api/formations'
        body = {'id': 'autotest', 'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        # scaling containers without backends should throw an error
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'web': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(json.loads(response.content), 'Must scale backends > 0 to host containers')
    
    def test_formation_actions(self):
        url = '/api/formations'
        body = {'id': 'autotest', 'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        # scale up
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'backends': 4, 'proxies': 2, 'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/backends'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 4)
        url = '/api/formations/{formation_id}/proxies'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        url = '/api/formations/{formation_id}/containers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
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