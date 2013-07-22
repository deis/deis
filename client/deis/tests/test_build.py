"""
Unit tests for config in the Deis CLI.

Run these tests with "python -m unittest deis.tests.test_config"
or all tests with "python -m unittest discover".
"""  # pylint: disable=C0103,R0201,R0904

import json
import os
import tempfile
import unittest
import uuid

import pexpect


# Locate the 'deis' executable script relative to this file.
CLI = os.path.abspath(
    os.path.join(os.path.dirname(__file__), '..', 'deis'))

CONTROLLER = 'http://localhost:8000/'


class TestConfig(unittest.TestCase):

    """Test configuration docs and config values."""

    def setUp(self):
        # TODO: set up the c3/api/fixtures/tests.json...somehow
        child = pexpect.spawn('{} login {}'.format(CLI, CONTROLLER))
        child.expect('username:')
        child.sendline('autotest')
        child.expect('password:')
        child.sendline('password')
        child.expect('Logged in as autotest')

    def tearDown(self):
        self.child = None

    def test_build(self):
        """Test that a user can publish a new build."""
        _, temp = tempfile.mkstemp()
        body = {
            'sha': uuid.uuid4().hex,
            'slug_size': 4096000,
            'procfile': json.dumps({'web': 'node server.js'}),
            'url':
            'http://deis.local/slugs/1c52739bbf3a44d3bfb9a58f7bbdd5fb.tar.gz',
            'checksum': uuid.uuid4().hex,
        }
        with open(temp, 'w') as f:
            f.write(json.dumps(body))
        child = pexpect.spawn(
            'cat {} | {} builds:create - --app=test-app'.format(temp, CLI))
        child.expect('Usage: ')
