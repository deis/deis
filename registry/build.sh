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
apk add --update-cache \
  build-base \
  git \
  openssl-dev \
  python-dev \
  libffi-dev \
  swig \
  libevent-dev \
  xz-dev

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.1.1/contrib/get-pip.py | python -

# workaround to python > 2.7.8 SSL issues
pip install pyopenssl ndg-httpsclient pyasn1

# create a registry user
adduser -D -s /bin/bash registry

# add the docker registry source from github
git clone -b new-repository-import-v091 --single-branch https://github.com/deis/docker-registry /docker-registry && \
  chown -R registry:registry /docker-registry

# install boto configuration
cp /docker-registry/config/boto.cfg /etc/boto.cfg
cd /docker-registry && pip install -r requirements/main.txt

# Install core
pip install /docker-registry/depends/docker-registry-core

# Install registry
pip install file:///docker-registry#egg=docker-registry[bugsnag,newrelic,cors]

patch \
  $(python -c 'import boto; import os; print os.path.dirname(boto.__file__)')/connection.py \
  < /docker-registry/contrib/boto_header_patch.diff

# cleanup. indicate that python is a required package.
apk del --purge \
  build-base \
  linux-headers \
  python-dev

rm -rf /var/cache/apk/*
