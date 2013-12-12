"""
Unit tests for the Deis example-[language] projects.

Run these tests with "python -m unittest client.tests.test_examples"
or with "./manage.py test client.ExamplesTest".
"""

from __future__ import unicode_literals
from unittest import TestCase
from uuid import uuid4

import pexpect
import time

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
        # scale the formation runtime=1
        child = pexpect.spawn("{} nodes:scale {} runtime=1".format(
            DEIS, cls.formation), timeout=10 * 60)
        child.expect('Scaling nodes...')
        child.expect(r'done in \d+s')
        child.expect(pexpect.EOF)

    @classmethod
    def tearDownClass(cls):
        # scale formation runtime=0
        child = pexpect.spawn("{} nodes:scale {} runtime=0".format(
            DEIS, cls.formation), timeout=3 * 60)
        child.expect('Scaling nodes...')
        child.expect(r'done in \d+s')
        child.expect(pexpect.EOF)
        # destroy the formation
        child = pexpect.spawn("{} formations:destroy {} --confirm={}".format(
            DEIS, cls.formation, cls.formation))
        child.expect('done in ', timeout=5 * 60)
        child.expect(pexpect.EOF)
        purge(cls.username, cls.password)

    def _test_example(self, repo_name, build_timeout=120, run_timeout=60):
        # `git clone` the example app repository
        _repo_type, repo_url = EXAMPLES[repo_name]
        # print repo_name, repo_type, repo_url
        clone(repo_url, repo_name)
        # create an App
        child = pexpect.spawn("{} create --formation={}".format(
            DEIS, self.formation))
        child.expect('done, created (?P<name>[-_\w]+)', timeout=60)
        app = child.match.group('name')
        try:
            child.expect('Git remote deis added')
            child.expect(pexpect.EOF)
            child = pexpect.spawn('git push deis master')
            # check git output for repo_type, e.g. "Clojure app detected"
            # TODO: for some reason, the next regex times out...
            # child.expect("{} app detected".format(repo_type), timeout=5 * 60)
            child.expect('Launching... ', timeout=build_timeout)
            child.expect('deployed to Deis(?P<url>.+)To learn more', timeout=run_timeout)
            url = child.match.group('url')
            child.expect(' -> master')
            child.expect(pexpect.EOF, timeout=10)
            # try to fetch the URL with curl a few times, ignoring 502's
            for _ in range(6):
                child = pexpect.spawn("curl -s {}".format(url))
                i = child.expect(['Powered by Deis', '502 Bad Gateway'], timeout=5)
                child.expect(pexpect.EOF)
                if i == 0:
                    break
                time.sleep(10)
            else:
                raise RuntimeError('Persistent 502 Bad Gateway')
            # `deis config:set POWERED_BY="Automated Testing"`
            child = pexpect.spawn(
                "{} config:set POWERED_BY='Automated Testing'".format(DEIS))
            child.expect(pexpect.EOF, timeout=3 * 60)
            # then re-fetch the URL with curl and recheck the output
            for _ in range(6):
                child = pexpect.spawn("curl -s {}".format(url))
                child.expect(['Powered by Automated Testing', '502 Bad Gateway'], timeout=5)
                child.expect(pexpect.EOF)
                if i == 0:
                    break
                time.sleep(10)
            else:
                raise RuntimeError('Config:set not working')
        finally:
            # destroy the app
            child = pexpect.spawn(
                "{} apps:destroy --app={} --confirm={}".format(DEIS, app, app),
                timeout=5 * 60)
            child.expect('Git remote deis removed')
            child.expect(pexpect.EOF)

    def test_clojure_ring(self):
        self._test_example('example-clojure-ring')

    def _test_dart(self):
        # TODO: fix broken buildpack / example app
        self._test_example('example-dart')

    def test_go(self):
        self._test_example('example-go')

    def test_java_jetty(self):
        self._test_example('example-java-jetty')

    def test_nodejs_express(self):
        self._test_example('example-nodejs-express')

    def test_perl(self):
        self._test_example('example-perl', build_timeout=600)

    def test_php(self):
        self._test_example('example-php')

    def _test_play(self):
        # TODO: fix broken buildpack / example app
        self._test_example('example-play', build_timeout=720)

    def test_python_flask(self):
        self._test_example('example-python-flask')

    def test_ruby_sinatra(self):
        self._test_example('example-ruby-sinatra')

    def test_scala(self):
        self._test_example('example-scala', build_timeout=720)
