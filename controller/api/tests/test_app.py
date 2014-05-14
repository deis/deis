"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import os.path

from django.test import TestCase
from django.test.utils import override_settings

from django.conf import settings


@override_settings(CELERY_ALWAYS_EAGER=True)
class AppTest(TestCase):

    """Tests creation of applications"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        body = {'id': 'autotest', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': {}}
        response = self.client.post('/api/clusters', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_app(self):
        """
        Test that a user can create, read, update and delete an application
        """
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        appname = response.data['name']  # noqa
        self.assertIn('cluster', response.data)
        self.assertIn('name', response.data)
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        url = '/api/apps/{appname}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        body = {'name': 'new'}
        response = self.client.patch(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/apps/new'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_app_override_name(self):
        body = {'cluster': 'autotest', 'name': 'myname'}
        response = self.client.post('/api/apps', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        body = {'cluster': response.data['cluster'], 'name': response.data['name']}
        response = self.client.post('/api/apps', json.dumps(body),
                                    content_type='application/json')
        self.assertContains(response, 'App with this Name already exists.', status_code=400)
        return response

    def test_app_actions(self):
        url = '/api/apps'
        body = {'cluster': 'autotest', 'name': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        appname = response.data['name']  # noqa
        # test logs
        if not os.path.exists(settings.DEIS_LOG_DIR):
            os.mkdir(settings.DEIS_LOG_DIR)
        path = os.path.join(settings.DEIS_LOG_DIR, appname + '.log')
        if os.path.exists(path):
            os.remove(path)
        url = '/api/apps/{appname}/logs'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, 'No logs for {}'.format(appname))
        # write out some fake log data and try again
        with open(path, 'w') as f:
            f.write(FAKE_LOG_DATA)
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, FAKE_LOG_DATA)
        # test run
        url = '/api/apps/{appname}/run'.format(**locals())
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data[0], 0)

    def test_app_errors(self):
        cluster_id, appname = 'autotest', 'autotest-errors'
        url = '/api/apps'
        body = {'cluster': cluster_id, 'name': 'camelCase'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'App names can only contain [a-z0-9-]', status_code=400)
        body = {'cluster': cluster_id, 'name': appname}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        appname = response.data['name']  # noqa
        url = '/api/apps/{appname}'.format(**locals())
        response = self.client.delete(url)
        self.assertEquals(response.status_code, 204)
        for endpoint in ('containers', 'config', 'releases', 'builds'):
            url = '/api/apps/{appname}/{endpoint}'.format(**locals())
            response = self.client.get(url)
            self.assertEquals(response.status_code, 404)


FAKE_LOG_DATA = """
2013-08-15 12:41:25 [33454] [INFO] Starting gunicorn 17.5
2013-08-15 12:41:25 [33454] [INFO] Listening at: http://0.0.0.0:5000 (33454)
2013-08-15 12:41:25 [33454] [INFO] Using worker: sync
2013-08-15 12:41:25 [33457] [INFO] Booting worker with pid 33457
"""
