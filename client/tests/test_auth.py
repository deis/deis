"""
Unit tests for the Deis CLI auth commands.

Run these tests with "python -m unittest client.tests.test_auth"
or with "./manage.py test client.AuthTest".
"""

from __future__ import unicode_literals
from unittest import TestCase

import pexpect

from .utils import DEIS
from .utils import DEIS_SERVER
from .utils import purge
from .utils import register


class AuthTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.username, cls.password = register()

    @classmethod
    def tearDownClass(cls):
        purge(cls.username, cls.password)

    def test_login(self):
        # log in the interactive way
        child = pexpect.spawn("{} login {}".format(DEIS, DEIS_SERVER))
        child.expect('username: ')
        child.sendline(self.username)
        child.expect('password: ')
        child.sendline(self.password)
        child.expect("Logged in as {}".format(self.username))
        child.expect(pexpect.EOF)

    def test_logout(self):
        child = pexpect.spawn("{} logout".format(DEIS))
        child.expect('Logged out')
        # log in the one-liner way
        child = pexpect.spawn("{} login {} --username={} --password={}".format(
            DEIS, DEIS_SERVER, self.username, self.password))
        child.expect("Logged in as {}".format(self.username))
        child.expect(pexpect.EOF)
