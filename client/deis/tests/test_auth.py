
"""
Unit tests for authentication in the CLI.

Run these tests with "python -m unittest deis.tests.test_auth"
or all tests with "python -m unittest discover".
"""  # pylint: disable=C0103,R0201,R0904

import os
import unittest

import pexpect


# Locate the 'of' executable script relative to this file.
CLI = os.path.abspath(
    os.path.join(os.path.dirname(__file__), '..', 'deis'))


class TestLogin(unittest.TestCase):

    """Test authentication in the CLI."""

    def setUp(self):
        # TODO: set up the CLI/api/fixtures/tests.json...somehow
        self.child = pexpect.spawn('{} login'.format(CLI))

    def tearDown(self):
        self.child = None

    def test_good_login(self):
        """Test that a valid login responds with a success message."""
        child = self.child
        child.expect('username:')
        child.sendline('autotest')
        child.expect('password:')
        child.sendline('password')
        child.expect('Logged in as autotest.')
        # call a protected API endpoint to ensure we were authenticated
        child = self.child = pexpect.spawn('{} apps'.format(CLI))
        child.expect('^\S+ +\(.+\)')

    def test_bad_login(self):
        """Test that an invalid login responds with a failure message."""
        child = self.child
        child.expect('username:')
        child.sendline('autotest')
        child.expect('password:')
        child.sendline('Pa55w0rd')
        child.expect('Login failed.')
        # call a protected API endpoint to ensure we get an unauth error
        child = self.child = pexpect.spawn('{} apps'.format(CLI))
        child.expect("\('Error")

    def test_logout(self):
        child = self.child = pexpect.spawn('{} logout'.format(CLI))
        child.expect('Logged out.')
        # call a protected API endpoint to ensure we get an unauth error
        child = self.child = pexpect.spawn('{} apps'.format(CLI))
        child.expect("\('Error")


class TestHelp(unittest.TestCase):

    """Test that the client can document its own behavior."""

    def test_version(self):
        """Test that the client reports its help message."""
        child = pexpect.spawn('{} --help'.format(CLI))
        child.expect(r'Usage: .*number of proxies\s+$')
        child = pexpect.spawn('{} -h'.format(CLI))
        child.expect(r'Usage: .*number of proxies\s+$')
        child = pexpect.spawn('{} help'.format(CLI))
        child.expect(r'Usage: .*number of proxies\s+$')


class TestVersion(unittest.TestCase):

    """Test that the client can report its version string."""

    def test_version(self):
        """Test that the client reports its version string."""
        child = pexpect.spawn('{} --version'.format(CLI))
        child.expect('Deis CLI 0.0.1')
        child = pexpect.spawn('{} -v'.format(CLI))
        child.expect('Deis CLI 0.0.1')
