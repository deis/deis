"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import uuid

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
        url = '/api/providers'
        creds = {'secret_key': 'x' * 64, 'access_key': 1 * 20}
        body = {'id': 'autotest', 'type': 'mock', 'creds': json.dumps(creds)}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/flavors'
        body = {
            'id': 'autotest',
            'provider': 'autotest',
            'params': json.dumps({'region': 'us-west-2'}),
        }
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        response = self.client.post('/api/formations', json.dumps(
            {'id': 'autotest', 'domain': 'localhost.localdomain'}),
            content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_build(self):
        """
        Test that a null build is created on a new formation, and that users
        can post new builds to a formation
        """
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check to see that no initial build was created
        url = "/api/apps/{app_id}/builds".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 1)
        # post a new build
        body = {
            'sha': uuid.uuid4().hex,
            'slug_size': 4096000,
            'procfile': json.dumps({'web': 'node server.js'}),
            'url': 'http://deis.local/slugs/1c52739bbf3a44d3bfb9a58f7bbdd5fb.tar.gz',
            'checksum': uuid.uuid4().hex,
        }
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        build_id = response.data['uuid']
        build1 = response.data
        self.assertEqual(response.data['url'], body['url'])
        # read the build
        url = "/api/apps/{app_id}/builds/{build_id}".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        build2 = response.data
        self.assertEqual(build1, build2)
        # post a new build
        url = "/api/apps/{app_id}/builds".format(**locals())
        body = {
            'sha': uuid.uuid4().hex,
            'slug_size': 4096000,
            'procfile': json.dumps({'web': 'node server.js'}),
            'url': 'http://deis.local/slugs/1c52739bbf3a44d3bfb9a58f7bbdd5fb.tar.gz',
            'checksum': uuid.uuid4().hex,
        }
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        build3 = response.data
        self.assertEqual(response.data['url'], body['url'])
        self.assertNotEqual(build2['uuid'], build3['uuid'])
        # disallow put/patch/delete
        self.assertEqual(self.client.put(url).status_code, 405)
        self.assertEqual(self.client.patch(url).status_code, 405)
        self.assertEqual(self.client.delete(url).status_code, 405)

    def test_build_str(self):
        """Test the text representation of a build."""
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check to see that no initial build was created
        url = "/api/apps/{app_id}/builds".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        data = response.data['results'][0]
        build = Build.objects.get(uuid=data['uuid'])
        sha = ''
        self.assertEqual(str(build), "{}-{}".format(data['app'], sha))
