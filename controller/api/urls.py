"""
URL routing patterns for the Deis API app.
"""
# pylint: disable=C0103

from __future__ import unicode_literals

from django.conf.urls import include
from django.conf.urls import patterns
from django.conf.urls import url

from api import routers
from api import views


router = routers.ApiRouter()

# Add the generated REST URLs and login/logout endpoint
urlpatterns = patterns(
    '',
    url(r'^', include(router.urls)),

    # key
    url(r'^keys/(?P<id>.+)/?',
        views.KeyViewSet.as_view({
            'get': 'retrieve', 'delete': 'destroy'})),
    url(r'^keys/?',
        views.KeyViewSet.as_view({'post': 'create', 'get': 'list'})),
    # provider
    url(r'^providers/(?P<id>[a-z0-9-]+)/?',
        views.ProviderViewSet.as_view({
            'get': 'retrieve', 'patch': 'partial_update', 'delete': 'destroy'})),
    url(r'^providers/?',
        views.ProviderViewSet.as_view({'post': 'create', 'get': 'list'})),
    # flavor
    url(r'^flavors/(?P<id>[a-z0-9-]+)/?',
        views.FlavorViewSet.as_view({
            'get': 'retrieve', 'patch': 'partial_update', 'delete': 'destroy'})),
    url(r'^flavors/?',
        views.FlavorViewSet.as_view({'post': 'create', 'get': 'list'})),

    # formation infrastructure
    url(r'^formations/(?P<id>[a-z0-9-]+)/nodes/?',
        views.FormationNodeViewSet.as_view({'get': 'list'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/backends/?',
        views.FormationBackendViewSet.as_view({'get': 'list'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/proxies/?',
        views.FormationProxyViewSet.as_view({'get': 'list'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/containers/?',
        views.FormationContainerViewSet.as_view({'get': 'list'})),
    # formation release components
    url(r'^formations/(?P<id>[a-z0-9-]+)/config/?',
        views.FormationConfigViewSet.as_view({'post': 'create', 'get': 'retrieve'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/image/?',
        views.FormationImageViewSet.as_view({'post': 'reset_image'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/build/?',
        views.FormationBuildViewSet.as_view({'post': 'create', 'get': 'retrieve'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/release/?',
        views.FormationReleaseViewSet.as_view({'get': 'retrieve'})),
    # formation actions
    url(r'^formations/(?P<id>[a-z0-9-]+)/scale/?',
        views.FormationViewSet.as_view({'post': 'scale'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/balance/?',
        views.FormationViewSet.as_view({'post': 'balance'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/calculate/?',
        views.FormationViewSet.as_view({'post': 'calculate'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/converge/?',
        views.FormationViewSet.as_view({'post': 'converge'})),
    # formation base endpoint
    url(r'^formations/(?P<id>[a-z0-9-]+)/?',
        views.FormationViewSet.as_view({'get': 'retrieve', 'delete': 'destroy'})),
    url(r'^formations/?',
        views.FormationViewSet.as_view({'post': 'create', 'get': 'list'})),
    
    # authn / authz
    url(r'^register/?',
        views.UserRegistrationView.as_view({'post': 'create'})),
    url(r'^auth/',
        include('rest_framework.urls', namespace='rest_framework')),
    url(r'^generate-api-key/',
        'rest_framework.authtoken.views.obtain_auth_token'),
)
