"""
RESTful view classes for presenting Deis API objects.
"""

from __future__ import absolute_import
from __future__ import unicode_literals
import json

from Crypto.PublicKey import RSA
from django.contrib.auth.models import AnonymousUser
from django.contrib.auth.models import User
from django.db import transaction
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

from api import models, serializers, tasks

from deis import settings


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
        elif hasattr(obj, 'formation'):
            return obj.formation.owner == request.user
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
    permission_classes = (IsAnonymous,)
    serializer_class = serializers.UserSerializer

    def post_save(self, user, created=False):
        """Seed both `Providers` and `Flavors` after registration."""
        if created:
            models.Provider.objects.seed(user)
            models.Flavor.objects.seed(user)

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


class KeyViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.Key`."""

    model = models.Key
    serializer_class = serializers.KeySerializer
    lookup_field = 'id'

    def post_save(self, key, created=False, **kwargs):
        tasks.converge_controller.apply_async().wait()


class ProviderViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.Provider`."""

    model = models.Provider
    serializer_class = serializers.ProviderSerializer
    lookup_field = 'id'


class FlavorViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.Flavor`."""

    model = models.Flavor
    serializer_class = serializers.FlavorSerializer
    lookup_field = 'id'

    def update(self, request, *args, **kwargs):
        """
        Override default update behavior to ensure that the params
        field is handled as a dict merge.
        """
        params = self.get_object().params
        new_params = json.loads(request.DATA.get('params', '{}'))
        params.update(new_params)
        # remove param if we provided a null value
        [params.pop(k) for k, v in params.items() if v is None]
        request.DATA['params'] = json.dumps(params)
        return super(FlavorViewSet, self).update(request, *args, **kwargs)


class FormationViewSet(viewsets.ModelViewSet):
    """RESTful views for :class:`~api.models.Formation`."""

    model = models.Formation
    serializer_class = serializers.FormationSerializer
    permission_classes = (permissions.IsAuthenticated, IsAdminOrSafeMethod)
    lookup_field = 'id'

    def pre_save(self, obj):
        if not hasattr(obj, 'owner'):
            obj.owner = self.request.user

    def post_save(self, formation, created=False, **kwargs):
        if created:
            formation.build()

    def scale(self, request, **kwargs):
        new_structure = {}
        try:
            for target, count in request.DATA.items():
                new_structure[target] = int(count)
        except ValueError:
            return Response('Invalid scaling format',
                            status=status.HTTP_400_BAD_REQUEST)
        # check for empty credentials
        for p in models.Provider.objects.filter(owner=request.user):
            if p.creds:
                break
        else:
            return Response('No provider credentials available',
                            status=status.HTTP_400_BAD_REQUEST)
        formation = self.get_object()
        try:
            databag = models.Node.objects.scale(formation, new_structure)
        except (models.Layer.DoesNotExist, EnvironmentError) as err:
            return Response(str(err),
                            status=status.HTTP_400_BAD_REQUEST)
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def balance(self, request, **kwargs):
        formation = self.get_object()
        databag = models.Container.objects.balance(formation)
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def calculate(self, request, **kwargs):
        formation = self.get_object()
        databag = formation.calculate()
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def converge(self, request, **kwargs):
        formation = self.get_object()
        databag = formation.converge()
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def destroy(self, request, **kwargs):
        formation = self.get_object()
        formation.destroy()
        return Response(status=status.HTTP_204_NO_CONTENT)


class FormationScopedViewSet(viewsets.ModelViewSet):

    permission_classes = (permissions.IsAuthenticated, IsAdmin)

    def pre_save(self, obj):
        if not hasattr(obj, 'owner'):
            obj.owner = self.request.user

    def get_queryset(self, **kwargs):
        formations = models.Formation.objects.all()
        formation = get_object_or_404(formations, id=self.kwargs['id'])
        return self.model.objects.filter(formation=formation)


class FormationLayerViewSet(FormationScopedViewSet):
    """RESTful views for :class:`~api.models.Layer`."""

    model = models.Layer
    serializer_class = serializers.LayerSerializer

    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = get_object_or_404(qs, id=self.kwargs['layer'])
        return obj

    def create(self, request, **kwargs):
        request._data = request.DATA.copy()
        formation = models.Formation.objects.get(id=self.kwargs['id'])
        request.DATA['formation'] = formation.id
        if not 'ssh_private_key' in request.DATA and not 'ssh_public_key' in request.DATA:
            # SECURITY: figure out best way to get keys with proper entropy
            key = RSA.generate(2048)
            request.DATA['ssh_private_key'] = key.exportKey('PEM')
            request.DATA['ssh_public_key'] = key.exportKey('OpenSSH')
        return super(FormationLayerViewSet, self).create(request, **kwargs)

    def post_save(self, layer, created=False, **kwargs):
        if created:
            layer.build()

    def destroy(self, request, **kwargs):
        layer = self.get_object()
        layer.destroy()
        return Response(status=status.HTTP_204_NO_CONTENT)


class FormationNodeViewSet(FormationScopedViewSet):
    """RESTful views for :class:`~api.models.Node`."""

    model = models.Node
    serializer_class = serializers.NodeSerializer

    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = get_object_or_404(qs, id=self.kwargs['node'])
        return obj

    def add(self, request, **kwargs):
        fqdn = request.DATA['fqdn']
        formation = models.Formation.objects.get(id=self.kwargs['id'])
        layer = models.Layer.objects.get(id=request.DATA['layer'])
        if self.model.objects.filter(fqdn=fqdn, formation=formation, layer=layer).exists():
            msg = "A node with fqdn={} already exists in the {} formation".format(fqdn, formation)
            return Response(data=msg, status=status.HTTP_409_CONFLICT)
        node = models.Node.objects.new(formation, layer, fqdn)
        node.build()
        return Response(status=status.HTTP_201_CREATED)

    def destroy(self, request, **kwargs):
        node = self.get_object()
        node.destroy()
        return Response(status=status.HTTP_204_NO_CONTENT)


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
        app.publish()
        tasks.converge_controller.apply_async().wait()
        models.log_event(app, "User {} was granted access to {}".format(user, app))
        return Response(status=status.HTTP_201_CREATED)

    def destroy(self, request, **kwargs):
        app = get_object_or_404(self.model, id=kwargs['id'])
        if request.user != app.owner:
            return Response(status=status.HTTP_403_FORBIDDEN)
        user = get_object_or_404(User, username=kwargs['username'])
        if user.has_perm(self.perm, app):
            remove_perm(self.perm, user, app)
            app.publish()
            tasks.converge_controller.apply_async().wait()
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


class NodeViewSet(FormationNodeViewSet):
    """RESTful views for :class:`~api.models.Node`."""

    def get_queryset(self, **kwargs):
        return self.model.objects.all()

    def converge(self, request, **kwargs):
        node = self.get_object()
        try:
            output, _ = node.converge()
        except RuntimeError as e:
            return Response(e.output, status=status.HTTP_500_INTERNAL_SERVER_ERROR,
                            content_type='text/plain')
        return Response(output, status=status.HTTP_200_OK, content_type='text/plain')


class AppViewSet(OwnerViewSet):
    """RESTful views for :class:`~api.models.App`."""

    model = models.App
    serializer_class = serializers.AppSerializer
    lookup_field = 'id'
    permission_classes = (permissions.IsAuthenticated, IsAppUser)

    def get_queryset(self, **kwargs):
        """
        Filter Apps by `owner` attribute or the
        `api.use_formation` permission.
        """
        return super(AppViewSet, self).get_queryset(**kwargs) | \
            get_objects_for_user(self.request.user, 'api.use_app')

    def post_save(self, app, created=False, **kwargs):
        if created:
            app.build()
        app.formation.converge(controller=True)

    def pre_save(self, app, created=False, **kwargs):
        if not app.pk and not app.formation.domain and app.formation.app_set.count() > 0:
            raise EnvironmentError('Formation does not support multiple apps')
        return super(AppViewSet, self).pre_save(app, **kwargs)

    def create(self, request, **kwargs):
        if not 'formation' in request.DATA:
            count = models.Formation.objects.count()
            if count == 1:
                request.DATA['formation'] = models.Formation.objects.first()
            elif count == 0:
                return Response('No formations available',
                                status=status.HTTP_400_BAD_REQUEST)
            else:
                return Response('Could not determine default formation',
                                status=status.HTTP_400_BAD_REQUEST)
        try:
            return OwnerViewSet.create(self, request, **kwargs)
        except EnvironmentError as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)

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
            models.Container.objects.scale(app, new_structure)
        except EnvironmentError as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)
        # save new structure now that scaling was successful
        app.containers.update(new_structure)
        app.save()
        databag = app.converge()
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def calculate(self, request, **kwargs):
        app = self.get_object()
        databag = app.calculate()
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def logs(self, request, **kwargs):
        app = self.get_object()
        try:
            logs = app.logs()
        except EnvironmentError:
            return Response("No logs for {}".format(app.id),
                            status=status.HTTP_404_NOT_FOUND,
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

    def destroy(self, request, **kwargs):
        app = self.get_object()
        app.destroy()
        app.formation.converge(controller=True)
        return Response(status=status.HTTP_204_NO_CONTENT)


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

    def post_save(self, obj, created=False):
        if created:
            models.release_signal.send(
                sender=self, config=obj, app=obj.app,
                user=self.request.user)
            # converge after each config update
            obj.app.formation.converge()

    def create(self, request, *args, **kwargs):
        request._data = request.DATA.copy()
        # assume an existing config object exists
        obj = self.get_object()
        # increment version and use the same formation
        request.DATA['version'] = obj.version + 1
        request.DATA['app'] = obj.app
        # merge config values
        values = obj.values.copy()
        provided = json.loads(request.DATA['values'])
        values.update(provided)
        # remove config keys if we provided a null value
        [values.pop(k) for k, v in provided.items() if v is None]
        request.DATA['values'] = values
        return super(AppConfigViewSet, self).create(request, *args, **kwargs)


class AppBuildViewSet(BaseAppViewSet):
    """RESTful views for :class:`~api.models.Build`."""

    model = models.Build
    serializer_class = serializers.BuildSerializer

    def post_save(self, obj, created=False):
        if created:
            models.release_signal.send(
                sender=self, build=obj, app=obj.app,
                user=self.request.user)

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        request._data = request.DATA.copy()
        request.DATA['app'] = app
        return super(AppBuildViewSet, self).create(request, *args, **kwargs)


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
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        last_version = app.release_set.latest().version
        version = int(request.DATA.get('version', last_version - 1))
        if version < 1:
            return Response(status=status.HTTP_404_NOT_FOUND)
        prev = app.release_set.get(version=version)
        with transaction.atomic():
            summary = "{} rolled back to v{}".format(request.user, version)
            app.release_set.create(owner=request.user, version=last_version + 1,
                                   build=prev.build, config=prev.config,
                                   summary=summary)
            app.converge()
        msg = "Rolled back to v{}".format(version)
        return Response(msg, status=status.HTTP_201_CREATED)


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
            return super(BuildHookViewSet, self).create(request, *args, **kwargs)
        raise PermissionDenied()

    def post_save(self, obj, created=False):
        if created:
            # create a new release
            models.release_signal.send(
                sender=self, build=obj, app=obj.app,
                user=obj.owner)
            models.release_signal.send(sender=self, build=obj, app=obj.app, user=obj.owner)
            # see if we need to scale an initial web container
            app = obj.app
            if len(app.formation.node_set.filter(layer__runtime=True)) > 0 and \
               len(app.container_set.filter(type='web')) < 1:
                # scale an initial web containers
                models.Container.objects.scale(app, {'web': 1})
            # publish and converge the application
            app.converge()
