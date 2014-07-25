import json

from django.http import HttpResponse
from rest_framework import status

from deis import __version__


class VersionMiddleware:

    def process_request(self, request):
        try:
            # server and client version must match "x.y"
            client_version = request.META['HTTP_X_DEIS_VERSION']
            server_version = __version__.rsplit('.', 1)[0]
            if client_version != server_version:
                message = {
                    'error': 'Client and server versions do not match.\n' +
                    'Client version: {}\n'.format(client_version) +
                    'Server version: {}'.format(server_version)
                }
                return HttpResponse(
                    json.dumps(message),
                    content_type='application/json',
                    status=status.HTTP_405_METHOD_NOT_ALLOWED
                )
        except KeyError:
            pass
