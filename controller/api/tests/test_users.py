
from __future__ import unicode_literals

from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token


class TestUsers(TestCase):
    """ Tests users endpoint"""

    fixtures = ['tests.json']

    def test_super_user_can_list(self):
        url = '/v1/users'

        user = User.objects.get(username='autotest')
        token = Token.objects.get(user=user)

        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(token))

        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 3)

    def test_non_super_user_cannot_list(self):
        url = '/v1/users'

        user = User.objects.get(username='autotest2')
        token = Token.objects.get(user=user)

        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(token))
        self.assertEqual(response.status_code, 403)
