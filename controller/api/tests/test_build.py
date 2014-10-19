"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import mock
import requests

from django.test import TransactionTestCase

from api.models import Build


def mock_import_repository_task(*args, **kwargs):
    resp = requests.Response()
    resp.status_code = 200
    resp._content_consumed = True
    return resp


class BuildTest(TransactionTestCase):

    """Tests build notification from build system"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))

    @mock.patch('requests.post', mock_import_repository_task)
    def test_build(self):
        """
        Test that a null build is created and that users can post new builds
        """
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check to see that an initial build was created
        url = "/api/apps/{app_id}/builds".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 1)
        # post a new build
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        build_id = response.data['uuid']
        build1 = response.data
        self.assertEqual(response.data['image'], body['image'])
        # read the build
        url = "/api/apps/{app_id}/builds/{build_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        build2 = response.data
        self.assertEqual(build1, build2)
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertIn('x-deis-release', response._headers)
        build3 = response.data
        self.assertEqual(response.data['image'], body['image'])
        self.assertNotEqual(build2['uuid'], build3['uuid'])
        # disallow put/patch/delete
        self.assertEqual(self.client.put(url).status_code, 405)
        self.assertEqual(self.client.patch(url).status_code, 405)
        self.assertEqual(self.client.delete(url).status_code, 405)

    @mock.patch('requests.post', mock_import_repository_task)
    def test_build_default_containers(self):
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # post an image as a build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = "/api/apps/{app_id}/containers/cmd".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        container = response.data['results'][0]
        self.assertEqual(container['type'], 'cmd')
        self.assertEqual(container['num'], 1)
        # start with a new app
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # post a new build with procfile
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example',
                'sha': 'a'*40,
                'dockerfile': "FROM scratch"}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = "/api/apps/{app_id}/containers/cmd".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        container = response.data['results'][0]
        self.assertEqual(container['type'], 'cmd')
        self.assertEqual(container['num'], 1)
        # start with a new app
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # post a new build with procfile
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example',
                'sha': 'a'*40,
                'dockerfile': "FROM scratch",
                'procfile': {'worker': 'node worker.js'}}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = "/api/apps/{app_id}/containers/cmd".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        container = response.data['results'][0]
        self.assertEqual(container['type'], 'cmd')
        self.assertEqual(container['num'], 1)
        # start with a new app
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # post a new build with procfile
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example',
                'sha': 'a'*40,
                'procfile': json.dumps({'web': 'node server.js',
                                        'worker': 'node worker.js'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = "/api/apps/{app_id}/containers/web".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        container = response.data['results'][0]
        self.assertEqual(container['type'], 'web')
        self.assertEqual(container['num'], 1)

    @mock.patch('requests.post', mock_import_repository_task)
    def test_build_str(self):
        """Test the text representation of a build."""
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        build = Build.objects.get(uuid=response.data['uuid'])
        self.assertEqual(str(build), "{}-{}".format(
                         response.data['app'], response.data['uuid'][:7]))

    @mock.patch('requests.post', mock_import_repository_task)
    def test_admin_can_create_builds_on_other_apps(self):
        """If a user creates an application, an administrator should be able
        to push builds.
        """
        # create app as non-admin
        self.client.login(username='autotest2', password='password')
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # post a new build as admin
        self.client.login(username='autotest', password='password')
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        build = Build.objects.get(uuid=response.data['uuid'])
        self.assertEqual(str(build), "{}-{}".format(
                         response.data['app'], response.data['uuid'][:7]))
