"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import uuid

from django.test import TestCase
from django.test.utils import override_settings

from deis import settings


@override_settings(CELERY_ALWAYS_EAGER=True)
class HookTest(TestCase):

    """Tests API hooks used to trigger actions from external components"""

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

    def test_push_hook(self):
        """Test creating a Push via the API"""
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # prepare a push body
        body = {
            'sha': 'df1e628f2244b73f9cdf944f880a2b3470a122f4',
            'fingerprint': '88:25:ed:67:56:91:3d:c6:1b:7f:42:c6:9b:41:24:80',
            'receive_user': 'autotest',
            'receive_repo': '{app_id}'.format(**locals()),
            'ssh_connection': '10.0.1.10 50337 172.17.0.143 22',
            'ssh_original_command': "git-receive-pack '{app_id}.git'".format(**locals()),
        }
        # post a request without the auth header
        url = "/api/hooks/push".format(**locals())
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)
        # now try with the builder key in the special auth header
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_X_DEIS_BUILDER_AUTH=settings.BUILDER_KEY)
        self.assertEqual(response.status_code, 201)
        for k in ('owner', 'app', 'sha', 'fingerprint', 'receive_repo', 'receive_user',
                  'ssh_connection', 'ssh_original_command'):
            self.assertIn(k, response.data)

    def test_push_abuse(self):
        """Test a user pushing to an unauthorized application"""
        # create a legit app as "autotest"
        url = '/api/apps'
        body = {'formation': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        # register an evil user
        username, password = 'eviluser', 'password'
        first_name, last_name = 'Evil', 'User'
        email = 'evil@deis.io'
        submit = {
            'username': username,
            'password': password,
            'first_name': first_name,
            'last_name': last_name,
            'email': email,
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # prepare a push body that simulates a git push
        body = {
            'sha': 'df1e628f2244b73f9cdf944f880a2b3470a122f4',
            'fingerprint': '88:25:ed:67:56:91:3d:c6:1b:7f:42:c6:9b:41:24:99',
            'receive_user': 'eviluser',
            'receive_repo': '{app_id}'.format(**locals()),
            'ssh_connection': '10.0.1.10 50337 172.17.0.143 22',
            'ssh_original_command': "git-receive-pack '{app_id}.git'".format(**locals()),
        }
        # try to push as "eviluser"
        url = "/api/hooks/push".format(**locals())
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_X_DEIS_BUILDER_AUTH=settings.BUILDER_KEY)
        self.assertEqual(response.status_code, 403)

    def test_build_hook(self):
        """Test creating a Build via the API"""
        formation_id = 'autotest'
        url = '/api/formations/{formation_id}/layers'.format(**locals())
        body = {'id': 'runtime', 'flavor': 'autotest', 'runtime': True, 'proxy': True}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/formations/{formation_id}/scale'.format(**locals())
        body = {'runtime': 2}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        url = '/api/apps'
        body = {'formation': formation_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        build = {'username': 'autotest', 'app': app_id}
        url = '/api/hooks/builds'.format(**locals())
        sha, checksum = uuid.uuid4().hex, uuid.uuid4().hex
        body = {'receive_user': 'autotest',
                'receive_repo': app_id,
                'sha': sha,
                'checksum': checksum,
                'procfile': {'web': 'node server.js'},
                'config': {'PATH': '/usr/local/bin:/usr/bin:/usr/sbin'},
                'url': 'http://deis-controller.local/slugs/{app_id}-{sha}.tar.gz'.format(**locals()),
                'size': 12345}
        # post the build without a session
        self.assertIsNone(self.client.logout())
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)
        # post the build with the builder auth key
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_X_DEIS_BUILDER_AUTH=settings.BUILDER_KEY)
        self.assertEqual(response.status_code, 201)
        build = response.data
        self.assertIn('sha', response.data)
        self.assertIn('procfile', response.data)
        procfile = json.loads(response.data['procfile'])
        self.assertIn('web', procfile)
        self.assertEqual(procfile['web'], 'node server.js')
        # calculate the databag
        self.assertTrue(
            self.client.login(username='autotest', password='password'))
        url = '/api/apps/{app_id}/calculate'.format(**locals())
        response = self.client.post(url)
        self.assertEqual(response.status_code, 200)
        databag = response.data
        self.assertIn('release', databag)
        self.assertIn('version', databag['release'])
        self.assertIn('containers', databag)
        self.assertIn('web', databag['containers'])
        self.assertIn('1', databag['containers']['web'])
        self.assertEqual(databag['containers']['web']['1'], 'up')

