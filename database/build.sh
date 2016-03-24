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
curl -sSL https://bootstrap.pypa.io/get-pip.py | python - pip==8.1.1

# install wal-e
pip install --disable-pip-version-check --no-cache-dir wal-e==0.8.1 oslo.config>=1.12.0

# python port of daemontools
pip install --disable-pip-version-check --no-cache-dir envdir==0.7

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
