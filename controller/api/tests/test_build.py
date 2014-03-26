"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase
from django.test.utils import override_settings

from api.models import Build


@override_settings(CELERY_ALWAYS_EAGER=True)
class BuildTest(TestCase):

    """Tests build notification from build system"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        body = {'id': 'autotest', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': {}}
        response = self.client.post('/api/clusters', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_build(self):
        """
        Test that a null build is created and that users can post new builds
        """
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
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
        build3 = response.data
        self.assertEqual(response.data['image'], body['image'])
        self.assertNotEqual(build2['uuid'], build3['uuid'])
        # disallow put/patch/delete
        self.assertEqual(self.client.put(url).status_code, 405)
        self.assertEqual(self.client.patch(url).status_code, 405)
        self.assertEqual(self.client.delete(url).status_code, 405)

    def test_build_str(self):
        """Test the text representation of a build."""
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
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
