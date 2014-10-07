"""
WSGI config for deis project.

This module contains the WSGI application used by Django's development server
and any production WSGI deployments. It should expose a module-level variable
named ``application``. Django's ``runserver`` and ``runfcgi`` commands discover
this application via the ``WSGI_APPLICATION`` setting.

"""

from __future__ import unicode_literals
import os

from django.core.wsgi import get_wsgi_application
import static

from registry.importer import forward_to_registry

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "deis.settings")


class Dispatcher(object):
    """
    Dispatches requests between two WSGI apps, a static file server and a
    Django server.
    """

    def __init__(self):
        self.django_handler = get_wsgi_application()
        self.static_handler = static.Cling(os.path.dirname(os.path.dirname(__file__)))

    def __call__(self, environ, start_response):
        if environ['PATH_INFO'].startswith('/static'):
            return self.static_handler(environ, start_response)
        elif environ['PATH_INFO'].startswith('/push'):
            environ['PATH_INFO'] = environ['PATH_INFO'][5:]
            return forward_to_registry(environ, start_response)
        else:
            return self.django_handler(environ, start_response)


application = Dispatcher()
