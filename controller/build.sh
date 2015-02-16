#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the controller locally"
  echo
  exit 1
fi

# install required system packages
# HACK: install git so we can install bacongobbler's fork of django-fsm
# install openssh-client for temporary fleetctl wrapper
echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories
apk-install libffi-dev \
    libpq \
    git \
    python \
    python-dev \
    postgresql-dev \
    build-base

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.0.8/contrib/get-pip.py | python -

# add a deis user
adduser deis -D -h /app -s /bin/bash

# create a /app directory for storing application data
mkdir -p /app && chown -R deis:deis /app

# create directory for confd templates
mkdir -p /templates && chown -R deis:deis /templates

# install dependencies
pip install -r /app/requirements.txt

apk del libffi-dev \
        python-dev \
        postgresql-dev \
        musl-dev \
        g++ \
        gcc \
        git \
        linux-headers \
        build-base

rm -rf /root/.cache \
  /usr/share/doc
