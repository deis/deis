"""
REST framework URL routing classes.
"""

from __future__ import unicode_literals

from rest_framework.routers import DefaultRouter
from rest_framework.routers import Route


class ApiRouter(DefaultRouter):
    """Generate URL patterns for list, detail, and viewset-specific
    HTTP routes.
    """

    routes = [
        # List route.
        Route(
            url=r"^{prefix}/?$",
            mapping={
                'get': 'list',
                'post': 'create'
            },
            name="{basename}-list",
            initkwargs={'suffix': 'List'}
        ),
        # Detail route.
        Route(
            url=r"^{prefix}/{lookup}/?$",
            mapping={
                'get': 'retrieve',
                'put': 'update',
                'patch': 'partial_update',
                'delete': 'destroy'
            },
            name="{basename}-detail",
            initkwargs={'suffix': 'Instance'}
        ),
        # Dynamically generated routes, from @action or @link decorators
        # on methods of the viewset.
        Route(
            url=r"^{prefix}/{lookup}/{methodname}/?$",
            mapping={
                "{httpmethod}": "{methodname}",
            },
            name="{basename}-{methodnamehyphen}",
            initkwargs={}
        ),
    ]
