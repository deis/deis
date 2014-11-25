import json

from django.http import HttpResponse
from rest_framework import status

from api import __version__


class APIVersionMiddleware:

    def process_request(self, request):
        try:
            # server and client version must match the major release point
            client_version = request.META['HTTP_X_DEIS_VERSION']
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
        # clients shouldn't care about the patch release
        response['X_DEIS_API_VERSION'] = __version__.rsplit('.', 1)[0]
        return response
