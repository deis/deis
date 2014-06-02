"""
RESTful view classes for presenting Deis API objects.
"""

from __future__ import absolute_import
from __future__ import unicode_literals
import json

from django.contrib.auth.models import AnonymousUser
from django.contrib.auth.models import User
from django.utils import timezone
from guardian.shortcuts import assign_perm
from guardian.shortcuts import get_objects_for_user
from guardian.shortcuts import get_users_with_perms
from guardian.shortcuts import remove_perm
from rest_framework import permissions
from rest_framework import status
from rest_framework import viewsets
from rest_framework.authentication import BaseAuthentication
from rest_framework.exceptions import PermissionDenied
from rest_framework.generics import get_object_or_404
from rest_framework.response import Response

from api import models, serializers

from django.conf import settings


class AnonymousAuthentication(BaseAuthentication):

    def authenticate(self, request):
        """
        Authenticate the request and return a two-tuple of (user, token).
        """
        user = AnonymousUser()
        return user, None


class IsAnonymous(permissions.BasePermission):
    """
    View permission to allow anonymous users.
    """

    def has_permission(self, request, view):
        """
        Return `True` if permission is granted, `False` otherwise.
        """
        return type(request.user) is AnonymousUser


class IsOwner(permissions.BasePermission):
    """
    Object-level permission to allow only owners of an object to access it.
    Assumes the model instance has an `owner` attribute.
    """

    def has_object_permission(self, request, view, obj):
        if hasattr(obj, 'owner'):
            return obj.owner == request.user
        else:
            return False


class IsAppUser(permissions.BasePermission):
    """
    Object-level permission to allow owners or collaborators to access
    an app-related model.
    """
    def has_object_permission(self, request, view, obj):
        if isinstance(obj, models.App) and obj.owner == request.user:
            return True
        elif hasattr(obj, 'app') and obj.app.owner == request.user:
            return True
        elif request.user.has_perm('use_app', obj):
            return request.method != 'DELETE'
        elif hasattr(obj, 'app') and request.user.has_perm('use_app', obj.app):
            return request.method != 'DELETE'
        else:
            return False


class IsAdmin(permissions.BasePermission):
    """
    View permission to allow only admins.
    """

    def has_permission(self, request, view):
        """
        Return `True` if permission is granted, `False` otherwise.
        """
        return request.user.is_superuser


class IsAdminOrSafeMethod(permissions.BasePermission):
    """
    View permission to allow only admins to use unsafe methods
    including POST, PUT, DELETE.

    This allows
    """

    def has_permission(self, request, view):
        """
        Return `True` if permission is granted, `False` otherwise.
        """
        return request.method in permissions.SAFE_METHODS or request.user.is_superuser


class HasRegistrationAuth(permissions.BasePermission):
    """
    Checks to see if registration is enabled
    """
    def has_permission(self, request, view):
        return settings.REGISTRATION_ENABLED


class HasBuilderAuth(permissions.BasePermission):
    """
    View permission to allow builder to perform actions
    with a special HTTP header
    """

    def has_permission(self, request, view):
        """
        Return `True` if permission is granted, `False` otherwise.
        """
        auth_header = request.environ.get('HTTP_X_DEIS_BUILDER_AUTH')
        if not auth_header:
            return False
        return auth_header == settings.BUILDER_KEY


class UserRegistrationView(viewsets.GenericViewSet,
                           viewsets.mixins.CreateModelMixin):
    model = User

    authentication_classes = (AnonymousAuthentication,)
    permission_classes = (IsAnonymous, HasRegistrationAuth)
    serializer_class = serializers.UserSerializer

    def pre_save(self, obj):
        """Replicate UserManager.create_user functionality."""
        now = timezone.now()
        obj.last_login = now
        obj.date_joined = now
        obj.is_active = True
        obj.email = User.objects.normalize_email(obj.email)
        obj.set_password(obj.password)
        # Make this first signup an admin / superuser
        if not User.objects.filter(is_superuser=True).exists():
            obj.is_superuser = obj.is_staff = True


class UserCancellationView(viewsets.GenericViewSet,
                           viewsets.mixins.DestroyModelMixin):
    model = User

    permission_classes = (permissions.IsAuthenticated,)

    def destroy(self, request, *args, **kwargs):
        obj = self.request.user
        obj.delete()
        return Response(status=status.HTTP_204_NO_CONTENT)


class OwnerViewSet(viewsets.ModelViewSet):
    """Scope views to an `owner` attribute."""

    permission_classes = (permissions.IsAuthenticated, IsOwner)

    def pre_save(self, obj):
        obj.owner = self.request.user

    def get_queryset(self, **kwargs):
        """Filter all querysets by an `owner` attribute.
        """
        return self.model.objects.filter(owner=self.request.user)


class ClusterViewSet(viewsets.ModelViewSet):
    """RESTful views for :class:`~api.models.Cluster`."""

    model = models.Cluster
    serializer_class = serializers.ClusterSerializer
    permission_classes = (permissions.IsAuthenticated, IsAdmin)
    lookup_field = 'id'

    def pre_save(self, obj):
        if not hasattr(obj, 'owner'):
            obj.owner = self.request.user

    def post_save(self, cluster, created=False, **kwargs):
        if created:
            cluster.create()

    def pre_delete(self, cluster):
        cluster.destroy()


class AppPermsViewSet(viewsets.ViewSet):
    """RESTful views for sharing apps with collaborators."""

    model = models.App  # models class
    perm = 'use_app'    # short name for permission

    def list(self, request, **kwargs):
        app = get_object_or_404(self.model, id=kwargs['id'])
        perm_name = "api.{}".format(self.perm)
        if request.user != app.owner and not request.user.has_perm(perm_name, app):
            return Response(status=status.HTTP_403_FORBIDDEN)
        usernames = [u.username for u in get_users_with_perms(app)
                     if u.has_perm(perm_name, app)]
        return Response({'users': usernames})

    def create(self, request, **kwargs):
        app = get_object_or_404(self.model, id=kwargs['id'])
        if request.user != app.owner:
            return Response(status=status.HTTP_403_FORBIDDEN)
        user = get_object_or_404(User, username=request.DATA['username'])
        assign_perm(self.perm, user, app)
        models.log_event(app, "User {} was granted access to {}".format(user, app))
        return Response(status=status.HTTP_201_CREATED)

    def destroy(self, request, **kwargs):
        app = get_object_or_404(self.model, id=kwargs['id'])
        if request.user != app.owner:
            return Response(status=status.HTTP_403_FORBIDDEN)
        user = get_object_or_404(User, username=kwargs['username'])
        if user.has_perm(self.perm, app):
            remove_perm(self.perm, user, app)
            models.log_event(app, "User {} was revoked access to {}".format(user, app))
            return Response(status=status.HTTP_204_NO_CONTENT)
        else:
            return Response(status=status.HTTP_404_NOT_FOUND)


class AdminPermsViewSet(viewsets.ModelViewSet):
    """RESTful views for sharing admin permissions with other users."""

    model = User
    serializer_class = serializers.AdminUserSerializer
    permission_classes = (IsAdmin,)

    def get_queryset(self, **kwargs):
        return self.model.objects.filter(is_active=True, is_superuser=True)

    def create(self, request, **kwargs):
        user = get_object_or_404(User, username=request.DATA['username'])
        user.is_superuser = user.is_staff = True
        user.save(update_fields=['is_superuser', 'is_staff'])
        return Response(status=status.HTTP_201_CREATED)

    def destroy(self, request, **kwargs):
        user = get_object_or_404(User, username=kwargs['username'])
        user.is_superuser = user.is_staff = False
        user.save(update_fields=['is_superuser', 'is_staff'])
        return Response(status=status.HTTP_204_NO_CONTENT)


class AppViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.App`."""

    model = models.App
    serializer_class = serializers.AppSerializer
    lookup_field = 'id'
    permission_classes = (permissions.IsAuthenticated, IsAppUser)

    def get_queryset(self, **kwargs):
        """
        Filter Apps by `owner` attribute or the `api.use_app` permission.
        """
        return super(AppViewSet, self).get_queryset(**kwargs) | \
            get_objects_for_user(self.request.user, 'api.use_app')

    def post_save(self, app, created=False, **kwargs):
        if created:
            app.create()

    def scale(self, request, **kwargs):
        new_structure = {}
        try:
            for target, count in request.DATA.items():
                new_structure[target] = int(count)
        except ValueError:
            return Response('Invalid scaling format',
                            status=status.HTTP_400_BAD_REQUEST)
        app = self.get_object()
        try:
            app.structure = new_structure
            app.scale()
        except EnvironmentError as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)
        return Response(status=status.HTTP_204_NO_CONTENT,
                        content_type='application/json')

    def logs(self, request, **kwargs):
        app = self.get_object()
        try:
            logs = app.logs()
        except EnvironmentError:
            return Response("No logs for {}".format(app.id),
                            status=status.HTTP_204_NO_CONTENT,
                            content_type='text/plain')
        return Response(logs, status=status.HTTP_200_OK,
                        content_type='text/plain')

    def run(self, request, **kwargs):
        app = self.get_object()
        command = request.DATA['command']
        try:
            output_and_rc = app.run(command)
        except EnvironmentError as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)
        return Response(output_and_rc, status=status.HTTP_200_OK,
                        content_type='text/plain')


class BaseAppViewSet(viewsets.ModelViewSet):

    permission_classes = (permissions.IsAuthenticated, IsAppUser)

    def pre_save(self, obj):
        obj.owner = self.request.user

    def get_queryset(self, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        return self.model.objects.filter(app=app)

    def get_object(self, *args, **kwargs):
        obj = self.get_queryset().latest('created')
        user = self.request.user
        if user == obj.app.owner or user in get_users_with_perms(obj.app):
            return obj
        raise PermissionDenied()


class AppBuildViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Build`."""

    model = models.Build
    serializer_class = serializers.BuildSerializer

    def post_save(self, build, created=False):
        if created:
            release = build.app.release_set.latest()
            self.release = release.new(self.request.user, build=build)
            initial = True if build.app.structure == {} else False
            build.app.deploy(self.release, initial=initial)

    def get_success_headers(self, data):
        headers = super(AppBuildViewSet, self).get_success_headers(data)
        headers.update({'X-Deis-Release': self.release.version})
        return headers

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        request._data = request.DATA.copy()
        request.DATA['app'] = app
        return super(AppBuildViewSet, self).create(request, *args, **kwargs)


class AppConfigViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Config`."""

    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def get_object(self, *args, **kwargs):
        """Return the Config associated with the App's latest Release."""
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        user = self.request.user
        if user == app.owner or user in get_users_with_perms(app):
            return app.release_set.latest().config
        raise PermissionDenied()

    def post_save(self, config, created=False):
        if created:
            release = config.app.release_set.latest()
            self.release = release.new(self.request.user, config=config)
            config.app.deploy(self.release)

    def get_success_headers(self, data):
        headers = super(AppConfigViewSet, self).get_success_headers(data)
        headers.update({'X-Deis-Release': self.release.version})
        return headers

    def create(self, request, *args, **kwargs):
        request._data = request.DATA.copy()
        # assume an existing config object exists
        obj = self.get_object()
        request.DATA['app'] = obj.app
        # merge config values
        values = obj.values.copy()
        provided = json.loads(request.DATA['values'])
        values.update(provided)
        # remove config keys if we provided a null value
        [values.pop(k) for k, v in provided.items() if v is None]
        request.DATA['values'] = values
        return super(AppConfigViewSet, self).create(request, *args, **kwargs)


class AppReleaseViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Release`."""

    model = models.Release
    serializer_class = serializers.ReleaseSerializer

    def get_object(self, *args, **kwargs):
        """Get Release by version always."""
        return self.get_queryset(**kwargs).get(version=self.kwargs['version'])

    # TODO: move logic into model
    def rollback(self, request, *args, **kwargs):
        """
        Create a new release as a copy of the state of the compiled slug and
        config vars of a previous release.
        """
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        release = app.release_set.latest()
        last_version = release.version
        version = int(request.DATA.get('version', last_version - 1))
        if version < 1:
            return Response(status=status.HTTP_404_NOT_FOUND)
        summary = "{} rolled back to v{}".format(request.user, version)
        prev = app.release_set.get(version=version)
        new_release = release.new(
            request.user,
            build=prev.build,
            config=prev.config,
            summary=summary)
        app.deploy(new_release)
        response = {'version': new_release.version}
        return Response(response, status=status.HTTP_201_CREATED)


class AppContainerViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.Container`."""

    model = models.Container
    serializer_class = serializers.ContainerSerializer

    def get_queryset(self, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        qs = self.model.objects.filter(app=app)
        container_type = self.kwargs.get('type')
        if container_type:
            qs = qs.filter(type=container_type)
        return qs

    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = qs.get(num=self.kwargs['num'])
        return obj


class KeyViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.Key`."""

    model = models.Key
    serializer_class = serializers.KeySerializer
    lookup_field = 'id'


class DomainViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.Domain`."""

    model = models.Domain
    serializer_class = serializers.DomainSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        request._data = request.DATA.copy()
        request.DATA['app'] = app
        return super(DomainViewSet, self).create(request, *args, **kwargs)

    def get_queryset(self, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        qs = self.model.objects.filter(app=app)
        return qs

    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = qs.get(domain=self.kwargs['domain'])
        return obj


class BaseHookViewSet(viewsets.ModelViewSet):

    permission_classes = (HasBuilderAuth,)

    def pre_save(self, obj):
        # SECURITY: we trust the username field to map to the owner
        obj.owner = self.request.DATA['owner']


class PushHookViewSet(BaseHookViewSet):
    """API hook to create new :class:`~api.models.Push`"""

    model = models.Push
    serializer_class = serializers.PushSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.DATA['receive_repo'])
        user = get_object_or_404(
            User, username=request.DATA['receive_user'])
        # check the user is authorized for this app
        if user == app.owner or user in get_users_with_perms(app):
            request._data = request.DATA.copy()
            request.DATA['app'] = app
            request.DATA['owner'] = user
            return super(PushHookViewSet, self).create(request, *args, **kwargs)
        raise PermissionDenied()


class BuildHookViewSet(BaseHookViewSet):
    """API hook to create new :class:`~api.models.Build`"""

    model = models.Build
    serializer_class = serializers.BuildSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.DATA['receive_repo'])
        user = get_object_or_404(
            User, username=request.DATA['receive_user'])
        # check the user is authorized for this app
        if user == app.owner or user in get_users_with_perms(app):
            request._data = request.DATA.copy()
            request.DATA['app'] = app
            request.DATA['owner'] = user
            super(BuildHookViewSet, self).create(request, *args, **kwargs)
            # return the application databag
            response = {'release': {'version': app.release_set.latest().version},
                        'domains': ['.'.join([app.id, app.cluster.domain])]}
            return Response(response, status=status.HTTP_200_OK)
        raise PermissionDenied()

    def post_save(self, build, created=False):
        if created:
            release = build.app.release_set.latest()
            new_release = release.new(build.owner, build=build)
            initial = True if build.app.structure == {} else False
            build.app.deploy(new_release, initial=initial)


class ConfigHookViewSet(BaseHookViewSet):
    """API hook to grab latest :class:`~api.models.Config`"""

    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.DATA['receive_repo'])
        user = get_object_or_404(
            User, username=request.DATA['receive_user'])
        # check the user is authorized for this app
        if user == app.owner or user in get_users_with_perms(app):
            config = app.release_set.latest().config
            serializer = self.get_serializer(config)
            return Response(serializer.data, status=status.HTTP_200_OK)
        raise PermissionDenied()
