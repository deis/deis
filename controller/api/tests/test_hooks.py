"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

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
        body = {'id': 'autotest', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': {}}
        response = self.client.post('/api/clusters', json.dumps(body),
                                    content_type='application/json')
        self.assertEqual(response.status_code, 201)

    def test_push_hook(self):
        """Test creating a Push via the API"""
        url = '/api/apps'
        body = {'cluster': 'autotest'}
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
        body = {'cluster': 'autotest'}
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
        """Test creating a Build via an API Hook"""
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']
        build = {'username': 'autotest', 'app': app_id}
        url = '/api/hooks/builds'.format(**locals())
        body = {'receive_user': 'autotest',
                'receive_repo': app_id,
                'image': 'registry.local:5000/autotest/{app_id}:v2'.format(**locals())}
        # post the build without a session
        self.assertIsNone(self.client.logout())
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)
        # post the build with the builder auth key
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_X_DEIS_BUILDER_AUTH=settings.BUILDER_KEY)
        self.assertEqual(response.status_code, 200)
        self.assertIn('release', response.data)
        self.assertIn('version', response.data['release'])
        self.assertIn('domains', response.data)
