import cStringIO
import hashlib
import json
import requests
import tarfile
import urlparse
import uuid

from django.conf import settings


def publish_release(repository_path, config, tag, source_tag='latest'):
    """
    Publish a new release as a Docker image

    Given a source repository path, a dictionary of environment variables
    and a target tag, create a new lightweight Docker image on the registry.

    source_tag is the name of the previous older tag that this image should
    be a child of. In most cases, this should be 'latest', but for rollbacks
    this should be an older tag.

    For example, publish_release('gabrtv/myapp', {'ENVVAR': 'values'}, 'v23')
    results in a new Docker image at: <registry_url>/gabrtv/myapp:v23
    which contains the new configuration as ENV entries.
    """
    try:
        image_id = _get_tag(repository_path, source_tag)
    except RuntimeError:
        if source_tag == 'latest':
            # no image exists yet, so let's build one!
            _put_first_image(repository_path)
            image_id = _get_tag(repository_path, 'latest')
        else:
            raise
    image = _get_image(image_id)
    # construct the new image
    image['parent'] = image['id']
    image['id'] = _new_id()
    image['config']['Env'] = _construct_env(image['config']['Env'], config)
    # update and tag the new image
    _commit(repository_path, image, _empty_tar_archive(), tag)


# registry access


def _commit(repository_path, image, layer, tag):
    _put_image(image)
    cookies = _put_layer(image['id'], layer)
    _put_checksum(image, cookies)
    _put_tag(image['id'], repository_path, tag)
    # point latest to the new tag
    _put_tag(image['id'], repository_path, 'latest')


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
    # FIXME: update API calls for docker 0.10.0+
    base_headers = {'user-agent': 'docker/0.9.0'}
    r = None
    if len(headers) > 0:
        for header, value in headers.iteritems():
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
    print r.text
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


def _put_checksum(image, cookies):
    path = "/v1/images/{id}/checksum".format(**image)
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    tarsum = TarSum(json.dumps(image)).compute()
    headers = {'X-Docker-Checksum': tarsum}
    r = _api_call(url, headers=headers, cookies=cookies, request_type='PUT')
    if not r.status_code == 200:
        raise RuntimeError("PUT Checksum Error ({}: {})".format(r.status_code, r.text))
    print r.json()


def _put_tag(image_id, repository_path, tag):
    path = "/v1/repositories/{repository_path}/tags/{tag}".format(**locals())
    url = urlparse.urljoin(settings.REGISTRY_URL, path)
    r = _api_call(url, data=json.dumps(image_id), request_type='PUT')
    if not r.status_code == 200:
        raise RuntimeError("PUT Tag Error ({}: {})".format(r.status_code, r.text))
    print r.json()


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
        new_env.append("{}={}".format(k, v))
    # add other config ENV items
    for k, v in config.items():
        new_env.append("{}={}".format(k, v))
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


#
# Below adapted from https://github.com/dotcloud/docker-registry/blob/master/lib/checksums.py
#

def sha256_file(fp, data=None):
    h = hashlib.sha256(data or '')
    if not fp:
        return h.hexdigest()
    while True:
        buf = fp.read(4096)
        if not buf:
            break
        h.update(buf)
    return h.hexdigest()


def sha256_string(s):
    return hashlib.sha256(s).hexdigest()


class TarSum(object):

    def __init__(self, json_data):
        self.json_data = json_data
        self.hashes = []
        self.header_fields = ('name', 'mode', 'uid', 'gid', 'size', 'mtime',
                              'type', 'linkname', 'uname', 'gname', 'devmajor',
                              'devminor')

    def append(self, member, tarobj):
        header = ''
        for field in self.header_fields:
            value = getattr(member, field)
            if field == 'type':
                field = 'typeflag'
            elif field == 'name':
                if member.isdir() and not value.endswith('/'):
                    value += '/'
            header += '{0}{1}'.format(field, value)
        h = None
        try:
            if member.size > 0:
                f = tarobj.extractfile(member)
                h = sha256_file(f, header)
            else:
                h = sha256_string(header)
        except KeyError:
            h = sha256_string(header)
        self.hashes.append(h)

    def compute(self):
        self.hashes.sort()
        data = self.json_data + ''.join(self.hashes)
        tarsum = 'tarsum+sha256:{0}'.format(sha256_string(data))
        return tarsum
