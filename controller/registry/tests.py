"""
Unit tests for the Deis registry app.

Run the tests with "./manage.py test registry"
"""

import unittest
try:
    from unittest import mock
except ImportError:
    import mock

from django.conf import settings
from rest_framework.exceptions import PermissionDenied
from registry.dockerclient import DockerClient
from registry.dockerclient import strip_prefix


@mock.patch('docker.Client')
class DockerClientTest(unittest.TestCase):
    """Test that the client makes appropriate Docker engine API calls."""

    def setUp(self):
        settings.REGISTRY_HOST, settings.REGISTRY_PORT = 'localhost', 5000

    def test_publish_release(self, mock_client):
        self.client = DockerClient()
        self.client.publish_release('ozzy/embryo:git-f2a8020',
                                    {'POWERED_BY': 'Deis'}, 'ozzy/embryo:v4', True)
        self.assertTrue(self.client.client.pull.called)
        self.assertTrue(self.client.client.tag.called)
        self.assertTrue(self.client.client.build.called)
        self.assertTrue(self.client.client.push.called)
        # Test that a registry host prefix is replaced with deis-registry for the target
        self.client.publish_release('ozzy/embryo:git-f2a8020',
                                    {'POWERED_BY': 'Deis'}, 'quay.io/ozzy/embryo:v4', True)
        docker_push = self.client.client.push
        docker_push.assert_called_with(
            'localhost:5000/ozzy/embryo', tag='v4', insecure_registry=True, stream=True)
        # Test that blacklisted image names can't be published
        with self.assertRaises(PermissionDenied):
            self.client.publish_release(
                'deis/controller:v1.11.1', {}, 'deis/controller:v1.11.1', True)
        with self.assertRaises(PermissionDenied):
            self.client.publish_release(
                'localhost:5000/deis/controller:v1.11.1', {}, 'deis/controller:v1.11.1', True)

    def test_build(self, mock_client):
        # test that self.client.build was called with proper arguments
        self.client = DockerClient()
        self.client.build('ozzy/embryo:git-f3a8020', {'POWERED_BY': 'Deis'}, 'ozzy/embryo', 'v4')
        docker_build = self.client.client.build
        self.assertTrue(docker_build.called)
        args = {"rm": True, "tag": u'localhost:5000/ozzy/embryo:v4', "stream": True}
        kwargs = docker_build.call_args[1]
        self.assertDictContainsSubset(args, kwargs)
        # test that the fileobj arg to "docker build" contains a correct Dockerfile
        f = kwargs['fileobj']
        self.assertEqual(f.read(), "FROM ozzy/embryo:git-f3a8020\nENV POWERED_BY=\"Deis\"")
        # Test that blacklisted image names can't be built
        with self.assertRaises(PermissionDenied):
            self.client.build('deis/controller:v1.11.1', {}, 'deis/controller', 'v1.11.1')
        with self.assertRaises(PermissionDenied):
            self.client.build(
                'localhost:5000/deis/controller:v1.11.1', {}, 'deis/controller', 'v1.11.1')

    def test_pull(self, mock_client):
        self.client = DockerClient()
        self.client.pull('alpine', '3.2')
        docker_pull = self.client.client.pull
        docker_pull.assert_called_once_with(
            'alpine', tag='3.2', insecure_registry=True, stream=True)
        # Test that blacklisted image names can't be pulled
        with self.assertRaises(PermissionDenied):
            self.client.pull('deis/controller', 'v1.11.1')
        with self.assertRaises(PermissionDenied):
            self.client.pull('localhost:5000/deis/controller', 'v1.11.1')

    def test_push(self, mock_client):
        self.client = DockerClient()
        self.client.push('ozzy/embryo', 'v4')
        docker_push = self.client.client.push
        docker_push.assert_called_once_with(
            'ozzy/embryo', tag='v4', insecure_registry=True, stream=True)

    def test_tag(self, mock_client):
        self.client = DockerClient()
        self.client.tag('ozzy/embryo:git-f2a8020', 'ozzy/embryo', 'v4')
        docker_tag = self.client.client.tag
        docker_tag.assert_called_once_with(
            'ozzy/embryo:git-f2a8020', 'ozzy/embryo', tag='v4', force=True)
        # Test that blacklisted image names can't be tagged
        with self.assertRaises(PermissionDenied):
            self.client.tag('deis/controller:v1.11.1', 'deis/controller', 'v1.11.1')
        with self.assertRaises(PermissionDenied):
            self.client.tag('localhost:5000/deis/controller:v1.11.1', 'deis/controller', 'v1.11.1')

    def test_strip_prefix(self, mock_client):
        self.assertEqual(strip_prefix('quay.io/boris/riotsugar'), 'boris/riotsugar')
        self.assertEqual(strip_prefix('127.0.0.1:5000/boris/galaxians'), 'boris/galaxians')
        self.assertEqual(strip_prefix('boris/jacksonhead'), 'boris/jacksonhead')
        self.assertEqual(strip_prefix(':8888/boris/pink'), 'boris/pink')
