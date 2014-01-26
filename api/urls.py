"""
RESTful URL patterns and routing for the Deis API app.

Keys
====

.. http:get:: /api/keys/(string:id)/

  Retrieve a :class:`~api.models.Key` by its `id`.

.. http:delete:: /api/keys/(string:id)/

  Destroy a :class:`~api.models.Key` by its `id`.

.. http:get:: /api/keys/

  List all :class:`~api.models.Key`\s.

.. http:post:: /api/keys/

  Create a new :class:`~api.models.Key`.


Providers
=========

.. http:get:: /api/providers/(string:id)/

  Retrieve a :class:`~api.models.Provider` by its `id`.

.. http:patch:: /api/providers/(string:id)/

  Update parts of a :class:`~api.models.Provider`.

.. http:delete:: /api/providers/(string:id)/

  Destroy a :class:`~api.models.Provider` by its `id`.

.. http:get:: /api/providers/

  List all :class:`~api.models.Provider`\s.

.. http:post:: /api/providers/

  Create a new :class:`~api.models.Provider`.


Flavors
=======

.. http:get:: /api/flavors/(string:id)/

  Retrieve a :class:`~api.models.Flavor` by its `id`.

.. http:patch:: /api/flavors/(string:id)/

  Update parts of a :class:`~api.models.Flavor`.

.. http:delete:: /api/flavors/(string:id)/

  Destroy a :class:`~api.models.Flavor` by its `id`.

.. http:get:: /api/flavors/

  List all :class:`~api.models.Flavor`\s.

.. http:post:: /api/flavors/

  Create a new :class:`~api.models.Flavor`.


Formations
==========

.. http:get:: /api/formations/(string:id)/

  Retrieve a :class:`~api.models.Formation` by its `id`.

.. http:patch:: /api/formations/(string:id)/

  Update parts of a :class:`~api.models.Formation`.

.. http:delete:: /api/formations/(string:id)/

  Destroy a :class:`~api.models.Formation` by its `id`.

.. http:get:: /api/formations/

  List all :class:`~api.models.Formation`\s.

.. http:post:: /api/formations/

  Create a new :class:`~api.models.Formation`.

  See also
  :meth:`FormationViewSet.post_save() <api.views.FormationViewSet.post_save>`


Formation Infrastructure
------------------------

.. http:get:: /api/formations/(string:id)/layers/(string:id)/

  Retrieve a :class:`~api.models.Layer` by its `id`.

.. http:patch:: /api/formations/(string:id)/layers/(string:id)/

  Update parts of a :class:`~api.models.Layer`.

.. http:delete:: /api/formations/(string:id)/layers/(string:id)/

  Destroy a :class:`~api.models.Layer` by its `id`.

  See also
  :meth:`FormationLayerViewSet.destroy() <api.views.FormationLayerViewSet.destroy>`

.. http:get:: /api/formations/(string:id)/layers/

  List all :class:`~api.models.Layer`\s.

.. http:post:: /api/formations/(string:id)/layers/

  Create a new :class:`~api.models.Layer`.

  See also
  :meth:`FormationLayerViewSet.create() <api.views.FormationLayerViewSet.create>`

.. http:get:: /api/formations/(string:id)/nodes/(string:id)/

  Retrieve a :class:`~api.models.Node` by its `id`.

.. http:post:: /api/formations/(string:id)/nodes/

  Create a new :class:`~api.models.Node` for an existing instance.

.. http:delete:: /api/formations/(string:id)/nodes/(string:id)/

  Destroy a :class:`~api.models.Node` by its `id`.

.. http:get:: /api/formations/(string:id)/nodes/

  List all :class:`~api.models.Node`\s.


Formation Actions
-----------------

.. http:post:: /api/formations/(string:id)/scale/

  See also
  :meth:`FormationViewSet.scale() <api.views.FormationViewSet.scale>`

.. http:post:: /api/formations/(string:id)/balance/

  See also
  :meth:`FormationViewSet.balance() <api.views.FormationViewSet.balance>`

.. http:post:: /api/formations/(string:id)/calculate/

  See also
  :meth:`FormationViewSet.calculate() <api.views.FormationViewSet.calculate>`

.. http:post:: /api/formations/(string:id)/converge/

  See also
  :meth:`FormationViewSet.converge() <api.views.FormationViewSet.converge>`


Applications
============

.. http:get:: /api/apps/(string:id)/

  Retrieve a :class:`~api.models.Application` by its `id`.

.. http:delete:: /api/apps/(string:id)/

  Destroy a :class:`~api.models.Formation` by its `id`.

.. http:get:: /api/apps/

  List all :class:`~api.models.Formation`\s.

.. http:post:: /api/apps/

  Create a new :class:`~api.models.Formation`.


Application Release Components
------------------------------

.. http:get:: /api/apps/(string:id)/config/

  List all :class:`~api.models.Config`\s.

.. http:post:: /api/apps/(string:id)/config/

  Create a new :class:`~api.models.Config`.

.. http:get:: /api/apps/(string:id)/builds/(string:uuid)/

  Retrieve a :class:`~api.models.Build` by its `uuid`.

.. http:get:: /api/apps/(string:id)/builds/

  List all :class:`~api.models.Build`\s.

.. http:post:: /api/apps/(string:id)/builds/

  Create a new :class:`~api.models.Build`.

.. http:get:: /api/apps/(string:id)/releases/(int:version)/

  Retrieve a :class:`~api.models.Release` by its `version`.

.. http:get:: /api/apps/(string:id)/releases/

  List all :class:`~api.models.Release`\s.

.. http:post:: /api/apps/(string:id)/releases/rollback/

  Rollback to a previous :class:`~api.models.Release`.


Application Infrastructure
--------------------------

.. http:get:: /api/apps/(string:id)/containers/(string:type)/(int:num)/

  List all :class:`~api.models.Container`\s.

.. http:get:: /api/apps/(string:id)/containers/(string:type)/

  List all :class:`~api.models.Container`\s.

.. http:get:: /api/apps/(string:id)/containers/

  List all :class:`~api.models.Container`\s.


Application Actions
-------------------

.. http:post:: /api/apps/(string:id)/scale/

  See also
  :meth:`AppViewSet.scale() <api.views.AppViewSet.scale>`

.. http:post:: /api/apps/(string:id)/logs/

  See also
  :meth:`AppViewSet.logs() <api.views.AppViewSet.logs>`

.. http:post:: /api/apps/(string:id)/run/

  See also
  :meth:`AppViewSet.run() <api.views.AppViewSet.run>`

.. http:post:: /api/apps/(string:id)/calculate/

  See also
  :meth:`AppViewSet.calculate() <api.views.AppViewSet.calculate>`


Application Sharing
===================

.. http:delete:: /api/apps/(string:id)/perms/(string:username)/

  Destroy an app permission by its `username`.

.. http:get:: /api/apps/(string:id)/perms/

  List all permissions granted to this app.

.. http:post:: /api/apps/(string:id)/perms/

  Create a new app permission.


API Hooks
=========

.. http:post:: /api/hooks/push/

  Create a new :class:`~api.models.Push`.

.. http:post:: /api/hooks/build/

  Create a new :class:`~api.models.Build`.


Nodes
=====

.. http:get:: /api/nodes/(string:id)/

  Retrieve a :class:`~api.models.Node` by its `id`.

.. http:patch:: /api/nodes/(string:id)/

  Update parts of a :class:`~api.models.Node`.

.. http:delete:: /api/nodes/(string:id)/

  Destroy a :class:`~api.models.Node` by its `id`.

.. http:get:: /api/nodes/

  List all :class:`~api.models.Node`\s.

.. http:post:: /api/nodes/

  Create a new :class:`~api.models.Node`.


Auth
====

.. http:post:: /api/auth/register/

  Create a new :class:`~django.contrib.auth.models.User`.

.. http:delete:: /api/auth/register/

  Destroy the logged-in :class:`~django.contrib.auth.models.User`.

.. http:post:: /api/auth/login

  Authenticate for the REST framework.

.. http:post:: /api/auth/logout

  Clear authentication for the REST framework.

.. http:get:: /api/generate-api-key/

  Generate an API key.


Admin Sharing
=============

.. http:delete:: /api/admin/perms/(string:username)/

  Destroy an admin permission by its `username`.

.. http:get:: /api/admin/perms/

  List all admin permissions granted.

.. http:post:: /api/admin/perms/

  Create a new admin permission.

"""

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
        views.KeyViewSet.as_view({'get': 'list', 'post': 'create'})),
    # provider
    url(r'^providers/(?P<id>[-_\w]+)/?',
        views.ProviderViewSet.as_view({
            'get': 'retrieve', 'patch': 'partial_update', 'delete': 'destroy'})),
    url(r'^providers/?',
        views.ProviderViewSet.as_view({'get': 'list', 'post': 'create'})),
    # flavor
    url(r'^flavors/(?P<id>[-_\w]+)/?',
        views.FlavorViewSet.as_view({
            'get': 'retrieve', 'patch': 'partial_update', 'delete': 'destroy'})),
    url(r'^flavors/?',
        views.FlavorViewSet.as_view({'get': 'list', 'post': 'create'})),
    # formation infrastructure
    url(r'^formations/(?P<id>[-_\w]+)/layers/(?P<layer>[-_\w]+)/?',
        views.FormationLayerViewSet.as_view({
            'get': 'retrieve', 'patch': 'partial_update', 'delete': 'destroy'})),
    url(r'^formations/(?P<id>[-_\w]+)/layers/?',
        views.FormationLayerViewSet.as_view({'get': 'list', 'post': 'create'})),
    url(r'^formations/(?P<id>[-_\w]+)/nodes/(?P<node>[-_\w]+)/?',
        views.FormationNodeViewSet.as_view({
            'get': 'retrieve', 'delete': 'destroy'})),
    url(r'^formations/(?P<id>[-_\w]+)/nodes/?',
        views.FormationNodeViewSet.as_view({'get': 'list', 'post': 'add'})),
    # formation actions
    url(r'^formations/(?P<id>[-_\w]+)/scale/?',
        views.FormationViewSet.as_view({'post': 'scale'})),
    url(r'^formations/(?P<id>[-_\w]+)/balance/?',
        views.FormationViewSet.as_view({'post': 'balance'})),
    url(r'^formations/(?P<id>[-_\w]+)/calculate/?',
        views.FormationViewSet.as_view({'post': 'calculate'})),
    url(r'^formations/(?P<id>[-_\w]+)/converge/?',
        views.FormationViewSet.as_view({'post': 'converge'})),
    # formation base endpoint
    url(r'^formations/(?P<id>[-_\w]+)/?',
        views.FormationViewSet.as_view({
            'get': 'retrieve', 'patch': 'partial_update', 'delete': 'destroy'})),
    url(r'^formations/?',
        views.FormationViewSet.as_view({'get': 'list', 'post': 'create'})),
    # application release components
    url(r'^apps/(?P<id>[-_\w]+)/config/?',
        views.AppConfigViewSet.as_view({'get': 'retrieve', 'post': 'create'})),
    url(r'^apps/(?P<id>[-_\w]+)/builds/(?P<uuid>[-_\w]+)/?',
        views.AppBuildViewSet.as_view({'get': 'retrieve'})),
    url(r'^apps/(?P<id>[-_\w]+)/builds/?',
        views.AppBuildViewSet.as_view({'get': 'list', 'post': 'create'})),
    url(r'^apps/(?P<id>[-_\w]+)/releases/v(?P<version>[0-9]+)/?',
        views.AppReleaseViewSet.as_view({'get': 'retrieve'})),
    url(r'^apps/(?P<id>[-_\w]+)/releases/rollback/?',
        views.AppReleaseViewSet.as_view({'post': 'rollback'})),
    url(r'^apps/(?P<id>[-_\w]+)/releases/?',
        views.AppReleaseViewSet.as_view({'get': 'list'})),
    # application infrastructure
    url(r'^apps/(?P<id>[-_\w]+)/containers/(?P<type>[-_\w]+)/(?P<num>[-_\w]+)/?',
        views.AppContainerViewSet.as_view({'get': 'retrieve'})),
    url(r'^apps/(?P<id>[-_\w]+)/containers/(?P<type>[-_\w.]+)/?',
        views.AppContainerViewSet.as_view({'get': 'list'})),
    url(r'^apps/(?P<id>[-_\w]+)/containers/?',
        views.AppContainerViewSet.as_view({'get': 'list'})),
    # application actions
    url(r'^apps/(?P<id>[-_\w]+)/scale/?',
        views.AppViewSet.as_view({'post': 'scale'})),
    url(r'^apps/(?P<id>[-_\w]+)/logs/?',
        views.AppViewSet.as_view({'post': 'logs'})),
    url(r'^apps/(?P<id>[-_\w]+)/run/?',
        views.AppViewSet.as_view({'post': 'run'})),
    url(r'^apps/(?P<id>[-_\w]+)/calculate/?',
        views.AppViewSet.as_view({'post': 'calculate'})),
    # apps sharing
    url(r'^apps/(?P<id>[-_\w]+)/perms/(?P<username>[-_\w]+)/?',
        views.AppPermsViewSet.as_view({'delete': 'destroy'})),
    url(r'^apps/(?P<id>[-_\w]+)/perms/?',
        views.AppPermsViewSet.as_view({'get': 'list', 'post': 'create'})),
    # apps base endpoint
    url(r'^apps/(?P<id>[-_\w]+)/?',
        views.AppViewSet.as_view({'get': 'retrieve', 'delete': 'destroy'})),
    url(r'^apps/?',
        views.AppViewSet.as_view({'get': 'list', 'post': 'create'})),
    # hooks
    url(r'^hooks/push/?',
        views.PushHookViewSet.as_view({'post': 'create'})),
    url(r'^hooks/build/?',
        views.BuildHookViewSet.as_view({'post': 'create'})),
    # nodes
    url(r'^nodes/(?P<node>[-_\w]+)/converge/?',
        views.NodeViewSet.as_view({'post': 'converge'})),
    url(r'^nodes/(?P<node>[-_\w]+)/?',
        views.NodeViewSet.as_view({
            'get': 'retrieve', 'delete': 'destroy'})),
    url(r'^nodes/?',
        views.NodeViewSet.as_view({'get': 'list'})),
    # authn / authz
    url(r'^auth/register/?',
        views.UserRegistrationView.as_view({'post': 'create'})),
    url(r'^auth/cancel/?',
        views.UserCancellationView.as_view({'delete': 'destroy'})),
    url(r'^auth/',
        include('rest_framework.urls', namespace='rest_framework')),
    url(r'^generate-api-key/',
        'rest_framework.authtoken.views.obtain_auth_token'),
    # admin sharing
    url(r'^admin/perms/(?P<username>[-_\w]+)/?',
        views.AdminPermsViewSet.as_view({'delete': 'destroy'})),
    url(r'^admin/perms/?',
        views.AdminPermsViewSet.as_view({'get': 'list', 'post': 'create'})),
)
