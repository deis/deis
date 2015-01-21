from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated

from api import permissions


class OwnerViewSet(viewsets.ModelViewSet):
    """
    A simple ViewSet for objects filtered by their 'owner' attribute.

    To use it, at minimum you'll need to provide the `serializer_class` attribute and
    the `model` attribute shortcut.
    """
    permission_classes = [IsAuthenticated, permissions.IsOwner]

    def get_queryset(self):
        return self.model.objects.filter(owner=self.request.user)

    def perform_create(self, serializer):
        obj = serializer.save(owner=self.request.user)
        self.post_save(obj)

    def post_save(self, obj):
        """A post_save hook for performing actions after the object has been pushed to the
        database.

        Leave it up to child classes to implement."""
        pass
