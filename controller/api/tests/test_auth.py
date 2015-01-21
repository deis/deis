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
        for key in response.data.keys():
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

    @override_settings(REGISTRATION_ENABLED=False)
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

    def test_cancel(self):
        """Test that a registered user can cancel her account."""
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
            # try to abuse superuser/staff level perms
            'is_superuser': True,
            'is_staff': True,
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
