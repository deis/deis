#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the registry locally"
  echo
  exit 1
fi

# install required packages (copied from dotcloud/docker-registry Dockerfile)
apk add --no-cache \
  build-base \
  git \
  openssl-dev \
  python-dev \
  libffi-dev \
  swig \
  libevent-dev \
  xz-dev

# install pip
curl -sSL https://bootstrap.pypa.io/get-pip.py | python - pip==8.1.1

# create a registry user
adduser -D -s /bin/bash registry

# add the docker registry source from github
git clone -b deis-v1-lts --single-branch https://github.com/deis/docker-registry /docker-registry && \
  chown -R registry:registry /docker-registry

# install boto configuration
cp /docker-registry/config/boto.cfg /etc/boto.cfg
cd /docker-registry && pip install --disable-pip-version-check --no-cache-dir -r requirements/main.txt

# Install core
pip install --disable-pip-version-check --no-cache-dir /docker-registry/depends/docker-registry-core

# Install registry
pip install --disable-pip-version-check --no-cache-dir "file:///docker-registry#egg=docker-registry[bugsnag,newrelic,cors]"

# patch boto
cd "$(python -c 'import boto; import os; print os.path.dirname(boto.__file__)')" \
  && patch -i /docker-registry/contrib/boto_header_patch.diff connection.py

# cleanup. indicate that python is a required package.
apk del --no-cache \
  build-base \
  linux-headers \
  python-dev
