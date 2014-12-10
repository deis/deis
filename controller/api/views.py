"""
RESTful view classes for presenting Deis API objects.
"""

from __future__ import absolute_import
from __future__ import unicode_literals

from django.conf import settings
from django.contrib.auth.models import AnonymousUser, User
from django.core.exceptions import ValidationError
from django.http import Http404
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
from api.permissions import IsAnonymous, IsOwner, IsAppUser, \
    IsAdmin, HasRegistrationAuth, HasBuilderAuth


class AnonymousAuthentication(BaseAuthentication):

    def authenticate(self, request):
        """
        Authenticate the request and return a two-tuple of (user, token).
        """
        user = AnonymousUser()
        return user, None


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


class UserManagementView(viewsets.GenericViewSet,
                         viewsets.mixins.CreateModelMixin,
                         viewsets.mixins.DestroyModelMixin):
    model = User
    permission_classes = (permissions.IsAuthenticated,)

    def passwd(self, request, *args, **kwargs):
        obj = self.request.user
        if not obj.check_password(request.DATA['password']):
            return Response("Current password did not match", status=status.HTTP_400_BAD_REQUEST)
        obj.set_password(request.DATA['new_password'])
        obj.save()
        return Response({'status': 'password set'})

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


class AppPermsViewSet(viewsets.ViewSet):
    """RESTful views for sharing apps with collaborators."""

    model = models.App  # models class
    perm = 'use_app'    # short name for permission

    def list(self, request, **kwargs):
        app = get_object_or_404(self.model, id=kwargs['id'])
        perm_name = "api.{}".format(self.perm)
        if request.user != app.owner and \
                not request.user.has_perm(perm_name, app) and \
                not request.user.is_superuser:
            return Response(status=status.HTTP_403_FORBIDDEN)
        usernames = [u.username for u in get_users_with_perms(app)
                     if u.has_perm(perm_name, app)]
        return Response({'users': usernames})

    def create(self, request, **kwargs):
        app = get_object_or_404(self.model, id=kwargs['id'])
        if request.user != app.owner and not request.user.is_superuser:
            return Response(status=status.HTTP_403_FORBIDDEN)
        user = get_object_or_404(User, username=request.DATA['username'])
        assign_perm(self.perm, user, app)
        models.log_event(app, "User {} was granted access to {}".format(user, app))
        return Response(status=status.HTTP_201_CREATED)

    def destroy(self, request, **kwargs):
        app = get_object_or_404(self.model, id=kwargs['id'])
        if request.user != app.owner and not request.user.is_superuser:
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
        except (TypeError, ValueError):
            return Response('Invalid scaling format',
                            status=status.HTTP_400_BAD_REQUEST)
        app = self.get_object()
        try:
            models.validate_app_structure(new_structure)
            app.scale(request.user, new_structure)
        except (EnvironmentError, ValidationError) as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError as e:
            return Response(str(e), status=status.HTTP_503_SERVICE_UNAVAILABLE)
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
            output_and_rc = app.run(self.request.user, command)
        except EnvironmentError as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError as e:
            return Response(str(e), status=status.HTTP_503_SERVICE_UNAVAILABLE)
        return Response(output_and_rc, status=status.HTTP_200_OK,
                        content_type='text/plain')

    def destroy(self, request, **kwargs):
        obj = get_object_or_404(self.model, id=kwargs['id'])
        obj.delete()
        return Response(status=status.HTTP_204_NO_CONTENT)


class BaseAppViewSet(viewsets.ModelViewSet):

    permission_classes = (permissions.IsAuthenticated, IsAppUser)

    def pre_save(self, obj):
        obj.owner = self.request.user

    def get_queryset(self, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        try:
            self.check_object_permissions(self.request, app)
        except PermissionDenied:
            raise Http404("No {} matches the given query.".format(
                self.model._meta.object_name))
        return self.model.objects.filter(app=app)

    def get_object(self, *args, **kwargs):
        obj = self.get_queryset().latest('created')
        self.check_object_permissions(self.request, obj)
        return obj


class AppBuildViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Build`."""

    model = models.Build
    serializer_class = serializers.BuildSerializer

    def post_save(self, build, created=False):
        if created:
            self.release = build.create(self.request.user)

    def get_success_headers(self, data):
        headers = super(AppBuildViewSet, self).get_success_headers(data)
        headers.update({'X-Deis-Release': self.release.version})
        return headers

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        request._data = request.DATA.copy()
        request.DATA['app'] = app
        try:
            return super(AppBuildViewSet, self).create(request, *args, **kwargs)
        except RuntimeError as e:
            return Response(str(e), status=status.HTTP_503_SERVICE_UNAVAILABLE)


class AppConfigViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Config`."""

    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def get_object(self, *args, **kwargs):
        """Return the Config associated with the App's latest Release."""
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        try:
            self.check_object_permissions(self.request, app)
            return app.release_set.latest().config
        except (PermissionDenied, models.Release.DoesNotExist):
            raise Http404("No {} matches the given query.".format(
                self.model._meta.object_name))

    def pre_save(self, config):
        """merge the old config with the new"""
        previous_config = config.app.config_set.latest()
        config.owner = self.request.user
        if previous_config:
            for attr in ['cpu', 'memory', 'tags', 'values']:
                # Guard against migrations from older apps without fixes to
                # JSONField encoding.
                try:
                    data = getattr(previous_config, attr).copy()
                except AttributeError:
                    data = {}
                try:
                    new_data = getattr(config, attr).copy()
                except AttributeError:
                    new_data = {}
                data.update(new_data)
                # remove config keys if we provided a null value
                [data.pop(k) for k, v in new_data.items() if v is None]
                setattr(config, attr, data)

    def post_save(self, config, created=False):
        if created:
            release = config.app.release_set.latest()
            self.release = release.new(self.request.user, config=config, build=release.build)
            try:
                config.app.deploy(self.request.user, self.release)
            except RuntimeError:
                self.release.delete()
                raise

    def get_success_headers(self, data):
        headers = super(AppConfigViewSet, self).get_success_headers(data)
        headers.update({'X-Deis-Release': self.release.version})
        return headers

    def create(self, request, *args, **kwargs):
        obj = self.get_object()
        request.DATA['app'] = obj.app
        try:
            return super(AppConfigViewSet, self).create(request, *args, **kwargs)
        except RuntimeError as e:
            return Response(str(e), status=status.HTTP_503_SERVICE_UNAVAILABLE)


class AppReleaseViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Release`."""

    model = models.Release
    serializer_class = serializers.ReleaseSerializer

    def get_object(self, *args, **kwargs):
        """Get Release by version always."""
        return self.get_queryset(**kwargs).get(version=self.kwargs['version'])

    def rollback(self, request, *args, **kwargs):
        """
        Create a new release as a copy of the state of the compiled slug and
        config vars of a previous release.
        """
        try:
            app = get_object_or_404(models.App, id=self.kwargs['id'])
            release = app.release_set.latest()
            version_to_rollback_to = release.version - 1
            if request.DATA.get('version'):
                version_to_rollback_to = int(request.DATA['version'])
            new_release = release.rollback(request.user, version_to_rollback_to)
            response = {'version': new_release.version}
            return Response(response, status=status.HTTP_201_CREATED)
        except EnvironmentError as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError as e:
            new_release.delete()
            return Response(str(e), status=status.HTTP_503_SERVICE_UNAVAILABLE)


class AppContainerViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Container`."""

    model = models.Container
    serializer_class = serializers.ContainerSerializer

    def get_queryset(self, **kwargs):
        qs = super(AppContainerViewSet, self).get_queryset(**kwargs)
        container_type = self.kwargs.get('type')
        if container_type:
            qs = qs.filter(type=container_type)
        else:
            qs = qs.exclude(type='run')
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
        if user == app.owner or \
           user in get_users_with_perms(app) or \
           user.is_superuser:
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
        self.user = get_object_or_404(
            User, username=request.DATA['receive_user'])
        # check the user is authorized for this app
        if self.user == app.owner or \
           self.user in get_users_with_perms(app) or \
           self.user.is_superuser:
            request._data = request.DATA.copy()
            request.DATA['app'] = app
            request.DATA['owner'] = self.user
            try:
                super(BuildHookViewSet, self).create(request, *args, **kwargs)
                # return the application databag
                response = {'release': {'version': app.release_set.latest().version},
                            'domains': ['.'.join([app.id, settings.DEIS_DOMAIN])]}
                return Response(response, status=status.HTTP_200_OK)
            except RuntimeError as e:
                return Response(str(e), status=status.HTTP_503_SERVICE_UNAVAILABLE)
        raise PermissionDenied()

    def post_save(self, build, created=False):
        if created:
            build.create(self.user)


class ConfigHookViewSet(BaseHookViewSet):
    """API hook to grab latest :class:`~api.models.Config`"""

    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.DATA['receive_repo'])
        user = get_object_or_404(
            User, username=request.DATA['receive_user'])
        # check the user is authorized for this app
        if user == app.owner or \
           user in get_users_with_perms(app) or \
           user.is_superuser:
            config = app.release_set.latest().config
            serializer = self.get_serializer(config)
            return Response(serializer.data, status=status.HTTP_200_OK)
        raise PermissionDenied()
