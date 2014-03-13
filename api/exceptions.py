"""
Deis API exception classes.
"""

from __future__ import unicode_literals

from rest_framework.exceptions import APIException
from rest_framework import status


class BuildNodeError(APIException):
    """
    Indicates a problem in building or bootstrapping a node.

    This exception is subclassed from rest_framework's APIException so it
    isn't reported as "500 SERVER ERROR."
    """

    status_code = status.HTTP_401_UNAUTHORIZED

    def __init__(self, detail=None):
        self.detail = detail


class BuildFormationError(APIException):
    """
    Indicates a problem in creating a formation.

    This exception is subclassed from rest_framework's APIException so it
    isn't reported as "500 SERVER ERROR."
    """

    status_code = status.HTTP_400_BAD_REQUEST

    def __init__(self, detail=None):
        self.detail = detail
