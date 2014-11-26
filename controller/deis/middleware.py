from deis import __version__


class PlatformVersionMiddleware:

    def process_response(self, request, response):
        response['X_DEIS_PLATFORM_VERSION'] = __version__
        return response
