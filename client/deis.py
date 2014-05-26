#!/usr/bin/env python
"""
The Deis command-line client issues API calls to a Deis controller.

Usage: deis <command> [<args>...]

Auth commands::

  register      register a new user with a controller
  login         login to a controller
  logout        logout from the current controller

Subcommands, use ``deis help [subcommand]`` to learn more::

  apps          manage applications used to provide services
  clusters      manage clusters used to host applications
  ps            manage processes inside an app container
  config        manage environment variables that define app config
  domains       manage and assign domain names to your applications
  builds        manage builds created using `git push`
  releases      manage releases of an application

  keys          manage ssh keys used for `git push` deployments
  perms         manage permissions for shared apps and clusters

Developer shortcut commands::

  create        create a new application
  scale         scale processes by type (web=2, worker=1)
  info          view information about the current app
  open          open a URL to the app in a browser
  logs          view aggregated log info for the app
  run           run a command in an ephemeral app container
  destroy       destroy an application

Use ``git push deis master`` to deploy to an application.

"""

from __future__ import print_function
from collections import namedtuple
from collections import OrderedDict
from cookielib import MozillaCookieJar
from datetime import datetime
from getpass import getpass
from itertools import cycle
from threading import Event
from threading import Thread
import base64
import glob
import json
import locale
import os.path
import re
import subprocess
import sys
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

__version__ = 'v0.9.0'


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
        m = re.match('\S+/(?P<app>[a-z0-9-]+)(.git)?$', url)
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
            var, val = arg.split('=', 1)
        except ValueError:
            raise DocoptExit()
        # Try to coerce the value to an int since that's a common use case
        try:
            data[var] = int(val)
        except ValueError:
            data[var] = val
    return data


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
    yesterday = now + relativedelta.relativedelta(days= -1)
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

    def apps_create(self, args):
        """
        Create a new application

        If no ID is provided, one will be generated automatically.
        If no cluster is provided, a cluster named "dev" will be used.

        Usage: deis apps:create [<id> --cluster=<cluster> --no-remote] [options]

        Options

        --cluster=CLUSTER      target cluster to host application [default: dev]
        --no-remote            do not create a 'deis' git remote
        """
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
        body = {}
        app_name = args.get('<id>')
        if app_name:
            body.update({'id': app_name})
        cluster = args.get('--cluster')
        if cluster:
            body.update({'cluster': cluster})
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
            # TODO: retrieve the hostname from service discovery
            hostname = urlparse.urlparse(self._settings['controller']).netloc.split(':')[0]
            git_remote = "ssh://git@{hostname}:2222/{app_id}.git".format(**locals())
            if args.get('--no-remote'):
                print('remote available at {}'.format(git_remote))
            else:
                try:
                    subprocess.check_call(
                        ['git', 'remote', 'add', '-f', 'deis', git_remote],
                        stdout=subprocess.PIPE)
                    print('Git remote deis added')
                except subprocess.CalledProcessError:
                    print('Could not create Deis remote')
                    sys.exit(1)
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
            # If the requested app is in the current dir, delete the git remote
            try:
                if app == self._session.app:
                    subprocess.check_call(
                        ['git', 'remote', 'rm', 'deis'],
                        stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                    print('Git remote deis removed')
            except (EnvironmentError, subprocess.CalledProcessError):
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
                print('{id}'.format(**item))
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
            self.ps_list(args)
            self.domains_list(args)
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
        # TODO: replace with a single API call to apps endpoint
        response = self._dispatch('get', "/api/apps/{}".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            cluster = response.json()['cluster']
        else:
            raise ResponseError(response)
        response = self._dispatch('get', "/api/clusters/{}".format(cluster))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            domain = response.json()['domain']
            # use the OS's default handler to open this URL
            webbrowser.open('http://{}.{}/'.format(app, domain))
            return domain
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
        app = args.get('--app')
        if not app:
            app = self._session.app
        body = {'command': ' '.join(sys.argv[2:])}
        response = self._dispatch('post',
                                  "/api/apps/{}/run".format(app),
                                  json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            rc, output = json.loads(response.content)
            sys.stdout.write(output)
            sys.stdout.flush()
            sys.exit(rc)
        else:
            raise ResponseError(response)

    def auth(self, args):
        """
        Valid commands for auth:

        auth:register          register a new user
        auth:cancel            remove the current account
        auth:login             authenticate against a controller
        auth:logout            clear the current user session

        Use `deis help [command]` to learn more
        """
        return

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
                sys.exit(1)
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
            sys.exit(1)
        return True

    def auth_cancel(self, args):
        """
        Cancel and remove the current account.

        Usage: deis auth:cancel
        """
        controller = self._settings.get('controller')
        if not controller:
            print('Not logged in to a Deis controller')
            sys.exit(1)
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
            self._session.cookies.clear()
            self._session.cookies.save()
            raise ResponseError(response)

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

        builds:list        list build history for an application
        builds:create      coming soon!

        Use `deis help [command]` to learn more
        """
        return self.builds_list(args)

    def builds_create(self, args):
        """
        Create a new build of an application

        Usage: deis builds:create <image> [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        body = {'image': args['<image>']}
        sys.stdout.write('Creating build... ')
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('post', "/api/apps/{}/builds".format(app), json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            version = response.headers['x-deis-release']
            print("done, v{}".format(version))
        else:
            raise ResponseError(response)

    def builds_list(self, args):
        """
        List build history for an application

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

    def clusters(self, args):
        """
        Valid commands for clusters:

        clusters:create        create a new cluster
        clusters:list          list accessible clusters
        clusters:update        update cluster fields
        clusters:info          print a represenation of the cluster
        clusters:destroy       destroy a cluster

        Use `deis help [command]` to learn more
        """
        return self.clusters_list(args)

    def clusters_create(self, args):
        """
        Create a new cluster

        A globally unique cluster ID must be provided.

        A domain field must also be provided to support multiple
        applications hosted on the cluster.  Note this requires
        wildcard DNS configuration on the domain.

        For example, a domain of "deisapp.com" requires that \\*.deisapp.com\\
        resolve to the cluster's router endpoints.

        Usage: deis clusters:create <id> <domain> --hosts=<hosts> --auth=<auth> [options]

        Parameters:

        <id>             a name for the cluster
        <domain>         a domain under which app hostnames will live
        <hosts>          a comma-separated list of cluster members
        <auth>           a path to an SSH private key used to connect to cluster members

        Options:

        --type=TYPE      cluster type [default: coreos]
        """
        body = {'id': args['<id>'], 'domain': args['<domain>'],
                'hosts': args['--hosts'], 'type': args['--type']}
        auth_path = os.path.expanduser(args['--auth'])
        if not os.path.exists(auth_path):
            print('Path to authentication credentials does not exist: {}'.format(auth_path))
            sys.exit(1)
        with open(auth_path) as f:
            data = f.read()
        body.update({'auth': base64.b64encode(data)})
        sys.stdout.write('Creating cluster... ')
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('post', '/api/clusters', json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            data = response.json()
            cluster = data['id']
            print("done, created {}".format(cluster))
        else:
            raise ResponseError(response)

    def clusters_info(self, args):
        """
        Print info about a cluster

        Usage: deis clusters:info <id>
        """
        cluster = args.get('<id>')
        response = self._dispatch('get', "/api/clusters/{}".format(cluster))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print("=== {} Cluster".format(cluster))
            print(json.dumps(response.json(), indent=2))
            print()
        else:
            raise ResponseError(response)

    def clusters_list(self, args):
        """
        List available clusters

        Usage: deis clusters:list
        """
        response = self._dispatch('get', '/api/clusters')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            print("=== Clusters")
            for item in data['results']:
                print("{id}".format(**item))
        else:
            raise ResponseError(response)

    def clusters_destroy(self, args):
        """
        Destroy a cluster

        Usage: deis clusters:destroy <id> [--confirm=<confirm>]
        """
        cluster = args.get('<id>')
        confirm = args.get('--confirm')
        if confirm == cluster:
            pass
        else:
            print("""
 !    WARNING: Potentially Destructive Action
 !    This command will destroy the cluster: {cluster}
 !    To proceed, type "{cluster}" or re-run this command with --confirm={cluster}
""".format(**locals()))
            confirm = raw_input('> ').strip('\n')
            if confirm != cluster:
                print('Destroy aborted')
                return
        sys.stdout.write("Destroying cluster... ".format(cluster))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            before = time.time()
            response = self._dispatch('delete', "/api/clusters/{}".format(cluster))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code in (requests.codes.no_content,  # @UndefinedVariable
                                    requests.codes.not_found):  # @UndefinedVariable
            print('done in {}s'.format(int(time.time() - before)))
        else:
            raise ResponseError(response)

    def clusters_update(self, args):
        """
        Update cluster fields

        Usage: deis clusters:update <id> [--domain=<domain> --hosts=<hosts> --auth=<auth>] [options]

        Options:

        --type=TYPE      cluster type [default: coreos]
        """
        cluster = args['<id>']
        body = {}
        for k, arg in (('domain', '--domain'), ('hosts', '--hosts'),
                       ('auth', '--auth'), ('type', '--type')):
            v = args.get(arg)
            if v:
                body.update({k: v})
        response = self._dispatch('patch', '/api/clusters/{}'.format(cluster),
                                  json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
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
        sys.argv[1] = 'config:list'
        args = docopt(self.config_list.__doc__)
        return self.config_list(args)

    def config_list(self, args):
        """
        List environment variables for an application

        Usage: deis config:list [--oneline] [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app

        oneline = args.get('--oneline')
        response = self._dispatch('get', "/api/apps/{}/config".format(app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            config = response.json()
            values = json.loads(config['values'])
            print("=== {} Config".format(app))
            items = values.items()
            if len(items) == 0:
                print('No configuration')
                return
            keys = sorted(values)

            if not oneline:
                width = max(map(len, keys)) + 5
                for k in keys:
                    v = values[k]
                    print(("{k:<" + str(width) + "} {v}").format(**locals()))
            else:
                output = []
                for k in keys:
                    v = values[k]
                    output.append("{k}={v}".format(**locals()))
                print(' '.join(output))
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
        sys.stdout.write('Creating config... ')
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('post', "/api/apps/{}/config".format(app), json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            version = response.headers['x-deis-release']
            print("done, v{}\n".format(version))
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
        sys.stdout.write('Creating config... ')
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('post', "/api/apps/{}/config".format(app), json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            version = response.headers['x-deis-release']
            print("done, v{}\n".format(version))
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

    def domains(self, args):
        """
        Valid commands for domains:

        domains:add           bind a domain to an application
        domains:list          list domains bound to an application
        domains:remove        unbind a domain from an application

        Use `deis help [command]` to learn more
        """
        return self.domains_list(args)

    def domains_add(self, args):
        """
        Bind a domain to an application

        Usage: deis domains:add <domain> [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        domain = args.get('<domain>')
        body = {'domain': domain}
        sys.stdout.write("Adding {domain} to {app}... ".format(**locals()))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('post', "/api/apps/{app}/domains".format(app=app), json.dumps(body))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print("done")
        else:
            raise ResponseError(response)

    def domains_remove(self, args):
        """
        Unbind a domain for an application

        Usage: deis domains:remove <domain> [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        domain = args.get('<domain>')
        sys.stdout.write("Removing {domain} from {app}... ".format(**locals()))
        sys.stdout.flush()
        try:
            progress = TextProgress()
            progress.start()
            response = self._dispatch('delete', "/api/apps/{app}/domains/{domain}".format(**locals()))
        finally:
            progress.cancel()
            progress.join()
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print("done")
        else:
            raise ResponseError(response)

    def domains_list(self, args):
        """
        List domains bound to an application

        Usage: deis domains:list [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.app
        response = self._dispatch(
            'get', "/api/apps/{app}/domains".format(app=app))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            domains = response.json()['results']
            print("=== {} Domains".format(app))
            if len(domains) == 0:
                print('No domains')
                return
            for domain in domains:
                print(domain['domain'])
        else:
            raise ResponseError(response)

    def ps(self, args):
        """
        Valid commands for processes:

        ps:list        list application processes
        ps:scale       scale processes (e.g. web=4 worker=2)

        Use `deis help [command]` to learn more
        """
        sys.argv[1] = 'ps:list'
        args = docopt(self.ps_list.__doc__)
        return self.ps_list(args)

    def ps_list(self, args, app=None):
        """
        List processes servicing an application

        Usage: deis ps:list [--app=<app>]
        """
        if not app:
            app = args.get('--app')
            if not app:
                app = self._session.app
        response = self._dispatch('get',
                                  "/api/apps/{}/containers".format(app))
        if response.status_code != requests.codes.ok:  # @UndefinedVariable
            raise ResponseError(response)
        processes = response.json()
        print("=== {} Processes".format(app))
        c_map = {}
        for item in processes['results']:
            c_map.setdefault(item['type'], []).append(item)
        print()
        for c_type in c_map.keys():
            print("--- {c_type}: ".format(**locals()))
            for c in c_map[c_type]:
                print("{type}.{num} {state} ({release})".format(**c))
            print()

    def ps_scale(self, args):
        """
        Scale an application's processes by type

        Example: deis ps:scale web=4 worker=2

        Usage: deis ps:scale <type=num>... [--app=<app>]
        """
        app = args.get('--app')
        if not app:
            app = self._session.get_app()
        body = {}
        for type_num in args.get('<type=num>'):
            typ, count = type_num.split('=')
            body.update({typ: int(count)})
        print('Scaling processes... but first, coffee!')
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
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print('done in {}s\n'.format(int(time.time() - before)))
            self.ps_list({}, app)
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

        path = args.get('<key>')
        if not path:
            selected_key = self._ask_pubkey_interactively()
        else:
            # check the specified key format
            selected_key = self._parse_key(path)
            if not selected_key:
                return
        # Upload the key to Deis
        body = {
            'id': selected_key.id,
            'public': "{} {}".format(selected_key.type, selected_key.str)
        }
        sys.stdout.write("Uploading {} to Deis...".format(selected_key.id))
        sys.stdout.flush()
        response = self._dispatch('post', '/api/keys', json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('done')
        else:
            raise ResponseError(response)

    def _parse_key(self, path):
        """Parse an SSH public key path into a Key namedtuple."""
        Key = namedtuple('Key', 'path name type str comment id')

        name = path.split(os.path.sep)[-1]
        with open(path) as f:
            data = f.read()
            match = re.match(r'^(ssh-...) ([^ ]+) ?(.*)', data)
            if not match:
                print("Could not parse SSH public key {0}".format(name))
                sys.exit(1)
            key_type, key_str, key_comment = match.groups()
            if key_comment:
                key_id = key_comment
            else:
                key_id = name.replace('.pub', '')
            return Key(path, name, key_type, key_str, key_comment, key_id)

    def _ask_pubkey_interactively(self):
        # find public keys and prompt the user to pick one
        ssh_dir = os.path.expanduser('~/.ssh')
        pubkey_paths = glob.glob(os.path.join(ssh_dir, '*.pub'))
        if not pubkey_paths:
            print('No SSH public keys found')
            return
        pubkeys_list = [self._parse_key(k) for k in pubkey_paths]
        print('Found the following SSH public keys:')
        for i, key_ in enumerate(pubkeys_list):
            print("{}) {} {}".format(i + 1, key_.name, key_.comment))
        print("0) Enter path to pubfile (or use keys:add <key_path>) ")
        inp = raw_input('Which would you like to use with Deis? ')
        try:
            if int(inp) != 0:
                selected_key = pubkeys_list[int(inp) - 1]
            else:
                selected_key_path = raw_input('Enter the path to the pubkey file: ')
                selected_key = self._parse_key(os.path.expanduser(selected_key_path))
        except:
            print('Aborting')
            return
        return selected_key

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

    def perms(self, args):
        """
        Valid commands for perms:

        perms:list            list permissions granted on an app or cluster
        perms:create          create a new permission for a user
        perms:delete          delete a permission for a user

        Use `deis help perms:[command]` to learn more
        """
        # perms:transfer        transfer ownership of an app or cluster
        sys.argv[1] = 'perms:list'
        args = docopt(self.perms_list.__doc__)
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

    def releases(self, args):
        """
        Valid commands for releases:

        releases:list        list an application's release history
        releases:info        print information about a specific release
        releases:rollback    return to a previous release

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

    def shortcuts(self, args):
        """
        Show valid shortcuts for client commands.

        Usage: deis shortcuts
        """
        print('Valid shortcuts are:\n')
        for shortcut, command in SHORTCUTS.items():
            if ':' not in shortcut:
                print("{:<10} -> {}".format(shortcut, command))
        print('\nUse `deis help [command]` to learn more')

SHORTCUTS = OrderedDict([
    ('create', 'apps:create'),
    ('destroy', 'apps:destroy'),
    ('init', 'clusters:create'),
    ('info', 'apps:info'),
    ('run', 'apps:run'),
    ('open', 'apps:open'),
    ('logs', 'apps:logs'),
    ('register', 'auth:register'),
    ('login', 'auth:login'),
    ('logout', 'auth:logout'),
    ('scale', 'ps:scale'),
    ('rollback', 'releases:rollback'),
    ('sharing', 'perms:list'),
    ('sharing:list', 'perms:list'),
    ('sharing:add', 'perms:create'),
    ('sharing:remove', 'perms:delete'),
])


def parse_args(cmd):
    """
    Parse command-line args applying shortcuts and looking for help flags
    """
    if cmd == 'help':
        cmd = sys.argv[-1]
        help_flag = True
    else:
        cmd = sys.argv[1]
        help_flag = False
    # swap cmd with shortcut
    if cmd in SHORTCUTS:
        cmd = SHORTCUTS[cmd]
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
    args = docopt(__doc__, version=__version__,
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
