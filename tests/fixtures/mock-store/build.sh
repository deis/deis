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
apt-get update \
  && apt-get install -y build-essential git python-dev curl net-tools

# install etcdctl
curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
  && chmod +x /usr/local/bin/etcdctl


git clone https://github.com/jserver/mock-s3 /app/mock-s3

cd /app/mock-s3

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

python setup.py install

# cleanup. indicate that python, libpq and libyanl are required packages.
apt-mark unmarkauto python && \
  apt-get remove -y --purge build-essential python-dev gcc cpp git && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*

