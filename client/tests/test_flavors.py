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
from .utils import purge
from .utils import register


class FlavorsTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.username, cls.password = register()

    @classmethod
    def tearDownClass(cls):
        purge(cls.username, cls.password)

    def test_create(self):
        # create a new flavor
        id_ = "test-flavor-{}".format(uuid4().hex[:4])
        child = pexpect.spawn("{} flavors:create {} --provider={} --params='{}'".format(
            DEIS, id_, 'ec2',
            '{"region":"ap-southeast-2","image":"ami-d5f66bef","zone":"any","size":"m1.medium"}'
        ))
        child.expect(id_)
        child.expect(pexpect.EOF)
        # list the flavors and make sure it's in there
        child = pexpect.spawn("{} flavors".format(DEIS))
        child.expect(pexpect.EOF)
        flavors = re.findall('([\w|-]+): .*', child.before)
        self.assertIn(id_, flavors)
        # delete the new flavor
        child = pexpect.spawn("{} flavors:delete {}".format(DEIS, id_))
        child.expect(pexpect.EOF)
        self.assertNotIn('Error', child.before)

    def test_update(self):
        # create a new flavor
        id_ = "test-flavor-{}".format(uuid4().hex[:4])
        child = pexpect.spawn("{} flavors:create {} --provider={} --params={}".format(
            DEIS, id_, 'mock', '{}'))
        child.expect(id_)
        child.expect(pexpect.EOF)
        # update the provider
        child = pexpect.spawn("{} flavors:update {} --provider={}".format(DEIS, id_, 'ec2'))
        child.expect(pexpect.EOF)
        # update the params
        child = pexpect.spawn("{} flavors:update {} {}".format(
            DEIS, id_, "'{\"key1\": \"val1\"}'"))
        child.expect(pexpect.EOF)
        # test the flavor contents
        child = pexpect.spawn("{} flavors:info {}".format(DEIS, id_))
        child.expect(pexpect.EOF)
        results = json.loads(child.before)
        self.assertEqual('ec2', results['provider'])
        self.assertIn('key1', results['params'])
        self.assertIn('val1', results['params'])
        # update the params and provider
        child = pexpect.spawn("{} flavors:update {} {} --provider={}".format(
            DEIS, id_, "'{\"key2\": \"val2\"}'", 'mock'))
        child.expect(pexpect.EOF)
        # test the flavor contents
        child = pexpect.spawn("{} flavors:info {}".format(DEIS, id_))
        child.expect(pexpect.EOF)
        results = json.loads(child.before)
        self.assertIn('key1', results['params'])
        self.assertIn('val1', results['params'])
        self.assertIn('key2', results['params'])
        self.assertIn('val2', results['params'])
        self.assertEqual('mock', results['provider'])
        # update the params to remove a value
        child = pexpect.spawn("{} flavors:update {} {}".format(
            DEIS, id_, "'{\"key1\": null}'"))
        child.expect(pexpect.EOF)
        # test the flavor contents
        child = pexpect.spawn("{} flavors:info {}".format(DEIS, id_))
        child.expect(pexpect.EOF)
        results = json.loads(child.before)
        self.assertNotIn('key1', results['params'])
        self.assertNotIn('val1', results['params'])
        self.assertIn('key2', results['params'])
        self.assertIn('val2', results['params'])
        # delete the new flavor
        child = pexpect.spawn("{} flavors:delete {}".format(DEIS, id_))
        child.expect(pexpect.EOF)
        self.assertNotIn('Error', child.before)

    def test_delete(self):
        # create a new flavor
        id_ = "test-flavor-{}".format(uuid4().hex[:4])
        child = pexpect.spawn("{} flavors:create {} --provider={} --params='{}'".format(
            DEIS, id_, 'ec2',
            '{"region":"ap-southeast-2","image":"ami-d5f66bef","zone":"any","size":"m1.medium"}'
        ))
        child.expect(id_)
        child.expect(pexpect.EOF)
        # delete the new flavor
        child = pexpect.spawn("{} flavors:delete {}".format(DEIS, id_))
        child.expect(pexpect.EOF)
        self.assertNotIn('Error', child.before)
        # list the flavors and make sure it's not in there
        child = pexpect.spawn("{} flavors".format(DEIS))
        child.expect(pexpect.EOF)
        flavors = re.findall('([\w|-]+): .*', child.before)
        self.assertNotIn(id_, flavors)

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
        results = json.loads(child.before)
        self.assertIn('created', results)
        self.assertIn('updated', results)
        self.assertIn('provider', results)
        self.assertIn('id', results)
        self.assertIn('params', results)
        self.assertEqual(results['owner'], self.username)
