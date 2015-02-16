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
apk add --update-cache \
  build-base \
  git \
  libffi-dev \
  libpq \
  openldap \
  openldap-dev \
  postgresql-dev \
  python \
  python-dev

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.1.1/contrib/get-pip.py | python -

# add a deis user
adduser deis -D -h /app -s /bin/bash

# create a /app directory for storing application data
mkdir -p /app && chown -R deis:deis /app

# create directory for confd templates
mkdir -p /templates && chown -R deis:deis /templates

# install dependencies
pip install -r /app/requirements.txt

# cleanup.
apk del --purge \
  build-base \
  git \
  libffi-dev \
  openldap-dev \
  postgresql-dev \
  python-dev
rm -rf /var/cache/apk/*
