import cStringIO
import hashlib
import json
import requests
import tarfile
import urlparse
import uuid

from django.conf import settings
from docker.utils import utils

from api.utils import encode


def publish_release(source, config, target):
    """
    Publish a new release as a Docker image

    Given a source image and dictionary of last-mile configuration,
    create a target Docker image on the registry.

    For example::

        publish_release('registry.local:5000/gabrtv/myapp:v22',
                        {'ENVVAR': 'values'},
                        'registry.local:5000/gabrtv/myapp:v23')

    results in a new Docker image at 'registry.local:5000/gabrtv/myapp:v23' which
    contains the new configuration as ENV entries.
    """
    try:
        repo, tag = utils.parse_repository_tag(source)
        src_image = repo
        src_tag = tag if tag is not None else 'latest'

        nameparts = repo.rsplit('/', 1)
        if len(nameparts) == 2:
            if '/' in nameparts[0]:
                # strip the hostname and just use the app name
                src_image = '{}/{}'.format(nameparts[0].rsplit('/', 1)[1],
                                           nameparts[1])
            elif '.' in nameparts[0]:
                # we got a name like registry.local:5000/registry
                src_image = nameparts[1]

        target_image = target.rsplit(':', 1)[0]
        target_tag = target.rsplit(':', 1)[1]
        image_id = _get_tag(src_image, src_tag)
    except RuntimeError:
        if src_tag == 'latest':
            # no image exists yet, so let's build one!
            _put_first_image(src_image)
            image_id = _get_tag(src_image, src_tag)
        else:
            raise
    image = _get_image(image_id)
    # construct the new image
    image['parent'] = image['id']
    image['id'] = _new_id()
    config['DEIS_APP'] = target_image
    config['DEIS_RELEASE'] = target_tag
    image['config']['Env'] = _construct_env(image['config']['Env'], config)
    # update and tag the new image
    _commit(target_image, image, _empty_tar_archive(), target_tag)


# registry access


def _commit(repository_path, image, layer, tag):
    _put_image(image)
    cookies = _put_layer(image['id'], layer)
    _put_checksum(image, layer, cookies)
    _put_tag(image['id'], repository_path, tag)


def _put_first_image(repository_path):
    image = {
        'id': _new_id(),
        'parent': '',
        'config': {
            'Env': []
        }
    }
    # tag as v0 in the registry
    _commit(repository_path, image, _empty_tar_archive(), 'v0')


def _api_call(endpoint, data=None, headers={}, cookies=None, request_type='GET'):
    base_headers = {'user-agent': 'docker/1.0.0'}
    r = None
    if len(headers) > 0:
        for header, value in headers.viewitems():
            base_headers[header] = value
    if request_type == 'GET':
        r = requests.get(endpoint, headers=base_headers)
    elif request_type == 'PUT':
        r = requests.put(endpoint, data=data, headers=base_headers, cookies=cookies)
    else:
        raise AttributeError("request type not supported: {}".format(request_type))
    return r


def _get_tag(repository, tag):
    path = "/v1/repositories/{repository}/tags/{tag}".format(**locals())
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    r = _api_call(url)
    if not r.status_code == 200:
        raise RuntimeError("GET Image Error ({}: {})".format(r.status_code, r.text))
    return r.json()


def _get_image(image_id):
    path = "/v1/images/{image_id}/json".format(**locals())
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    r = _api_call(url)
    if not r.status_code == 200:
        raise RuntimeError("GET Image Error ({}: {})".format(r.status_code, r.text))
    return r.json()


def _put_image(image):
    path = "/v1/images/{id}/json".format(**image)
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    r = _api_call(url, data=json.dumps(image), request_type='PUT')
    if not r.status_code == 200:
        raise RuntimeError("PUT Image Error ({}: {})".format(r.status_code, r.text))
    return r.json()


def _put_layer(image_id, layer_fileobj):
    path = "/v1/images/{image_id}/layer".format(**locals())
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    r = _api_call(url, data=layer_fileobj.read(), request_type='PUT')
    if not r.status_code == 200:
        raise RuntimeError("PUT Layer Error ({}: {})".format(r.status_code, r.text))
    return r.cookies


def _put_checksum(image, layer, cookies):
    path = "/v1/images/{id}/checksum".format(**image)
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    h = hashlib.sha256(json.dumps(image) + '\n')
    h.update(layer.getvalue())
    layer_checksum = "sha256:{0}".format(h.hexdigest())
    headers = {'X-Docker-Checksum-Payload': layer_checksum}
    r = _api_call(url, headers=headers, cookies=cookies, request_type='PUT')
    if not r.status_code == 200:
        raise RuntimeError("PUT Checksum Error ({}: {})".format(r.status_code, r.text))


def _put_tag(image_id, repository_path, tag):
    path = "/v1/repositories/{repository_path}/tags/{tag}".format(**locals())
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    r = _api_call(url, data=json.dumps(image_id), request_type='PUT')
    if not r.status_code == 200:
        raise RuntimeError("PUT Tag Error ({}: {})".format(r.status_code, r.text))


# utility functions

def _construct_env(env, config):
    "Update current environment with latest config"
    new_env = []
    # see if we need to update existing ENV vars
    for e in env:
        k, v = e.split('=', 1)
        if k in config:
            # update values defined by config
            v = config.pop(k)
        new_env.append("{}={}".format(encode(k), encode(v)))
    # add other config ENV items
    for k, v in config.viewitems():
        new_env.append("{}={}".format(encode(k), encode(v)))
    return new_env


def _new_id():
    "Return 64-char UUID for use as Image ID"
    return ''.join(uuid.uuid4().hex * 2)


def _empty_tar_archive():
    "Return an empty tar archive (in memory)"
    data = cStringIO.StringIO()
    tar = tarfile.open(mode="w", fileobj=data)
    tar.close()
    data.seek(0)
    return data
