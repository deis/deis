"""
HTTP middleware for the Deis REST API.

See https://docs.djangoproject.com/en/1.6/topics/http/middleware/
"""

import json

from django.http import HttpResponse
from rest_framework import status

from api import __version__


class APIVersionMiddleware(object):
    """
    Return an error if a client request is incompatible with this REST API
    version, and include that REST API version with each response.
    """

    def process_request(self, request):
        """
        Return a 405 "Not Allowed" if the request's client major version
        doesn't match this controller's REST API major version (currently "1").
        """
        try:
            client_version = request.META['HTTP_DEIS_VERSION']
            server_version = __version__.rsplit('.', 2)[0]
            if client_version != server_version:
                message = {
                    'error': 'Client and server versions do not match. ' +
                             'Client version: {} '.format(client_version) +
                             'Server version: {}'.format(server_version)
                }
                return HttpResponse(
                    json.dumps(message),
                    content_type='application/json',
                    status=status.HTTP_405_METHOD_NOT_ALLOWED
                )
        except KeyError:
            pass

    def process_response(self, request, response):
        """
        Include the controller's REST API major and minor version in
        a response header.
        """
        # clients shouldn't care about the patch release
        response['DEIS_API_VERSION'] = __version__.rsplit('.', 1)[0]
        response['X_DEIS_API_VERSION'] = response['DEIS_API_VERSION']  # DEPRECATED
        return response
