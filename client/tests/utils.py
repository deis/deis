"""
Common code used by the Deis CLI unit tests.
"""

from __future__ import unicode_literals
import os
import os.path
import random
import shutil
import tempfile
from uuid import uuid4

import pexpect


DEIS = os.path.abspath(
    os.path.join(os.path.dirname(__file__), '..', 'deis.py'))
try:
    DEIS_SERVER = os.environ['DEIS_SERVER']
except KeyError:
    raise EnvironmentError(
        'Environment variable DEIS_SERVER must point to a Deis controller URL.')
DEIS_TEST_FLAVOR = os.environ.get('DEIS_TEST_FLAVOR', 'ec2-us-west-2')
REPOSITORIES = {
    'Clojure': 'https://github.com/opdemand/example-clojure-ring.git',
    # 'Python': 'https://github.com/opdemand/example-python-django.git',
    'Python': 'https://github.com/opdemand/example-python-flask.git',
    'Java': 'https://github.com/opdemand/example-java-jetty.git',
    'Go': 'https://github.com/opdemand/example-go.git',
    'Node.js': 'https://github.com/opdemand/example-nodejs-express.git',
    # 'Ruby/Rails': 'https://github.com/opdemand/example-rails-helloworld.git',
    'Ruby/Rails': 'https://github.com/opdemand/example-rails-todo.git',
    'Ruby/Rack': 'https://github.com/opdemand/example-ruby-sinatra.git',
}


def purge(username, password):
    """Purge an existing Deis user."""
    child = pexpect.spawn("{} auth:cancel".format(DEIS))
    child.expect('username: ')
    child.sendline(username)
    child.expect('password: ')
    child.sendline(password)
    child.expect('\? \(y/n\) ')
    child.sendline('y')
    child.expect(pexpect.EOF)


def register():
    """Register a new Deis user from the command line."""
    username = "autotester-{}".format(uuid4().hex[:4])
    password = 'password'
    child = pexpect.spawn("{} register {}".format(DEIS, DEIS_SERVER))
    child.expect('username: ')
    child.sendline(username)
    child.expect('password: ')
    child.sendline(password)
    child.expect('password \(confirm\): ')
    child.sendline(password)
    child.expect('email: ')
    child.sendline('autotest@opdemand.com')
    child.expect('Which would you like to use with Deis')
    child.sendline('1')
    child.expect('Import these credentials\? \(y/n\) :')
    child.sendline('y')
    child.expect(pexpect.EOF)
    return username, password


def random_repo():
    """Return an example Heroku-style repository name and URL."""
    name = random.choice(REPOSITORIES.keys())
    return name, REPOSITORIES[name]


def setup(repo_url=None):
    """Do user and fixture setup for CLI tests."""
    # create an autotest user
    username, password = register()
    repo_dir = None
    if repo_url:
        # clone an example repository
        repo_dir = tempfile.mkdtemp()
        child = pexpect.spawn("git clone {} {}".format(repo_url, repo_dir))
        child.expect(', done')
        child.expect(pexpect.EOF)
        # cd to repo dir
        os.chdir(repo_dir)
    return (username, password, repo_dir)


def teardown(username, password, repo_dir=None):
    """Undo user and fixture setup for CLI tests."""
    # destroy the example repository
    if repo_dir:
        shutil.rmtree(repo_dir)
    # destroy the autotest user
    purge(username, password)
