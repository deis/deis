#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the registry locally"
  echo
  exit 1
fi

DEBIAN_FRONTEND=noninteractive

sed -i 's/main$/main universe/' /etc/apt/sources.list

# install required packages (copied from dotcloud/docker-registry Dockerfile)
apt-get update && \
    apt-get install -y git-core build-essential python-dev \
    libevent-dev python-openssl liblzma-dev

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

# create a registry user
useradd -s /bin/bash registry

# add the docker registry source from github
git clone https://github.com/deis/docker-registry /docker-registry && \
    cd /docker-registry && \
    git checkout 54fa9a1 && \
    chown -R registry:registry /docker-registry

# install boto configuration
cp /docker-registry/config/boto.cfg /etc/boto.cfg
cd /docker-registry && pip install -r requirements/main.txt

# Install core
pip install /docker-registry/depends/docker-registry-core

# Install registry
pip install file:///docker-registry#egg=docker-registry[bugsnag,newrelic,cors]

# cleanup. indicate that python is a required package.
apt-mark unmarkauto python python-openssl && \
  apt-get remove -y --purge git-core build-essential python-dev && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
