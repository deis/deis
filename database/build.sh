#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the database locally"
  echo
  exit 1
fi

apk-install curl \
  file \
  gcc \
  git \
  libxml2-dev \
  libxslt-dev \
  postgresql \
  postgresql-client \
  build-base \
  libffi-dev \
  openssl-dev \
  python-dev

# pv port.
curl http://dl-3.alpinelinux.org/alpine/edge/testing/x86_64/pv-1.5.3-r0.apk -o /tmp/pv-1.5.3-r0.apk
apk add /tmp/pv-1.5.3-r0.apk

/etc/init.d/postgresql stop || true

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.0.8/contrib/get-pip.py | python -

# install wal-e
# https://github.com/wal-e/wal-e/issues/165
cd /tmp
git clone https://github.com/aledbf/wal-e.git

cd /tmp/wal-e
git checkout 5d3fb74

pip install .

# python port of daemontools
pip install envdir

mkdir -p /etc/wal-e.d/env /etc/postgresql/main /var/run/postgresql

chown -R root:postgres /etc/wal-e.d /etc/postgresql/main /var/run/postgresql

# cleanup.
apk del gcc \
  git \
  musl-dev \
  linux-headers \
  libc-dev \
  build-base

rm -rf /root/.cache \
  /usr/share/doc \
  /tmp/*