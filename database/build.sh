#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the database locally"
  echo
  exit 1
fi

apk add --update-cache \
  build-base \
  curl \
  file \
  gcc \
  git \
  libffi-dev \
  libxml2-dev \
  libxslt-dev \
  openssl-dev \
  postgresql \
  postgresql-client \
  python-dev

# pv port.
curl http://dl-3.alpinelinux.org/alpine/edge/testing/x86_64/pv-1.6.0-r0.apk -o /tmp/pv-1.6.0-r0.apk
apk add /tmp/pv-1.6.0-r0.apk

/etc/init.d/postgresql stop || true

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.1.1/contrib/get-pip.py | python -

# install wal-e
cd /tmp
git clone https://github.com/wal-e/wal-e.git

# get a post-v0.8.0 commit which includes a busybox fix
cd /tmp/wal-e
git checkout c6dd4b1

pip install . oslo.config>=1.12.0

# python port of daemontools
pip install envdir

mkdir -p /etc/wal-e.d/env /etc/postgresql/main /var/lib/postgresql

chown -R root:postgres /etc/wal-e.d /etc/postgresql/main /var/lib/postgresql

# cleanup.
apk del --purge \
  build-base \
  gcc \
  git
rm -rf /root/.cache \
  /usr/share/doc \
  /tmp/* \
  /var/cache/apk/*
