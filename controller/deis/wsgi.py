"""
WSGI config for deis project.

It exposes the WSGI callable as a module-level variable named ``application``.

For more information on this file, see
https://docs.djangoproject.com/en/1.7/howto/deployment/wsgi/
"""

import os

from django.core.wsgi import get_wsgi_application
import static


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
        else:
            return self.django_handler(environ, start_response)


application = Dispatcher()
