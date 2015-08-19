
from __future__ import unicode_literals
import json

from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token


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
        url = '/v1/auth/register'
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
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])

    def test_list(self):
        submit = {
            'username': 'firstuser',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        user = User.objects.get(username='firstuser')
        token = Token.objects.get(user=user).key
        response = self.client.get('/v1/admin/perms', content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(token))
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
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        user = User.objects.get(username='seconduser')
        token = Token.objects.get(user=user).key
        response = self.client.get('/v1/admin/perms', content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 403)
        self.assertIn('You do not have permission', response.data['detail'])

    def test_create(self):
        submit = {
            'username': 'first',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        submit = {
            'username': 'second',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        user = User.objects.get(username='first')
        token = Token.objects.get(user=user).key
        # grant user 2 the superuser perm
        url = '/v1/admin/perms'
        body = {'username': 'second'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 201)
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        self.assertIn('second', str(response.data['results']))

    def test_delete(self):
        submit = {
            'username': 'first',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertTrue(response.data['is_superuser'])
        submit = {
            'username': 'second',
            'password': 'password',
            'email': 'autotest@deis.io',
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertFalse(response.data['is_superuser'])
        user = User.objects.get(username='first')
        token = Token.objects.get(user=user).key
        # grant user 2 the superuser perm
        url = '/v1/admin/perms'
        body = {'username': 'second'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 201)
        # revoke the superuser perm
        response = self.client.delete(url + '/second', HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 204)
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        self.assertNotIn('two', str(response.data['results']))


class TestAppPerms(TestCase):

    fixtures = ['test_sharing.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest-1')
        self.token = Token.objects.get(user=self.user).key
        self.user2 = User.objects.get(username='autotest-2')
        self.token2 = Token.objects.get(user=self.user2).key
        self.user3 = User.objects.get(username='autotest-3')
        self.token3 = Token.objects.get(user=self.user3).key

    def test_create(self):
        # check that user 1 sees her lone app and user 2's app
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        app_id = response.data['results'][0]['id']
        # check that user 2 can only see his app
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(len(response.data['results']), 1)
        # check that user 2 can't see any of the app's builds, configs,
        # containers, limits, or releases
        for model in ['builds', 'config', 'containers', 'releases']:
            response = self.client.get("/v1/apps/{}/{}/".format(app_id, model),
                                       HTTP_AUTHORIZATION='token {}'.format(self.token2))
            msg = "Failed: status '%s', and data '%s'" % (response.status_code, response.data)
            self.assertEqual(response.status_code, 403, msg=msg)
            self.assertEqual(response.data['detail'],
                             'You do not have permission to perform this action.', msg=msg)
        # TODO: test that git pushing to the app fails
        # give user 2 permission to user 1's app
        url = "/v1/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # check that user 2 can see the app
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        # check that user 2 sees (empty) results now for builds, containers,
        # and releases. (config and limit will still give 404s since we didn't
        # push a build here.)
        for model in ['builds', 'containers', 'releases']:
            response = self.client.get("/v1/apps/{}/{}/".format(app_id, model),
                                       HTTP_AUTHORIZATION='token {}'.format(self.token2))
            self.assertEqual(len(response.data['results']), 0)
        # TODO:  check that user 2 can git push the app

    def test_create_errors(self):
        # check that user 1 sees her lone app
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token))
        app_id = response.data['results'][0]['id']
        # check that user 2 can't create a permission
        url = "/v1/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 403)

    def test_delete(self):
        # give user 2 permission to user 1's app
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token))
        app_id = response.data['results'][0]['id']
        url = "/v1/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # check that user 2 can see the app as well as his own
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        # delete permission to user 1's app
        url = "/v1/apps/{}/perms/{}".format(app_id, 'autotest-2')
        response = self.client.delete(url, content_type='application/json',
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)
        self.assertIsNone(response.data)
        # check that user 2 can only see his app
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(len(response.data['results']), 1)
        # delete permission to user 1's app again, expecting an error
        response = self.client.delete(url, content_type='application/json',
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 403)

    def test_list(self):
        # check that user 1 sees her lone app and user 2's app
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)
        app_id = response.data['results'][0]['id']
        # create a new object permission
        url = "/v1/apps/{}/perms".format(app_id)
        body = {'username': 'autotest-2'}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # list perms on the app
        response = self.client.get(
            "/v1/apps/{}/perms".format(app_id), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.data, {'users': ['autotest-2']})

    def test_admin_can_list(self):
        """Check that an administrator can list an app's perms"""
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 2)

    def test_list_errors(self):
        response = self.client.get('/v1/apps', HTTP_AUTHORIZATION='token {}'.format(self.token))
        app_id = response.data['results'][0]['id']
        # login as user 2, list perms on the app
        response = self.client.get(
            "/v1/apps/{}/perms".format(app_id), content_type='application/json',
            HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 403)

    def test_unauthorized_user_cannot_modify_perms(self):
        """
        An unauthorized user should not be able to modify other apps' permissions.

        Since an unauthorized user should not know about the application at all, these
        requests should return a 404.
        """
        app_id = 'autotest'
        url = '/v1/apps'
        body = {'id': app_id}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        unauthorized_user = self.user2
        unauthorized_token = self.token2
        url = '{}/{}/perms'.format(url, app_id)
        body = {'username': unauthorized_user.username}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(unauthorized_token))
        self.assertEqual(response.status_code, 403)

    def test_collaborator_cannot_share(self):
        """
        An collaborator should not be able to modify the app's permissions.
        """
        app_id = "autotest-1-app"
        owner_token = self.token
        collab = self.user2
        collab_token = self.token2
        url = '/v1/apps/{}/perms'.format(app_id)
        # Share app with collaborator
        body = {'username': collab.username}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(owner_token))
        self.assertEqual(response.status_code, 201)
        # Collaborator should fail to share app
        body = {'username': self.user3.username}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(collab_token))
        self.assertEqual(response.status_code, 403)
        # Collaborator can list
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(collab_token))
        self.assertEqual(response.status_code, 200)
        # Share app with user 3 for rest of tests
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(owner_token))
        self.assertEqual(response.status_code, 201)
        response = self.client.get(url, content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(collab_token))
        self.assertEqual(response.status_code, 200)
        # Collaborator cannot delete other collaborator
        url += "/{}".format(self.user3.username)
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(collab_token))
        self.assertEqual(response.status_code, 403)
        # Collaborator can delete themselves
        url = '/v1/apps/{}/perms/{}'.format(app_id, collab.username)
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(collab_token))
        self.assertEqual(response.status_code, 204)
