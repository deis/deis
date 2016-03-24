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
apk add --no-cache \
  build-base \
  curl \
  file \
  gcc \
  git \
  python-dev

# install etcdctl
curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/get-deis/etcdctl-v0.4.9 \
  && chmod +x /usr/local/bin/etcdctl

git clone https://github.com/jserver/mock-s3 /app/mock-s3
cd /app/mock-s3
#FIXME: This is a gisted patch to a known "good" version of mock-s3 to enable pseudo-handling of POST requests, otherwise wal-e crashes attempting to delete old wal segments
git checkout 4c3c3752f990db97e8969c00666251a3b427ef4c
git apply /tmp/mock-s3-patch.diff

# install pip
curl -sSL https://bootstrap.pypa.io/get-pip.py | python - pip==8.1.1

python setup.py install

# cleanup.
apk del --no-cache \
  build-base \
  gcc \
  git
rm -rf /tmp/*
