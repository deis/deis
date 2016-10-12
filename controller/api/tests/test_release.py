"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.contrib.auth.models import User
from django.test import TransactionTestCase
import mock
from django.conf import settings
from django.test.utils import override_settings
from rest_framework.authtoken.models import Token

from api.models import Release
from . import mock_status_ok


@mock.patch('api.models.publish_release', lambda *args: None)
class ReleaseTest(TransactionTestCase):

    """Tests push notification from build system"""

    fixtures = ['tests.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key

        settings.DEFAULT_PERMISSIONS_APP_MANAGEMENT = False

    @mock.patch('requests.post', mock_status_ok)
    def test_release(self):
        """
        Test that a release is created when an app is created, and
        that updating config or build or triggers a new release
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check that updating config rolls a new release
        url = '/v1/apps/{app_id}/config'.format(**locals())
        body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('NEW_URL1', response.data['values'])
        # check to see that an initial release was created
        url = '/v1/apps/{app_id}/releases'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        # account for the config release as well
        self.assertEqual(response.data['count'], 2)
        url = '/v1/apps/{app_id}/releases/v1'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release1 = response.data
        self.assertIn('config', response.data)
        self.assertIn('build', response.data)
        self.assertEquals(release1['version'], 1)
        # check to see that a new release was created
        url = '/v1/apps/{app_id}/releases/v2'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release2 = response.data
        self.assertNotEqual(release1['uuid'], release2['uuid'])
        self.assertNotEqual(release1['config'], release2['config'])
        self.assertEqual(release1['build'], release2['build'])
        self.assertEquals(release2['version'], 2)
        # check that updating the build rolls a new release
        url = '/v1/apps/{app_id}/builds'.format(**locals())
        build_config = json.dumps({'PATH': 'bin:/usr/local/bin:/usr/bin:/bin'})
        body = {'image': 'autotest/example'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data['image'], body['image'])
        # check to see that a new release was created
        url = '/v1/apps/{app_id}/releases/v3'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release3 = response.data
        self.assertNotEqual(release2['uuid'], release3['uuid'])
        self.assertNotEqual(release2['build'], release3['build'])
        self.assertEquals(release3['version'], 3)
        # check that we can fetch a previous release
        url = '/v1/apps/{app_id}/releases/v2'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release2 = response.data
        self.assertNotEqual(release2['uuid'], release3['uuid'])
        self.assertNotEqual(release2['build'], release3['build'])
        self.assertEquals(release2['version'], 2)
        # disallow post/put/patch/delete
        url = '/v1/apps/{app_id}/releases'.format(**locals())
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.put(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.patch(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        return release3

    @mock.patch('requests.post', mock_status_ok)
    def test_response_data(self):
        body = {'id': 'test'}
        response = self.client.post('/v1/apps', json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        body = {'values': json.dumps({'NEW_URL': 'http://localhost:8080/'})}
        config_response = self.client.post('/v1/apps/test/config', json.dumps(body),
                                           content_type='application/json',
                                           HTTP_AUTHORIZATION='token {}'.format(self.token))
        url = '/v1/apps/test/releases/v2'
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        for key in response.data.keys():
            self.assertIn(key, ['uuid', 'owner', 'created', 'updated', 'app', 'build', 'config',
                                'summary', 'version'])
        expected = {
            'owner': self.user.username,
            'app': 'test',
            'build': None,
            'config': config_response.data['uuid'],
            'summary': '{} added NEW_URL'.format(self.user.username),
            'version': 2
        }
        self.assertDictContainsSubset(expected, response.data)

    @mock.patch('requests.post', mock_status_ok)
    def test_release_rollback(self):
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # try to rollback with only 1 release extant, expecting 400
        url = "/v1/apps/{app_id}/releases/rollback/".format(**locals())
        response = self.client.post(url, content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {'detail': 'version cannot be below 0'})
        self.assertEqual(response.get('content-type'), 'application/json')
        # update config to roll a new release
        url = '/v1/apps/{app_id}/config'.format(**locals())
        body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # update the build to roll a new release
        url = '/v1/apps/{app_id}/builds'.format(**locals())
        body = {'image': 'autotest/example'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # rollback and check to see that a 4th release was created
        # with the build and config of release #2
        url = "/v1/apps/{app_id}/releases/rollback/".format(**locals())
        response = self.client.post(url, content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        url = '/v1/apps/{app_id}/releases'.format(**locals())
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 4)
        url = '/v1/apps/{app_id}/releases/v2'.format(**locals())
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release2 = response.data
        self.assertEquals(release2['version'], 2)
        url = '/v1/apps/{app_id}/releases/v4'.format(**locals())
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release4 = response.data
        self.assertEquals(release4['version'], 4)
        self.assertNotEqual(release2['uuid'], release4['uuid'])
        self.assertEqual(release2['build'], release4['build'])
        self.assertEqual(release2['config'], release4['config'])
        # rollback explicitly to release #1 and check that a 5th release
        # was created with the build and config of release #1
        url = "/v1/apps/{app_id}/releases/rollback/".format(**locals())
        body = {'version': 1}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        url = '/v1/apps/{app_id}/releases'.format(**locals())
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 5)
        url = '/v1/apps/{app_id}/releases/v1'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release1 = response.data
        url = '/v1/apps/{app_id}/releases/v5'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        release5 = response.data
        self.assertEqual(release5['version'], 5)
        self.assertNotEqual(release1['uuid'], release5['uuid'])
        self.assertEqual(release1['build'], release5['build'])
        self.assertEqual(release1['config'], release5['config'])
        # check to see that the current config is actually the initial one
        url = "/v1/apps/{app_id}/config".format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['values'], {})
        # rollback to #3 and see that it has the correct config
        url = "/v1/apps/{app_id}/releases/rollback/".format(**locals())
        body = {'version': 3}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        url = "/v1/apps/{app_id}/config".format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        values = response.data['values']
        self.assertIn('NEW_URL1', values)
        self.assertEqual('http://localhost:8080/', values['NEW_URL1'])

    @mock.patch('requests.post', mock_status_ok)
    def test_release_str(self):
        """Test the text representation of a release."""
        release3 = self.test_release()
        release = Release.objects.get(uuid=release3['uuid'])
        self.assertEqual(str(release), "{}-v3".format(release3['app']))

    @mock.patch('requests.post', mock_status_ok)
    def test_release_summary(self):
        """Test the text summary of a release."""
        release3 = self.test_release()
        release = Release.objects.get(uuid=release3['uuid'])
        # check that the release has push and env change messages
        self.assertIn('autotest deployed ', release.summary)

    @mock.patch('requests.post', mock_status_ok)
    @override_settings(DEFAULT_PERMISSIONS_APPS=True)
    def test_admin_can_create_release(self):
        """If a non-user creates an app, an admin should be able to create releases."""
        user = User.objects.get(username='autotest2')
        token = Token.objects.get(user=user).key
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # check that updating config rolls a new release
        url = '/v1/apps/{app_id}/config'.format(**locals())
        body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertIn('NEW_URL1', response.data['values'])
        # check to see that an initial release was created
        url = '/v1/apps/{app_id}/releases'.format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        # account for the config release as well
        self.assertEqual(response.data['count'], 2)

    @mock.patch('requests.post', mock_status_ok)
    def test_unauthorized_user_cannot_modify_release(self):
        """
        An unauthorized user should not be able to modify other releases.

        Since an unauthorized user should not know about the application at all, these
        requests should return a 404.
        """
        app_id = 'autotest'
        base_url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(base_url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        # push a new build
        url = '{base_url}/{app_id}/builds'.format(**locals())
        body = {'image': 'test'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        # update config to roll a new release
        url = '{base_url}/{app_id}/config'.format(**locals())
        body = {'values': json.dumps({'NEW_URL1': 'http://localhost:8080/'})}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        unauthorized_user = User.objects.get(username='autotest2')
        unauthorized_token = Token.objects.get(user=unauthorized_user).key
        # try to rollback
        url = '{base_url}/{app_id}/releases/rollback/'.format(**locals())
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(unauthorized_token))
        self.assertEqual(response.status_code, 403)
