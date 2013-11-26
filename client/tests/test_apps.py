"""
Unit tests for the Deis CLI apps commands.

Run these tests with "python -m unittest client.tests.test_apps"
or with "./manage.py test client.AppsTest".
"""

from __future__ import unicode_literals
from unittest import TestCase
from uuid import uuid4
import json
import re

import pexpect

from .utils import DEIS
from .utils import DEIS_TEST_FLAVOR
from .utils import clone
from .utils import purge
from .utils import random_repo
from .utils import register


class AppsTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.username, cls.password = register()
        # create a new formation
        cls.formation = "{}-test-formation-{}".format(
            cls.username, uuid4().hex[:4])
        child = pexpect.spawn("{} formations:create {} --flavor={}".format(
            DEIS, cls.formation, DEIS_TEST_FLAVOR))
        child.expect("created {}.*to scale a basic formation".format(
            cls.formation))
        child.expect(pexpect.EOF)
        repo_name, (repo_type, repo_url) = random_repo()
        # print repo_name, repo_type, repo_url
        clone(repo_url, repo_name)

    @classmethod
    def tearDownClass(cls):
        # delete the formation
        child = pexpect.spawn("{} formations:destroy {} --confirm={}".format(
            DEIS, cls.formation, cls.formation))
        child.expect('done in ', timeout=5 * 60)
        child.expect(pexpect.EOF)
        purge(cls.username, cls.password)

    def test_create(self):
        # create the app
        self.assertIsNotNone(self.formation)
        child = pexpect.spawn("{} create --formation={}".format(
            DEIS, self.formation))
        child.expect('done, created (?P<name>[-_\w]+)', timeout=5 * 60)
        app = child.match.group('name')
        child.expect('Git remote deis added')
        child.expect(pexpect.EOF)
        # check that it's in the list of apps
        child = pexpect.spawn("{} apps".format(DEIS))
        child.expect('=== Apps')
        child.expect(pexpect.EOF)
        apps = re.findall(r'([-_\w]+) {\w?}', child.before)
        self.assertIn(app, apps)
        # destroy the app
        child = pexpect.spawn("{} apps:destroy --confirm={}".format(DEIS, app),
                              timeout=5 * 60)
        child.expect('Git remote deis removed')
        child.expect(pexpect.EOF)

    def test_destroy(self):
        # create the app
        self.assertIsNotNone(self.formation)
        child = pexpect.spawn("{} apps:create --formation={}".format(
            DEIS, self.formation))
        child.expect('done, created ([-_\w]+)', timeout=5 * 60)
        app = child.match.group(1)
        child.expect(pexpect.EOF)
        # check that it's in the list of apps
        child = pexpect.spawn("{} apps".format(DEIS))
        child.expect('=== Apps')
        child.expect(pexpect.EOF)
        apps = re.findall(r'([-_\w]+) {\w?}', child.before)
        self.assertIn(app, apps)
        # destroy the app
        child = pexpect.spawn("{} destroy --confirm={}".format(DEIS, app))
        child.expect("Destroying {}".format(app))
        child.expect('done in \d+s')
        child.expect('Git remote deis removed')
        child.expect(pexpect.EOF)

    def test_list(self):
        # list apps and get their names
        child = pexpect.spawn("{} apps".format(DEIS))
        child.expect('=== Apps')
        child.expect(pexpect.EOF)
        apps_before = re.findall(r'([-_\w]+) {\w?}', child.before)
        # create a new app
        self.assertIsNotNone(self.formation)
        child = pexpect.spawn("{} apps:create --formation={}".format(
            DEIS, self.formation))
        child.expect('done, created ([-_\w]+)')
        app = child.match.group(1)
        child.expect(pexpect.EOF)
        # list apps and get their names
        child = pexpect.spawn("{} apps".format(DEIS))
        child.expect('=== Apps')
        child.expect(pexpect.EOF)
        apps = re.findall(r'([-_\w]+) {\w?}', child.before)
        # test that the set of names contains the previous set
        self.assertLess(set(apps_before), set(apps))
        # delete the app
        child = pexpect.spawn("{} apps:destroy --app={} --confirm={}".format(
            DEIS, app, app))
        child.expect('done in ', timeout=5 * 60)
        child.expect(pexpect.EOF)
        # list apps and get their names
        child = pexpect.spawn("{} apps:list".format(DEIS))
        child.expect('=== Apps')
        child.expect(pexpect.EOF)
        apps = re.findall(r'([-_\w]+) {\w?}', child.before)
        # test that the set of names is equal to the original set
        self.assertEqual(set(apps_before), set(apps))

    def test_info(self):
        # create a new app
        self.assertIsNotNone(self.formation)
        child = pexpect.spawn("{} create --formation={}".format(
            DEIS, self.formation))
        child.expect('done, created (?P<name>[-_\w]+)')
        app = child.match.group('name')
        child.expect('Git remote deis added')
        child.expect(pexpect.EOF)
        # get app info
        child = pexpect.spawn("{} info".format(DEIS))
        child.expect("=== {} Application".format(app))
        child.expect("=== {} Containers".format(app))
        response = json.loads(child.before)
        child.expect(pexpect.EOF)
        self.assertEqual(response['id'], app)
        self.assertEqual(response['formation'], self.formation)
        self.assertEqual(response['owner'], self.username)
        self.assertIn('uuid', response)
        self.assertIn('created', response)
        self.assertIn('containers', response)
        # delete the app
        child = pexpect.spawn("{} apps:destroy --app={} --confirm={}".format(
            DEIS, app, app))
        child.expect('done in ', timeout=5 * 60)
        child.expect(pexpect.EOF)

    # def test_calculate(self):
    #     pass

    # def test_open(self):
    #     pass

    # def test_logs(self):
    #     pass

    # def test_run(self):
    #     pass
