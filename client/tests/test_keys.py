"""
Unit tests for the Deis CLI keys commands.

Run these tests with "python -m unittest client.tests.test_keys"
or with "./manage.py test client.KeysTest".
"""

from __future__ import unicode_literals
from unittest import TestCase
import glob
import os.path

import pexpect

from .utils import DEIS
from .utils import purge
from .utils import register


class KeysTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.username, cls.password = register(False, False)

    @classmethod
    def tearDownClass(cls):
        purge(cls.username, cls.password)

    def test_add(self):
        # test adding a specified key--the "choose a key" path is well
        # covered in utils.register()
        ssh_dir = os.path.expanduser('~/.ssh')
        pubkey = glob.glob(os.path.join(ssh_dir, '*.pub'))[0]
        child = pexpect.spawn("{} keys:add {}".format(DEIS, pubkey))
        child.expect('Uploading')
        child.expect('...done')
        child.expect(pexpect.EOF)
