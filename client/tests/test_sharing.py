"""
Unit tests for the Deis example-[language] projects.

Run these tests with "python -m unittest client.tests.test_sharing"
or with "./manage.py test client.SharingTest".
"""

from __future__ import unicode_literals
from unittest import TestCase
from uuid import uuid4
import os

import pexpect
from .utils import DEIS
from .utils import DEIS_TEST_FLAVOR
from .utils import EXAMPLES
from .utils import clone
from .utils import login
from .utils import purge
from .utils import register


class SharingTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.username2, cls.password2 = register()
        cls.username, cls.password = register()
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
        purge(cls.username2, cls.password2)
        login(cls.username, cls.password)
        child = pexpect.spawn("{} formations:destroy {} --confirm={}".format(
            DEIS, cls.formation, cls.formation))
        child.expect('done in ', timeout=5 * 60)
        child.expect(pexpect.EOF)
        purge(cls.username, cls.password)

    def _test_sharing(self, repo_name):
        # `git clone` the example app repository
        repo_type, repo_url = EXAMPLES[repo_name]
        clone(repo_url, repo_name)
        # create an App
        child = pexpect.spawn("{} create --formation={}".format(
            DEIS, self.formation))
        child.expect('done, created (?P<name>[-_\w]+)', timeout=3 * 60)
        app = child.match.group('name')
        try:
            child.expect('Git remote deis added')
            child.expect(pexpect.EOF)
            home = os.environ['HOME']
            login(self.username2, self.password2)
            os.chdir(os.path.join(home, repo_name))
            child = pexpect.spawn('git push deis master')
            child.expect('access denied')
            child.expect(pexpect.EOF)
            login(self.username, self.password)
            os.chdir(os.path.join(home, repo_name))
            child = pexpect.spawn("{} sharing:add {} --app={}".format(
                DEIS, self.username2, app))
            child.expect('done')
            child.expect(pexpect.EOF)
            login(self.username2, self.password2)
            os.chdir(os.path.join(home, repo_name))
            child = pexpect.spawn('git push deis master')
            # check git output for repo_type, e.g. "Clojure app detected"
            # TODO: for some reason, the next regex times out...
            # child.expect("{} app detected".format(repo_type), timeout=5 * 60)
            child.expect('Launching... ', timeout=10 * 60)
            child.expect('deployed to Deis(?P<url>.+)To learn more', timeout=3 * 60)
            url = child.match.group('url')  # noqa
            child.expect(' -> master')
            child.expect(pexpect.EOF, timeout=2 * 60)
        finally:
            login(self.username, self.password)
            os.chdir(os.path.join(home, repo_name))
            # destroy the app
            child = pexpect.spawn(
                "{} apps:destroy --app={} --confirm={}".format(DEIS, app, app),
                timeout=5 * 60)
            child.expect('Git remote deis removed')
            child.expect(pexpect.EOF)

    def test_go(self):
        self._test_sharing('example-go')
