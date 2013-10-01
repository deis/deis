"""
Unit tests for the Deis CLI flavors commands.

Run these tests with "python -m unittest client.tests.test_flavors"
or with "./manage.py test client.FlavorsTest".
"""

from __future__ import unicode_literals
from unittest import TestCase
import re

import json
import pexpect
import random
from uuid import uuid4

from .utils import DEIS
from .utils import setup
from .utils import teardown


class FlavorsTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.username, cls.password, _ = setup()

    @classmethod
    def tearDownClass(cls):
        teardown(cls.username, cls.password)

    def test_create(self):
        id_ = "test-flavor-{}".format(uuid4().hex[:4])
        child = pexpect.spawn("{} flavors:create {} --provider={} --params='{}'".format(
            DEIS, id_, 'ec2',
            '{"region":"ap-southeast-2","image":"ami-d5f66bef","zone":"any","size":"m1.medium"}'
        ))
        child.expect(id_)
        child.expect(pexpect.EOF)
        child = pexpect.spawn("{} flavors:delete {}".format(DEIS, id_))
        child.expect(pexpect.EOF)
        self.assertNotIn('Error', child.before)

    # def test_update(self):
    #     pass

    # def test_delete(self):
    #     pass

    def test_list(self):
        child = pexpect.spawn("{} flavors".format(DEIS))
        child.expect(pexpect.EOF)
        before = child.before
        flavors = re.findall('([\w|-]+): .*', before)
        # test that there were at least 3 flavors seeded
        self.assertGreaterEqual(len(flavors), 3)
        # test that "flavors" and "flavors:list" are equivalent
        child = pexpect.spawn("{} flavors:list".format(DEIS))
        child.expect(pexpect.EOF)
        self.assertEqual(before, child.before)

    def test_info(self):
        child = pexpect.spawn("{} flavors".format(DEIS))
        child.expect(pexpect.EOF)
        flavor = random.choice(re.findall('([\w|-]+): .*', child.before))
        child = pexpect.spawn("{} flavors:info {}".format(DEIS, flavor))
        child.expect(pexpect.EOF)
        # test that we received JSON results
        # TODO: There's some error here, but only when run as part of the
        # entire test suite?
        print child.before
        results = json.loads(child.before)
        self.assertIn('created', results)
        self.assertIn('updated', results)
        self.assertIn('provider', results)
        self.assertIn('id', results)
        self.assertIn('params', results)
        self.assertEqual(results['owner'], self.username)
