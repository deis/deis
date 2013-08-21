"""
Unit tests for the Deis CLI flavors commands.

Run these tests with "python -m unittest client.tests.test_git"
or with "./manage.py test client.GitTest".
"""

from __future__ import unicode_literals
from unittest import TestCase

# import pexpect

# from .utils import DEIS
from .utils import random_repo
from .utils import setup
from .utils import teardown


class GitTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.repo_name, repo_url = random_repo()
        cls.username, cls.password, cls.repo_dir = setup(repo_url)

    @classmethod
    def tearDownClass(cls):
        teardown(cls.username, cls.password, cls.repo_dir)

    # def test_push(self):
    #     pushes = {
    #         'clojure': 'Clojure app detected.*\[new branch\]      master -> master',
    #         'django': 'Python app detected.*\[new branch\]      master -> master',
    #         'flask': 'Python app detected.*\[new branch\]      master -> master',
    #         'go': 'Go app detected.*\[new branch\]      master -> master',
    #         'java': 'Java app detected.*\[new branch\]      master -> master',
    #         'nodejs': 'Node.js app detected.*\[new branch\]      master -> master',
    #         'rails': 'Ruby/Rails app detected.*\[new branch\]      master -> master',
    #         'rails-todo': 'Ruby/Rails app detected.*\[new branch\]      master -> master',
    #         'sinatra': 'Ruby/Rack app detected.*\[new branch\]      master -> master',
    #     }
    #     child = pexpect.spawn("python {} create --flavor=ec2-us-west-2".format(DEIS))
    #     child.expect('created (?P<name>[a-z]{6}-[a-z]{8}).*to scale a basic formation')
    #     formation = child.match.group('name')
    #     child = pexpect.spawn('git push deis master')
    #     child.expect(pushes[self.repo_name], timeout=180)
    #     child.expect(pexpect.EOF)
    #     # destroy formation the one-liner way
    #     child = pexpect.spawn("{} destroy --confirm={}".format(DEIS, formation))
    #     child.expect('Git remote deis removed')
    #     child.expect(pexpect.EOF)
