"""
Common code used by the Deis CLI unit tests.
"""

from __future__ import unicode_literals
import os
import os.path
import random
import re
import shutil
import stat
import tempfile
from urllib2 import urlparse
from uuid import uuid4

import pexpect


DEIS = os.path.abspath(
    os.path.join(os.path.dirname(__file__), '..', 'deis.py'))
try:
    DEIS_SERVER = os.environ['DEIS_SERVER']
except KeyError:
    DEIS_SERVER = None
    print 'Error: env var DEIS_SERVER must point to a Deis controller URL.'
DEIS_TEST_FLAVOR = os.environ.get('DEIS_TEST_FLAVOR', 'ec2-us-west-2')
REPOSITORIES = {
    # 'example-clojure-ring': ('Clojure', 'https://github.com/opdemand/example-clojure-ring.git'),
    'Clojure': 'https://github.com/opdemand/example-clojure-ring.git',
    # 'Python': 'https://github.com/opdemand/example-python-django.git',
    'Python': 'https://github.com/opdemand/example-python-flask.git',
    'Java': 'https://github.com/opdemand/example-java-jetty.git',
    'Go': 'https://github.com/opdemand/example-go.git',
    'Node.js': 'https://github.com/opdemand/example-nodejs-express.git',
    # 'Ruby/Rails': 'https://github.com/opdemand/example-rails-helloworld.git',
    # 'Ruby/Rails': 'https://github.com/opdemand/example-rails-todo.git',
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
    ssh_path = os.path.expanduser('~/.ssh')
    child = pexpect.spawn(os.path.expanduser(
        "rm -f {}/{}*".format(ssh_path, username)))
    child.expect(pexpect.EOF)


def register():
    """Register a new Deis user from the command line."""
    username = "autotester-{}".format(uuid4().hex[:4])
    password = 'password'
    home = os.path.join('/tmp', username)
    os.environ['HOME'] = home
    os.mkdir(home)
    os.chdir(home)
    # generate an SSH key
    ssh_path = os.path.expandvars('$HOME/.ssh')
    os.mkdir(ssh_path, 0700)
    key_path = "/{}/{}".format(ssh_path, username)
    child = pexpect.spawn("ssh-keygen -f {} -t rsa -N '' -C {}".format(
        key_path, username))
    child.expect("Your public key has been saved")
    child.expect(pexpect.EOF)
    # write out ~/.ssh/config
    ssh_config_path = os.path.expandvars("$HOME/.ssh/config")
    with open(ssh_config_path, 'w') as ssh_config:
        # get hostname from DEIS_SERVER
        server = urlparse.urlparse(DEIS_SERVER).netloc
        ssh_config.write("""\
    Hostname {}
    IdentitiesOnly yes
    IdentityFile {}/.ssh/{}
""".format(server, home, username))
    # make a GIT_SSH script to enforce use of our key
    git_ssh_path = os.path.expandvars("$HOME/git_ssh.sh")
    with open(git_ssh_path, 'w') as git_ssh:
        git_ssh.write("""\
#!/bin/sh

SSH_ORIGINAL_COMMAND="ssh $@"
ssh -F {} "$@"
""".format(ssh_config_path))
    os.chmod(git_ssh_path, stat.S_IRUSR | stat.S_IRGRP | stat.S_IROTH
             | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)
    os.environ['GIT_SSH'] = git_ssh_path
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
    for index, key in re.findall('(\d)\) ([ \S]+)', child.before):
        if username in key:
            child.sendline(index)
            break
    opt = child.expect(['Import EC2 credentials\? \(y/n\) :',
                       'No EC2 credentials discovered.'])
    if opt == 0:
        child.sendline('y')
    opt = child.expect(['Import Rackspace credentials\? \(y/n\) :',
                       'No Rackspace credentials discovered.'])
    if opt == 0:
        child.sendline('y')
    opt = child.expect(['Import DigitalOcean credentials\? \(y/n\) :',
                       'No DigitalOcean credentials discovered.'])
    if opt == 0:
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
