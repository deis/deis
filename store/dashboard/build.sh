#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the store dashboard locally"
  echo
  exit 1
fi

DEBIAN_FRONTEND=noninteractive

apt-get update && \
    apt-get install -yq python-dev apache2 libapache2-mod-wsgi python gunicorn git

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

pip install Flask werkzeug

git clone https://github.com/Crapworks/ceph-dash /app/cephdash && \
    cd /app/cephdash && \
    git checkout e2a5f6a

# cleanup. indicate that python is a required package.
apt-mark unmarkauto python python-openssl && \
  apt-get remove -y --purge git git-core build-essential python-dev && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*

