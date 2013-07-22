
from __future__ import unicode_literals

# add patch support to built-in django test client

from django.test.client import RequestFactory, Client


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

from .auth import *
from .backend import *
from .build import *
from .config import *
from .container import *
from .flavor import *
from .formation import *
from .node import *
from .provider import *
from .proxy import *
from .release import *
