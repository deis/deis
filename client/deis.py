#!/usr/bin/env python
"""
The Deis command-line client issues API calls to a Deis controller.

Usage: deis <command> [<args>...]

Auth commands::

  register      register a new user with a controller
  login         login to a controller
  logout        logout from the current controller

Subcommands, use ``deis help [subcommand]`` to learn more::

  formations    manage formations used to host applications
  layers        manage layers used for node configuration
  nodes         manage nodes used to host containers and proxies

  apps          manage applications used to provide services
  containers    manage containers used to handle requests and jobs
  config        manage environment variables that define app config
  builds        manage builds created using `git push`
  releases      manage releases of an application

  providers     manage credentials used to access cloud providers
  flavors       manage flavors of nodes including size and location
  keys          manage ssh keys used for `git push` deployments

  perms         manage permissions for shared apps and formations

Developer shortcut commands::

  create        create a new application
  scale         scale containers by type (web=2, worker=1)
  info          view information about the current app
  open          open a URL to the app in a browser
  logs          view aggregated log info for the app
  run           run a command in an ephemeral app container
  destroy       destroy an application

Use ``git push deis master`` to deploy to an application.

"""

from __future__ import print_function
from collections import namedtuple
from cookielib import MozillaCookieJar
from datetime import datetime
from getpass import getpass
from itertools import cycle
from threading import Event
from threading import Thread
import glob
import json
import locale
import os.path
import random
import re
import subprocess
import sys
import tempfile
import time
import urlparse
import webbrowser

from dateutil import parser
from dateutil import relativedelta
from dateutil import tz
from docopt import docopt
from docopt import DocoptExit
import requests
import yaml

__version__ = '0.4.0'


locale.setlocale(locale.LC_ALL, '')


class Session(requests.Session):
    """
    Session for making API requests and interacting with the filesystem
    """

    def __init__(self):
        super(Session, self).__init__()
        self.trust_env = False
        cookie_file = os.path.expanduser('~/.deis/cookies.txt')
        cookie_dir = os.path.dirname(cookie_file)
        self.cookies = MozillaCookieJar(cookie_file)
        # Create the $HOME/.deis dir if it doesn't exist
        if not os.path.isdir(cookie_dir):
            os.mkdir(cookie_dir, 0700)
        # Load existing cookies if the cookies.txt exists
        if os.path.isfile(cookie_file):
            self.cookies.load()
            self.cookies.clear_expired_cookies()

    def clear(self, domain):
        """Clear cookies for the specified domain."""
        try:
            self.cookies.clear(domain)
            self.cookies.save()
        except KeyError:
            pass

    def git_root(self):
        """
        Return the absolute path from the git repository root

        If no git repository exists, raise an EnvironmentError
        """
        try:
            git_root = subprocess.check_output(
                ['git', 'rev-parse', '--show-toplevel'],
                stderr=subprocess.PIPE).strip('\n')
        except subprocess.CalledProcessError:
            raise EnvironmentError('Current directory is not a git repository')
        return git_root

    def get_app(self):
        """
        Return the application name for the current directory

        The application is determined by parsing `git remote -v` output.
        If no application is found, raise an EnvironmentError.
        """
        git_root = self.git_root()
        # try to match a deis remote
        remotes = subprocess.check_output(['git', 'remote', '-v'],
                                          cwd=git_root)
        m = re.search(r'^deis\W+(?P<url>\S+)\W+\(', remotes, re.MULTILINE)
        if not m:
            raise EnvironmentError(
                'Could not find deis remote in `git remote -v`')
        url = m.groupdict()['url']
        m = re.match('\S+:(?P<app>[a-z0-9-]+)(.git)?', url)
        if not m:
            raise EnvironmentError("Could not parse: {url}".format(**locals()))
        return m.groupdict()['app']

    app = property(get_app)

    def request(self, *args, **kwargs):
        """
        Issue an HTTP request with proper cookie handling including
        `Django CSRF tokens <https://docs.djangoproject.com/en/dev/ref/contrib/csrf/>`
        """
        for cookie in self.cookies:
            if cookie.name == 'csrftoken':
                if 'headers' in kwargs:
                    kwargs['headers']['X-CSRFToken'] = cookie.value
                else:
                    kwargs['headers'] = {'X-CSRFToken': cookie.value}
                break
        response = super(Session, self).request(*args, **kwargs)
        self.cookies.save()
        return response


class Settings(dict):
    """
    Settings backed by a file in the user's home directory

    On init, settings are loaded from ~/.deis/client.yaml
    """

    def __init__(self):
        path = os.path.expanduser('~/.deis')
        if not os.path.exists(path):
            os.mkdir(path)
        self._path = os.path.join(path, 'client.yaml')
        if not os.path.exists(self._path):
            with open(self._path, 'w') as f:
                f.write(yaml.safe_dump({}))
        # load initial settings
        self.load()

    def load(self):
        """
        Deserialize and load settings from the filesystem
        """
        with open(self._path) as f:
            data = f.read()
        settings = yaml.safe_load(data)
        self.update(settings)
        return settings

    def save(self):
        """
        Serialize and save settings to the filesystem
        """
        data = yaml.safe_dump(dict(self))
        with open(self._path, 'w') as f:
            f.write(data)
        return data


_counter = 0


def _newname(template="Thread-{}"):
    """Generate a new thread name."""
    global _counter
    _counter += 1
    return template.format(_counter)


FRAMES = {
    'arrow': ['^', '>', 'v', '<'],
    'dots': ['...', 'o..', '.o.', '..o'],
    'ligatures': ['bq', 'dp', 'qb', 'pd'],
    'lines': [' ', '-', '=', '#', '=', '-'],
    'slash': ['-', '\\', '|', '/'],
}


class TextProgress(Thread):
    """Show progress for a long-running operation on the command-line."""

    def __init__(self, group=None, target=None, name=None, args=(), kwargs={}):
        name = name or _newname("TextProgress-Thread-{}")
        style = kwargs.get('style', 'dots')
        super(TextProgress, self).__init__(
            group, target, name, args, kwargs)
        self.daemon = True
        self.cancelled = Event()
        self.frames = cycle(FRAMES[style])

    def run(self):
        """Write ASCII progress animation frames to stdout."""
        if not os.environ.get('DEIS_HIDE_PROGRESS'):
            time.sleep(0.5)
            self._write_frame(self.frames.next(), erase=False)
            while not self.cancelled.is_set():
                time.sleep(0.4)
                self._write_frame(self.frames.next())
            # clear the animation
            sys.stdout.write('\b' * (len(self.frames.next()) + 2))
            sys.stdout.flush()

    def cancel(self):
        """Set the animation thread as cancelled."""
        self.cancelled.set()

    def _write_frame(self, frame, erase=True):
        if erase:
            backspaces = '\b' * (len(frame) + 2)
        else:
            backspaces = ''
        sys.stdout.write("{} {} ".format(backspaces, frame))
        # flush stdout or we won't see the frame
        sys.stdout.flush()


def dictify(args):
    """Converts a list of key=val strings into a python dict.

    >>> dictify(['MONGODB_URL=http://mongolabs.com/test', 'scale=5'])
    {'MONGODB_URL': 'http://mongolabs.com/test', 'scale': 5}
    """
    data = {}
    for arg in args:
        try:
            var, val = arg.split('=')
        except ValueError:
            raise DocoptExit()
        # Try to coerce the value to an int since that's a common use case
        try:
            data[var] = int(val)
        except ValueError:
            data[var] = val
    return data


def get_provider_creds(provider, raise_error=False):
    """Query environment variables and return a provider's creds if found.
    """
    cred_types = {
        'ec2': [[('AWS_ACCESS_KEY_ID', 'access_key', None),
                 ('AWS_SECRET_ACCESS_KEY', 'secret_key', None)],
                [('AWS_ACCESS_KEY', 'access_key', None),
                 ('AWS_SECRET_KEY', 'secret_key', None)]],
        'rackspace': [[('RACKSPACE_USERNAME', 'username', None),
                       ('RACKSPACE_API_KEY', 'api_key', None),
                       ('CLOUD_ID_TYPE', 'identity_type', 'rackspace')]],
        'digitalocean': [[('DIGITALOCEAN_CLIENT_ID', 'client_id', None),
                          ('DIGITALOCEAN_API_KEY', 'api_key', None)]]
    }
    missing = None
    for cred_set in cred_types[provider]:
        creds = {}
        for envvar, key, default in cred_set:
            val = os.environ.get(envvar, default)
            if not val:
                missing = envvar
                break
            else:
                creds[key] = val
        if creds:
            return creds
    if raise_error:
        raise EnvironmentError(
            "Missing environment variable: {}".format(missing))


def readable_datetime(datetime_str):
    """
    Return a human-readable datetime string from an ECMA-262 (JavaScript)
    datetime string.
    """
    timezone = tz.tzlocal()
    dt = parser.parse(datetime_str).astimezone(timezone)
    now = datetime.now(timezone)
    delta = relativedelta.relativedelta(now, dt)
    # if it happened today, say "2 hours and 1 minute ago"
    if delta.days <= 1 and dt.day == now.day:
        if delta.hours == 0:
            hour_str = ''
        elif delta.hours == 1:
            hour_str = '1 hour '
        else:
            hour_str = "{} hours ".format(delta.hours)
        if delta.minutes == 0:
            min_str = ''
        elif delta.minutes == 1:
            min_str = '1 minute '
        else:
            min_str = "{} minutes ".format(delta.minutes)
        if not any((hour_str, min_str)):
            return 'Just now'
        else:
            return "{}{}ago".format(hour_str, min_str)
    # if it happened yesterday, say "yesterday at 3:23 pm"
    yesterday = now + relativedelta.relativedelta(days=-1)
    if delta.days <= 2 and dt.day == yesterday.day:
        return dt.strftime("Yesterday at %X")
    # otherwise return locale-specific date/time format
    else:
        return dt.strftime('%c %Z')


def trim(docstring):
    """
    Function to trim whitespace from docstring

    c/o PEP 257 Docstring Conventions
    <http://www.python.org/dev/peps/pep-0257/>
    """
    if not docstring:
        return ''
    # Convert tabs to spaces (following the normal Python rules)
    # and split into a list of lines:
    lines = docstring.expandtabs().splitlines()
    # Determine minimum indentation (first line doesn't count):
    indent = sys.maxint
    for line in lines[1:]:
        stripped = line.lstrip()
        if stripped:
            indent = min(indent, len(line) - len(stripped))
    # Remove indentation (first line is special):
    trimmed = [lines[0].strip()]
    if indent < sys.maxint:
        for line in lines[1:]:
            trimmed.append(line[indent:].rstrip())
    # Strip off trailing and leading blank lines:
    while trimmed and not trimmed[-1]:
        trimmed.pop()
    while trimmed and not trimmed[0]:
        trimmed.pop(0)
    # Return a single string:
    return '\n'.join(trimmed)


class ResponseError(Exception):
    pass


class DeisClient(object):
    """
    A client which interacts with a Deis controller.
    """

    def __init__(self):
        self._session = Session()
        self._settings = Settings()

    def _dispatch(self, method, path, body=None,
                  headers={'content-type': 'application/json'}, **kwargs):
        """
        Dispatch an API request to the active Deis controller
        """
        func = getattr(self._session, method.lower())
        controller = self._settings['controller']
        if not controller:
            raise EnvironmentError(
                'No active controller. Use `deis login` or `deis register` to get started.')
        url = urlparse.urljoin(controller, path, **kwargs)
        response = func(url, data=body, headers=headers)
        return response

    def apps(self, args):
        """
        Valid commands for apps:

        apps:create        create a new application
        apps:list          list accessible applications
        apps:info          view info about an application
        apps:open          open the application in a browser
        apps:logs          view aggregated application logs
        apps:run           run a command in an ephemeral app container
        apps:destroy       destroy an application

        Use `deis help [command]` to learn more
        """
        return self.apps_list(args)

    def apps_calculate(self, args, quiet=False):
        """
        Calculate the application's JSON representation

        Usage: deis apps:calculate [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch('post',
                                  "/api/apps/{}/calculate".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            databag = json.loads(response.content)
            if quiet is False:
                print(json.dumps(databag, indent=2))
            return databag
        else:
            raise ResponseError(response)

    def apps_create(self, args):
        """
        Create a new application

        If no ID is provided, one will be generated automatically.
        If no formation is provided, the first available will be used.

        Usage: deis apps:create [--id=<id> --formation=<formation>]
        """
        body = {}
        try:
            self._session.git_root()  # check for a git repository
        except EnvironmentError:
            print('No git repository found, use `git init` to create one')
            sys.exit(1)
        try:
            self._session.get_app()
            print('Deis remote already exists')
            sys.exit(1)
        except EnvironmentError:
            pass
        for opt in ('--id', '--formation'):
            o = args.get(opt)
            if o:
                body.update({opt.strip('-'): o})
        sys.stdout.write('Creating application... ')
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('post', '/api/apps',
                                      json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            data = response.json()
            app_id = data['id']
            print("done, created {}".format(app_id))
            # add a git remote
            hostname = urlparse.urlparse(self._settings['controller']).netloc
            git_remote = "git@{hostname}:{app_id}.git".format(**locals())
            try:
                subprocess.check_call(
                    ['git', 'remote', 'add', '-f', 'deis', git_remote],
                    stdout=subprocess.PIPE)
            except subprocess.CalledProcessError:
                sys.exit(1)
            print('Git remote deis added')
        else:
            raise ResponseError(response)

    def apps_destroy(self, args):
        """
        Destroy an application

        Usage: deis apps:destroy [--app=<id> --confirm=<confirm>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        confirm = args.get('--confirm')
        if confirm == app:
            pass
        else:
            print("""
 !    WARNING: Potentially Destructive Action
 !    This command will destroy the application: {app}
 !    To proceed, type "{app}" or re-run this command with --confirm={app}
""".format(**locals()))
            confirm = raw_input('> ').strip('\n')
            if confirm != app:
                print('Destroy aborted')
                return
        sys.stdout.write("Destroying {}... ".format(app))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('delete', "/api/apps/{}".format(app))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code in (requests.codes.no_content,  # @UndefinedVariable
                                    requests.codes.not_found):  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
            try:
                subprocess.check_call(
                    ['git', 'remote', 'rm', 'deis'],
                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                print('Git remote deis removed')
            except subprocess.CalledProcessError:
                pass  # ignore error
        else:
            raise ResponseError(response)

    def apps_list(self, args):
        """
        List applications visible to the current user

        Usage: deis apps:list
        """
        response = self._dispatch('get', '/api/apps')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            print('=== Apps')
            for item in data['results']:
                print('{id} {containers}'.format(**item))
        else:
            raise ResponseError(response)

    def apps_info(self, args):
        """
        Print info about the current application

        Usage: deis apps:info [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch('get', "/api/apps/{}".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print("=== {} Application".format(app))
            print(json.dumps(response.json(), indent=2))
            print()
            self.containers_list(args)
            print()
        else:
            raise ResponseError(response)

    def apps_open(self, args):
        """
        Open a URL to the application in a browser

        Usage: deis apps:open [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        # TODO: replace with a proxy lookup that doesn't have any side effects
        # this currently recalculates and updates the databag
        response = self._dispatch('post',
                                  "/api/apps/{}/calculate".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            databag = json.loads(response.content)
            domains = databag.get('domains', [])
            if domains:
                domain = random.choice(domains)
                # use the OS's default handler to open this URL
                webbrowser.open('http://{}/'.format(domain))
                return domain
            else:
                print('No proxies found. Use `deis nodes:scale myformation proxy=1` to scale up.')
        else:
            raise ResponseError(response)

    def apps_logs(self, args):
        """
        Retrieve the most recent log events

        Usage: deis apps:logs [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch('post',
                                  "/api/apps/{}/logs".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(response.json())
        else:
            raise ResponseError(response)

    def apps_run(self, args):
        """
        Run a command inside an ephemeral app container

        Usage: deis apps:run <command>...
        """
        app = self._session.app
        body = {'command': ' '.join(sys.argv[2:])}
        response = self._dispatch('post',
                                  "/api/apps/{}/run".format(app),
                                  json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            output, rc = json.loads(response.content)
            if rc != 0:
                print('Warning: non-zero return code {}'.format(rc))
            sys.stdout.write(output)
            sys.stdout.flush()
        else:
            raise ResponseError(response)

    def auth_register(self, args):
        """
        Register a new user with a Deis controller

        Usage: deis auth:register <controller> [options]

        Options:

        --username=USERNAME    provide a username for the new account
        --password=PASSWORD    provide a password for the new account
        --email=EMAIL          provide an email address
        """
        controller = args['<controller>']
        if not urlparse.urlparse(controller).scheme:
            controller = "http://{}".format(controller)
        username = args.get('--username')
        if not username:
            username = raw_input('username: ')
        password = args.get('--password')
        if not password:
            password = getpass('password: ')
            confirm = getpass('password (confirm): ')
            if password != confirm:
                print('Password mismatch, aborting registration.')
                return False
        email = args.get('--email')
        if not email:
            email = raw_input('email: ')
        url = urlparse.urljoin(controller, '/api/auth/register')
        payload = {'username': username, 'password': password, 'email': email}
        response = self._session.post(url, data=payload, allow_redirects=False)
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            self._settings['controller'] = controller
            self._settings.save()
            print("Registered {}".format(username))
            login_args = {'--username': username, '--password': password,
                          '<controller>': controller}
            if self.auth_login(login_args) is False:
                print('Login failed')
        else:
            print('Registration failed', response.content)
            return False
        return True

    def auth_cancel(self, args):
        """
        Cancel and remove the current account.

        Usage: deis auth:cancel
        """
        controller = self._settings.get('controller')
        if not controller:
            print('Not logged in to a Deis controller')
            return False
        print('Please log in again in order to cancel this account')
        username = self.auth_login({'<controller>': controller})
        if username:
            confirm = raw_input("Cancel account \"{}\" at {}? (y/n) ".format(username, controller))
            if confirm == 'y':
                self._dispatch('delete', '/api/auth/cancel')
                self._session.cookies.clear()
                self._session.cookies.save()
                self._settings['controller'] = None
                self._settings.save()
                print('Account cancelled')
            else:
                print('Accont not changed')

    def auth_login(self, args):
        """
        Login by authenticating against a controller

        Usage: deis auth:login <controller> [--username=<username> --password=<password>]
        """
        controller = args['<controller>']
        if not urlparse.urlparse(controller).scheme:
            controller = "http://{}".format(controller)
        username = args.get('--username')
        headers = {}
        if not username:
            username = raw_input('username: ')
        password = args.get('--password')
        if not password:
            password = getpass('password: ')
        url = urlparse.urljoin(controller, '/api/auth/login/')
        payload = {'username': username, 'password': password}
        # clear any cookies for this controller's domain
        self._session.clear(urlparse.urlparse(url).netloc)
        # prime cookies for login
        self._session.get(url, headers=headers)
        # post credentials to the login URL
        response = self._session.post(url, data=payload, allow_redirects=False)
        if response.status_code == requests.codes.found:  # @UndefinedVariable
            self._settings['controller'] = controller
            self._settings.save()
            print("Logged in as {}".format(username))
            return username
        else:
            print('Login failed')
            self._session.cookies.clear()
            self._session.cookies.save()
            return False

    def auth_logout(self, args):
        """
        Logout from a controller and clear the user session

        Usage: deis auth:logout
        """
        controller = self._settings.get('controller')
        if controller:
            self._dispatch('get', '/api/auth/logout/')
        self._session.cookies.clear()
        self._session.cookies.save()
        self._settings['controller'] = None
        self._settings.save()
        print('Logged out')

    def builds(self, args):
        """
        Valid commands for builds:

        builds:list        list build history for a formation
        builds:create      coming soon!

        Use `deis help [command]` to learn more
        """
        return self.builds_list(args)

    def builds_list(self, args):
        """
        List build history for a formation

        Usage: deis builds:list [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch('get', "/api/apps/{}/builds".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print("=== {} Builds".format(app))
            data = response.json()
            for item in data['results']:
                print("{0[uuid]:<23} {0[created]}".format(item))
        else:
            raise ResponseError(response)

    def config(self, args):
        """
        Valid commands for config:

        config:list        list environment variables for an app
        config:set         set environment variables for an app
        config:unset       unset environment variables for an app

        Use `deis help [command]` to learn more
        """
        return self.config_list(args)

    def config_list(self, args):
        """
        List environment variables for an application

        Usage: deis config:list [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch('get', "/api/apps/{}/config".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            config = response.json()
            values = json.loads(config['values'])
            print("=== {} Config".format(app))
            items = values.items()
            if len(items) == 0:
                print('No configuration')
                return
            for k, v in values.items():
                print("{k}: {v}".format(**locals()))
        else:
            raise ResponseError(response)

    def config_set(self, args):
        """
        Set environment variables for an application

        Usage: deis config:set <var>=<value>... [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        body = {'values': json.dumps(dictify(args['<var>=<value>']))}
        response = self._dispatch('post',
                                  "/api/apps/{}/config".format(app),
                                  json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            config = response.json()
            values = json.loads(config['values'])
            print("=== {}".format(app))
            items = values.items()
            if len(items) == 0:
                print('No configuration')
                return
            for k, v in values.items():
                print("{k}: {v}".format(**locals()))
        else:
            raise ResponseError(response)

    def config_unset(self, args):
        """
        Unset an environment variable for an application

        Usage: deis config:unset <key>... [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        values = {}
        for k in args.get('<key>'):
            values[k] = None
        body = {'values': json.dumps(values)}
        response = self._dispatch('post',
                                  "/api/apps/{}/config".format(app),
                                  json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            config = response.json()
            values = json.loads(config['values'])
            print("=== {}".format(app))
            items = values.items()
            if len(items) == 0:
                print('No configuration')
                return
            for k, v in values.items():
                print("{k}: {v}".format(**locals()))
        else:
            raise ResponseError(response)

    def containers(self, args):
        """
        Valid commands for containers:

        containers:list        list application containers
        containers:scale       scale app containers (e.g. web=4 worker=2)

        Use `deis help [command]` to learn more
        """
        return self.containers_list(args)

    def containers_list(self, args):
        """
        List containers servicing an application

        Usage: deis containers:list [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.get_app()
        response = self._dispatch('get',
                                  "/api/apps/{}/containers".format(app))
        if response.status_code != requests.codes.ok:  # @UndefinedVariable
            raise ResponseError(response)
        containers = response.json()
        response = self._dispatch('get', "/api/apps/{}/builds".format(app))
        if response.status_code != requests.codes.ok:  # @UndefinedVariable
            raise ResponseError(response)
        txt = response.json()['results'][0]['procfile']
        procfile = json.loads(txt) if txt else {}
        print("=== {} Containers".format(app))
        c_map = {}
        for item in containers['results']:
            c_map.setdefault(item['type'], []).append(item)
        print()
        for c_type in c_map.keys():
            command = procfile.get(c_type, '<none>')
            print("--- {c_type}: `{command}`".format(**locals()))
            for c in c_map[c_type]:
                print("{type}.{num} up {created} ({node})".format(**c))
            print()

    def containers_scale(self, args):
        """
        Scale an application's containers by type

        Example: deis containers:scale web=4 worker=2

        Usage: deis containers:scale <type=num>... [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.get_app()
        body = {}
        for type_num in args.get('<type=num>'):
            typ, count = type_num.split('=')
            body.update({typ: int(count)})
        print('Scaling containers... but first, coffee!')
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('post',
                                      "/api/apps/{}/scale".format(app),
                                      json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('done in {}s\n'.format(int(time.time() - before)))
            self.containers_list({})
        else:
            raise ResponseError(response)

    def flavors(self, args):
        """
        Valid commands for flavors:

        flavors:create        create a new node flavor
        flavors:info          print information about a node flavor
        flavors:list          list available flavors
        flavors:update        update an existing node flavor
        flavors:delete        delete a node flavor

        Use `deis help [command]` to learn more
        """
        return self.flavors_list(args)

    def flavors_create(self, args):
        """
        Create a new node flavor

        Usage: deis flavors:create <id> --provider=<provider> --params=<params>
        """
        body = {'id': args.get('<id>'), 'provider': args.get('--provider'),
                'params': args.get('--params', json.dumps({}))}
        response = self._dispatch('post', '/api/flavors', json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print("{0[id]}".format(response.json()))
        else:
            raise ResponseError(response)

    def flavors_delete(self, args):
        """
        Delete a node flavor

        Usage: deis flavors:delete <id>
        """
        flavor = args.get('<id>')
        response = self._dispatch('delete', "/api/flavors/{}".format(flavor))
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            pass
        else:
            print('Error!', response.status_code, response.text)

    def flavors_info(self, args):
        """
        Print information about a node flavor

        Usage: deis flavors:info <flavor>
        """
        flavor = args.get('<flavor>')
        response = self._dispatch('get', "/api/flavors/{}".format(flavor))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def flavors_list(self, args):
        """
        List available node flavors

        Usage: deis flavors:list
        """
        response = self._dispatch('get', '/api/flavors')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            if data['count'] == 0:
                print('No flavors found')
                return
            print("=== {owner} Flavors".format(**data['results'][0]))
            for item in data['results']:
                print("{id}: params => {params}".format(**item))
        else:
            raise ResponseError(response)

    def flavors_update(self, args):
        """
        Update an existing node flavor

        Usage: deis flavors:update <id> [<params>] [--provider=<provider>]
        """
        id_ = args.get('<id>')
        body = {'id': id_}
        params = args.get('<params>')
        if params:
            body['params'] = params
        provider = args.get('--provider')
        if provider:
            body['provider'] = provider
        response = self._dispatch(
            'patch', '/api/flavors/{}'.format(id_), json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def formations(self, args):
        """
        Valid commands for formations:

        formations:create        create a new container formation
        formations:list          list accessible formations
        formations:update        update formation fields
        formations:info          print a represenation of the formation
        formations:converge      force-converge all nodes in the formation
        formations:calculate     calculate and display the formation databag
        formations:destroy       destroy a container formation

        Use `deis help [command]` to learn more
        """
        return self.formations_list(args)

    def formations_calculate(self, args, quiet=False):
        """
        Calculate and display the formation's databag

        This command will calculate the databag and return
        its JSON representation.

        Usage: deis formations:calculate <formation>
        """
        formation = args.get('<formation>')
        response = self._dispatch('post',
                                  "/api/formations/{}/calculate".format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            databag = json.loads(response.content)
            if quiet is False:
                print(json.dumps(databag, indent=2))
            return databag
        else:
            raise ResponseError(response)

    def formations_converge(self, args):
        """
        Force converge a formation

        Converging a formation will force a converge on all nodes in the
        formation, ensuring it is completely up-to-date.

        Usage: deis formations:converge <id>
        """
        formation = args.get('<id>')
        sys.stdout.write('Converging {} formation... '.format(formation))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('post',
                                      "/api/formations/{}/converge".format(formation))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
            databag = json.loads(response.content)
            print(json.dumps(databag, indent=2))
        else:
            raise ResponseError(response)

    def formations_create(self, args):
        """
        Create a new container formation

        A globally unique formation ID must be provided.

        If a flavor is provided, a default layer will be initialized
        with dual proxy and runtime capability, faciliating a simple
        single-node formation with a `deis nodes:scale` command.

        The name of the default layer is "runtime" unless overriden
        with the --layer=<layer> option.

        The domain field is required for a single formation to host
        multiple applications.  Note this requires wildcard DNS
        configuration on the provided domain.

        For example: --domain=deisapp.com requires that \\*.deisapp.com\\
        resolve to the formation's proxy nodes.

        Usage: deis formations:create <id> [--flavor=<flavor>] [--domain=<domain> --layer=<layer>]
        """
        body = {'id': args['<id>'], 'domain': args.get('--domain')}
        flavor = args.get('--flavor')
        if flavor:
            response = self._dispatch('get', '/api/flavors/{}'.format(flavor))
            if response.status_code != 200:
                print('Flavor not found')
                return
        sys.stdout.write('Creating formation... ')
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('post', '/api/formations',
                                      json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            data = response.json()
            formation = data['id']
            print("done, created {}".format(formation))
            if flavor:
                layer = args.get('--layer') or 'runtime'
                self.layers_create({'<formation>': formation, '<id>': layer,
                                    '<flavor>': flavor, '--proxy': True, '--runtime': True})
                print('\nUse `deis nodes:scale {formation} {layer}=1` '
                      'to scale a basic formation'.format(**locals()))
            else:
                print('\nSee `deis help layers:create` to begin '
                      'building your formation'.format(**locals()))
        else:
            raise ResponseError(response)

    def formations_info(self, args):
        """
        Print info about a formation

        Usage: deis formations:info <id>
        """
        formation = args.get('<id>')
        response = self._dispatch('get', "/api/formations/{}".format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            print("=== {} Formation".format(formation))
            print(json.dumps(response.json(), indent=2))
            print()
            args = {'<formation>': data['id']}
            self.layers_list(args)
            print()
            self.nodes_list(args)
            print()
        else:
            raise ResponseError(response)

    def formations_list(self, args):
        """
        List available formations

        Usage: deis formations:list
        """
        response = self._dispatch('get', '/api/formations')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            print("=== Formations")
            for item in data['results']:
                print("{id} {nodes}".format(**item))
        else:
            raise ResponseError(response)

    def formations_destroy(self, args):
        """
        Destroy a formation

        Usage: deis formations:destroy <id> [--confirm=<confirm>]
        """
        formation = args.get('<id>')
        confirm = args.get('--confirm')
        if confirm == formation:
            pass
        else:
            print("""
 !    WARNING: Potentially Destructive Action
 !    This command will destroy the formation: {formation}
 !    To proceed, type "{formation}" or re-run this command with --confirm={formation}
""".format(**locals()))
            confirm = raw_input('> ').strip('\n')
            if confirm != formation:
                print('Destroy aborted')
                return
        sys.stdout.write("Destroying formation... ".format(formation))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('delete', "/api/formations/{}".format(formation))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code in (requests.codes.no_content,  # @UndefinedVariable
                                    requests.codes.not_found):  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
        else:
            raise ResponseError(response)

    def formations_update(self, args):
        """
        Update formation fields

        This is typically used to add a "domain" to to host
        multiple applications on a single formation

        Usage: deis formations:update <id> [--domain=<domain>]
        """
        formation = args['<id>']
        domain = args.get('--domain')
        body = {'domain': domain}
        response = self._dispatch('patch', '/api/formations/{}'.format(formation),
                                  json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def keys(self, args):
        """
        Valid commands for SSH keys:

        keys:list        list SSH keys for the logged in user
        keys:add         add an SSH key
        keys:remove      remove an SSH key

        Use `deis help [command]` to learn more
        """
        return self.keys_list(args)

    def keys_add(self, args):
        """
        Add SSH keys for the logged in user

        Usage: deis keys:add [<key>]
        """

        Key = namedtuple('Key', 'path name type str comment')

        def parse_key(path):
            """Parse an SSH public key path into a Key namedtuple."""
            name = path.split(os.path.sep)[-1]
            with open(path) as f:
                data = f.read()
                match = re.match(r'^(ssh-...) ([^ ]+) ?(.*)', data)
                if not match:
                    print("Could not parse SSH public key {0}".format(name))
                    return
                key_type, key_str, key_comment = match.groups()
                return Key(path, name, key_type, key_str, key_comment)

        path = args.get('<key>')
        if not path:
            # find public keys and prompt the user to pick one
            ssh_dir = os.path.expanduser('~/.ssh')
            pubkey_paths = glob.glob(os.path.join(ssh_dir, '*.pub'))
            if not pubkey_paths:
                print('No SSH public keys found')
                return
            pubkeys_list = [parse_key(k) for k in pubkey_paths]
            print('Found the following SSH public keys:')
            for i, key_ in enumerate(pubkeys_list):
                print("{}) {} {}".format(i + 1, key_.name, key_.comment))
            inp = raw_input('Which would you like to use with Deis? ')
            try:
                selected_key = pubkeys_list[int(inp) - 1]
            except:
                print('Aborting')
                return
        else:
            # check the specified key format
            selected_key = parse_key(path)
            if not selected_key:
                return
        # Upload the key to Deis
        if selected_key.comment:
            key_id = selected_key.comment
        else:
            key_id = selected_key.name.replace('.pub', '')
        body = {
            'id': key_id,
            'public': "{} {}".format(selected_key.type, selected_key.str)
        }
        sys.stdout.write("Uploading {} to Deis...".format(key_id))
        sys.stdout.flush()
        response = self._dispatch('post', '/api/keys', json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('done')
        else:
            raise ResponseError(response)

    def keys_list(self, args):
        """
        List SSH keys for the logged in user

        Usage: deis keys:list
        """
        response = self._dispatch('get', '/api/keys')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            if data['count'] == 0:
                print('No keys found')
                return
            print("=== {owner} Keys".format(**data['results'][0]))
            for key in data['results']:
                public = key['public']
                print("{0} {1}...{2}".format(
                    key['id'], public[0:16], public[-10:]))
        else:
            raise ResponseError(response)

    def keys_remove(self, args):
        """
        Remove an SSH key for the logged in user

        Usage: deis keys:remove <key>
        """
        key = args.get('<key>')
        sys.stdout.write("Removing {} SSH Key... ".format(key))
        sys.stdout.flush()
        response = self._dispatch('delete', "/api/keys/{}".format(key))
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print('done')
        else:
            raise ResponseError(response)

    def layers(self, args):
        """
        Valid commands for node layers:

        layers:create        create a layer of nodes for a formation
        layers:list          list layers in a formation
        layers:info          print info about a particular layer
        layers:destroy       destroy a layer of nodes in a formation

        Use `deis help [command]` to learn more
        """
        sys.argv[1] = 'layers:list'
        args = docopt(self.layers_list.__doc__)
        return self.layers_list(args)

    def layers_create(self, args):
        """
        Create a layer of nodes

        Usage: deis layers:create <formation> <id> <flavor> [options]

        Options:
        --proxy=<yn>                    layer can be used for proxy [default: y]
        --runtime=<yn>                  layer can be used for runtime [default: y]
        --ssh_username=USERNAME         username for ssh connections [default: ubuntu]
        --ssh_private_key=PRIVATE_KEY   private key for ssh comm (default: auto-gen)
        --ssh_public_key=PUBLIC_KEY     public key for ssh comm (default: auto-gen)
        --ssh_port=<port>               port number for ssh comm (default: 22)

        """
        formation = args.get('<formation>')
        body = {'id': args['<id>'], 'flavor': args['<flavor>']}
        for opt in ('--formation', '--ssh_username', '--ssh_private_key',
                    '--ssh_public_key'):
            o = args.get(opt)
            if o:
                body.update({opt.strip('-'): o})
        o = args.get('--ssh_port')
        if o:
            body.update({'ssh_port': int(o)})
        for opt in ('--proxy', '--runtime'):
            o = args.get(opt)
            if o and str(o).lower() in ['n', 'no', 'f', 'false', '0', 'off']:
                body.update({opt.strip('-'): False})
            else:
                body.update({opt.strip('-'): True})
        sys.stdout.write("Creating {} layer... ".format(args['<id>']))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('post', "/api/formations/{}/layers".format(formation),
                                      json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
        else:
            raise ResponseError(response)

    def layers_destroy(self, args):
        """
        Destroy a layer of nodes

        Usage: deis layers:destroy <formation> <id>
        """
        formation = args.get('<formation>')
        layer = args['<id>']
        sys.stdout.write("Destroying {layer} layer... ".format(**locals()))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch(
                'delete', "/api/formations/{formation}/layers/{layer}".format(**locals()))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
        else:
            raise ResponseError(response)

    def layers_info(self, args):
        """
        Print info about a layer of nodes

        Usage: deis layers:info <formation> <id>
        """
        formation = args.get('<formation>')
        layer = args.get('<id>')
        response = self._dispatch('get', "/api/formations/{}/layers/{}".format(formation, layer))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def layers_list(self, args):
        """
        List a formation's layers

        Usage: deis layers:list <formation>
        """
        formation = args.get('<formation>')
        response = self._dispatch('get',
                                  "/api/formations/{}/layers".format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print("=== {} Layers".format(formation))
            data = response.json()
            format_str = "{id} => flavor: {flavor}, proxy: {proxy}, runtime: {runtime}"
            for item in data['results']:
                print(format_str.format(**item))
        else:
            raise ResponseError(response)

    def layers_update(self, args):
        """
        Create a layer of nodes

        Usage: deis layers:update <formation> <id> [options]

        Options:

        --proxy=<yn>                    layer can be used for proxy [default: y]
        --runtime=<yn>                  layer can be used for runtime [default: y]
        --ssh_username=USERNAME         username for ssh connections [default: ubuntu]
        --ssh_private_key=PRIVATE_KEY   private key for ssh comm (default: auto-gen)
        --ssh_public_key=PUBLIC_KEY     public key for ssh comm (default: auto-gen)
        --ssh_port=<port>               port number for ssh comm (default: 22)

        """
        formation = args.get('<formation>')
        layer = args['<id>']
        body = {'id': args['<id>']}
        for opt in ('--ssh_username', '--ssh_private_key', '--ssh_public_key',
                    '--ssh_port'):
            o = args.get(opt)
            if o:
                body.update({opt.strip('-'): o})
        o = args.get('--ssh_port')
        if o:
            body.update({'ssh_port': int(o)})
        for opt in ('--proxy', '--runtime'):
            o = args.get(opt)
            if o is not None:
                if str(o).lower() in ['n', 'no', 'f', 'false', '0', 'off']:
                    body.update({opt.strip('-'): False})
                else:
                    body.update({opt.strip('-'): True})
        sys.stdout.write("Updating {} layer... ".format(args['<id>']))
        sys.stdout.flush()
        response = self._dispatch(
            'patch', "/api/formations/{formation}/layers/{layer}".format(**locals()),
            json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('done.')
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def nodes(self, args):
        """
        Valid commands for nodes:

        nodes:list            list nodes in a formation
        nodes:info            print info for a given node
        nodes:scale           scale nodes by layer (e.g. runtime=4)
        nodes:converge        force-converge a node and return the output
        nodes:ssh             ssh directly into a node
        nodes:destroy         destroy a node by ID

        Use `deis help [command]` to learn more
        """
        sys.argv[1] = 'nodes:list'
        args = docopt(self.nodes_list.__doc__)
        return self.nodes_list(args)

    def nodes_create(self, args):
        """
        Add an existing node to a formation.

        Usage: deis nodes:create <formation> <fqdn> --layer=<layer>
        """
        formation = args.get('<formation>')
        fqdn, layer = args.get('<fqdn>'), args.get('--layer')
        body = {'fqdn': fqdn, 'layer': layer}
        sys.stdout.write("Creating node for {}... ".format(fqdn))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('post', "/api/formations/{}/nodes".format(formation),
                                      json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
        else:
            raise ResponseError(response)

    def nodes_info(self, args):
        """
        Print info about a particular node

        Usage: deis nodes:info <node>
        """
        node = args.get('<node>')
        response = self._dispatch('get', "/api/nodes/{}".format(node))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def nodes_list(self, args):
        """
        List nodes in a formation

        Usage: deis nodes:list <formation>
        """
        formation = args.get('<formation>')
        response = self._dispatch('get',
                                  "/api/formations/{}/nodes".format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print("=== {} Nodes".format(formation))
            data = response.json()
            format_str = "{id} {provider_id} {fqdn}"
            for item in data['results']:
                print(format_str.format(**item))
        else:
            raise ResponseError(response)

    def nodes_destroy(self, args):
        """
        Destroy a node by ID

        Usage: deis nodes:destroy <id>
        """
        node = args['<id>']
        sys.stdout.write("Destroying {}... ".format(node))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch(
                'delete', "/api/nodes/{node}".format(**locals()))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
        else:
            raise ResponseError(response)

    def nodes_scale(self, args):
        """
        Scale nodes in a formation

        Scaling nodes will launch or terminate nodes to meet the
        requested structure.  For example, to scale up to 4 nodes
        in the "dev" formation's runtime layer:

        ``deis nodes:scale dev runtime=4``

        Usage: deis nodes:scale <formation> <type=num>...
        """
        formation = args.get('<formation>')
        body = {}
        runtimes = True
        for type_num in args.get('<type=num>'):
            typ, count = type_num.split('=')
            if (typ, count) == ('runtime', '0'):
                runtimes = False
            body.update({typ: int(count)})
        print('Scaling nodes... but first, coffee!')
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('post',
                                      "/api/formations/{}/scale".format(formation),
                                      json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('done in {}s\n'.format(int(time.time() - before)))
            if runtimes:
                print('Use `deis create --formation={}` to create an application'.format(
                      formation))
        else:
            raise ResponseError(response)

    def nodes_converge(self, args):
        """
        Force converge a node

        Converging a node will force a chef-client run and
        return its output

        Usage: deis nodes:converge <id>
        """
        node = args.get('<id>')
        sys.stdout.write('Converging {} node... '.format(node))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('post',
                                      "/api/nodes/{}/converge".format(node))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
            output = json.loads(response.content)
            print(output)
        else:
            raise ResponseError(response)

    def nodes_ssh(self, args):
        """
        SSH into a node and optionally run a command

        Usage: deis nodes:ssh <node> [<command>...]
        """
        node = args.get('<node>')
        response = self._dispatch('get',
                                  "/api/nodes/{node}".format(**locals()))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            node = response.json()
            response = self._dispatch('get',
                                      '/api/formations/{formation}/layers/{layer}'.format(**node))
            if response.status_code != requests.codes.ok:  # @UndefinedVariable
                raise ResponseError(response)
            layer = response.json()
            _, key_path = tempfile.mkstemp()
            os.chmod(key_path, 0600)
            with open(key_path, 'w') as f:
                f.write(layer['ssh_private_key'])
            ssh_args = ['-o UserKnownHostsFile=/dev/null', '-o StrictHostKeyChecking=no',
                        '-i', key_path, '{}@{}'.format(layer['ssh_username'], node['fqdn'])]
            command = args.get('<command>')
            if command:
                ssh_args.extend(command)
            os.execvp('ssh', ssh_args)
        else:
            raise ResponseError(response)

    def perms(self, args):
        """
        Valid commands for perms:

        perms:list            list permissions granted on an app or formation
        perms:create          create a new permission for a user
        perms:delete          delete a permission for a user

        Use `deis help perms:[command]` to learn more
        """
        # perms:transfer        transfer ownership of an app or formation
        return self.perms_list(args)

    def perms_list(self, args):
        """
        List all users with permission to use an app, or list all users
        with system administrator privileges.

        Usage: deis perms:list [--app=<app>|--admin]
        """
        app, url = self._parse_perms_args(args)
        response = self._dispatch('get', url)
        if response.status_code == requests.codes.ok:
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def perms_create(self, args):
        """
        Give another user permission to use an app, or give another user
        system administrator privileges.

        Usage: deis perms:create <username> [--app=<app>|--admin]
        """
        app, url = self._parse_perms_args(args)
        username = args.get('<username>')
        body = {'username': username}
        if app:
            msg = "Adding {} to {} collaborators... ".format(username, app)
        else:
            msg = "Adding {} to system administrators... ".format(username)
        sys.stdout.write(msg)
        sys.stdout.flush()
        response = self._dispatch('post', url, json.dumps(body))
        if response.status_code == requests.codes.created:
            print('done')
        else:
            raise ResponseError(response)

    def perms_delete(self, args):
        """
        Revoke another user's permission to use an app, or revoke another
        user's system administrator privileges.

        Usage: deis perms:delete <username> [--app=<app>|--admin]
        """
        app, url = self._parse_perms_args(args)
        username = args.get('<username>')
        url = "{}/{}".format(url, username)
        if app:
            msg = "Removing {} from {} collaborators... ".format(username, app)
        else:
            msg = "Remove {} from system administrators... ".format(username)
        sys.stdout.write(msg)
        sys.stdout.flush()
        response = self._dispatch('delete', url)
        if response.status_code == requests.codes.no_content:
            print('done')
        else:
            raise ResponseError(response)

    def _parse_perms_args(self, args):
        app = args.get('--app'),
        admin = args.get('--admin')
        if admin:
            app = None
            url = '/api/admin/perms'
        else:
            app = app[0] or self._session.app
            url = "/api/apps/{}/perms".format(app)
        return app, url

    def providers(self, args):
        """
        Valid commands for providers:

        providers:list        list available providers for the logged in user
        providers:discover    discover provider credentials using envvars
        providers:create      create a new provider for use by deis
        providers:info        print information about a specific provider

        Use `deis help [command]` to learn more
        """
        return self.providers_list(args)

    def providers_create(self, args):  # noqa
        """
        Create a provider for use by Deis

        This command is only necessary when adding a duplicate set of
        credentials for a provider. User accounts start with empty providers,
        EC2, Rackspace, and DigitalOcean by default, which should be updated
        in place.

        Use `providers:discover` to update the credentials for the default
        providers created with your account.

        Usage: deis providers:create <id> <type> <creds>
        """
        type = args.get('<type>')  # @ReservedAssignment
        if type in ['ec2', 'rackspace', 'digitalocean']:
            creds = get_provider_creds(type, raise_error=True)
        else:
            creds = json.loads(args.get('<creds>'))
        id = args.get('<id>')  # @ReservedAssignment
        if not id:
            id = type  # @ReservedAssignment
        body = {'id': id, 'type': type, 'creds': json.dumps(creds)}
        response = self._dispatch('post', '/api/providers',
                                  json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print("{0[id]}".format(response.json()))
        else:
            raise ResponseError(response)

    def providers_discover(self, args):  # noqa
        """
        Discover and update provider credentials

        This command will discover provider credentials using
        standard environment variables like AWS_ACCESS_KEY and
        AWS_SECRET_KEY.  It will use those credentials to update
        the existing provider record, allowing you to use
        pre-installed node flavors.

        Usage: deis providers:discover
        """
        provider_data = [
            # Provider, human-redable provider name, sample field to display
            ('ec2', 'EC2', 'access_key'),
            ('rackspace', 'Rackspace', 'api_key'),
            ('digitalocean', 'DigitalOcean', 'api_key'),
        ]
        for provider, name, field in provider_data:
            creds = get_provider_creds(provider)
            if creds:
                print ("Discovered {} credentials: {}".format(name, creds[field]))
                inp = raw_input("Import {} credentials? (y/n) : ".format(name))
                if inp.lower().strip('\n') != 'y':
                    print('Aborting.')
                else:
                    body = {'creds': json.dumps(creds)}
                    sys.stdout.write("Uploading {} credentials... ".format(name))
                    sys.stdout.flush()
                    endpoint = "/api/providers/{}".format(provider)
                    response = self._dispatch('patch', endpoint, json.dumps(body))
                    if response.status_code == requests.codes.ok:  # @UndefinedVariable
                        print('done')
                    else:
                        raise ResponseError(response)
            else:
                print("No {} credentials discovered.".format(name))

        # Check for locally booted Deis Controller VM
        try:
            running_vms = subprocess.check_output(
                ['vboxmanage', 'list', 'runningvms'],
                stderr=subprocess.PIPE
            )
        except subprocess.CalledProcessError:
            running_vms = ""
        # Vagrant internally names a running VM using the folder name in which the Vagrantfile
        # resides, eg; my-deis-code-folder_default_1383326629
        try:
            deis_codebase_folder = self._session.git_root().split('/')[-1]
        except EnvironmentError:
            deis_codebase_folder = 'deis'
        if deis_codebase_folder and deis_codebase_folder in running_vms:
            print("Discovered locally running Deis Controller VM")
            # In order for the Controller to be able to boot Vagrant VMs it needs to run commands
            # on the host machine. It does this via an SSH server. In order to access that server
            # we need to send the current user's name and host.
            try:
                user = subprocess.check_output(
                    "whoami",
                    stderr=subprocess.PIPE
                ).strip()
                hostname = subprocess.check_output(
                    "hostname",
                    stderr=subprocess.PIPE,
                    shell=True
                ).strip()
            except subprocess.CalledProcessError:
                print("Error detecting username and host address.")
                sys.exit(1)
            if not hostname.endswith('.local'):
                hostname += '.local'
            creds = {
                'user': user,
                'host': hostname
            }
            body = {'creds': json.dumps(creds)}
            sys.stdout.write('Activating Vagrant as a provider... ')
            sys.stdout.flush()
            response = self._dispatch('patch', '/api/providers/vagrant',
                                      json.dumps(body))
            if response.status_code == requests.codes.ok:  # @UndefinedVariable
                print('done')
            else:
                raise ResponseError(response)
        else:
            print("No Vagrant VMs discovered.")

    def providers_info(self, args):
        """
        Print information about a specific provider

        Usage: deis providers:info <provider>
        """
        provider = args.get('<provider>')
        response = self._dispatch('get', "/api/providers/{}".format(provider))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def providers_list(self, args):
        """
        List providers for the logged in user

        Usage: deis providers:list
        """
        response = self._dispatch('get', '/api/providers')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            if data['count'] == 0:
                print('No providers found')
                return
            print("=== {owner} Providers".format(**data['results'][0]))
            for item in data['results']:
                creds = json.loads(item['creds'])
                if 'secret_key' in creds:
                    creds.pop('secret_key')
                print("{} => {}".format(item['id'], creds))
        else:
            raise ResponseError(response)

    def releases(self, args):
        """
        Valid commands for releases:

        releases:list        list an application's release history
        releases:info        print information about a specific release
        releases:rollback    coming soon!

        Use `deis help [command]` to learn more
        """
        return self.releases_list(args)

    def releases_info(self, args):
        """
        Print info about a particular release

        Usage: deis releases:info <version> [--app=<app>]
        """
        version = args.get('<version>')
        if not version.startswith('v'):
            version = 'v' + version
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch(
            'get', "/api/apps/{app}/releases/{version}".format(**locals()))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            raise ResponseError(response)

    def releases_list(self, args):
        """
        List release history for an application

        Usage: deis releases:list [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch('get', "/api/apps/{app}/releases".format(**locals()))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print("=== {} Releases".format(app))
            data = response.json()
            for item in data['results']:
                item['created'] = readable_datetime(item['created'])
                print("v{version:<6} {created:<33} {summary}".format(**item))
        else:
            raise ResponseError(response)

    def releases_rollback(self, args):
        """
        Roll back to a previous application release.

        Usage: deis releases:rollback [--app=<app>] [<version>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        version = args.get('<version>')
        if version:
            if version.startswith('v'):
                version = version[1:]
            body = {'version': int(version)}
        else:
            body = {}
        url = "/api/apps/{app}/releases/rollback".format(**locals())
        response = self._dispatch('post', url, json.dumps(body))
        if response.status_code == requests.codes.created:
            print(response.json())
        else:
            raise ResponseError(response)


def parse_args(cmd):
    """
    Parse command-line args applying shortcuts and looking for help flags
    """
    shortcuts = {
        'register': 'auth:register',
        'login': 'auth:login',
        'logout': 'auth:logout',
        'create': 'apps:create',
        'destroy': 'apps:destroy',
        'ps': 'containers:list',
        'info': 'apps:info',
        'scale': 'containers:scale',
        'converge': 'formations:converge',
        'calculate': 'apps:calculate',
        'ssh': 'nodes:ssh',
        'open': 'apps:open',
        'logs': 'apps:logs',
        'rollback': 'releases:rollback',
        'run': 'apps:run',
        'sharing': 'perms:list',
        'sharing:list': 'perms:list',
        'sharing:add': 'perms:create',
        'sharing:remove': 'perms:delete',
    }
    if cmd == 'help':
        cmd = sys.argv[-1]
        help_flag = True
    else:
        cmd = sys.argv[1]
        help_flag = False
    # swap cmd with shortcut
    if cmd in shortcuts:
        cmd = shortcuts[cmd]
        # change the cmdline arg itself for docopt
        if not help_flag:
            sys.argv[1] = cmd
        else:
            sys.argv[2] = cmd
    # convert : to _ for matching method names and docstrings
    if ':' in cmd:
        cmd = '_'.join(cmd.split(':'))
    return cmd, help_flag


def _dispatch_cmd(method, args):
    try:
        method(args)
    except requests.exceptions.ConnectionError as err:
        print("Couldn't connect to the Deis Controller. Make sure that the Controller URI is \
correct and the server is running.")
        sys.exit(1)
    except EnvironmentError as err:
        raise DocoptExit(err.message)
    except ResponseError as err:
        resp = err.message
        print('{} {}'.format(resp.status_code, resp.reason))
        try:
            msg = resp.json()
            if 'detail' in msg:
                msg = "Detail:\n{}".format(msg['detail'])
        except:
            msg = resp.text
        print(msg)
        sys.exit(1)


def main():
    """
    Create a client, parse the arguments received on the command line, and
    call the appropriate method on the client.
    """
    cli = DeisClient()
    args = docopt(__doc__, version='Deis CLI {}'.format(__version__),
                  options_first=True)
    cmd = args['<command>']
    cmd, help_flag = parse_args(cmd)
    # print help if it was asked for
    if help_flag:
        if cmd != 'help' and cmd in dir(cli):
            print(trim(getattr(cli, cmd).__doc__))
            return
        docopt(__doc__, argv=['--help'])
    # unless cmd needs to use sys.argv directly
    if hasattr(cli, cmd):
        method = getattr(cli, cmd)
    else:
        raise DocoptExit('Found no matching command, try `deis help`')
    # re-parse docopt with the relevant docstring unless it needs sys.argv
    if cmd not in ('apps_run',):
        docstring = trim(getattr(cli, cmd).__doc__)
        if 'Usage: ' in docstring:
            args.update(docopt(docstring))
    # dispatch the CLI command
    _dispatch_cmd(method, args)


if __name__ == '__main__':
    main()
    sys.exit(0)
