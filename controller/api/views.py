"""
RESTful view classes for presenting Deis API objects.
"""

from django.conf import settings
from django.core.exceptions import ValidationError
from django.contrib.auth.models import User
from django.shortcuts import get_object_or_404
from guardian.shortcuts import assign_perm, get_objects_for_user, \
    get_users_with_perms, remove_perm
from rest_framework import mixins, renderers, status
from rest_framework.exceptions import PermissionDenied
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.viewsets import GenericViewSet

from api import authentication, models, permissions, serializers, viewsets


class UserRegistrationViewSet(GenericViewSet,
                              mixins.CreateModelMixin):
    """ViewSet to handle registering new users. The logic is in the serializer."""
    authentication_classes = [authentication.AnonymousAuthentication]
    permission_classes = [permissions.IsAnonymous, permissions.HasRegistrationAuth]
    serializer_class = serializers.UserSerializer


class UserManagementViewSet(GenericViewSet,
                            mixins.DestroyModelMixin):
    serializer_class = serializers.UserSerializer

    def get_queryset(self):
        return User.objects.filter(pk=self.request.user.pk)

    def get_object(self):
        return self.get_queryset()[0]

    def passwd(self, request, **kwargs):
        obj = self.get_object()
        if not obj.check_password(request.DATA['password']):
            return Response({'detail': 'Current password does not match'},
                            status=status.HTTP_400_BAD_REQUEST)
        obj.set_password(request.DATA['new_password'])
        obj.save()
        return Response({'status': 'password set'})


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
        request.DATA['app'] = self.get_app()
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
        headers.update({'X-Deis-Release': self.release.version})
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
            for target, count in request.DATA.items():
                new_structure[target] = int(count)
            models.validate_app_structure(new_structure)
            app.scale(request.user, new_structure)
        except (TypeError, ValueError):
            return Response({'detail': 'Invalid scaling format'},
                            status=status.HTTP_400_BAD_REQUEST)
        except (ValidationError, EnvironmentError) as e:
            return Response({'detail': str(e)}, status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError as e:
            return Response({'detail': str(e)}, status=status.HTTP_503_SERVICE_UNAVAILABLE)
        return Response(status=status.HTTP_204_NO_CONTENT)

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
            return Response({'detail': str(e)}, status=status.HTTP_400_BAD_REQUEST)
        except RuntimeError as e:
            return Response({'detail': str(e)}, status=status.HTTP_503_SERVICE_UNAVAILABLE)
        return Response(output_and_rc, status=status.HTTP_200_OK,
                        content_type='text/plain')


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


class DomainViewSet(AppResourceViewSet):
    """A viewset for interacting with Domain objects."""
    model = models.Domain
    serializer_class = serializers.DomainSerializer


class KeyViewSet(BaseDeisViewSet):
    """A viewset for interacting with Key objects."""
    model = models.Key
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
        try:
            app = self.get_app()
            release = app.release_set.latest()
            version_to_rollback_to = release.version - 1
            if request.DATA.get('version'):
                version_to_rollback_to = int(request.DATA['version'])
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
        self.user = get_object_or_404(User, username=request.data['receive_user'])
        # check the user is authorized for this app
        if self.user == app.owner or \
           self.user in get_users_with_perms(app) or \
           self.user.is_superuser:
                request.data['app'] = app
                request.data['owner'] = self.user
                return super(PushHookViewSet, self).create(request, *args, **kwargs)
        raise PermissionDenied()

    def perform_create(self, serializer, **kwargs):
        serializer.save(owner=self.user)


class BuildHookViewSet(BaseHookViewSet):
    """API hook to create new :class:`~api.models.Build`"""
    model = models.Build
    serializer_class = serializers.BuildSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.data['receive_repo'])
        self.user = get_object_or_404(User, username=request.data['receive_user'])
        # check the user is authorized for this app
        if self.user == app.owner or \
           self.user in get_users_with_perms(app) or \
           self.user.is_superuser:
            request._data = request.data.copy()
            request.data['app'] = app
            request.data['owner'] = self.user
            super(BuildHookViewSet, self).create(request, *args, **kwargs)
            # return the application databag
            response = {'release': {'version': app.release_set.latest().version},
                        'domains': ['.'.join([app.id, settings.DEIS_DOMAIN])]}
            return Response(response, status=status.HTTP_200_OK)
        raise PermissionDenied()

    def perform_create(self, serializer, **kwargs):
        build = serializer.save(owner=self.user)
        self.post_save(build)

    def post_save(self, build):
        build.create(self.user)


class ConfigHookViewSet(BaseHookViewSet):
    """API hook to grab latest :class:`~api.models.Config`"""
    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def create(self, request, *args, **kwargs):
        app = get_object_or_404(models.App, id=request.DATA['receive_repo'])
        user = get_object_or_404(User, username=request.DATA['receive_user'])
        # check the user is authorized for this app
        if user == app.owner or \
           user in get_users_with_perms(app) or \
           user.is_superuser:
            config = app.release_set.latest().config
            serializer = self.get_serializer(config)
            return Response(serializer.data, status=status.HTTP_200_OK)
        raise PermissionDenied()


class AppPermsViewSet(BaseDeisViewSet):
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
            return Response(status=status.HTTP_403_FORBIDDEN)


class AdminPermsViewSet(BaseDeisViewSet):
    """RESTful views for sharing admin permissions with other users."""

    model = User
    serializer_class = serializers.AdminUserSerializer
    permission_classes = (permissions.IsAdmin,)

    def get_queryset(self, **kwargs):
        self.check_object_permissions(self.request, self.request.user)
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
