"""
RESTful URL patterns and routing for the Deis API app.

Keys
====

.. http:get:: /api/keys/(string:id)/

  Retrieve a :class:`~api.models.Key` by its id.

.. http:delete:: /api/keys/(string:id)/

  Delete a :class:`~api.models.Key` by its id.

.. http:get:: /api/keys/

  List all :class:`~api.models.Key`\s.

.. http:post:: /api/keys/

  Create a new :class:`~api.models.Key`.


Providers
=========

.. http:get:: /api/providers/(string:id)/

  Retrieve a :class:`~api.models.Provider` by its id.

.. http:patch:: /api/providers/(string:id)/

  Update parts of a :class:`~api.models.Provider`.

.. http:delete:: /api/providers/(string:id)/

  Delete a :class:`~api.models.Provider` by its id.

.. http:get:: /api/providers/

  List all :class:`~api.models.Provider`\s.

.. http:post:: /api/providers/

  Create a new :class:`~api.models.Provider`.


Flavors
=======

.. http:get:: /api/flavors/(string:id)/

  Retrieve a :class:`~api.models.Flavor` by its id.

.. http:patch:: /api/flavors/(string:id)/

  Update parts of a :class:`~api.models.Flavor`.

.. http:delete:: /api/flavors/(string:id)/

  Delete a :class:`~api.models.Flavor` by its id.

.. http:get:: /api/flavors/

  List all :class:`~api.models.Flavor`\s.

.. http:post:: /api/flavors/

  Create a new :class:`~api.models.Flavor`.


Formations
==========

.. http:get:: /api/formations/(string:id)/

  Retrieve a :class:`~api.models.Formation` by its id.

.. http:delete:: /api/formations/(string:id)/

  Delete a :class:`~api.models.Formation` by its id.

.. http:get:: /api/formations/

  List all :class:`~api.models.Formation`\s.

.. http:post:: /api/formations/

  Create a new :class:`~api.models.Formation`.

Infrastructure
--------------

.. http:get:: /api/formations/(string:id)/layers/(string:id)/

  Retrieve a :class:`~api.models.Layer` by its id.

.. http:patch:: /api/formations/(string:id)/layers/(string:id)/

  Update parts of a :class:`~api.models.Layer`.

.. http:delete:: /api/formations/(string:id)/layers/(string:id)/

  Delete a :class:`~api.models.Layer` by its id.

.. http:get:: /api/formations/(string:id)/layers/

  List all :class:`~api.models.Layer`\s.

.. http:post:: /api/formations/(string:id)/layers/

  Create a new :class:`~api.models.Layer`.

.. http:get:: /api/formations/(string:id)/nodes/(string:id)/

  Retrieve a :class:`~api.models.Node` by its id.

.. http:delete:: /api/formations/(string:id)/nodes/(string:id)/

  Delete a :class:`~api.models.Node` by its id.

.. http:get:: /api/formations/(string:id)/nodes/

  List all :class:`~api.models.Node`\s.

.. http:get:: /api/formations/(string:id)/containers/

  List all :class:`~api.models.Container`\s.


Release Components
------------------

.. http:get:: /api/formations/(string:id)/config/

  List all :class:`~api.models.Config`\s.

.. http:post:: /api/formations/(string:id)/config/

  Create a new :class:`~api.models.Config`.

.. http:get:: /api/formations/(string:id)/builds/(string:uuid)/

  Retrieve a :class:`~api.models.Build` by its uuid.

.. http:post:: /api/formations/(string:id)/builds/

  Create a new :class:`~api.models.Build`.

.. http:get:: /api/formations/(string:id)/releases/(int:version)/

  Retrieve a :class:`~api.models.Release` by its version.

.. http:get:: /api/formations/(string:id)/releases/

  List all :class:`~api.models.Release`\s.

Actions
-------

.. http:post:: /api/formations/(string:id)/scale/layers/

    scale_layers

.. http:post:: /api/formations/(string:id)/scale/containers/

    scale_containers

.. http:post:: /api/formations/(string:id)/balance/

    balance

.. http:post:: /api/formations/(string:id)/calculate/

    calculate

.. http:post:: /api/formations/(string:id)/converge/

    converge


Auth
====

.. http:post:: /api/auth/register/

  Create a new :class:`~api.models.UserRegistration`.

.. http:post:: /api/auth/???

  TODO: document the important rest_framework login URLs

.. http:get:: /api/generate-api-key/

  Generate an API key.

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
    url(r'^formations/(?P<id>[a-z0-9-]+)/layers/(?P<layer>[a-z0-9-]+)/?',
        views.FormationLayerViewSet.as_view({
            'get': 'retrieve', 'patch': 'partial_update', 'delete': 'destroy'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/layers/?',
        views.FormationLayerViewSet.as_view({'post': 'create', 'get': 'list'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/nodes/(?P<node>[a-z0-9-]+)/?',
        views.FormationNodeViewSet.as_view({
            'get': 'retrieve', 'delete': 'destroy'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/nodes/?',
        views.FormationNodeViewSet.as_view({'get': 'list'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/containers/?',
        views.FormationContainerViewSet.as_view({'get': 'list'})),
    # formation release components
    url(r'^formations/(?P<id>[a-z0-9-]+)/config/?',
        views.FormationConfigViewSet.as_view({'post': 'create', 'get': 'retrieve'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/builds/(?P<uuid>[a-z0-9-]+)/?',
        views.FormationBuildViewSet.as_view({'get': 'retrieve'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/builds/?',
        views.FormationBuildViewSet.as_view({'post': 'create', 'get': 'list'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/releases/(?P<version>[0-9]+)/?',
        views.FormationReleaseViewSet.as_view({'get': 'retrieve'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/releases/?',
        views.FormationReleaseViewSet.as_view({'get': 'list'})),
    # formation actions
    url(r'^formations/(?P<id>[a-z0-9-]+)/scale/layers/?',
        views.FormationViewSet.as_view({'post': 'scale_layers'})),
    url(r'^formations/(?P<id>[a-z0-9-]+)/scale/containers/?',
        views.FormationViewSet.as_view({'post': 'scale_containers'})),
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
    url(r'^auth/register/?',
        views.UserRegistrationView.as_view({'post': 'create'})),
    url(r'^auth/',
        include('rest_framework.urls', namespace='rest_framework')),
    url(r'^generate-api-key/',
        'rest_framework.authtoken.views.obtain_auth_token'),
)
