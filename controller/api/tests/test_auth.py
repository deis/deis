"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json
import urllib

from django.contrib.auth.models import User
from django.test import TestCase
from django.test.utils import override_settings
from rest_framework.authtoken.models import Token


class AuthTest(TestCase):

    fixtures = ['test_auth.json']

    """Tests user registration, authentication and authorization"""

    def setUp(self):
        self.admin = User.objects.get(username='autotest')
        self.admin_token = Token.objects.get(user=self.admin).key
        self.user1 = User.objects.get(username='autotest2')
        self.user1_token = Token.objects.get(user=self.user1).key
        self.user2 = User.objects.get(username='autotest3')
        self.user2_token = Token.objects.get(user=self.user2).key

    def test_auth(self):
        """
        Test that a user can register using the API, login and logout
        """
        # test registration workflow
        username, password = 'newuser', 'password'
        first_name, last_name = 'Otto', 'Test'
        email = 'autotest@deis.io'
        submit = {
            'username': username,
            'password': password,
            'first_name': first_name,
            'last_name': last_name,
            'email': email,
            # try to abuse superuser/staff level perms (not the first signup!)
            'is_superuser': True,
            'is_staff': True,
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        for key in response.data:
            self.assertIn(key, ['id', 'last_login', 'is_superuser', 'username', 'first_name',
                                'last_name', 'email', 'is_active', 'is_superuser', 'is_staff',
                                'date_joined', 'groups', 'user_permissions'])
        expected = {
            'username': username,
            'email': email,
            'first_name': first_name,
            'last_name': last_name,
            'is_active': True,
            'is_superuser': False,
            'is_staff': False
        }
        self.assertDictContainsSubset(expected, response.data)
        # test login
        url = '/v1/auth/login/'
        payload = urllib.urlencode({'username': username, 'password': password})
        response = self.client.post(url, data=payload,
                                    content_type='application/x-www-form-urlencoded')
        self.assertEqual(response.status_code, 200)

    @override_settings(REGISTRATION_MODE="disabled")
    def test_auth_registration_disabled(self):
        """test that a new user cannot register when registration is disabled."""
        url = '/v1/auth/register'
        submit = {
            'username': 'testuser',
            'password': 'password',
            'first_name': 'test',
            'last_name': 'user',
            'email': 'test@user.com',
            'is_superuser': False,
            'is_staff': False,
        }
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 403)

    @override_settings(REGISTRATION_MODE="admin_only")
    def test_auth_registration_admin_only_fails_if_not_admin(self):
        """test that a non superuser cannot register when registration is admin only."""
        url = '/v1/auth/register'
        submit = {
            'username': 'testuser',
            'password': 'password',
            'first_name': 'test',
            'last_name': 'user',
            'email': 'test@user.com',
            'is_superuser': False,
            'is_staff': False,
        }
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 403)

    @override_settings(REGISTRATION_MODE="admin_only")
    def test_auth_registration_admin_only_works(self):
        """test that a superuser can register when registration is admin only."""
        url = '/v1/auth/register'

        username, password = 'newuser_by_admin', 'password'
        first_name, last_name = 'Otto', 'Test'
        email = 'autotest@deis.io'

        submit = {
            'username': username,
            'password': password,
            'first_name': first_name,
            'last_name': last_name,
            'email': email,
            # try to abuse superuser/staff level perms (not the first signup!)
            'is_superuser': True,
            'is_staff': True,
        }
        response = self.client.post(url, json.dumps(submit), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.admin_token))

        self.assertEqual(response.status_code, 201)
        for key in response.data:
            self.assertIn(key, ['id', 'last_login', 'is_superuser', 'username', 'first_name',
                                'last_name', 'email', 'is_active', 'is_superuser', 'is_staff',
                                'date_joined', 'groups', 'user_permissions'])
        expected = {
            'username': username,
            'email': email,
            'first_name': first_name,
            'last_name': last_name,
            'is_active': True,
            'is_superuser': False,
            'is_staff': False
        }
        self.assertDictContainsSubset(expected, response.data)
        # test login
        url = '/v1/auth/login/'
        payload = urllib.urlencode({'username': username, 'password': password})
        response = self.client.post(url, data=payload,
                                    content_type='application/x-www-form-urlencoded')
        self.assertEqual(response.status_code, 200)

    @override_settings(REGISTRATION_MODE="not_a_mode")
    def test_auth_registration_fails_with_nonexistant_mode(self):
        """test that a registration should fail with a nonexistant mode"""
        url = '/v1/auth/register'
        submit = {
            'username': 'testuser',
            'password': 'password',
            'first_name': 'test',
            'last_name': 'user',
            'email': 'test@user.com',
            'is_superuser': False,
            'is_staff': False,
        }

        try:
            self.client.post(url, json.dumps(submit), content_type='application/json')
        except Exception, e:
            self.assertEqual(str(e), 'not_a_mode is not a valid registation mode')

    def test_cancel(self):
        """Test that a registered user can cancel her account."""
        # test registration workflow
        username, password = 'newuser', 'password'
        submit = {
            'username': username,
            'password': password,
            'first_name': 'Otto',
            'last_name': 'Test',
            'email': 'autotest@deis.io',
            # try to abuse superuser/staff level perms
            'is_superuser': True,
            'is_staff': True,
        }

        other_username, other_password = 'newuser2', 'password'
        other_submit = {
            'username': other_username,
            'password': other_password,
            'first_name': 'Test',
            'last_name': 'Tester',
            'email': 'autotest-2@deis.io',
            'is_superuser': False,
            'is_staff': False,
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)

        # cancel the account
        url = '/v1/auth/cancel'
        user = User.objects.get(username=username)
        token = Token.objects.get(user=user).key
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 204)

        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(other_submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)

        # normal user can't delete another user
        url = '/v1/auth/cancel'
        other_user = User.objects.get(username=other_username)
        other_token = Token.objects.get(user=other_user).key
        response = self.client.delete(url, json.dumps({'username': self.admin.username}),
                                      content_type='application/json',
                                      HTTP_AUTHORIZATION='token {}'.format(other_token))
        self.assertEqual(response.status_code, 403)

        # admin can delete another user
        response = self.client.delete(url, json.dumps({'username': other_username}),
                                      content_type='application/json',
                                      HTTP_AUTHORIZATION='token {}'.format(self.admin_token))
        self.assertEqual(response.status_code, 204)
        # user can not be deleted if it has an app attached to it
        response = self.client.post(
            '/v1/apps',
            HTTP_AUTHORIZATION='token {}'.format(self.admin_token)
        )
        self.assertEqual(response.status_code, 201)
        app_id = response.data['id']  # noqa
        self.assertIn('id', response.data)

        response = self.client.delete(url, json.dumps({'username': str(self.admin)}),
                                      content_type='application/json',
                                      HTTP_AUTHORIZATION='token {}'.format(self.admin_token))
        self.assertEqual(response.status_code, 409)

    def test_passwd(self):
        """Test that a registered user can change the password."""
        # test registration workflow
        username, password = 'newuser', 'password'
        first_name, last_name = 'Otto', 'Test'
        email = 'autotest@deis.io'
        submit = {
            'username': username,
            'password': password,
            'first_name': first_name,
            'last_name': last_name,
            'email': email,
        }
        url = '/v1/auth/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        # change password
        url = '/v1/auth/passwd'
        user = User.objects.get(username=username)
        token = Token.objects.get(user=user).key
        submit = {
            'password': 'password2',
            'new_password': password,
        }
        response = self.client.post(url, json.dumps(submit), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {'detail': 'Current password does not match'})
        self.assertEqual(response.get('content-type'), 'application/json')
        submit = {
            'password': password,
            'new_password': 'password2',
        }
        response = self.client.post(url, json.dumps(submit), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 200)
        # test login with old password
        url = '/v1/auth/login/'
        payload = urllib.urlencode({'username': username, 'password': password})
        response = self.client.post(url, data=payload,
                                    content_type='application/x-www-form-urlencoded')
        self.assertEqual(response.status_code, 400)
        # test login with new password
        payload = urllib.urlencode({'username': username, 'password': 'password2'})
        response = self.client.post(url, data=payload,
                                    content_type='application/x-www-form-urlencoded')
        self.assertEqual(response.status_code, 200)

    def test_change_user_passwd(self):
        """
        Test that an administrator can change a user's password, while a regular user cannot.
        """
        # change password
        url = '/v1/auth/passwd'
        old_password = self.user1.password
        new_password = 'password'
        submit = {
            'username': self.user1.username,
            'new_password': new_password,
        }
        response = self.client.post(url, json.dumps(submit), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.admin_token))
        self.assertEqual(response.status_code, 200)
        # test login with old password
        url = '/v1/auth/login/'
        payload = urllib.urlencode({'username': self.user1.username, 'password': old_password})
        response = self.client.post(url, data=payload,
                                    content_type='application/x-www-form-urlencoded')
        self.assertEqual(response.status_code, 400)
        # test login with new password
        payload = urllib.urlencode({'username': self.user1.username, 'password': new_password})
        response = self.client.post(url, data=payload,
                                    content_type='application/x-www-form-urlencoded')
        self.assertEqual(response.status_code, 200)
        # Non-admins can't change another user's password
        submit['password'], submit['new_password'] = submit['new_password'], old_password
        url = '/v1/auth/passwd'
        response = self.client.post(url, json.dumps(submit), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.user2_token))
        self.assertEqual(response.status_code, 403)
        # change back password with a regular user
        response = self.client.post(url, json.dumps(submit), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.user1_token))
        self.assertEqual(response.status_code, 200)
        # test login with new password
        url = '/v1/auth/login/'
        payload = urllib.urlencode({'username': self.user1.username, 'password': old_password})
        response = self.client.post(url, data=payload,
                                    content_type='application/x-www-form-urlencoded')
        self.assertEqual(response.status_code, 200)

    def test_regenerate(self):
        """ Test that token regeneration works"""

        url = '/v1/auth/tokens/'

        response = self.client.post(url, '{}', content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.admin_token))

        self.assertEqual(response.status_code, 200)
        self.assertNotEqual(response.data['token'], self.admin_token)

        self.admin_token = Token.objects.get(user=self.admin)

        response = self.client.post(url, '{"username" : "autotest2"}',
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.admin_token))

        self.assertEqual(response.status_code, 200)
        self.assertNotEqual(response.data['token'], self.user1_token)

        response = self.client.post(url, '{"all" : "true"}',
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.admin_token))
        self.assertEqual(response.status_code, 200)

        response = self.client.post(url, '{}', content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.admin_token))

        self.assertEqual(response.status_code, 401)
