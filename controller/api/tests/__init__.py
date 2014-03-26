
from __future__ import unicode_literals
import logging

from django.test.client import RequestFactory, Client
from django.test.simple import DjangoTestSuiteRunner


# add patch support to built-in django test client

def construct_patch(self, path, data='',
                    content_type='application/octet-stream', **extra):
    """Construct a PATCH request."""
    return self.generic('PATCH', path, data, content_type, **extra)


def send_patch(self, path, data='', content_type='application/octet-stream',
               follow=False, **extra):
    """Send a resource to the server using PATCH."""
    # FIXME: figure out why we need to reimport Client (otherwise NoneType)
    from django.test.client import Client  # @Reimport
    response = super(Client, self).patch(
        path, data=data, content_type=content_type, **extra)
    if follow:
        response = self._handle_redirects(response, **extra)
    return response


RequestFactory.patch = construct_patch
Client.patch = send_patch


class SilentDjangoTestSuiteRunner(DjangoTestSuiteRunner):
    """Prevents api log messages from cluttering the console during tests."""

    def run_tests(self, test_labels, extra_tests=None, **kwargs):
        """Run tests with all but critical log messages disabled."""
        # hide any log messages less than critical
        logging.disable(logging.CRITICAL)
        return super(SilentDjangoTestSuiteRunner, self).run_tests(
            test_labels, extra_tests, **kwargs)


from .test_app import *  # noqa
from .test_auth import *  # noqa
from .test_build import *  # noqa
from .test_cluster import *  # noqa
from .test_config import *  # noqa
from .test_container import *  # noqa
from .test_hooks import *  # noqa
from .test_key import *  # noqa
from .test_perm import *  # noqa
from .test_release import *  # noqa
