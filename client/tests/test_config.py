"""
Unit tests for the Deis CLI config commands.

Run these tests with "python -m unittest client.tests.test_config"
or with "./manage.py test client.ConfigTest".
"""

from __future__ import unicode_literals
from unittest import TestCase


class ConfigTest(TestCase):

    pass


# class TestConfig(unittest.TestCase):

#     """Test configuration docs and config values."""

#     def setUp(self):
#         # TODO: set up the c3/api/fixtures/tests.json...somehow
#         child = pexpect.spawn('{} login'.format(CLI))
#         child.expect('username:')
#         child.sendline('autotest')
#         child.expect('password:')
#         child.sendline('password')
#         child.expect('Logged in as autotest.')

#     def tearDown(self):
#         self.child = None

#     def test_config_syntax(self):
#         key, value = 'MONGODB_URL', 'http://mongolab.com/test'
#         # Test some invalid command line input
#         child = pexpect.spawn('{} config:set {}'.format(
#             CLI, key))
#         child.expect('Usage: ')
#         child = pexpect.spawn('{} config:set {} {}'.format(
#             CLI, key, value))
#         child.expect('Usage: ')
#         child = pexpect.spawn('{} config set {}={}'.format(
#             CLI, key, value))
#         child.expect('Usage: ')

#     def test_config(self):
#         """Test that a user can set a config value."""
#         key, value = 'MONGODB_URL', 'http://mongolab.com/test'
#         child = pexpect.spawn('{} config:set {}={}'.format(
#             CLI, key, value))
#         child.expect(pexpect.EOF)
#         child = pexpect.spawn('{} config:set {}={} DEBUG=True'.format(
#             CLI, key, value))
#         child.expect(pexpect.EOF)
