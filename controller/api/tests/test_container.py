"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase
from django.test.utils import override_settings

from api.models import Container
from deis import settings


def get_allocations(container_dict):
    counts = {}
    for container in container_dict.values():
        name, _id = container.split(':')
        if name in counts:
            counts[name] += 1
        else:
            counts[name] = 1
    return sorted(counts.values())


@override_settings(CELERY_ALWAYS_EAGER=True)
class ContainerTest(TestCase):

    """Tests creation of containers on nodes"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        url = '/api/providers'
        creds = {'access_key': getattr(settings, 'EC2_ACCESS_KEY', 'x' * 32),
                 'secret_key': getattr(settings, 'EC2_SECRET_KEY', 'x' * 64)}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {'id': 'autotest', 'provider': 'autotest',
                'params': json.dumps({'region': 'us-west-2', 'instance_size': 'm1.medium'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        response = self.client.post('/api/formations', json.dumps(
            {'id': 'autotest', 'domain': 'localhost.localdomain'}),
            content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create & scale a basic formation
        formation_id = 'autotest'
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'proxy', 'flavor': 'autotest', 'proxy': True,
                'run_list': 'recipe[deis::proxy]'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'runtime': True,
                'run_list': 'recipe[deis::runtime]'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'proxy': 2, 'runtime': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)

    def test_container(self):
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['containers'], json.dumps(body))
        # test listing/retrieving container info
        url = "/api/apps/{app_id}/containers/web".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 4)
        num = response.data['results'][0]['num']
        url = "/api/apps/{app_id}/containers/web/{num}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['num'], num)
        # scale down
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 2, 'worker': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['containers'], json.dumps(body))
        # scale down to 0
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 0, 'worker': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['containers'], json.dumps(body))

    def test_container_errors(self):
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 'not_an_int'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'Invalid scaling format', status_code=400)

    def test_container_single_layer(self):
        # create & scale a single layer formation
        response = self.client.post('/api/formations', json.dumps(
            {'id': 'single-layer', 'domain': 'localhost.localdomain'}),
            content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = 'single-layer'
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'default', 'flavor': 'autotest', 'proxy': True, 'runtime': True,
                'run_list': 'recipe[deis::runtime],recipe[deis::proxy]'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'default': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/apps'
        body = {'formation': 'single-layer'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['containers'], json.dumps(body))
        # scale down
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 2, 'worker': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['containers'], json.dumps(body))
        # scale down to 0
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 0, 'worker': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['containers'], json.dumps(body))

    def test_container_multiple_apps(self):
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app1_id = response.data['id']
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app2_id = response.data['id']
        # scale up
        url = "/api/apps/{app1_id}/scale".format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app2_id}/scale".format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app1_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        url = "/api/apps/{app2_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        # check port assignments
        url = '/api/formations/autotest/calculate'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        databag = response.data.copy()
        ports = []
        for app in databag['apps'].values():
            for containers_set in app['containers'].values():
                for node_port in containers_set.values():
                    _, port = node_port.split(':')
                    ports.append(int(port))
        ports.sort()
        self.assertEqual(ports, range(10001, 10013))
        # scale down
        url = "/api/apps/{app1_id}/scale".format(**locals())
        body = {'web': 2, 'worker': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app2_id}/scale".format(**locals())
        body = {'web': 2, 'worker': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app1_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        url = "/api/apps/{app2_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        # check port assignments
        url = '/api/formations/autotest/calculate'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        databag = response.data.copy()
        ports = []
        for app in databag['apps'].values():
            for containers_set in app['containers'].values():
                for node_port in containers_set.values():
                    _, port = node_port.split(':')
                    ports.append(int(port))
        ports.sort()
        self.assertEqual(len(set(ports)), 6)

    def test_container_allocation(self):
        url = '/api/apps'
        formation_id = 'autotest'
        body = {'formation': formation_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # With 4 nodes and 13 web containers
        url = "/api/formations/{formation_id}/scale".format(**locals())
        body = {'runtime': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 13}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # test that one node has 4 and 3 nodes have 3 containers
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(get_allocations(response.data['apps'][app_id]['containers']['web']),
                         [3, 3, 3, 4])
        # With 1 node
        url = "/api/formations/{formation_id}/scale".format(**locals())
        body = {'runtime': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # test that the node has all 13 containers
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(get_allocations(response.data['apps'][app_id]['containers']['web']),
                         [13])
        # With 2 nodes
        url = "/api/formations/{formation_id}/scale".format(**locals())
        body = {'runtime': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # test that one has 6 and the other has 7 containers
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(get_allocations(response.data['apps'][app_id]['containers']['web']),
                         [6, 7])
        # With 8 containers
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 8}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # test that both have 4 containers
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(get_allocations(response.data['apps'][app_id]['containers']['web']),
                         [4, 4])
        # With 0 nodes
        url = "/api/formations/{formation_id}/scale".format(**locals())
        body = {'runtime': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # test that there are no containers
        self.assertNotIn('web', response.data['apps'][app_id]['containers'])
        # With 5 containers
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 5}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        # test that we get an error message about runtime nodes
        self.assertEqual(response.status_code, 400)
        self.assertIn('No nodes available for containers', response.data)
        # With 1 node
        url = "/api/formations/{formation_id}/scale".format(**locals())
        body = {'runtime': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # test that it gets all 8 containers
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(get_allocations(response.data['apps'][app_id]['containers']['web']),
                         [8])

    def test_container_balance(self):
        url = '/api/apps'
        formation_id = 'autotest'
        body = {'formation': formation_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # scale layer
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'runtime': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # scale up
        url = '/api/apps/{app_id}/scale'.format(**locals())
        body = {'web': 8, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # scale layer up
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'runtime': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # calculate the formation
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url)
        containers = response.data['apps'][app_id]['containers']
        # check balance of web types
        by_backend = {}
        for c in containers['web'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([len(by_backend[b]) for b in by_backend.keys()])
        b_max = max([len(by_backend[b]) for b in by_backend.keys()])
        self.assertLess(b_max - b_min, 2)
        # check balance of worker types
        by_backend = {}
        for c in containers['worker'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([len(by_backend[b]) for b in by_backend.keys()])
        b_max = max([len(by_backend[b]) for b in by_backend.keys()])
        self.assertLess(b_max - b_min, 2)
        # scale up more
        url = '/api/apps/{app_id}/scale'.format(**locals())
        body = {'web': 6, 'worker': 4}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # calculate the formation
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url)
        containers = response.data['apps'][app_id]['containers']
        # check balance of web types
        by_backend = {}
        for c in containers['web'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([len(by_backend[b]) for b in by_backend.keys()])
        b_max = max([len(by_backend[b]) for b in by_backend.keys()])
        self.assertLess(b_max - b_min, 2)
        # check balance of worker types
        by_backend = {}
        for c in containers['worker'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([len(by_backend[b]) for b in by_backend.keys()])
        b_max = max([len(by_backend[b]) for b in by_backend.keys()])
        self.assertLess(b_max - b_min, 2)
        # scale down
        url = '/api/apps/{app_id}/scale'.format(**locals())
        body = {'web': 2, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 4)
        # calculate the formation
        url = "/api/formations/{formation_id}/calculate".format(**locals())
        response = self.client.post(url)
        containers = response.data['apps'][app_id]['containers']
        # check balance of web types
        by_backend = {}
        for c in containers['web'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([len(by_backend[b]) for b in by_backend.keys()])
        b_max = max([len(by_backend[b]) for b in by_backend.keys()])
        self.assertLess(b_max - b_min, 2)
        # check balance of worker types
        by_backend = {}
        for c in containers['worker'].values():
            backend, port = c.split(':')
            by_backend.setdefault(backend, []).append(port)
        b_min = min([len(by_backend[b]) for b in by_backend.keys()])
        b_max = max([len(by_backend[b]) for b in by_backend.keys()])
        self.assertLess(b_max - b_min, 2)

    def test_container_str(self):
        """Test the text representation of a container."""
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # scale up
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        uuid = response.data['results'][0]['uuid']
        container = Container.objects.get(uuid=uuid)
        self.assertEqual(container.short_name(),
                         "{}.{}".format(container.type, container.num))
        self.assertEqual(str(container),
                         "{} {}".format(container.formation.id, container.short_name()))
