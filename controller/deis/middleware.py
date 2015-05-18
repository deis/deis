from deis import __version__


class PlatformVersionMiddleware:

    def process_response(self, request, response):
        response['DEIS_PLATFORM_VERSION'] = __version__
        response['X_DEIS_PLATFORM_VERSION'] = __version__  # DEPRECATED
        return response
