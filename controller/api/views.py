"""
RESTful view classes for presenting Deis API objects.
"""
# pylint: disable=R0901,R0904

from __future__ import unicode_literals
from Crypto.PublicKey import RSA
from api import models, serializers
from django.conf import settings
from django.contrib.auth.models import Group, AnonymousUser, User
from django.db.utils import IntegrityError
from django.utils import timezone
from rest_framework import permissions, status, viewsets
from rest_framework.authentication import BaseAuthentication
from rest_framework.response import Response
from rest_framework.status import HTTP_400_BAD_REQUEST, HTTP_201_CREATED
import json
import os.path
import yaml


class AnonymousAuthentication(BaseAuthentication):

    def authenticate(self, request):
        """
        Authenticate the request and return a two-tuple of (user, token).
        """
        user = AnonymousUser()
        return user, None


class IsAnonymous(permissions.BasePermission):
    """
    Object-level permission to allow anonymous users
    """

    def has_permission(self, request, view):
        """
        Return `True` if permission is granted, `False` otherwise.
        """
        if type(request.user) == AnonymousUser:
            return True
        return False


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


class GroupViewSet(viewsets.ModelViewSet):

    model = Group


class UserViewSet(viewsets.ModelViewSet):

    model = settings.AUTH_USER_MODEL


class UserRegistrationView(viewsets.GenericViewSet,
                           viewsets.mixins.CreateModelMixin):

    model = User

    authentication_classes = (AnonymousAuthentication,)
    permission_classes = (IsAnonymous,)
    serializer_class = serializers.UserSerializer

    def pre_save(self, obj):
        "Replicate UserManager.create_user functionality"
        now = timezone.now()
        obj.last_login = now
        obj.date_joined = now
        obj.email = User.objects.normalize_email(obj.email)
        obj.set_password(obj.password)


class OwnerViewSet(viewsets.ModelViewSet):
    """
    Base ViewSet for views scoped to a particular Owner
    """
    permission_classes = (permissions.IsAuthenticated, IsOwner)

    def pre_save(self, obj):
        obj.owner = self.request.user

    def get_queryset(self, **kwargs):
        return self.model.objects.filter(owner=self.request.user)


class KeyViewSet(OwnerViewSet):

    model = models.Key
    serializer_class = serializers.KeySerializer
    lookup_field = 'id'

    def post_save(self, obj, created=False, **kwargs):
        # update gitosis
        models.Formation.objects.publish()

    def destroy(self, request, **kwargs):
        resp = super(KeyViewSet, self).destroy(self, request, **kwargs)
        # publish gitosis updates
        models.Formation.objects.publish()
        return resp


class ProviderViewSet(OwnerViewSet):

    model = models.Provider
    serializer_class = serializers.ProviderSerializer
    lookup_field = 'id'


def _load_cloud_config_base():
    # load cloud-config-base yaml_
    _cloud_config_path = os.path.abspath(
            os.path.join(__file__, '..', 'files', 'cloud-config-base.yml'))
    with open(_cloud_config_path) as f:
        _data = f.read()
    return yaml.safe_load(_data)

class FlavorViewSet(OwnerViewSet):

    model = models.Flavor
    serializer_class = serializers.FlavorSerializer
    lookup_field = 'id'

    def create(self, request, **kwargs):
        request._data = request.DATA.copy()
        if not 'ssh_private_key' in request.DATA and not 'ssh_public_key' in request.DATA:
            # SECURITY: figure out best way to get keys with proper entropy
            key = RSA.generate(2048)
            request.DATA['ssh_private_key'] = key.exportKey('PEM')
            request.DATA['ssh_public_key'] = key.exportKey('OpenSSH')
        # set default cloud-init configuration
        if not 'init' in request.DATA:
            request.DATA['init'] = _load_cloud_config_base()
        return viewsets.ModelViewSet.create(self, request, **kwargs)


class FormationViewSet(OwnerViewSet):

    model = models.Formation
    serializer_class = serializers.FormationSerializer
    lookup_field = 'id'

    def create(self, request, **kwargs):
        try:
            return OwnerViewSet.create(self, request, **kwargs)
        except IntegrityError as _e:
            return Response("Formation with this Id already exists.",
                            status=HTTP_400_BAD_REQUEST)
    
    def post_save(self, formation, created=False, **kwargs):
        if created:
            config = models.Config.objects.create(
                owner=formation.owner, formation=formation, values={})
            build = models.Build.objects.create(
                owner=formation.owner, formation=formation)
            _release = models.Release.objects.create(
                owner=formation.owner, formation=formation, config=config,
                image=formation.image, build=build)
            # prepare the formation's required infrastruture at the provider
            formation.prepare_provider().delay().wait()
        # update gitosis
        models.Formation.objects.publish()

    def scale(self, request, **kwargs):
        new_structure = {}
        try:
            for target, count in request.DATA.items():
                new_structure[target] = int(count)
        except ValueError:
            return Response('Invalid scaling format', status=HTTP_400_BAD_REQUEST)
        formation = self.get_object()
        formation.structure.update(new_structure)
        formation.save()
        try:
            databag = formation.scale()
        except models.ScalingError as e:
            return Response(str(e), status=status.HTTP_400_BAD_REQUEST)
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def balance(self, request, **kwargs):
        formation = self.get_object()
        databag = formation.balance()
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')
    
    def calculate(self, request, **kwargs):
        formation = self.get_object()
        databag = formation.calculate()
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def converge(self, request, **kwargs):
        formation = self.get_object()
        databag = formation.converge(formation.calculate())
        return Response(databag, status=status.HTTP_200_OK,
                        content_type='application/json')

    def destroy(self, request, **kwargs):
        formation = self.get_object()
        formation.destroy()
        formation.delete()
        # publish gitosis updates
        models.Formation.objects.publish()
        return Response(status=status.HTTP_204_NO_CONTENT)


class FormationNodeViewSet(OwnerViewSet):

    model = models.Node
    serializer_class = serializers.NodeSerializer
    
    def get_queryset(self, **kwargs):
        formation = models.Formation.objects.get(
                owner=self.request.user, id=self.kwargs['id'])
        return self.model.objects.filter(owner=self.request.user, formation=formation)
            
    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = qs.get(id=self.kwargs['id'])
        return obj


class FormationBackendViewSet(OwnerViewSet):

    model = models.Backend
    serializer_class = serializers.BackendSerializer
    
    def get_queryset(self, **kwargs):
        formation = models.Formation.objects.get(
                owner=self.request.user, id=self.kwargs['id'])
        return self.model.objects.filter(owner=self.request.user, formation=formation)

    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = qs.get(pk=self.kwargs['id'])
        return obj


class FormationProxyViewSet(OwnerViewSet):

    model = models.Proxy
    serializer_class = serializers.ProxySerializer
    
    def get_queryset(self, **kwargs):
        formation = models.Formation.objects.get(
                owner=self.request.user, id=self.kwargs['id'])
        return self.model.objects.filter(owner=self.request.user, formation=formation)

    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = qs.get(pk=self.kwargs['id'])
        return obj


class FormationContainerViewSet(OwnerViewSet):

    model = models.Container
    serializer_class = serializers.ContainerSerializer
    
    def get_queryset(self, **kwargs):
        formation = models.Formation.objects.get(
                owner=self.request.user, id=self.kwargs['id'])
        return self.model.objects.filter(owner=self.request.user, formation=formation)
    
    def get_object(self, *args, **kwargs):
        qs = self.get_queryset(**kwargs)
        obj = qs.get(pk=self.kwargs['id'])
        return obj


class FormationImageViewSet(OwnerViewSet):

    model = models.Release
    serializer_class = serializers.ReleaseSerializer

    def get_queryset(self, **kwargs):
        return self.model.objects.filter(owner=self.request.user)
            
    def get_object(self, *args, **kwargs):        
        formation = models.Formation.objects.get(
                owner=self.request.user, id=self.kwargs['id'])
        return formation.release

    def reset_image(self, request, *args, **kwargs):
        formation = models.Formation.objects.get(
                owner=self.request.user, id=self.kwargs['id'])
        models.release_signal.send(sender=self, image=request.DATA['image'], 
                                   formation=formation, user=self.request.user)
        return Response(status=HTTP_201_CREATED)


class FormationConfigViewSet(OwnerViewSet):

    model = models.Config
    serializer_class = serializers.ConfigSerializer

    def get_queryset(self, **kwargs):
        return self.model.objects.filter(owner=self.request.user)
    
    def get_object(self, *args, **kwargs):
        formation = models.Formation.objects.get(id=self.kwargs['id'])
        config = self.model.objects.filter(
                formation=formation).order_by('-created')[0]
        return config

    def post_save(self, obj, created=False):
        if created:
            models.release_signal.send(sender=self,
                config=obj, formation=obj.formation, user=self.request.user)
            # recalculate and converge after each config update
            databag = obj.formation.calculate()
            obj.formation.converge(databag)

    def create(self, request, *args, **kwargs):
        request._data = request.DATA.copy()
        # assume an existing config object exists
        obj = self.get_object()
        # increment version and use the same formation
        request.DATA['version'] = obj.version + 1
        request.DATA['formation'] = obj.formation
        # merge config values
        values = obj.values.copy()
        provided = json.loads(request.DATA['values'])
        values.update(provided)
        # remove config keys if we provided a null value
        [ values.pop(k) for k, v in provided.items() if v is None ]
        request.DATA['values'] = values
        return super(OwnerViewSet, self).create(request, *args, **kwargs)


class FormationBuildViewSet(OwnerViewSet):

    model = models.Build
    serializer_class = serializers.BuildSerializer

    def get_queryset(self, **kwargs):
        return self.model.objects.filter(owner=self.request.user)
    
    def get_object(self, *args, **kwargs):
        formation = models.Formation.objects.get(id=self.kwargs['id'])
        build = self.model.objects.filter(
                formation=formation).order_by('-created')[0]
        return build

    def post_save(self, obj, created=False):
        if created:
            models.release_signal.send(sender=self,
                build=obj, formation=obj.formation, user=self.request.user)

    def create(self, request, *args, **kwargs):
        request._data = request.DATA.copy()
        obj = self.get_object()
        # increment version and use the same formation
        request.DATA['version'] = obj.version + 1
        request.DATA['formation'] = obj.formation
        return super(OwnerViewSet, self).create(request, *args, **kwargs)


class FormationReleaseViewSet(OwnerViewSet):

    model = models.Release
    serializer_class = serializers.ReleaseSerializer

    def get_queryset(self, **kwargs):
        return self.model.objects.filter(owner=self.request.user)
            
    def get_object(self, *args, **kwargs):
        formation = models.Formation.objects.get(id=self.kwargs['id'])
        release = self.model.objects.filter(
                formation=formation).order_by('-created')[0]
        return release
