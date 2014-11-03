#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the controller locally"
  echo
  exit 1
fi

DEBIAN_FRONTEND=noninteractive

# install required system packages
# HACK: install git so we can install bacongobbler's fork of django-fsm
# install openssh-client for temporary fleetctl wrapper
apt-get update && \
    apt-get install -yq python-dev libffi-dev libpq-dev libyaml-dev git

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

# add a deis user that has passwordless sudo (for now)
useradd deis --groups sudo --home-dir /app --shell /bin/bash
sed -i -e 's/%sudo\tALL=(ALL:ALL) ALL/%sudo\tALL=(ALL:ALL) NOPASSWD:ALL/' /etc/sudoers

# create a /app directory for storing application data
mkdir -p /app && chown -R deis:deis /app

# create directory for confd templates
mkdir -p /templates && chown -R deis:deis /templates

# install dependencies
pip install -r /app/requirements.txt

# cleanup. indicate that python, libpq and libyanl are required packages.
apt-mark unmarkauto python python-openssl libpq5 libpython2.7 libffi6 libyaml-0-2 && \
  apt-get remove -y --purge python-dev gcc cpp libffi-dev libpq-dev libyaml-dev git && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
