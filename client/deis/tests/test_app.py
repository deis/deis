"""
Unit tests for applications in the Deis CLI.

Run these tests with "python -m unittest deis.tests.test_app"
or all tests with "python -m unittest discover".
"""  # pylint: disable=C0103,R0201,R0904

import os
import unittest

import pexpect


# Locate the 'of' executable script relative to this file.
CLI = os.path.abspath(
    os.path.join(os.path.dirname(__file__), '..', 'deis'))


class TestLifecycle(unittest.TestCase):

    """Test basic lifecycle methods in the Deis client."""

    def setUp(self):
        # TODO: set up the deis/api/fixtures/tests.json...somehow
        child = pexpect.spawn('{} login'.format(CLI))
        child.expect('username:')
        child.sendline('autotest')
        child.expect('password:')
        child.sendline('password')
        child.expect('Logged in as autotest.')

    def tearDown(self):
        self.child = None

    def test_create(self):
        """Test that a user can create an application."""
        child = pexpect.spawn('{} create'.format(CLI))
        child.expect('Created (.+)')
        app_name = child.match.group(1).strip()
        # check that the new app shows up in the app list
        child = pexpect.spawn('{} apps'.format(CLI))
        child.expect('.*{}.*'.format(app_name))
