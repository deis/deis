#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the database locally"
  echo
  exit 1
fi

<<<<<<< HEAD
# install postgresql 9.3 from postgresql.org repository as well as requirements for building wal-e
echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main" > /etc/apt/sources.list.d/pgdg.list
curl -sk https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
apt-get update && apt-get install -yq \
                                      curl \
                                      daemontools \
                                      file \
                                      gcc \
                                      g++ \
                                      git \
                                      libffi-dev \
                                      libxml2-dev \
                                      libxslt1-dev \
                                      libssl-dev \
                                      lzop \
                                      postgresql-9.3 \
                                      pv \
                                      python-dev

/etc/init.d/postgresql stop
=======
apk-install curl \
  file \
  gcc \
  git \
  libxml2-dev \
  libxslt-dev \
  postgresql \
  postgresql-client \
  build-base \
  python-dev

/etc/init.d/postgresql stop || true
>>>>>>> feat(*): use alpine linux.

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.0.8/contrib/get-pip.py | python -

# install wal-e
# https://github.com/wal-e/wal-e/issues/165
cd /tmp
git clone https://github.com/aledbf/wal-e.git

cd /tmp/wal-e
<<<<<<< HEAD
git checkout v0.8c2
=======
git checkout 5d3fb74
>>>>>>> feat(*): use alpine linux.

pip install .

# python port of daemontools
pip install envdir

mkdir -p /etc/wal-e.d/env /etc/postgresql/main /var/run/postgresql

<<<<<<< HEAD
# cleanup. indicate python, curl, and others as required packages.
apt-mark unmarkauto python curl daemontools file libffi-dev libxml2-dev \
  libxslt1-dev libssl-dev lzop postgresql-9.3 pv && \
  apt-get remove -y --purge gcc g++ git python-dev && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
=======
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
>>>>>>> feat(*): use alpine linux.
