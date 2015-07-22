#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the store mock component locally"
  echo
  exit 1
fi

# install required packages to build
apk add --update-cache \
  build-base \
  curl \
  file \
  gcc \
  git \
  python-dev

# install etcdctl
ETCD_VERSION=2.0.13
curl -sSL https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz | \
    tar zxv -C /usr/local/bin etcd-v${ETCD_VERSION}-linux-amd64/etcdctl && ln -s /usr/local/bin/etcd-v${ETCD_VERSION}-linux-amd64/etcdctl /usr/local/bin && chmod +x /usr/local/bin/etcdctl

git clone https://github.com/jserver/mock-s3 /app/mock-s3 --depth 1
cd /app/mock-s3
#FIXME: This is a gisted patch to enable pseudo-handling of POST requests, otherwise wal-e crashes attempting to delete old wal segments
curl https://gist.githubusercontent.com/anonymous/c565f11a8d90d6e2d92b/raw/c5815f6c83aa5c2cfb7b0a34cfab4a075c97be16/mock-s3-post.diff|git apply

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/7.0.3/contrib/get-pip.py | python -

python setup.py install

# cleanup.
apk del --purge \
  build-base \
  gcc \
  git
rm -rf /var/cache/apk/* /tmp/*
