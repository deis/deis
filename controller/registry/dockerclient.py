# -*- coding: utf-8 -*-
"""Support the Deis workflow by manipulating and publishing Docker images."""

from __future__ import unicode_literals
import io
import logging

from django.conf import settings
from rest_framework.exceptions import PermissionDenied
from simpleflock import SimpleFlock
import docker

logger = logging.getLogger(__name__)


class DockerClient(object):
    """Use the Docker API to pull, tag, build, and push images to deis-registry."""

    FLOCKFILE = '/tmp/controller-pull'

    def __init__(self):
        self.client = docker.Client(version='auto')
        self.registry = settings.REGISTRY_HOST + ':' + str(settings.REGISTRY_PORT)

    def publish_release(self, source, config, target, deis_registry):
        """Update a source Docker image with environment config and publish it to deis-registry."""
        # get the source repository name and tag
        src_name, src_tag = docker.utils.parse_repository_tag(source)
        # get the target repository name and tag
        name, tag = docker.utils.parse_repository_tag(target)
        # strip any "http://host.domain:port" prefix from the target repository name,
        # since we always publish to the Deis registry
        name = strip_prefix(name)

        # pull the source image from the registry
        # NOTE: this relies on an implementation detail of deis-builder, that
        # the image has been uploaded already to deis-registry
        if deis_registry:
            repo = "{}/{}".format(self.registry, src_name)
        else:
            repo = src_name
        self.pull(repo, src_tag)

        # tag the image locally without the repository URL
        image = "{}:{}".format(repo, src_tag)
        self.tag(image, src_name, tag=src_tag)

        # build a Docker image that adds a "last-mile" layer of environment
        config.update({'DEIS_APP': name, 'DEIS_RELEASE': tag})
        self.build(source, config, name, tag)

        # push the image to deis-registry
        self.push("{}/{}".format(self.registry, name), tag)

    def build(self, source, config, repo, tag):
        """Add a "last-mile" layer of environment config to a Docker image for deis-registry."""
        check_blacklist(repo)
        env = ' '.join('{}="{}"'.format(
            k, v.encode('unicode-escape').replace('"', '\\"')) for k, v in config.viewitems())
        dockerfile = "FROM {}\nENV {}".format(source, env)
        f = io.BytesIO(dockerfile.encode('utf-8'))
        target_repo = "{}/{}:{}".format(self.registry, repo, tag)
        logger.info("Building Docker image {}".format(target_repo))
        with SimpleFlock(self.FLOCKFILE, timeout=1200):
            stream = self.client.build(fileobj=f, tag=target_repo, stream=True, rm=True)
            log_output(stream)

    def pull(self, repo, tag):
        """Pull a Docker image into the local storage graph."""
        check_blacklist(repo)
        logger.info("Pulling Docker image {}:{}".format(repo, tag))
        with SimpleFlock(self.FLOCKFILE, timeout=1200):
            stream = self.client.pull(repo, tag=tag, stream=True, insecure_registry=True)
            log_output(stream)

    def push(self, repo, tag):
        """Push a local Docker image to a registry."""
        logger.info("Pushing Docker image {}:{}".format(repo, tag))
        stream = self.client.push(repo, tag=tag, stream=True, insecure_registry=True)
        log_output(stream)

    def tag(self, image, repo, tag):
        """Tag a local Docker image with a new name and tag."""
        check_blacklist(repo)
        logger.info("Tagging Docker image {} as {}:{}".format(image, repo, tag))
        if not self.client.tag(image, repo, tag=tag, force=True):
            raise docker.errors.DockerException("tagging failed")


def check_blacklist(repo):
    """Check a Docker repository name for collision with deis/* components."""
    blacklisted = [  # NOTE: keep this list up to date!
        'builder', 'cache', 'controller', 'database', 'logger', 'logspout',
        'publisher', 'registry', 'router', 'store-admin', 'store-daemon',
        'store-gateway', 'store-metadata', 'store-monitor',
    ]
    if any("deis/{}".format(c) in repo for c in blacklisted):
        raise PermissionDenied("Repository name {} is not allowed".format(repo))


def log_output(stream):
    """Log a stream at DEBUG level, and raise DockerException if it contains "error"."""
    for chunk in stream:
        logger.debug(chunk)
        # error handling requires looking at the response body
        if '"error"' in chunk.lower():
            raise docker.errors.DockerException(chunk)


def strip_prefix(name):
    """Strip the schema and host:port from a Docker repository name."""
    paths = name.split('/')
    return '/'.join(p for p in paths if p and '.' not in p and ':' not in p)


def publish_release(source, config, target, deis_registry):

    client = DockerClient()
    return client.publish_release(source, config, target, deis_registry)
