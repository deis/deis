"""
Unit tests for the Deis CLI flavors commands.

Run these tests with "python -m unittest client.tests.test_git"
or with "./manage.py test client.GitTest".
"""

from __future__ import unicode_literals
from unittest import TestCase
from uuid import uuid4

import pexpect

from .utils import DEIS
from .utils import DEIS_TEST_FLAVOR
from .utils import random_repo
from .utils import setup
from .utils import teardown


class GitTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.repo_name, repo_url = random_repo()
        cls.username, cls.password, cls.repo_dir = setup(repo_url)
        # create a new formation
        cls.formation = "{}-test-formation-{}".format(
            cls.username, uuid4().hex[:4])
        child = pexpect.spawn("{} formations:create {} --flavor={}".format(
            DEIS, cls.formation, DEIS_TEST_FLAVOR))
        child.expect("created {}.*to scale a basic formation".format(
            cls.formation))
        child.expect(pexpect.EOF)
        # create an app
        child = pexpect.spawn("{} create --formation={}".format(
            DEIS, cls.formation))
        child.expect('done, created (?P<name>[-_\w]+)')
        cls.app = child.match.group('name')
        child.expect(pexpect.EOF)

    @classmethod
    def tearDownClass(cls):
        # destroy the formation
        child = pexpect.spawn("{} formations:destroy {} --confirm={}".format(
            DEIS, cls.formation, cls.formation))
        child.expect('done in ', timeout=5*60)
        child.expect(pexpect.EOF)
        teardown(cls.username, cls.password, cls.repo_dir)

    def test_push(self):
        child = pexpect.spawn('git push deis master')
        # check git output for "Clojure app detected", for example
        print self.repo_name
        child.expect('----->')
        print child.before
        # child.expect("{} app detected".format(self.repo_name))
        # print child.before
        child.expect(' -> master', timeout=10*60)
        child.expect(pexpect.EOF, timeout=2*60)
        print child.before
