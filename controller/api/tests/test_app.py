"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import os.path

from django.test import TestCase
from django.conf import settings

from api.models import App


class AppTest(TestCase):
    """Tests creation of applications"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        # provide mock authentication used for run commands
        settings.SSH_PRIVATE_KEY = '<some-ssh-private-key>'

    def tearDown(self):
        # reset global vars for other tests
        settings.SSH_PRIVATE_KEY = ''

    def test_app(self):
        """
        Test that a user can create, read, update and delete an application
        """
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        self.assertIn('id', response.data)
        self.assertIn('url', response.data)
        self.assertEqual(response.data['url'], '{app_id}.deisapp.local'.format(**locals()))
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/api/apps/{app_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        body = {'id': 'new'}
        response = self.client.patch(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 405)
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_app_override_id(self):
        body = {'id': 'myid'}
        response = self.client.post('/api/apps', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        body = {'id': response.data['id']}
        response = self.client.post('/api/apps', json.dumps(body),
                                    content_type='application/json')
        self.assertContains(response, 'App with this Id already exists.', status_code=400)
        return response

    def test_app_actions(self):
        url = '/api/apps'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        # test logs
        if not os.path.exists(settings.DEIS_LOG_DIR):
            os.mkdir(settings.DEIS_LOG_DIR)
        path = os.path.join(settings.DEIS_LOG_DIR, app_id + '.log')
        # HACK: remove app lifecycle logs
        if os.path.exists(path):
            os.remove(path)
        url = '/api/apps/{app_id}/logs'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 204)
        self.assertEqual(response.data, 'No logs for {}'.format(app_id))
        # write out some fake log data and try again
        with open(path, 'a') as f:
            f.write(FAKE_LOG_DATA)
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, FAKE_LOG_DATA)
        os.remove(path)
        # test run
        url = '/api/apps/{app_id}/run'.format(**locals())
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data[0], 0)
        # delete file for future runs
        os.remove(path)

    def test_app_release_notes_in_logs(self):
        """Verifies that an app's release summary is dumped into the logs."""
        url = '/api/apps'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        path = os.path.join(settings.DEIS_LOG_DIR, app_id + '.log')
        url = '/api/apps/{app_id}/logs'.format(**locals())
        response = self.client.get(url)
        self.assertIn('autotest created initial release', response.data)
        self.assertEqual(response.status_code, 200)
        # delete file for future runs
        os.remove(path)

    def test_app_errors(self):
        app_id = 'autotest-errors'
        url = '/api/apps'
        body = {'id': 'camelCase'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'App IDs can only contain [a-z0-9-]', status_code=400)
        url = '/api/apps'
        body = {'id': 'deis'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, "App IDs cannot be 'deis'", status_code=400)
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        url = '/api/apps/{app_id}'.format(**locals())
        response = self.client.delete(url)
        self.assertEquals(response.status_code, 204)
        for endpoint in ('containers', 'config', 'releases', 'builds'):
            url = '/api/apps/{app_id}/{endpoint}'.format(**locals())
            response = self.client.get(url)
            self.assertEquals(response.status_code, 404)

    def test_app_structure_is_valid_json(self):
        """Application structures should be valid JSON objects."""
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        self.assertIn('structure', response.data)
        self.assertEqual(response.data['structure'], {})
        app = App.objects.get(id=app_id)
        app.structure = {'web': 1}
        app.save()
        url = '/api/apps/{}'.format(app_id)
        response = self.client.get(url)
        self.assertIn('structure', response.data)
        self.assertEqual(response.data['structure'], {"web": 1})

    def test_admin_can_manage_other_apps(self):
        """Administrators of Deis should be able to manage all applications.
        """
        # log in as non-admin user and create an app
        self.assertTrue(
            self.client.login(username='autotest2', password='password'))
        app_id = 'autotest'
        url = '/api/apps'
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        # log in as admin, check to see if they have access
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        url = '/api/apps/{}'.format(app_id)
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        # check app logs
        url = '/api/apps/{app_id}/logs'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertIn('autotest2 created initial release', response.data)
        # run one-off commands
        url = '/api/apps/{app_id}/run'.format(**locals())
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data[0], 0)
        # delete the app
        url = '/api/apps/{}'.format(app_id)
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_admin_can_see_other_apps(self):
        """If a user creates an application, the administrator should be able
        to see it.
        """
        # log in as non-admin user and create an app
        self.assertTrue(
            self.client.login(username='autotest2', password='password'))
        app_id = 'autotest'
        url = '/api/apps'
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        # log in as admin
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        response = self.client.get(url)
        self.assertEqual(response.data['count'], 1)

    def test_run_without_auth(self):
        """If the administrator has not provided SSH private key for run commands,
        make sure a friendly error message is provided on run"""
        settings.SSH_PRIVATE_KEY = ''
        url = '/api/apps'
        body = {'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        # test run
        url = '/api/apps/{app_id}/run'.format(**locals())
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEquals(response.status_code, 400)
        self.assertEquals(response.data, 'Support for admin commands is not configured')


FAKE_LOG_DATA = """
2013-08-15 12:41:25 [33454] [INFO] Starting gunicorn 17.5
2013-08-15 12:41:25 [33454] [INFO] Listening at: http://0.0.0.0:5000 (33454)
2013-08-15 12:41:25 [33454] [INFO] Using worker: sync
2013-08-15 12:41:25 [33457] [INFO] Booting worker with pid 33457
"""
