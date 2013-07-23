"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase

from deis import settings


class ContainerTest(TestCase):

    """Tests creation of containers on nodes"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        url = '/api/providers'
        creds = {'access_key': getattr(settings, 'EC2_ACCESS_KEY', 'x'*32),
                 'secret_key': getattr(settings, 'EC2_SECRET_KEY', 'x'*64)}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {'id': 'autotest', 'provider': 'autotest', 'ssh_username': 'ubuntu',
                'params': json.dumps({'region': 'us-west-2', 'instance_size': 'm1.medium'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
    
    def test_container_scale(self):
        url = '/api/formations'
        body = {'id': 'autotest', 'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        # scale backends
        url = '/api/formations/{formation_id}/backends'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'backends': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/backends'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 4)
        # should start with zero
        url = '/api/formations/{formation_id}/containers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/containers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        # scale down
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'web': 2, 'worker': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        url = '/api/formations/{formation_id}/containers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        # scale down to 0
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'web': 0, 'worker': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/containers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)

    def test_container_balance(self):
        url = '/api/formations'
        body = {'id': 'autotest', 'flavor': 'autotest', 'image': 'deis/autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        # scale backends
        url = '/api/formations/{formation_id}/backends'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'backends': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # should start with zero
        url = '/api/formations/{formation_id}/containers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'web': 8, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # scale backends
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'backends': 4 }
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # calculate the formation
        url = '/api/formations/{formation_id}/calculate'.format(**locals())
        response = self.client.post(url)
        containers = response.data['containers']
        # check balance of web types
        by_backend = {}
        for c in containers['web'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([ len(by_backend[b]) for b in by_backend.keys() ])
        b_max = max([ len(by_backend[b]) for b in by_backend.keys() ])
        self.assertLess(b_max - b_min, 2)
        # check balance of worker types
        by_backend = {}
        for c in containers['worker'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([ len(by_backend[b]) for b in by_backend.keys() ])
        b_max = max([ len(by_backend[b]) for b in by_backend.keys() ])
        self.assertLess(b_max - b_min, 2)
        # scale up more
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'web': 6, 'worker': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # calculate the formation
        url = '/api/formations/{formation_id}/calculate'.format(**locals())
        response = self.client.post(url)
        containers = response.data['containers']
        # check balance of web types
        by_backend = {}
        for c in containers['web'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([ len(by_backend[b]) for b in by_backend.keys() ])
        b_max = max([ len(by_backend[b]) for b in by_backend.keys() ])
        self.assertLess(b_max - b_min, 2)
        # check balance of worker types
        by_backend = {}
        for c in containers['worker'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([ len(by_backend[b]) for b in by_backend.keys() ])
        b_max = max([ len(by_backend[b]) for b in by_backend.keys() ])
        self.assertLess(b_max - b_min, 2)
        # scale down
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'web': 2, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/containers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 4)
        # calculate the formation
        url = '/api/formations/{formation_id}/calculate'.format(**locals())
        response = self.client.post(url)
        containers = response.data['containers']
        # check balance of web types
        by_backend = {}
        for c in containers['web'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([ len(by_backend[b]) for b in by_backend.keys() ])
        b_max = max([ len(by_backend[b]) for b in by_backend.keys() ])
        self.assertLess(b_max - b_min, 2)
        # check balance of worker types
        by_backend = {}
        for c in containers['worker'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([ len(by_backend[b]) for b in by_backend.keys() ])
        b_max = max([ len(by_backend[b]) for b in by_backend.keys() ])
        self.assertLess(b_max - b_min, 2)
