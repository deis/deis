"""
Tests for the command-line client for the Deis system.
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
