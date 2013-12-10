
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
            'username': 'one',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        submit = {
            'username': 'two',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        self.assertTrue(
            self.client.login(username='one', password='password'))
        # grant user 2 the superuser perm
        url = '/api/admin/perms'
        body = {'username': 'two'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        self.assertIn('two', str(response.data['results']))

    def test_delete(self):
        submit = {
            'username': 'uno',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        submit = {
            'username': 'dos',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/api/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        self.assertTrue(
            self.client.login(username='uno', password='password'))
        # grant user 2 the superuser perm
        url = '/api/admin/perms'
        body = {'username': 'dos'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # revoke the superuser perm
        response = self.client.delete(url + '/dos')
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
        # check that user 1 sees her lone app
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        app_id = response.data['results'][0]['id']
        # check that user 2 can't see any apps
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/apps')
        self.assertEqual(len(response.data['results']), 0)
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
        self.assertEqual(len(response.data['results']), 1)
        # TODO:  check that user 2 can git push the app

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
        # check that user 2 can see the app
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
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
        # check that user 2 can't see any apps
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/apps')
        self.assertEqual(len(response.data['results']), 0)

    def test_list(self):
        # check that user 1 sees her lone app
        response = self.client.get('/api/apps')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
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


class TestFormationPerms(TestCase):

    fixtures = ['test_sharing.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))

    def test_create(self):
        # check that user 1 sees her lone formation
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        formation_id = response.data['results'][0]['id']
        # check that user 2 can't see any formations
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)
        # give user 2 permission to user 1's formation
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))
        url = '/api/formations/{formation_id}/perms'.format(**locals())
        body = {'username': 'autotest-2'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # check that user 1 can see a formation
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertEqual(response.data['results'][0]['id'], formation_id)
        # revoke user 2's permission to user 1's formation
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))
        username = 'autotest-2'
        url = '/api/formations/{formation_id}/perms/{username}/'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_create_errors(self):
        response = self.client.get('/api/formations')
        formation_id = response.data['results'][0]['id']
        # try to give user 2 permission to user 1's formation as user 2
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        url = '/api/formations/{formation_id}/perms'.format(**locals())
        body = {'username': 'autotest-2'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)
        # try to give user 1 permission to user 1's formation as user 2
        url = '/api/formations/{formation_id}/perms'.format(**locals())
        body = {'username': 'autotest-1'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)

    def test_delete(self):
        # give user 2 permission to user 1's formation
        response = self.client.get('/api/formations')
        formation_id = response.data['results'][0]['id']
        url = '/api/formations/{formation_id}/perms'.format(**locals())
        body = {'username': 'autotest-2'}
        response = self.client.post(
            url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # check that user 1 can see a formation
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertEqual(response.data['results'][0]['id'], formation_id)
        # revoke user 2's permission to user 1's formation
        self.assertTrue(
            self.client.login(username='autotest-1', password='password'))
        username = 'autotest-2'
        url = '/api/formations/{formation_id}/perms/{username}/'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)
        # check that user 1 can't see any formations
        self.assertTrue(
            self.client.login(username='autotest-2', password='password'))
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 0)

    def test_delete_errors(self):
        # revoke user 2's permission to user 1's formation
        response = self.client.get('/api/formations')
        formation_id = response.data['results'][0]['id']
        username = 'autotest-2'
        url = '/api/formations/{formation_id}/perms/{username}/'.format(**locals())
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 404)

    def test_list(self):
        # check that user 1 sees her lone formation
        response = self.client.get('/api/formations')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        formation_id = response.data['results'][0]['id']
        # create a new object permission
        url = "/api/formations/{}/perms".format(formation_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # list perms on the app
        response = self.client.get(
            "/api/formations/{}/perms".format(formation_id), content_type='application/json')
        self.assertEqual(response.data, {'users': ['autotest-2']})
