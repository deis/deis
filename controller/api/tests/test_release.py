"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase
from django.test.utils import override_settings

from api.models import Release


@override_settings(CELERY_ALWAYS_EAGER=True)
class ReleaseTest(TestCase):

    """Tests push notification from build system"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        body = {'id': 'autotest', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': {}}
        response = self.client.post('/api/clusters', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_release(self):
        """
        Test that a release is created when a cluster is created, and
        that updating config or build or triggers a new release
        """
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check that updating config rolls a new release
        url = '/api/apps/{app_id}/config'.format(**locals())
        body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertIn('NEW_URL1', json.loads(response.data['values']))
        # check to see that an initial release was created
        url = '/api/apps/{app_id}/releases'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        # account for the config release as well
        self.assertEqual(response.data['count'], 2)
        url = '/api/apps/{app_id}/releases/v1'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        release1 = response.data
        self.assertIn('config', response.data)
        self.assertIn('build', response.data)
        self.assertEquals(release1['version'], 1)
        # check to see that a new release was created
        url = '/api/apps/{app_id}/releases/v2'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        release2 = response.data
        self.assertNotEqual(release1['uuid'], release2['uuid'])
        self.assertNotEqual(release1['config'], release2['config'])
        self.assertEqual(release1['build'], release2['build'])
        self.assertEquals(release2['version'], 2)
        # check that updating the build rolls a new release
        url = '/api/apps/{app_id}/builds'.format(**locals())
        build_config = json.dumps({'PATH': 'bin:/usr/local/bin:/usr/bin:/bin'})
        body = {'image': 'autotest/example'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data['image'], body['image'])
        # check to see that a new release was created
        url = '/api/apps/{app_id}/releases/v3'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        release3 = response.data
        self.assertNotEqual(release2['uuid'], release3['uuid'])
        self.assertNotEqual(release2['build'], release3['build'])
        self.assertEquals(release3['version'], 3)
        # check that we can fetch a previous release
        url = '/api/apps/{app_id}/releases/v2'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        release2 = response.data
        self.assertNotEqual(release2['uuid'], release3['uuid'])
        self.assertNotEqual(release2['build'], release3['build'])
        self.assertEquals(release2['version'], 2)
        # disallow post/put/patch/delete
        url = '/api/apps/{app_id}/releases'.format(**locals())
        self.assertEqual(self.client.post(url).status_code, 405)
        self.assertEqual(self.client.put(url).status_code, 405)
        self.assertEqual(self.client.patch(url).status_code, 405)
        self.assertEqual(self.client.delete(url).status_code, 405)
        return release3

    def test_release_rollback(self):
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # try to rollback with only 1 release extant, expecting 404
        url = "/api/apps/{app_id}/releases/rollback/".format(**locals())
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 404)
        # update config to roll a new release
        url = '/api/apps/{app_id}/config'.format(**locals())
        body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # update the build to roll a new release
        url = '/api/apps/{app_id}/builds'.format(**locals())
        build_config = json.dumps({'PATH': 'bin:/usr/local/bin:/usr/bin:/bin'})
        body = {'image': 'autotest/example'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # rollback and check to see that a 4th release was created
        # with the build and config of release #2
        url = "/api/apps/{app_id}/releases/rollback/".format(**locals())
        response = self.client.post(url, content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/apps/{app_id}/releases'.format(**locals())
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 4)
        url = '/api/apps/{app_id}/releases/v2'.format(**locals())
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        release2 = response.data
        self.assertEquals(release2['version'], 2)
        url = '/api/apps/{app_id}/releases/v4'.format(**locals())
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        release4 = response.data
        self.assertEquals(release4['version'], 4)
        self.assertNotEqual(release2['uuid'], release4['uuid'])
        self.assertEqual(release2['build'], release4['build'])
        self.assertEqual(release2['config'], release4['config'])
        # rollback explicitly to release #1 and check that a 5th release
        # was created with the build and config of release #1
        url = "/api/apps/{app_id}/releases/rollback/".format(**locals())
        body = {'version': 1}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/apps/{app_id}/releases'.format(**locals())
        response = self.client.get(url, content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 5)
        url = '/api/apps/{app_id}/releases/v1'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        release1 = response.data
        url = '/api/apps/{app_id}/releases/v5'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        release5 = response.data
        self.assertEqual(release5['version'], 5)
        self.assertNotEqual(release1['uuid'], release5['uuid'])
        self.assertEqual(release1['build'], release5['build'])
        self.assertEqual(release1['config'], release5['config'])
        # check to see that the current config is actually the initial one
        url = "/api/apps/{app_id}/config".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['values'], json.dumps({}))
        # rollback to #3 and see that it has the correct config
        url = "/api/apps/{app_id}/releases/rollback/".format(**locals())
        body = {'version': 3}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = "/api/apps/{app_id}/config".format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        values = json.loads(response.data['values'])
        self.assertIn('NEW_URL1', values)
        self.assertEqual('http://localhost:8080/', values['NEW_URL1'])

    def test_release_str(self):
        """Test the text representation of a release."""
        release3 = self.test_release()
        release = Release.objects.get(uuid=release3['uuid'])
        self.assertEqual(str(release), "{}-v3".format(release3['app']))

    def test_release_summary(self):
        """Test the text summary of a release."""
        release3 = self.test_release()
        release = Release.objects.get(uuid=release3['uuid'])
        # check that the release has push and env change messages
        self.assertIn('autotest deployed ', release.summary)
