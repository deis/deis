"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import os.path

from django.test import TestCase
from django.test.utils import override_settings

from deis import settings


@override_settings(CELERY_ALWAYS_EAGER=True)
class AppTest(TestCase):

    """Tests creation of applications"""

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
        body = {'id': 'autotest', 'provider': 'autotest',
                'params': json.dumps({'region': 'us-west-2', 'instance_size': 'm1.medium'})}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        formation_id = 'autotest'
        response = self.client.post('/api/formations', json.dumps({'id': formation_id}),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create & scale a basic formation
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'proxy', 'flavor': 'autotest', 'proxy': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'runtime': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'proxy': 1, 'runtime': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)

    def test_app(self):
        """
        Test that a user can create, read, update and delete an application
        """
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        self.assertIn('formation', response.data)
        self.assertIn('id', response.data)
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

    def test_app_cm(self):
        """
        Test that configuration management is updated on app changes
        """
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        path = os.path.join(settings.TEMPDIR, 'app-{}'.format(app_id))
        with open(path) as f:
            data = json.loads(f.read())
        self.assertIn('id', data)
        self.assertEquals(data['id'], app_id)
        url = '/api/apps/{app_id}'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)
        self.assertFalse(os.path.exists(path))
        formation_id = 'autotest'
        path = os.path.join(settings.TEMPDIR, 'formation-{}'.format(formation_id))
        with open(path) as f:
            data = json.loads(f.read())
        self.assertNotIn(app_id, data['apps'])

    def test_app_override_id(self):
        body = {'formation': 'autotest', 'id': 'myid'}
        response = self.client.post('/api/apps', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        body = {'formation': response.data['formation'], 'id': response.data['id']}
        response = self.client.post('/api/apps', json.dumps(body),
                                    content_type='application/json')
        self.assertContains(response, 'App with this Id already exists.', status_code=400)
        return response

    def test_app_default_formation(self):
        # delete formation created in setUp
        response = self.client.delete('/api/formations/autotest')
        self.assertEqual(response.status_code, 204)
        # try creating an app with no formation specified
        url = '/api/apps'
        response = self.client.post(url)
        self.assertContains(response, 'No formations available', status_code=400)
        # create a formation
        formation1 = 'autotest'
        response = self.client.post('/api/formations', json.dumps({'id': formation1}),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # try again to create an app with no formation specified
        url = '/api/apps'
        response = self.client.post(url)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(formation1, response.data['formation'])
        # create a second formation
        formation2 = 'autotest2'
        response = self.client.post('/api/formations', json.dumps({'id': formation2}),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # create another app with no formation specified
        url = '/api/apps'
        response = self.client.post(url)
        self.assertContains(response, 'Could not determine default formation', status_code=400)

    def test_multiple_apps(self):
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app1_id = response.data['id']
        # test single app domain
        url = "/api/apps/{app1_id}/calculate".format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['domains'], ['localhost.localdomain.local'])
        # create second app without multi-app support
        url = '/api/apps'
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'Formation does not support multiple apps', status_code=400)
        # add domain for multi-app support
        url = '/api/formations/autotest'
        response = self.client.patch(url, json.dumps({'domain': 'deisapp.local'}),
                                     content_type='application/json')
        self.assertEqual(response.status_code, 200)
        # create second app
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app2_id = response.data['id']
        # test multiple app domains
        url = "/api/apps/{app1_id}/calculate".format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['domains'], ['{}.deisapp.local'.format(app1_id)])
        url = "/api/apps/{app2_id}/calculate".format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['domains'], ['{}.deisapp.local'.format(app2_id)])

    def test_app_actions(self):
        url = '/api/apps'
        body = {'formation': 'autotest', 'id': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        # test logs
        if not os.path.exists(settings.DEIS_LOG_DIR):
            os.mkdir(settings.DEIS_LOG_DIR)
        path = os.path.join(settings.DEIS_LOG_DIR, app_id + '.log')
        if os.path.exists(path):
            os.remove(path)
        url = '/api/apps/{app_id}/logs'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, 'No logs for {}'.format(app_id))
        # write out some fake log data and try again
        with open(path, 'w') as f:
            f.write(FAKE_LOG_DATA)
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, FAKE_LOG_DATA)
        # test run with mock error
        url = '/api/apps/{app_id}/run'.format(**locals())
        body = {'command': 'error'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertIn('run `git push deis master` first', response.data[0])
        self.assertEqual(response.data[1], 1)
        # test run
        url = '/api/apps/{app_id}/run'.format(**locals())
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertIn('drwx------  2 deis deis 4096 Dec 21 10:00 .chef', response.data[0])
        self.assertEqual(response.data[1], 0)
        # test calculate
        url = '/api/apps/{app_id}/calculate'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        databag = response.data
        self.assertIn('release', databag)
        self.assertIn('version', databag['release'])
        self.assertIn('containers', databag)

    def test_app_errors(self):
        formation_id, app_id = 'autotest', 'autotest-errors'
        url = '/api/apps'
        body = {'formation': formation_id, 'id': 'camelCase'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'App IDs can only contain [a-z0-9-]', status_code=400)
        body = {'formation': formation_id, 'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'proxy': 0, 'runtime': 0}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/apps/{app_id}/run'.format(**locals())
        body = {'command': 'ls -al'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertContains(response, 'No nodes available to run command', status_code=400)
        url = '/api/apps/{app_id}'.format(**locals())
        response = self.client.delete(url)
        self.assertEquals(response.status_code, 204)
        for endpoint in ('containers', 'config', 'releases', 'builds'):
            url = '/api/apps/{app_id}/{endpoint}'.format(**locals())
            response = self.client.get(url)
            self.assertEquals(response.status_code, 404)


FAKE_LOG_DATA = """
2013-08-15 12:41:25 [33454] [INFO] Starting gunicorn 17.5
2013-08-15 12:41:25 [33454] [INFO] Listening at: http://0.0.0.0:5000 (33454)
2013-08-15 12:41:25 [33454] [INFO] Using worker: sync
2013-08-15 12:41:25 [33457] [INFO] Booting worker with pid 33457
"""
