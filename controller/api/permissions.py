from rest_framework import permissions
from django.conf import settings
from django.contrib.auth.models import AnonymousUser

from api import models


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
