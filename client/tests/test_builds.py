"""
Unit tests for the Deis CLI build commands.

Run these tests with "python -m unittest client.tests.test_builds"
or with "./manage.py test client.BuildsTest".
"""

from __future__ import unicode_literals
from unittest import TestCase


class BuildsTest(TestCase):

    pass


# class TestBuild(unittest.TestCase):

#     """Test builds."""

#     def setUp(self):
#         # TODO: set up the c3/api/fixtures/tests.json...somehow
#         child = pexpect.spawn('{} login {}'.format(CLI, CONTROLLER))
#         child.expect('username:')
#         child.sendline('autotest')
#         child.expect('password:')
#         child.sendline('password')
#         child.expect('Logged in as autotest')

#     def tearDown(self):
#         self.child = None

#     def test_build(self):
#         """Test that a user can publish a new build."""
#         _, temp = tempfile.mkstemp()
#         body = {
#             'sha': uuid.uuid4().hex,
#             'slug_size': 4096000,
#             'procfile': json.dumps({'web': 'node server.js'}),
#             'url':
#             'http://deis.local/slugs/1c52739bbf3a44d3bfb9a58f7bbdd5fb.tar.gz',
#             'checksum': uuid.uuid4().hex,
#         }
#         with open(temp, 'w') as f:
#             f.write(json.dumps(body))
#         child = pexpect.spawn(
#             'cat {} | {} builds:create - --app=test-app'.format(temp, CLI))
#         child.expect('Usage: ')
