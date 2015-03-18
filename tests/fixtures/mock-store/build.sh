#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the store mock component locally"
  echo
  exit 1
fi

apk add --update curl \
  file \
  gcc \
  git \
  build-base \
  python-dev

# install etcdctl
curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
  && chmod +x /usr/local/bin/etcdctl

git clone https://github.com/jserver/mock-s3 /app/mock-s3

cd /app/mock-s3

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.0.8/contrib/get-pip.py | python -

python setup.py install

# cleanup.
apk del gcc \
  git \
  build-base

rm -rf /root/.cache \
  /usr/share/doc \
  /tmp/* \
  /var/cache/apk/*