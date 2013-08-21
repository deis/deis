"""
Unit tests for the Deis CLI apps commands.

Run these tests with "python -m unittest client.tests.test_apps"
or with "./manage.py test client.AppsTest".
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


class AppsTest(TestCase):

    @classmethod
    def setUpClass(cls):
        repo_name, repo_url = random_repo()
        cls.username, cls.password, cls.repo_dir = setup(repo_url)
        # create a new formation
        cls.formation = "{}-test-formation-{}".format(
            cls.username, uuid4().hex[:4])
        child = pexpect.spawn("{} formations:create {} --flavor={}".format(
            DEIS, cls.formation, DEIS_TEST_FLAVOR))
        child.expect("created {}.*to scale a basic formation".format(
            cls.formation))
        child.expect(pexpect.EOF)

    @classmethod
    def tearDownClass(cls):
        # delete the formation
        child = pexpect.spawn("{} formations:destroy {} --confirm={}".format(
            DEIS, cls.formation, cls.formation))
        child.expect('done in ', timeout=5*60)
        child.expect(pexpect.EOF)
        teardown(cls.username, cls.password, cls.repo_dir)

    def test_app_create(self):
        self.assertIsNotNone(self.formation)
