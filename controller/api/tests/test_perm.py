
from __future__ import unicode_literals
import json

from django.test import TestCase


class TestAdminPerms(TestCase):

    def test_first_signup(self):
        # register a first user
        username, password = 'firstuser', 'password'
        email = 'autotest@deis.io'
        submit = {
            'username': username,
            'password': password,
            'email': email,
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        # register a second user
        username, password = 'seconduser', 'password'
        email = 'autotest@deis.io'
        submit = {
            'username': username,
            'password': password,
            'email': email,
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])

    def test_list(self):
        submit = {
            'username': 'firstuser',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        self.assertTrue(
            self.client.login(username='firstuser', password='password'))
        response = self.client.get('/api/admin/perms', content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertEqual(response.data['results'][0]['username'], 'firstuser')
        self.assertTrue(response.data['results'][0]['is_superuser'])
        # register a non-superuser
        submit = {
            'username': 'seconduser',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        self.assertTrue(
            self.client.login(username='seconduser', password='password'))
        response = self.client.get('/api/admin/perms', content_type='application/json')
        self.assertEqual(response.status_code, 403)
        self.assertIn('You do not have permission', response.data['detail'])

    def test_create(self):
        submit = {
            'username': 'first',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        submit = {
            'username': 'second',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        self.assertTrue(
            self.client.login(username='first', password='password'))
        # grant user 2 the superuser perm
        url = '/api/admin/perms'
        body = {'username': 'second'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        self.assertIn('second', str(response.data['results']))

    def test_delete(self):
        submit = {
            'username': 'first',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        submit = {
            'username': 'second',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        self.assertTrue(
            self.client.login(username='first', password='password'))
        # grant user 2 the superuser perm
        url = '/api/admin/perms'
        body = {'username': 'second'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # revoke the superuser perm
        response = self.client.delete(url + '/second')
        self.assertEqual(response.status_code, 204)
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertNotIn('two', str(response.data['results']))


class TestAppPerms(TestCase):

    fixtures = ['test_sharing.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))

    def test_create(self):
        # check that user 1 sees her lone app and user 2's app
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        app_id = response.data['results'][0]['id']
        # check that user 2 can only see his app
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/apps')
        self.assertEqual(len(response.data['results']), 1)
        # check that user 2 can't see any of the app's builds, configs,
        # containers, limits, or releases
        for model in ['builds', 'config', 'containers', 'limits', 'releases']:
            response = self.client.get("/api/apps/{}/{}/".format(app_id, model))
            self.assertEqual(response.data['detail'], 'Not found')
        # TODO: test that git pushing to the app fails
        # give user 2 permission to user 1's app
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))
        url = "/api/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # check that user 2 can see the app
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        # check that user 2 sees (empty) results now for builds, containers,
        # and releases. (config and limit will still give 404s since we didn't
        # push a build here.)
        for model in ['builds', 'containers', 'releases']:
            response = self.client.get("/api/apps/{}/{}/".format(app_id, model))
            self.assertEqual(len(response.data['results']), 0)
        # TODO:  check that user 2 can git push the app

    def test_create_errors(self):
        # check that user 1 sees her lone app
        response = self.client.get('/api/apps')
        app_id = response.data['results'][0]['id']
        # check that user 2 can't create a permission
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        url = "/api/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)

    def test_delete(self):
        # give user 2 permission to user 1's app
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))
        response = self.client.get('/api/apps')
        app_id = response.data['results'][0]['id']
        url = "/api/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # check that user 2 can see the app as well as his own
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        # try to delete the permission as user 2
        url = "/api/apps/{}/perms/{}".format(app_id, 'autotest-2')
        response = self.client.delete(url, content_type='application/json')
        self.assertEqual(response.status_code, 403)
        self.assertIsNone(response.data)
        # delete permission to user 1's app
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))
        response = self.client.delete(url, content_type='application/json')
        self.assertEqual(response.status_code, 204)
        self.assertIsNone(response.data)
        # check that user 2 can only see his app
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/apps')
        self.assertEqual(len(response.data['results']), 1)
        # delete permission to user 1's app again, expecting an error
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))
        response = self.client.delete(url, content_type='application/json')
        self.assertEqual(response.status_code, 404)

    def test_list(self):
        # check that user 1 sees her lone app and user 2's app
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        app_id = response.data['results'][0]['id']
        # create a new object permission
        url = "/api/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # list perms on the app
        response = self.client.get(
            "/api/apps/{}/perms".format(app_id), content_type='application/json')
        self.assertEqual(response.data, {'users': ['autotest-2']})

    def test_admin_can_list(self):
        """Check that an administrator can list an app's perms"""
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)

    def test_list_errors(self):
        response = self.client.get('/api/apps')
        app_id = response.data['results'][0]['id']
        # login as user 2
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        # list perms on the app
        response = self.client.get(
            "/api/apps/{}/perms".format(app_id), content_type='application/json')
        self.assertEqual(response.status_code, 403)
