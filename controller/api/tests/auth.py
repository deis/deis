"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase


class AuthTest(TestCase):

    """Tests user registration, authentication and authorization"""

    def test_auth(self):
        """
        Test that a user can register using the API, login and logout
        """
        # make sure logging in with an invalid username/password
        # results in a 404
        # post credentials to the login URL
        url = '/api/auth/login'
        body = {'username': 'fail', 'password': 'this'}
        response = self.client.post(url, data=json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 404)
        # test registration workflow
        username, password = 'newuser', 'password'
        first_name, last_name = 'Otto', 'Test'
        email = 'autotest@deis.io'
        submit = {'username': username, 'password': password, 
                  'first_name': first_name, 'last_name': last_name,
                  'email': email,
                  # try to abuse superuser/staff level perms
                  'is_superuser': True, 'is_staff': True}
        url = '/api/register'
        response = self.client.post(url, json.dumps(submit), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data['username'], username)
        self.assertNotIn('password', response.data)
        self.assertEqual(response.data['email'], email)
        self.assertEqual(response.data['first_name'], first_name)
        self.assertEqual(response.data['last_name'], last_name)
        self.assertFalse(response.data['is_superuser'])
        self.assertFalse(response.data['is_staff'])
        self.assertTrue(
            self.client.login(username=username, password=password))
