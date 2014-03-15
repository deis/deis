"""
Tests for the command-line client for the Deis system.

These tests can only be run against an already provisioned Deis installation. It's recommended not
to run them on the controller, but from where you would normally run the command line client.

If you're running the tests for an installation that doesn't have any registered users yet, you can
run this command from the root of the project:

`DEIS_SERVER=<controller domain> python -m unittest client.tests`

Where <controller domain> is the FQDN for the Deis Controller, eg. 'deis-controller.local' for a
vagrant installation.

If you're running the tests for an existing installation with users already registered then you
will need to use the following:

`DEIS_SERVER=<controller domain> DEIS_SUPER_USER=<user> DEIS_SUPER_PASS=<password> \
python -m unittest client.tests`

Where <user> and <password> are replaced with the first registered user or a known superuser.
"""

try:
    import pexpect  # noqa
except ImportError:
    print('Please install the python pexpect library.')
    raise

from .test_apps import *  # noqa
from .test_auth import *  # noqa
from .test_builds import *  # noqa
from .test_config import *  # noqa
from .test_containers import *  # noqa
from .test_examples import *  # noqa
from .test_flavors import *  # noqa
from .test_formations import *  # noqa
from .test_keys import *  # noqa
from .test_layers import *  # noqa
from .test_misc import *  # noqa
from .test_nodes import *  # noqa
from .test_providers import *  # noqa
from .test_releases import *  # noqa
from .test_sharing import *  # noqa
