#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

DEBIAN_FRONTEND=noninteractive

sed -i 's/main$/main universe/' /etc/apt/sources.list

# install required packages (copied from dotcloud/docker-registry Dockerfile)
apt-get update && \
    apt-get install -y git-core build-essential python-dev \
    libevent-dev libssl-dev libyaml-0-2 libyaml-dev python-openssl liblzma-dev swig

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.0.8/contrib/get-pip.py | python -

# create a registry user
useradd -s /bin/bash registry

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
apt-mark unmarkauto python python-openssl libyaml-0-2 && \
  apt-get remove -y --purge git-core build-essential libssl-dev libyaml-dev python-dev swig && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
