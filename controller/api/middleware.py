import json

from django.http import HttpResponse

from deis import __version__


class VersionMiddleware:

    def process_request(self, request):
        # server and client version must match "x.y"
        server_version = __version__.rsplit('.', 1)[0]
        try:
            if request.META['HTTP_X_DEIS_VERSION'] != server_version:
                message = {
                    'error': 'Client and server versions do not match.\n' +
                    'Client version: {}\n'.format(server_version) +
                    'Server version: {}'.format(request.META['HTTP_X_DEIS_VERSION'])
                }
                return HttpResponse(
                    json.dumps(message),
                    content_type='application/json'
                )
        except KeyError:
            pass
