"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""
from __future__ import unicode_literals

from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token

from api import __version__


class APIMiddlewareTest(TestCase):

    """Tests middleware.py's business logic"""

    fixtures = ['tests.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key

    def test_deis_version_header_good(self):
        """
        Test that when the version header is sent, the request is accepted.
        """
        response = self.client.get(
            '/v1/apps',
            HTTP_DEIS_VERSION=__version__.rsplit('.', 2)[0],
            HTTP_AUTHORIZATION='token {}'.format(self.token),
        )
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.has_header('DEIS_API_VERSION'), True)
        self.assertEqual(response['DEIS_API_VERSION'], __version__.rsplit('.', 1)[0])

    def test_deis_version_header_bad(self):
        """
        Test that when an improper version header is sent, the request is declined.
        """
        response = self.client.get(
            '/v1/apps',
            HTTP_DEIS_VERSION='1234.5678',
            HTTP_AUTHORIZATION='token {}'.format(self.token),
        )
        self.assertEqual(response.status_code, 405)

    def test_deis_version_header_not_present(self):
        """
        Test that when the version header is not present, the request is accepted.
        """
        response = self.client.get('/v1/apps',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
