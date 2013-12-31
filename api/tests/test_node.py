"""
Unit tests for the Deis api app.

Run these tests with "./manage.py test api.tests.test_node"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase
from django.test.utils import override_settings

from api.models import Node


@override_settings(CELERY_ALWAYS_EAGER=True)
class NodeTest(TestCase):

    """Tests creation of nodes"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        url = '/api/providers'
        creds = {'secret_key': 'x' * 64, 'access_key': 'A' * 20}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {'id': 'autotest', 'provider': 'autotest',
                'params': json.dumps({'region': 'us-west-2', 'instance_size': 'm1.medium'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_node(self):
        """
        Test that a user can create, read, update and delete a node
        """
        url = '/api/formations'
        body = {'id': 'autotest', 'domain': 'localhost.localdomain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'run_list': 'recipe[deis::runtime]'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # should start with zero
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'runtime': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        node = response.data['results'][0]['id']
        url = '/api/formations/{formation_id}/nodes/{node}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        # query alternate nodes endpoints
        url = '/api/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        node = response.data['results'][0]['id']
        url = '/api/nodes/{node}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertIn('fqdn', response.data)
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_node_scale(self):
        url = '/api/formations'
        body = {'id': 'autotest', 'domain': 'localhost.localdomain'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'proxy', 'flavor': 'autotest', 'run_list': 'recipe[deis::proxy]',
                'runtime': False, 'proxy': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'run_list': 'recipe[deis::runtime]',
                'runtime': True, 'proxy': False}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        # should start with zero
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'proxy': 2, 'runtime': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['nodes'], json.dumps(body))
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        # scale down
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'proxy': 1, 'runtime': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        url = '/api/formations/{formation_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['nodes'], json.dumps(body))
        # scale down to 0
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'proxy': 0, 'runtime': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        url = '/api/formations/{formation_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['nodes'], json.dumps(body))

    def test_node_scale_errors(self):
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']  # noqa
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'runtime': 'not_an_int'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'Invalid scaling format', status_code=400)
        body = {'runtime': '1'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'Layer matching query does not exist', status_code=400)
        url = '/api/providers/autotest'
        body = {'creds': json.dumps({})}
        response = self.client.patch(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'run_list': 'recipe[deis::runtime]'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'runtime': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'No provider credentials available', status_code=400)

    def test_node_actions(self):
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'runtime': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'runtime': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # get our node
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 1)
        node_id = response.data['results'][0]['id']
        url = '/api/formations/{formation_id}/nodes/{node_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(node_id, response.data['id'])
        node = response.data
        url = '/api/nodes/{id}/converge'.format(**node)
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)

    def test_node_create(self):
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'runtime': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create a node for an existing instance
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        body = {'fqdn': 'example.com', 'layer': 'runtime'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create it again, expecting an error
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        body = {'fqdn': 'example.com', 'layer': 'runtime'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 409)
        # get our node
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 1)
        node_id = response.data['results'][0]['id']
        url = '/api/formations/{formation_id}/nodes/{node_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(node_id, response.data['id'])
        node = response.data
        # delete our node
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)
        # check the node is gone
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 0)

    def test_node_create_errors(self):
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'runtime': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create a node for an existing instance
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        body = {'fqdn': 'error', 'layer': 'runtime'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 401)

    def test_node_str(self):
        """Test the text representation of a node."""
        url = '/api/formations'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = response.data['id']
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'runtime': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create a node for an existing instance
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        body = {'fqdn': 'example.com', 'layer': 'runtime'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # get our node
        url = '/api/formations/{formation_id}/nodes'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 1)
        node_id = response.data['results'][0]['id']
        node = Node.objects.get(id=node_id)
        self.assertEqual(str(node), 'autotest-runtime-1')
