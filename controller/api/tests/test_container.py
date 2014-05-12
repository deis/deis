"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.contrib.auth.models import User
from django.test import TransactionTestCase
from django.test.utils import override_settings

from django_fsm import TransitionNotAllowed

from api.models import Container, App


@override_settings(CELERY_ALWAYS_EAGER=True)
class ContainerTest(TransactionTestCase):

    """Tests creation of containers on nodes"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        body = {'id': 'autotest', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': {}}
        response = self.client.post('/api/clusters', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create a malicious scheduler as well
        body['id'] = 'autotest2'
        body['type'] = 'faulty'
        response = self.client.post('/api/clusters', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_container_state_good(self):
        """Test that the finite state machine transitions with a good scheduler"""
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # create a container
        c = Container.objects.create(owner=User.objects.get(username='autotest'),
                                     app=App.objects.get(id=app_id),
                                     release=App.objects.get(id=app_id).release_set.latest(),
                                     type='web',
                                     num=1)
        self.assertEqual(c.state, 'initialized')
        # test an illegal transition
        self.assertRaises(TransitionNotAllowed, lambda: c.start())
        c.create()
        self.assertEqual(c.state, 'created')
        c.start()
        self.assertEqual(c.state, 'up')
        c.deploy(App.objects.get(id=app_id).release_set.latest())
        self.assertEqual(c.state, 'up')
        c.destroy()
        self.assertEqual(c.state, 'destroyed')

    def test_container_state_bad(self):
        """Test that the finite state machine transitions with a faulty scheduler"""
        url = '/api/apps'
        body = {'cluster': 'autotest2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # create a container
        c = Container.objects.create(owner=User.objects.get(username='autotest'),
                                     app=App.objects.get(id=app_id),
                                     release=App.objects.get(id=app_id).release_set.latest(),
                                     type='web',
                                     num=1)
        self.assertEqual(c.state, 'initialized')
        self.assertRaises(Exception, lambda: c.create())
        self.assertEqual(c.state, 'initialized')
        # test an illegal transition
        self.assertRaises(TransitionNotAllowed, lambda: c.start())
        self.assertEqual(c.state, 'initialized')
        self.assertRaises(
            Exception,
            lambda: c.deploy(
                App.objects.get(id=app_id).release_set.latest()
            )
        )
        self.assertEqual(c.state, 'down')
        self.assertRaises(Exception, lambda: c.destroy())
        self.assertEqual(c.state, 'down')
        self.assertRaises(Exception, lambda: c.run('echo hello world'))
        self.assertEqual(c.state, 'down')

    def test_container_state_protected(self):
        """Test that you cannot directly modify the state"""
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        c = Container.objects.create(owner=User.objects.get(username='autotest'),
                                     app=App.objects.get(id=app_id),
                                     release=App.objects.get(id=app_id).release_set.latest(),
                                     type='web',
                                     num=1)
        self.assertRaises(AttributeError, lambda: setattr(c, 'state', 'up'))

    def test_container_api_heroku(self):
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example', 'procfile': json.dumps({'web': 'node server.js',
                                                                     'worker': 'node worker.js'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # scale up
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 204)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
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
        self.assertEqual(response.status_code, 204)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        self.assertEqual(max(c['num'] for c in response.data['results']), 2)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        # scale down to 0
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 0, 'worker': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 204)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)

    def test_container_api_docker(self):
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example', 'dockerfile': "FROM busybox\nCMD /bin/true"}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # scale up
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'cmd': 6}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 204)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        # test listing/retrieving container info
        url = "/api/apps/{app_id}/containers/cmd".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        # scale down
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'cmd': 3}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 204)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)
        self.assertEqual(max(c['num'] for c in response.data['results']), 3)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        # scale down to 0
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'cmd': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 204)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        url = "/api/apps/{app_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)

    def test_container_release(self):
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example', 'procfile': json.dumps({'web': 'node server.js',
                                                                     'worker': 'node worker.js'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # scale up
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 204)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertEqual(response.data['results'][0]['release'], 'v2')
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data['image'], body['image'])
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertEqual(response.data['results'][0]['release'], 'v3')
        # post new config
        url = "/api/apps/{app_id}/config".format(**locals())
        body = {'values': json.dumps({'KEY': 'value'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertEqual(response.data['results'][0]['release'], 'v4')

    def test_container_errors(self):
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 'not_an_int'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'Invalid scaling format', status_code=400)
        body = {'invalid': 1}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'Container type invalid', status_code=400)

    def test_container_str(self):
        """Test the text representation of a container."""
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example', 'procfile': json.dumps({'web': 'node server.js',
                                                                     'worker': 'node worker.js'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # scale up
        url = "/api/apps/{app_id}/scale".format(**locals())
        body = {'web': 4, 'worker': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 204)
        # should start with zero
        url = "/api/apps/{app_id}/containers".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 6)
        uuid = response.data['results'][0]['uuid']
        container = Container.objects.get(uuid=uuid)
        self.assertEqual(container.short_name(),
                         "{}.{}.{}".format(container.app, container.type, container.num))
        self.assertEqual(str(container),
                         "{}.{}.{}".format(container.app, container.type, container.num))
