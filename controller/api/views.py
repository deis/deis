"""
RESTful view classes for presenting Deis API objects.
"""

import os.path
import tempfile

from django.conf import settings
from django.core.exceptions import ValidationError
from django.contrib.auth.models import User
from django.http import HttpResponse
from django.shortcuts import get_object_or_404
from guardian.shortcuts import assign_perm, get_objects_for_user, \
    get_users_with_perms, remove_perm
from django.views.generic import View
from rest_framework import mixins, renderers, status
from rest_framework.exceptions import PermissionDenied
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.viewsets import GenericViewSet
from rest_framework.authtoken.models import Token
from simpleflock import SimpleFlock

from api import authentication, models, permissions, serializers, viewsets

import requests


class HealthCheckView(View):
    """Simple health check view to determine if the server
       is responding to HTTP requests.
    """

    def get(self, request):
        return HttpResponse("OK")
    head = get


class UserRegistrationViewSet(GenericViewSet,
                              mixins.CreateModelMixin):
    """ViewSet to handle registering new users. The logic is in the serializer."""
    authentication_classes = [authentication.AnonymousOrAuthenticatedAuthentication]
    permission_classes = [permissions.HasRegistrationAuth]
    serializer_class = serializers.UserSerializer


class UserManagementViewSet(GenericViewSet):
    serializer_class = serializers.UserSerializer

    def get_queryset(self):
        return User.objects.filter(pk=self.request.user.pk)

    def get_object(self):
        return self.get_queryset()[0]

    def destroy(self, request, **kwargs):
        calling_obj = self.get_object()
        target_obj = calling_obj

        if request.data.get('username'):
            # if you "accidentally" target yourself, that should be fine
            if calling_obj.username == request.data['username'] or calling_obj.is_superuser:
                target_obj = get_object_or_404(User, username=request.data['username'])
            else:
                raise PermissionDenied()

        # A user can not be removed without apps changing ownership first
        if len(models.App.objects.filter(owner=target_obj)) > 0:
            msg = '{} still has applications assigned. Delete or transfer ownership'.format(str(target_obj))  # noqa
            return Response({'detail': msg}, status=status.HTTP_409_CONFLICT)

        target_obj.delete()
        return Response(status=status.HTTP_204_NO_CONTENT)

    def passwd(self, request, **kwargs):
        caller_obj = self.get_object()
        target_obj = self.get_object()
        if request.data.get('username'):
            # if you "accidentally" target yourself, that should be fine
            if caller_obj.username == request.data['username'] or caller_obj.is_superuser:
                target_obj = get_object_or_404(User, username=request.data['username'])
            else:
                raise PermissionDenied()
        if request.data.get('password') or not caller_obj.is_superuser:
            if not target_obj.check_password(request.data['password']):
                return Response({'detail': 'Current password does not match'},
                                status=status.HTTP_400_BAD_REQUEST)
        target_obj.set_password(request.data['new_password'])
        target_obj.save()
        return Response({'status': 'password set'})


class TokenManagementViewSet(GenericViewSet,
                             mixins.DestroyModelMixin):
    serializer_class = serializers.UserSerializer
    permission_classes = [permissions.CanRegenerateToken]

    def get_queryset(self):
        return User.objects.filter(pk=self.request.user.pk)

    def get_object(self):
        return self.get_queryset()[0]

    def regenerate(self, request, **kwargs):
        obj = self.get_object()

        if 'all' in request.data:
            for user in User.objects.all():
                if not user.is_anonymous():
                    token = Token.objects.get(user=user)
                    token.delete()
                    Token.objects.create(user=user)
            return Response("")

        if 'username' in request.data:
            obj = get_object_or_404(User,
                                    username=request.data['username'])
            self.check_object_permissions(self.request, obj)

        token = Token.objects.get(user=obj)
        token.delete()
        token = Token.objects.create(user=obj)
        return Response({'token': token.key})


class BaseDeisViewSet(viewsets.OwnerViewSet):
    """
    A generic ViewSet for objects related to Deis.

    To use it, at minimum you'll need to provide the `serializer_class` attribute and
    the `model` attribute shortcut.
    """
    lookup_field = 'id'
    permission_classes = [IsAuthenticated, permissions.IsAppUser]
    renderer_classes = [renderers.JSONRenderer]

    def create(self, request, *args, **kwargs):
        try:
            return super(BaseDeisViewSet, self).create(request, *args, **kwargs)
        # If the scheduler oopsie'd
        except RuntimeError as e:
            return Response({'detail': str(e)}, status=status.HTTP_503_SERVICE_UNAVAILABLE)


class AppResourceViewSet(BaseDeisViewSet):
    """A viewset for objects which are attached to an application."""

    def get_app(self):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        self.check_object_permissions(self.request, app)
        return app

    def get_queryset(self, **kwargs):
        app = self.get_app()
        return self.model.objects.filter(app=app)

    def get_object(self, **kwargs):
        return self.get_queryset(**kwargs).latest('created')

    def create(self, request, **kwargs):
        request.data['app'] = self.get_app()
        return super(AppResourceViewSet, self).create(request, **kwargs)


class ReleasableViewSet(AppResourceViewSet):
    """A viewset for application resources which affect the release cycle.

    When a resource is created, a new release is created for the application
    and it returns some success headers regarding the new release.

    To use it, at minimum you'll need to provide a `release` attribute tied to your class before
    calling post_save().
    """
    def get_object(self):
        """Retrieve the object based on the latest release's value"""
        return getattr(self.get_app().release_set.latest(), self.model.__name__.lower())

    def get_success_headers(self, data, **kwargs):
        headers = super(ReleasableViewSet, self).get_success_headers(data)
        headers.update({'Deis-Release': self.release.version})
        headers.update({'X-Deis-Release': self.release.version})  # DEPRECATED
        return headers


class AppViewSet(BaseDeisViewSet):
    """A viewset for interacting with App objects."""
    model = models.App
    serializer_class = serializers.AppSerializer

    def get_queryset(self, *args, **kwargs):
        return self.model.objects.all(*args, **kwargs)

    def list(self, request, *args, **kwargs):
        """
        HACK: Instead of filtering by the queryset, we limit the queryset to list only the apps
        which are owned by the user as well as any apps they have been given permission to
        interact with.
        """
        queryset = super(AppViewSet, self).get_queryset(**kwargs) | \
            get_objects_for_user(self.request.user, 'api.use_app')
        instance = self.filter_queryset(queryset)
        page = self.paginate_queryset(instance)
        if page is not None:
            serializer = self.get_pagination_serializer(page)
        else:
            serializer = self.get_serializer(instance, many=True)
        return Response(serializer.data)

    def post_save(self, app):
        app.create()

    def scale(self, request, **kwargs):
        new_structure = {}
        app = self.get_object()
        try:
            for target, count in request.data.viewitems():
                new_structure[target] = int(count)
            models.validate_app_structure(new_structure)
            app.scale(request.user, new_structure)
        except (TypeError, ValueError) as e:
            return Response({'detail': 'Invalid scaling format: {}'.format(e)},
                            status=status.HTTP_400_BAD_REQUEST)
        except (EnvironmentError, ValidationError) as e:
            return Response({'detail': str(e)}, status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError as e:
            return Response({'detail': str(e)}, status=status.HTTP_503_SERVICE_UNAVAILABLE)
        return Response(status=status.HTTP_204_NO_CONTENT)

    def logs(self, request, **kwargs):
        app = self.get_object()
        try:
            return HttpResponse(app.logs(request.query_params.get('log_lines',
                                         str(settings.LOG_LINES))),
                                status=status.HTTP_200_OK, content_type='text/plain')
        except requests.exceptions.RequestException:
            return HttpResponse("Error accessing logs for {}".format(app.id),
                                status=status.HTTP_500_INTERNAL_SERVER_ERROR,
                                content_type='text/plain')
        except EnvironmentError as e:
            if e.message == 'Error accessing deis-logger':
                return HttpResponse("Error accessing logs for {}".format(app.id),
                                    status=status.HTTP_500_INTERNAL_SERVER_ERROR,
                                    content_type='text/plain')
            else:
                return HttpResponse(status=status.HTTP_204_NO_CONTENT)

    def run(self, request, **kwargs):
        app = self.get_object()
        try:
            output_and_rc = app.run(self.request.user, request.data['command'])
        except EnvironmentError as e:
            return Response({'detail': str(e)}, status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError as e:
            return Response({'detail': str(e)}, status=status.HTTP_503_SERVICE_UNAVAILABLE)
        return Response(output_and_rc, status=status.HTTP_200_OK,
                        content_type='text/plain')

    def update(self, request, **kwargs):
        app = self.get_object()

        if request.data.get('owner'):
            if self.request.user != app.owner and not self.request.user.is_superuser:
                raise PermissionDenied()
            new_owner = get_object_or_404(User, username=request.data['owner'])
            app.owner = new_owner
        app.save()
        return Response(status=status.HTTP_200_OK)


class BuildViewSet(ReleasableViewSet):
    """A viewset for interacting with Build objects."""
    model = models.Build
    serializer_class = serializers.BuildSerializer

    def post_save(self, build):
        self.release = build.create(self.request.user)
        super(BuildViewSet, self).post_save(build)


class ConfigViewSet(ReleasableViewSet):
    """A viewset for interacting with Config objects."""
    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def create(self, request, **kwargs):
        # Guard against overlapping config changes, using a filesystem lock so that
        # multiple controller processes can be coordinated.
        # Use a tempfile such as "/tmp/violet-valkyrie-config".
        lockfile = os.path.join(tempfile.gettempdir(), kwargs['id'] + '-config')
        try:
            with SimpleFlock(lockfile, timeout=5):
                return super(ConfigViewSet, self).create(request, **kwargs)
        except IOError as err:
            msg = "Config changes already in progress.\n{}".format(err)
            return Response(status=status.HTTP_409_CONFLICT, data={'error': msg})

    def post_save(self, config):
        release = config.app.release_set.latest()
        self.release = release.new(self.request.user, config=config, build=release.build)
        try:
            config.app.deploy(self.request.user, self.release)
        except RuntimeError:
            self.release.delete()
            raise


class ContainerViewSet(AppResourceViewSet):
    """A viewset for interacting with Container objects."""
    model = models.Container
    serializer_class = serializers.ContainerSerializer

    def get_queryset(self, **kwargs):
        qs = super(ContainerViewSet, self).get_queryset(**kwargs)
        container_type = self.kwargs.get('type')
        if container_type:
            qs = qs.filter(type=container_type)
        else:
            qs = qs.exclude(type='run')
        return qs

    def get_object(self, **kwargs):
        qs = self.get_queryset(**kwargs)
        return qs.get(num=self.kwargs['num'])

    def restart(self, *args, **kwargs):
        try:
            containers = self.get_app().restart(**kwargs)
            serializer = self.get_serializer(containers, many=True)
            return Response(serializer.data, status=status.HTTP_200_OK)
        except Exception as e:
            return Response({'detail': str(e)}, status=status.HTTP_503_SERVICE_UNAVAILABLE)


class DomainViewSet(AppResourceViewSet):
    """A viewset for interacting with Domain objects."""
    model = models.Domain
    serializer_class = serializers.DomainSerializer

    def get_object(self, **kwargs):
        qs = self.get_queryset(**kwargs)
        return qs.get(domain=self.kwargs['domain'])


class CertificateViewSet(BaseDeisViewSet):
    """A viewset for interacting with Domain objects."""
    model = models.Certificate
    serializer_class = serializers.CertificateSerializer

    def get_object(self, **kwargs):
        """Retrieve domain certificate by common name"""
        qs = self.get_queryset(**kwargs)
        return qs.get(common_name=self.kwargs['common_name'])


class KeyViewSet(BaseDeisViewSet):
    """A viewset for interacting with Key objects."""
    model = models.Key
    permission_classes = [IsAuthenticated, permissions.IsOwner]
    serializer_class = serializers.KeySerializer


class ReleaseViewSet(AppResourceViewSet):
    """A viewset for interacting with Release objects."""
    model = models.Release
    serializer_class = serializers.ReleaseSerializer

    def get_object(self, **kwargs):
        """Get release by version always"""
        return self.get_queryset(**kwargs).get(version=self.kwargs['version'])

    def rollback(self, request, **kwargs):
        """
        Create a new release as a copy of the state of the compiled slug and config vars of a
        previous release.
        """
        app = self.get_app()
        try:
            release = app.release_set.latest()
            version_to_rollback_to = release.version - 1
            if request.data.get('version'):
                version_to_rollback_to = int(request.data['version'])
            new_release = release.rollback(request.user, version_to_rollback_to)
            response = {'version': new_release.version}
            return Response(response, status=status.HTTP_201_CREATED)
        except EnvironmentError as e:
            return Response({'detail': str(e)}, status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError:
            new_release.delete()
            raise


class BaseHookViewSet(BaseDeisViewSet):
    permission_classes = [permissions.HasBuilderAuth]


class PushHookViewSet(BaseHookViewSet):
    """API hook to create new :class:`~api.models.Push`"""
    model = models.Push
    serializer_class = serializers.PushSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.data['receive_repo'])
        request.user = get_object_or_404(User, username=request.data['receive_user'])
        # check the user is authorized for this app
        if not permissions.is_app_user(request, app):
            raise PermissionDenied()
        request.data['app'] = app
        request.data['owner'] = request.user
        return super(PushHookViewSet, self).create(request, *args, **kwargs)


class BuildHookViewSet(BaseHookViewSet):
    """API hook to create new :class:`~api.models.Build`"""
    model = models.Build
    serializer_class = serializers.BuildSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.data['receive_repo'])
        self.user = request.user = get_object_or_404(User, username=request.data['receive_user'])
        # check the user is authorized for this app
        if not permissions.is_app_user(request, app):
            raise PermissionDenied()
        request.data['app'] = app
        request.data['owner'] = self.user
        super(BuildHookViewSet, self).create(request, *args, **kwargs)
        # return the application databag
        response = {'release': {'version': app.release_set.latest().version},
                    'domains': ['.'.join([app.id, settings.DEIS_DOMAIN])]}
        return Response(response, status=status.HTTP_200_OK)

    def post_save(self, build):
        build.create(self.user)


class ConfigHookViewSet(BaseHookViewSet):
    """API hook to grab latest :class:`~api.models.Config`"""
    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.data['receive_repo'])
        request.user = get_object_or_404(User, username=request.data['receive_user'])
        # check the user is authorized for this app
        if not permissions.is_app_user(request, app):
            raise PermissionDenied()
        config = app.release_set.latest().config
        serializer = self.get_serializer(config)
        return Response(serializer.data, status=status.HTTP_200_OK)


class AppPermsViewSet(BaseDeisViewSet):
    """RESTful views for sharing apps with collaborators."""

    model = models.App  # models class
    perm = 'use_app'    # short name for permission

    def get_queryset(self):
        return self.model.objects.all()

    def list(self, request, **kwargs):
        app = self.get_object()
        perm_name = "api.{}".format(self.perm)
        usernames = [u.username for u in get_users_with_perms(app)
                     if u.has_perm(perm_name, app)]
        return Response({'users': usernames})

    def create(self, request, **kwargs):
        app = self.get_object()
        if not permissions.IsOwnerOrAdmin.has_object_permission(permissions.IsOwnerOrAdmin(),
                                                                request, self, app):
            raise PermissionDenied()

        user = get_object_or_404(User, username=request.data['username'])
        assign_perm(self.perm, user, app)
        models.log_event(app, "User {} was granted access to {}".format(user, app))
        return Response(status=status.HTTP_201_CREATED)

    def destroy(self, request, **kwargs):
        app = get_object_or_404(models.App, id=self.kwargs['id'])
        user = get_object_or_404(User, username=kwargs['username'])

        perm_name = "api.{}".format(self.perm)
        if not user.has_perm(perm_name, app):
            raise PermissionDenied()

        if (user != request.user and
            not permissions.IsOwnerOrAdmin.has_object_permission(permissions.IsOwnerOrAdmin(),
                                                                 request, self, app)):
            raise PermissionDenied()
        remove_perm(self.perm, user, app)
        models.log_event(app, "User {} was revoked access to {}".format(user, app))
        return Response(status=status.HTTP_204_NO_CONTENT)


class AdminPermsViewSet(BaseDeisViewSet):
    """RESTful views for sharing admin permissions with other users."""

    model = User
    serializer_class = serializers.AdminUserSerializer
    permission_classes = [permissions.IsAdmin]

    def get_queryset(self, **kwargs):
        self.check_object_permissions(self.request, self.request.user)
        return self.model.objects.filter(is_active=True, is_superuser=True)

    def create(self, request, **kwargs):
        user = get_object_or_404(User, username=request.data['username'])
        user.is_superuser = user.is_staff = True
        user.save(update_fields=['is_superuser', 'is_staff'])
        return Response(status=status.HTTP_201_CREATED)

    def destroy(self, request, **kwargs):
        user = get_object_or_404(User, username=kwargs['username'])
        user.is_superuser = user.is_staff = False
        user.save(update_fields=['is_superuser', 'is_staff'])
        return Response(status=status.HTTP_204_NO_CONTENT)


class UserView(BaseDeisViewSet):
    """A Viewset for interacting with User objects."""
    model = User
    serializer_class = serializers.UserSerializer
    permission_classes = [permissions.IsAdmin]

    def get_queryset(self):
        return self.model.objects.exclude(username='AnonymousUser')
