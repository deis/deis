#!/usr/bin/env python

"""Usage: deis <command> [--formation <formation>] [<args>...]

Options:
  -h --help       Show this help screen
  -v --version    Show the CLI version

Auth commands:

  register      register a new user with a controller
  login         login to a controller
  logout        logout from the current controller

Common commands:

  create        create a new container formation
  info          print a represenation of the formation
  scale         scale container types (web=2, worker=1)
  balance       rebalance the container formation
  converge      force-converge all nodes in the formation
  calculate     recalculate and update the formation databag
  destroy       destroy a container formation

Use `deis help [subcommand]` to learn about these subcommands:

  formations    manage container formations
  layers        manage layers of nodes
  nodes         manage nodes of all types
  containers    manage the containers running on backends

  providers     manage cloud provider credentials
  flavors       manage node flavors on a provider
  keys          manage ssh keys

  config        manage environment variables for a formation
  builds        manage git-push builds for a formation
  releases      manage a formation's release history

Use `git push deis master` to deploy to the container formation.

"""

from cookielib import MozillaCookieJar
from getpass import getpass
import glob
import json
import os.path
import re
import subprocess
import sys
import urlparse

from docopt import docopt
from docopt import DocoptExit
import requests
import yaml


__version__ = '0.0.4'


class Session(requests.Session):

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

    def git_root(self):
        try:
            git_root = subprocess.check_output(
                ['git', 'rev-parse', '--show-toplevel'],
                stderr=subprocess.PIPE).strip('\n')
        except subprocess.CalledProcessError:
            raise EnvironmentError('Current directory is not a git repository')
        return git_root

    def get_formation(self):
        git_root = self.git_root()
        # try to match a deis remote
        remotes = subprocess.check_output(['git', 'remote', '-v'],
                                          cwd=git_root)
        m = re.match(r'^deis\W+(?P<url>\S+)\W+\(', remotes, re.MULTILINE)
        if not m:
            raise EnvironmentError(
                'Could not find deis remote in `git remote -v`')
        url = m.groupdict()['url']
        m = re.match('\S+:(?P<formation>[a-z0-9-]+)(.git)?', url)
        if not m:
            raise EnvironmentError('Could not parse: {url}'.format(**locals()))
        return m.groupdict()['formation']

    formation = property(get_formation)

    def request(self, *args, **kwargs):
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
        with open(self._path) as f:
            data = f.read()
        settings = yaml.safe_load(data)
        self.update(settings)
        return settings

    def save(self):
        data = yaml.safe_dump(dict(self))
        with open(self._path, 'w') as f:
            f.write(data)
        return data


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


def trim(docstring):
    """
    Function to trim whitespace from docstring

    c/o PEP 257 Docstring Conventions
    http://www.python.org/dev/peps/pep-0257/
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


class DeisClient(object):

    """
    A client which interacts with a Deis server.
    """

    def __init__(self):
        self._session = Session()
        self._settings = Settings()

    def _dispatch(self, method, path, body=None,
                  headers={'content-type': 'application/json'}, **kwargs):
        func = getattr(self._session, method.lower())
        url = urlparse.urljoin(self._settings['controller'], path, **kwargs)
        response = func(url, data=body, headers=headers)
        return response

    def auth_register(self, args):
        """
        Register a new user with a Deis controller

        Usage: deis auth:register <controller> [--username=<username> --password=<password> --email=<email>]
        """
        controller = args['<controller>']
        username = args.get('--username')
        if not username:
            username = raw_input('username: ')
        password = args.get('--password')
        if not password:
            password = getpass('password: ')
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
            # login after registering
            if self.auth_login(login_args) is False:
                print('Login failed')
                return
            # add ssh keys before formations are created
            print
            self.keys_add({})
            print
            self.providers_discover({})
            print
            print 'Use `deis create --flavor=ec2-us-east-1` to create a new formation'
        else:
            print('Registration failed', response.content)
            return False

    def auth_login(self, args):
        """
        Login by authenticating against a controller

        Usage: deis auth:login <controller> [--username=<username> --password=<password>]
        """
        controller = args['<controller>']
        username = args.get('--username')
        headers = {}
        if not username:
            username = raw_input('username: ')
        password = args.get('--password')
        if not password:
            password = getpass('password: ')
        url = urlparse.urljoin(controller, '/api/auth/login/')
        payload = {'username': username, 'password': password}
        # prime cookies for login
        self._session.get(url, headers=headers)
        # post credentials to the login URL
        response = self._session.post(url, data=payload, allow_redirects=False)
        if response.status_code == requests.codes.found:  # @UndefinedVariable
            self._settings['controller'] = controller
            self._settings.save()
            print("Logged in as {}".format(username))
            return True
        else:
            print('Login failed')
            self._session.cookies.clear()
            self._session.cookies.save()
            return False

    def auth_logout(self, args):
        """
        Logout from a controller, clearing the user session

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
        Builds help would be nice
        """
        return self.builds_list(args)

    def builds_create(self, args):
        """
        Usage: deis builds:create - [--formation=<formation>]
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        data = sys.stdin.read()
        # url / sha / slug_size / procfile / checksum
        j = json.loads(data)
        response = self._dispatch('post',
                                  '/api/formations/{}/builds'.format(formation),
                                  body=json.dumps(j))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('Build created.')
        else:
            print('Error!', response.text)

    def builds_list(self, args):
        """
        Usage: deis builds:list
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('get', '/api/formations/{}/builds'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('=== {0}'.format(formation))
            data = response.json()
            for item in data['results']:
                print('{0[uuid]:<23} {0[created]}'.format(item))
        else:
            print('Error!', response.text)

    def config(self, args):
        """
        Config help would be nice
        """
        return self.config_list(args)

    def config_list(self, args):
        """
        Usage: deis config:list
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('get', '/api/formations/{}/config'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            config = response.json()
            values = json.loads(config['values'])
            print('=== {0}'.format(formation))
            items = values.items()
            if len(items) == 0:
                print('No configuration')
                return
            for k, v in values.items():
                print('{k}: {v}'.format(**locals()))
        else:
            print('Error!', response.text)

    def config_set(self, args):
        """
        Usage: deis config:set <var>=<value>...
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        body = {'values': json.dumps(dictify(args['<var>=<value>']))}
        response = self._dispatch('post',
                                  '/api/formations/{}/config'.format(formation),
                                  json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            config = response.json()
            values = json.loads(config['values'])
            print('=== {0}'.format(formation))
            items = values.items()
            if len(items) == 0:
                print('No configuration')
                return
            for k, v in values.items():
                print('{k}: {v}'.format(**locals()))
        else:
            print('Error!', response.text)

    def config_unset(self, args):
        """
        Usage: deis config:unset <key>...
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        values = {}
        for k in args.get('<key>'):
            values[k] = None
        body = {'values': json.dumps(values)}
        response = self._dispatch('post',
                                  '/api/formations/{}/config'.format(formation),
                                  data=json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            config = response.json()
            values = json.loads(config['values'])
            print('=== {0}'.format(formation))
            items = values.items()
            if len(items) == 0:
                print('No configuration')
                return
            for k, v in values.items():
                print('{k}: {v}'.format(**locals()))
        else:
            print('Error!', response.text)

    def containers(self, args):
        """
        Containers help would be nice
        """
        return self.containers_list(args)

    def containers_list(self, args):
        """
        Usage: deis containers:list
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('get',
                                  '/api/formations/{}/containers'.format(formation))
        databag = self.formations_calculate({}, quiet=True)
        procfile = databag['release']['build'].get('procfile', {})
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            c_map = {}
            for item in data['results']:
                c_map.setdefault(item['type'], []).append(item)
            for c_type in c_map.keys():
                command = procfile.get(c_type, '<none>')
                print('=== {c_type}: `{command}`'.format(**locals()))
                for c in c_map[c_type]:
                    print('{type}.{num} up {created}'.format(**c))
                print
        else:
            print('Error!', response.text)

    def containers_scale(self, args):
        """
        Usage: deis containers:scale <type=num>...
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        body = {}
        for type_num in args.get('<type=num>'):
            typ, count = type_num.split('=')
            body.update({typ: int(count)})
        response = self._dispatch('post',
                                  '/api/formations/{}/scale/containers'.format(formation),
                                  json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            databag = json.loads(response.content)
            print(json.dumps(databag, indent=2))
        else:
            print('Error!', response.text)

    def flavors(self, args):
        """
        Flavors help would be nice
        """
        return self.flavors_list(args)

    def flavors_create(self, args):
        """
        Usage: deis flavors:create --id=<id> --provider=<provider> --params=<params> [options]

        Options:

        --params=PARAMS    provider-specific parameters (size, region, zone, etc.)
        --init=INIT        override Ubuntu cloud-init with custom YAML
        """
        body = {'id': args.get('--id'), 'provider': args.get('--provider')}
        fields = ('params', 'init', 'ssh_username', 'ssh_private_key',
                  'ssh_public_key')
        for fld in fields:
            opt = args.get('--' + fld)
            if opt:
                body.update({fld: opt})
        response = self._dispatch('post', '/api/flavors', json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('{0[id]}'.format(response.json()))
        else:
            print('Error!', response.text)

    def flavors_delete(self, args):
        """
        Usage: deis flavors:delete <id>
        """
        flavor = args.get('<id>')
        response = self._dispatch('delete', '/api/flavors/{}'.format(flavor))
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            pass
        else:
            print('Error!', response.status_code, response.text)

    def flavors_info(self, args):
        """
        Usage: deis flavors:info <flavor>
        """
        flavor = args.get('<flavor>')
        response = self._dispatch('get', '/api/flavors/{}'.format(flavor))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            print('Error!', response.text)

    def flavors_list(self, args):
        """
        Usage: deis flavors:list
        """
        response = self._dispatch('get', '/api/flavors')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            for item in data['results']:
                print('{0[id]:<23}'.format(item))
        else:
            print('Error!', response.text)

    def formations(self, args):
        """
        Formations help would be nice
        """
        return self.formations_list(args)

    def formations_create(self, args):
        """
        Usage: deis formations:create [--id=<id> --flavor=<flavor>]
        """
        body = {}
        try:
            self._session.git_root()  # check for a git repository
        except EnvironmentError:
            print 'No git repository found, use `git init` to create one'
            return
        for opt in ('--id',):
            o = args.get(opt)
            if o:
                body.update({opt.strip('-'): o})
        sys.stdout.write('Creating formation... ')
        sys.stdout.flush()
        response = self._dispatch('post', '/api/formations',
                                  json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            data = response.json()
            formation = data['id']
            print('done, created {}'.format(formation))
            # add a git remote
            hostname = urlparse.urlparse(self._settings['controller']).netloc
            git_remote = 'git@{hostname}:{formation}.git'.format(**locals())
            try:
                subprocess.check_call(
                    ['git', 'remote', 'add', '-f', 'deis', git_remote],
                    stdout=subprocess.PIPE)
            except subprocess.CalledProcessError:
                sys.exit(1)
            print('Git remote deis added')
            # create default layers if a flavor was provided
            flavor = args.get('--flavor')
            if flavor:
                print
                self.layers_create({'<id>': 'runtime', '<flavor>': flavor})
                self.layers_create({'<id>': 'proxy', '<flavor>': flavor})
                print('\nUse `deis layers:scale runtime=1 proxy=1` to scale a basic formation')
        else:
            print('Error!', response.text)

    def formations_info(self, args):
        """
        Usage: deis formations:info
        """
        formation = args.get('<formation>')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('get', '/api/formations/{}'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            print('Error!', response.text)

    def formations_list(self, args):
        """
        Usage: deis formations:list
        """
        response = self._dispatch('get', '/api/formations')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            for item in data['results']:
                print('{0[id]:<23}'.format(item))
        else:
            print('Error!', response.text)

    def formations_destroy(self, args):
        """
        Usage: deis formations:destroy [<formation>] [--confirm=<confirm>]
        """
        formation = args.get('<formation>')
        if not formation:
            formation = self._session.formation
        confirm = args.get('--confirm')
        if confirm == formation:
            pass
        else:
            print """
 !    WARNING: Potentially Destructive Action
 !    This command will destroy: {formation}
 !    To proceed, type "{formation}" or re-run this command with --confirm={formation}
""".format(**locals())
            confirm = raw_input('> ').strip('\n')
            if confirm != formation:
                print('Destroy aborted')
                return
        sys.stdout.write('Destroying {}... '.format(formation))
        sys.stdout.flush()
        response = self._dispatch('delete', '/api/formations/{}'.format(formation))
        if response.status_code in (requests.codes.no_content,  # @UndefinedVariable
                                    requests.codes.not_found):  # @UndefinedVariable
            print('done')
            try:
                subprocess.check_call(
                    ['git', 'remote', 'rm', 'deis'],
                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                print('Git remote deis removed')
            except subprocess.CalledProcessError:
                pass  # ignore error
        else:
            print('Error!', response.text)

    def formations_calculate(self, args, quiet=False):
        """
        Usage: deis formations:calculate
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('post',
                                  '/api/formations/{}/calculate'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            databag = json.loads(response.content)
            if quiet is False:
                print(json.dumps(databag, indent=2))
            return databag
        else:
            print('Error!', response.text)

    def formations_balance(self, args):
        """
        Usage: deis formations:balance
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('post',
                                  '/api/formations/{}/balance'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            databag = json.loads(response.content)
            print(json.dumps(databag, indent=2))
        else:
            print('Error!', response.text)

    def formations_converge(self, args):
        """
        Usage: deis formations:converge
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('post',
                                  '/api/formations/{}/converge'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            databag = json.loads(response.content)
            print(json.dumps(databag, indent=2))
        else:
            print('Error!', response.text)

    def keys(self, args):
        """
        Keys help would be nice
        """
        return self.keys_list(args)

    def keys_add(self, args):
        """
        Usage: deis keys:add [<key>]
        """
        path = args.get('<key>')
        if not path:
            ssh_dir = os.path.expanduser('~/.ssh')
            pubkeys = glob.glob(os.path.join(ssh_dir, '*.pub'))
            print('Found the following SSH public keys:')
            for i, k in enumerate(pubkeys):
                key = k.split(os.path.sep)[-1]
                print('{0}) {1}'.format(i + 1, key))
            inp = raw_input('Which would you like to use with Deis? ')
            try:
                path = pubkeys[int(inp) - 1]
                key_id = path.split(os.path.sep)[-1].replace('.pub', '')
            except:
                print 'Aborting'
                return
        with open(path) as f:
            data = f.read()
        match = re.match(r'^(ssh-...) ([^ ]+) (.+)', data)
        if not match:
            print 'Could not parse public key material'
            return
        key_type, key_str, _key_comment = match.groups()
        body = {'id': key_id, 'public': '{0} {1}'.format(key_type, key_str)}
        sys.stdout.write('Uploading {0} to Deis... '.format(path))
        sys.stdout.flush()
        response = self._dispatch('post', '/api/keys', json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('done')
        else:
            print('Error!', response.text)

    def keys_list(self, args):
        """
        List SSH keys for the logged in user.

        Usage: deis keys:list
        """
        response = self._dispatch('get', '/api/keys')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            if data['count'] == 0:
                print 'No keys found'
                return
            print('=== {0} Keys'.format(data['results'][0]['owner']))
            for key in data['results']:
                public = key['public']
                print('{0} {1}...{2}'.format(
                    key['id'], public[0:16], public[-10:]))
        else:
            print('Error!', response.text)

    def keys_remove(self, args):
        """
        Remove a specific SSH key for the logged in user.

        Usage: deis keys:remove <key>
        """
        key = args.get('<key>')
        sys.stdout.write('Removing {0} SSH Key... '.format(key))
        sys.stdout.flush()
        response = self._dispatch('delete', '/keys/{}'.format(key))
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print('done')
        else:
            print('Error!', response.text)

    def layers(self, args):
        """
        Layers help would be nice
        """
        return self.layers_list(args)

    def layers_create(self, args):
        """
        Create a layer of nodes.

        Usage: deis layers:create <id> <flavor> [options]

        Chef Options:

        --run_list=RUN_LIST         run-list to use when bootstrapping nodes
        --environment=ENVIRONMENT   chef environment to place nodes [default: _default]
        --attributes=INITIAL_ATTRS  initial attributes for nodes

        SSH Options:

        --ssh_username=USERNAME         username for ssh connections [default: ubuntu]
        --ssh_private_key=PRIVATE_KEY   private key for ssh comm (default: auto-gen)
        --ssh_public_key=PUBLIC_KEY     public key for ssh comm (default: auto-gen)

        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        body = {'id': args['<id>'], 'flavor': args['<flavor>']}
        for opt in ('--environment', '--initial_attributes', '--run_list',
                    '--ssh_username', '--ssh_private_key', '--ssh_public_key'):
            o = args.get(opt)
            if o:
                body.update({opt.strip('-'): o})
        # provide default run_list for runtime and proxy
        if not 'run_list' in body:
            if body['id'] == 'runtime':
                body['run_list'] = 'recipe[deis],recipe[deis::runtime]'
            elif body['id'] == 'proxy':
                body['run_list'] = 'recipe[deis],recipe[deis::proxy]'
        sys.stdout.write('Creating {} layer... '.format(args['<id>']))
        sys.stdout.flush()
        response = self._dispatch('post', '/api/formations/{}/layers'.format(formation),
                                  json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('done')
        else:
            print('Error!', response.text)

    def layers_destroy(self, args):
        """
        Destroy a layer of nodes.

        Usage: deis layers:destroy <id>
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        layer = args['<id>']  # noqa
        sys.stdout.write('Destroying {layer} layer... '.format(**locals()))
        sys.stdout.flush()
        response = self._dispatch(
            'delete', '/api/formations/{formation}/layers/{layer}'.format(**locals()))
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print('done')
        else:
            print('Error!', response.text)

    def layers_list(self, args):
        """
        List layers for this formation.

        Usage deis layers:list
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('get',
                                  '/api/formations/{}/layers'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('=== {0}'.format(formation))
            data = response.json()
            format_str = '{0[id]} => {0[run_list]}'
            for item in data['results']:
                print(format_str.format(item))
        else:
            print('Error!', response.text)

    def layers_scale(self, args):
        """
        Usage: deis layers:scale <type=num>...
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        body = {}
        for type_num in args.get('<type=num>'):
            typ, count = type_num.split('=')
            body.update({typ: int(count)})
        print('Scaling layers... but first, coffee!')
        # TODO: add threaded spinner to print dots
        response = self._dispatch('post',
                                  '/api/formations/{}/scale/layers'.format(formation),
                                  json.dumps(body))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('...done\n')
            print('Use `git push deis master` to deploy to your formation')
        else:
            print('Error!', response.text)

    def nodes(self, args):
        """
        Nodes help would be nice
        """
        return self.nodes_list(args)

    def nodes_info(self, args):
        """
        Usage: deis nodes:info <node>
        """
        node = args.get('<node>')
        response = self._dispatch('get', '/api/nodes/{}'.format(node))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            print('Error!', response.text)

    def nodes_list(self, args):
        """
        List nodes for this formation.

        Usage: deis nodes:list
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('get',
                                  '/api/formations/{}/nodes'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('=== {0}'.format(formation))
            data = response.json()
            format_str = '{0[id]:<23} {0[layer]} {0[provider_id]} {0[fqdn]}'
            for item in data['results']:
                print(format_str.format(item))
        else:
            print('Error!', response.text)

    def nodes_destroy(self, args):
        """
        Destroy a node by ID.

        Usage: deis nodes:destroy <id>
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        node = args['<id>']
        response = self._dispatch('delete',
                                  '/api/formations/{formation}/nodes/{node}'.format(**locals()))
        sys.stdout.write('Destroying {}... '.format(node))
        sys.stdout.flush()
        if response.status_code == requests.codes.no_content:  # @UndefinedVariable
            print('done')
        else:
            print('Error!', response.status_code, response.text)

    def providers(self, args):
        """
        Providers help would be nice
        """
        return self.providers_list(args)

    def providers_create(self, args):
        """
        Create a provider for use by Deis

        Usage: deis providers:create --type=<type> [--id=<id> --creds=<creds>]
        """
        type = args.get('--type')  # @ReservedAssignment
        if type == 'ec2':
            # read creds from envvars
            for k in ('AWS_ACCESS_KEY', 'AWS_SECRET_KEY'):
                if not k in os.environ:
                    msg = 'Missing environment variable: {}'.format(k)
                    raise EnvironmentError(msg)
            creds = {'access_key': os.environ['AWS_ACCESS_KEY'],
                     'secret_key': os.environ['AWS_SECRET_KEY']}
        id = args.get('--id')  # @ReservedAssignment
        if not id:
            id = type  # @ReservedAssignment
        body = {'id': id, 'type': type, 'creds': json.dumps(creds)}
        response = self._dispatch('post', '/api/providers',
                                  json.dumps(body))
        if response.status_code == requests.codes.created:  # @UndefinedVariable
            print('{0[id]}'.format(response.json()))
        else:
            print('Error!', response.text)

    def providers_discover(self, args):
        """
        Discover and update providers

        Usage: deis providers:discover
        """
        # look for ec2 credentials
        if 'AWS_ACCESS_KEY' in os.environ and 'AWS_SECRET_KEY' in os.environ:
            print('Found EC2 credentials: {0}'.format(os.environ['AWS_ACCESS_KEY']))
            inp = raw_input('Import these credentials? (y/n) : ')
            if inp.lower().strip('\n') != 'y':
                print 'Aborting.'
                return
            creds = {'access_key': os.environ['AWS_ACCESS_KEY'],
                     'secret_key': os.environ['AWS_SECRET_KEY']}
            body = {'creds': json.dumps(creds)}
            sys.stdout.write('Uploading EC2 credentials... ')
            sys.stdout.flush()
            response = self._dispatch('patch', '/api/providers/ec2',
                                      json.dumps(body))
            if response.status_code == requests.codes.ok:  # @UndefinedVariable
                print 'done'
            else:
                print('Error!', response.text)
        else:
            print 'No credentials discovered, did you install the EC2 Command Line tools?'
            return

    def providers_list(self, args):
        """
        List providers

        Usage: deis providers:list
        """
        response = self._dispatch('get', '/api/providers')
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            data = response.json()
            for item in data['results']:
                print(item['id'])
        else:
            print('Error!', response.text)

    def providers_info(self, args):
        """
        Show detail of a provider.

        Usage: deis providers:info <provider>
        """
        provider = args.get('<provider>')
        response = self._dispatch('get', '/api/providers/{}'.format(provider))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print(json.dumps(response.json(), indent=2))
        else:
            print('Error!', response.text)

    def releases_list(self, args):
        """
        Usage: deis releases:list
        """
        formation = args.get('--formation')
        if not formation:
            formation = self._session.formation
        response = self._dispatch('get', '/api/formations/{}/release'.format(formation))
        if response.status_code == requests.codes.ok:  # @UndefinedVariable
            print('=== {0}'.format(formation))
            data = response.json()
            for item in data['results']:
                print('{0[uuid]:<23} {0[created]}'.format(item))
        else:
            print('Error!', response.text)

def main():
    """
    Create a client, parse the arguments received on the command line, and
    call the appropriate method on the client.
    """
    # create a client instance
    cli = DeisClient()
    # parse base command-line arguments
    args = docopt(__doc__, version='Deis CLI {}'.format(__version__),
                  options_first=True)
    cmd = args['<command>']
    # split cmd with _ if it contains a :
    shortcuts = {
        'register': 'auth:register',
        'login': 'auth:login',
        'logout': 'auth:logout',
        'create': 'formations:create',
        'info': 'formations:info',
        'balance': 'formations:balance',
        'calculate': 'formations:calculate',
        'converge': 'formations:converge',
        'destroy': 'formations:destroy',
        'scale': 'containers:scale',
        'ps': 'containers:list',
    }
    # lookup cmd shortcut
    if cmd in shortcuts:
        cmd = shortcuts[cmd]
        sys.argv[1] = cmd  # change the cmdline arg itself
    # convert : to _ for matching method names and docstrings
    if ':' in cmd:
        cmd = '_'.join(cmd.split(':'))
    # re-parse docopt with the relevant docstring
    if cmd in dir(cli):
        docstring = trim(getattr(cli, cmd).__doc__)
        if 'Usage: ' in docstring:
            args.update(docopt(docstring))
    # find the right method for dispatching
    if cmd == 'help':
        if len(sys.argv) == 3 and sys.argv[2] in dir(cli):
            print trim(getattr(cli, sys.argv[2]).__doc__)
            return
        docopt(__doc__, argv=['--help'])
    elif hasattr(cli, cmd):
        method = getattr(cli, cmd)
    else:
        print 'Found no matching command'
        raise DocoptExit()
    # dispatch the CLI command
    try:
        method(args)
    except EnvironmentError:
        print 'Could not find git remote for deis'
        raise DocoptExit()


if __name__ == '__main__':
    main()
    sys.exit(0)
