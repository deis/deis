"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import logging
import mock
import requests

from django.conf import settings
from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token

from api.models import App
from . import mock_status_ok


class AppTest(TestCase):
    """Tests creation of applications"""

    fixtures = ['tests.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        # provide mock authentication used for run commands
        settings.SSH_PRIVATE_KEY = '<some-ssh-private-key>'

    def tearDown(self):
        # reset global vars for other tests
        settings.SSH_PRIVATE_KEY = ''

    def test_app(self):
        """
        Test that a user can create, read, update and delete an application
        """
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        self.assertIn('id', response.data)
        response = self.client.get('/v1/apps',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/v1/apps/{app_id}'.format(**locals())
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        body = {'id': 'new'}
        response = self.client.patch(url, json.dumps(body), content_type='application/json',
                                     HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)

    def test_response_data(self):
        """Test that the serialized response contains only relevant data."""
        body = {'id': 'test'}
        response = self.client.post('/v1/apps', json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        for key in response.data:
            self.assertIn(key, ['uuid', 'created', 'updated', 'id', 'owner', 'url', 'structure'])
        expected = {
            'id': 'test',
            'owner': self.user.username,
            'url': 'test.deisapp.local',
            'structure': {}
        }
        self.assertDictContainsSubset(expected, response.data)

    def test_app_override_id(self):
        body = {'id': 'myid'}
        response = self.client.post('/v1/apps', json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        body = {'id': response.data['id']}
        response = self.client.post('/v1/apps', json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertContains(response, 'This field must be unique.', status_code=400)
        return response

    @mock.patch('requests.get')
    def test_app_actions(self, mock_get):
        url = '/v1/apps'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa

        # test logs - 204 from deis-logger
        mock_response = mock.Mock()
        mock_response.status_code = 204
        mock_get.return_value = mock_response
        url = "/v1/apps/{app_id}/logs".format(**locals())
        response = self.client.get(url, HTTP_AUTHORIZATION="token {}".format(self.token))
        self.assertEqual(response.status_code, 204)

        # test logs - 404 from deis-logger
        mock_response.status_code = 404
        response = self.client.get(url, HTTP_AUTHORIZATION="token {}".format(self.token))
        self.assertEqual(response.status_code, 204)

        # test logs - unanticipated status code from deis-logger
        mock_response.status_code = 400
        response = self.client.get(url, HTTP_AUTHORIZATION="token {}".format(self.token))
        self.assertEqual(response.status_code, 500)
        self.assertEqual(response.content, "Error accessing logs for {}".format(app_id))

        # test logs - success accessing deis-logger
        mock_response.status_code = 200
        mock_response.content = FAKE_LOG_DATA
        response = self.client.get(url, HTTP_AUTHORIZATION="token {}".format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.content, FAKE_LOG_DATA)

        # test logs - HTTP request error while accessing deis-logger
        mock_get.side_effect = requests.exceptions.RequestException('Boom!')
        response = self.client.get(url, HTTP_AUTHORIZATION="token {}".format(self.token))
        self.assertEqual(response.status_code, 500)
        self.assertEqual(response.content, "Error accessing logs for {}".format(app_id))

        # TODO: test run needs an initial build

    @mock.patch('api.models.logger')
    def test_app_release_notes_in_logs(self, mock_logger):
        """Verifies that an app's release summary is dumped into the logs."""
        url = '/v1/apps'
        app_name = 'autotest'
        body = {'id': app_name}

        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app = App.objects.get(id=app_name)
        # check app logs
        exp_msg = "[{app_name}]: {self.user.username} created initial release".format(**locals())
        mock_logger.log.assert_called_with(logging.INFO, exp_msg)
        app.log('hello world')
        exp_msg = "[{app_name}]: hello world".format(**locals())
        mock_logger.log.assert_called_with(logging.INFO, exp_msg)
        app.log('goodbye world', logging.WARNING)
        # assert logging with a different log level
        exp_msg = "[{app_name}]: goodbye world".format(**locals())
        mock_logger.log.assert_called_with(logging.WARNING, exp_msg)

    def test_app_errors(self):
        app_id = 'autotest-errors'
        url = '/v1/apps'
        body = {'id': 'camelCase'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertContains(response, 'App IDs can only contain [a-z0-9-]', status_code=400)
        url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        url = '/v1/apps/{app_id}'.format(**locals())
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEquals(response.status_code, 204)
        for endpoint in ('containers', 'config', 'releases', 'builds'):
            url = '/v1/apps/{app_id}/{endpoint}'.format(**locals())
            response = self.client.get(url,
                                       HTTP_AUTHORIZATION='token {}'.format(self.token))
            self.assertEquals(response.status_code, 404)

    def test_app_reserved_names(self):
        """Nobody should be able to create applications with names which are reserved."""
        url = '/v1/apps'
        reserved_names = ['foo', 'bar']
        with self.settings(DEIS_RESERVED_NAMES=reserved_names):
            for name in reserved_names:
                body = {'id': name}
                response = self.client.post(url, json.dumps(body), content_type='application/json',
                                            HTTP_AUTHORIZATION='token {}'.format(self.token))
                self.assertContains(
                    response,
                    '{} is a reserved name.'.format(name),
                    status_code=400)

    def test_app_structure_is_valid_json(self):
        """Application structures should be valid JSON objects."""
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        self.assertIn('structure', response.data)
        self.assertEqual(response.data['structure'], {})
        app = App.objects.get(id=app_id)
        app.structure = {'web': 1}
        app.save()
        url = '/v1/apps/{}'.format(app_id)
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertIn('structure', response.data)
        self.assertEqual(response.data['structure'], {"web": 1})

    @mock.patch('requests.post', mock_status_ok)
    @mock.patch('api.models.logger')
    def test_admin_can_manage_other_apps(self, mock_logger):
        """Administrators of Deis should be able to manage all applications.
        """
        # log in as non-admin user and create an app
        user = User.objects.get(username='autotest2')
        token = Token.objects.get(user=user)
        app_id = 'autotest'
        url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(token))
        # log in as admin, check to see if they have access
        url = '/v1/apps/{}'.format(app_id)
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        # check app logs
        exp_msg = "autotest2 created initial release"
        exp_log_call = mock.call(logging.INFO, exp_msg)
        mock_logger.log.has_calls(exp_log_call)
        # TODO: test run needs an initial build
        # delete the app
        url = '/v1/apps/{}'.format(app_id)
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)

    def test_admin_can_see_other_apps(self):
        """If a user creates an application, the administrator should be able
        to see it.
        """
        # log in as non-admin user and create an app
        user = User.objects.get(username='autotest2')
        token = Token.objects.get(user=user)
        app_id = 'autotest'
        url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(token))
        # log in as admin
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.data['count'], 1)

    def test_run_without_auth(self):
        """If the administrator has not provided SSH private key for run commands,
        make sure a friendly error message is provided on run"""
        settings.SSH_PRIVATE_KEY = ''
        url = '/v1/apps'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        # test run
        url = '/v1/apps/{app_id}/run'.format(**locals())
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEquals(response.status_code, 400)
        self.assertEquals(response.data, {'detail': 'Support for admin commands '
                                                    'is not configured'})

    def test_run_without_release_should_error(self):
        """
        A user should not be able to run a one-off command unless a release
        is present.
        """
        app_id = 'autotest'
        url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        url = '/v1/apps/{}/run'.format(app_id)
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {'detail': 'No build associated with this '
                                                   'release to run this command'})

    def test_unauthorized_user_cannot_see_app(self):
        """
        An unauthorized user should not be able to access an app's resources.

        Since an unauthorized user can't access the application, these
        tests should return a 403, but currently return a 404. FIXME!
        """
        app_id = 'autotest'
        base_url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(base_url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        unauthorized_user = User.objects.get(username='autotest2')
        unauthorized_token = Token.objects.get(user=unauthorized_user).key
        url = '{}/{}/run'.format(base_url, app_id)
        body = {'command': 'foo'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(unauthorized_token))
        self.assertEqual(response.status_code, 403)
        url = '{}/{}/logs'.format(base_url, app_id)
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(unauthorized_token))
        self.assertEqual(response.status_code, 403)
        url = '{}/{}'.format(base_url, app_id)
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(unauthorized_token))
        self.assertEqual(response.status_code, 403)
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(unauthorized_token))
        self.assertEqual(response.status_code, 403)

    def test_app_info_not_showing_wrong_app(self):
        app_id = 'autotest'
        base_url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(base_url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        url = '{}/foo'.format(base_url)
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 404)

    def test_app_transfer(self):
        owner = User.objects.get(username='autotest2')
        owner_token = Token.objects.get(user=owner).key
        app_id = 'autotest'
        base_url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(base_url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(owner_token))
        # Transfer App
        url = '{}/{}'.format(base_url, app_id)
        new_owner = User.objects.get(username='autotest3')
        new_owner_token = Token.objects.get(user=new_owner).key
        body = {'owner': new_owner.username}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(owner_token))
        self.assertEqual(response.status_code, 200)

        # Original user can no longer access it
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(owner_token))
        self.assertEqual(response.status_code, 403)

        # New owner can access it
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(new_owner_token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['owner'], new_owner.username)

        # Collaborators can't transfer
        body = {'username': owner.username}
        perms_url = url+"/perms/"
        response = self.client.post(perms_url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(new_owner_token))
        self.assertEqual(response.status_code, 201)
        body = {'owner': self.user.username}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(owner_token))
        self.assertEqual(response.status_code, 403)

        # Admins can transfer
        body = {'owner': self.user.username}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['owner'], self.user.username)


FAKE_LOG_DATA = """
2013-08-15 12:41:25 [33454] [INFO] Starting gunicorn 17.5
2013-08-15 12:41:25 [33454] [INFO] Listening at: http://0.0.0.0:5000 (33454)
2013-08-15 12:41:25 [33454] [INFO] Using worker: sync
2013-08-15 12:41:25 [33457] [INFO] Booting worker with pid 33457
"""
