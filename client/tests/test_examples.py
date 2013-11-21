"""
Unit tests for the Deis example-[language] projects.

Run these tests with "python -m unittest client.tests.test_examples"
or with "./manage.py test client.ExamplesTest".
"""

from __future__ import unicode_literals
from unittest import TestCase
from uuid import uuid4

import pexpect
from .utils import DEIS
from .utils import DEIS_TEST_FLAVOR
from .utils import EXAMPLES
from .utils import clone
from .utils import purge
from .utils import register


class ExamplesTest(TestCase):

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
        # TODO: scale the formation runtime=1

    @classmethod
    def tearDownClass(cls):
        # TODO: scale formation runtime=0
        # destroy the formation
        child = pexpect.spawn("{} formations:destroy {} --confirm={}".format(
            DEIS, cls.formation, cls.formation))
        child.expect('done in ', timeout=5*60)
        child.expect(pexpect.EOF)
        purge(cls.username, cls.password)

    def _test_example(self, repo_name):
        # `git clone` the example app repository
        repo_type, repo_url = EXAMPLES[repo_name]
        # print repo_name, repo_type, repo_url
        clone(repo_url, repo_name)
        # create an App
        child = pexpect.spawn("{} create --formation={}".format(
            DEIS, self.formation))
        child.expect('done, created (?P<name>[-_\w]+)')
        app = child.match.group('name')
        try:
            child.expect('Git remote deis added')
            child.expect(pexpect.EOF)
            child = pexpect.spawn('git push deis master')
            # check git output for repo_type, e.g. "Clojure app detected"
            child.expect("{} app detected".format(repo_type), timeout=2*60)
            child.expect(' -> master', timeout=10*60)
            child.expect(pexpect.EOF, timeout=2*60)
            # TODO: scale up runtime nodes in setUpClass, then
            # actually fetch the URL with curl and check the output
            # TODO: `deis config:set POWERED_BY="Automated Testing"`
            # then re-fetch the URL with curl and recheck the output
        finally:
            # destroy the app
            child = pexpect.spawn(
                "{} apps:destroy --app={} --confirm={}".format(DEIS, app, app),
                timeout=5*60)
            child.expect('Git remote deis removed')
            child.expect(pexpect.EOF)

    def test_clojure_ring(self):
        self._test_example('example-clojure-ring')

    def test_dart(self):
        self._test_example('example-dart')

    def test_go(self):
        self._test_example('example-go')

    def test_java_jetty(self):
        self._test_example('example-java-jetty')

    def test_nodejs_express(self):
        self._test_example('example-nodejs-express')

    def test_perl(self):
        self._test_example('example-perl')

    def test_php(self):
        self._test_example('example-php')

    def test_play(self):
        self._test_example('example-play')

    def test_python_flask(self):
        self._test_example('example-python-flask')

    def test_ruby_sinatra(self):
        self._test_example('example-ruby-sinatra')

    def test_scala(self):
        self._test_example('example-scala')
