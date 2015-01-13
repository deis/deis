#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the database locally"
  echo
  exit 1
fi

# install postgresql 9.3 from postgresql.org repository as well as requirements for building wal-e
echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main" > /etc/apt/sources.list.d/pgdg.list
curl -sk https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
apt-get update && apt-get install -yq \
                                      curl \
                                      daemontools \
                                      file \
                                      gcc \
                                      git \
                                      libxml2-dev \
                                      libxslt1-dev \
                                      lzop \
                                      postgresql-9.3 \
                                      pv \
                                      python-dev

/etc/init.d/postgresql stop

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

# install wal-e
cd /tmp
git clone https://github.com/wal-e/wal-e.git

cd /tmp/wal-e
git checkout c16e58a

pip install .

mkdir -p /etc/wal-e.d/env

chown -R root:postgres /etc/wal-e.d

# cleanup. indicate that python, libpq and libyanl are required packages.
apt-mark unmarkauto python curl daemontools file libxml2-dev \
  libxslt1-dev lzop postgresql-9.3 pv && \
  apt-get remove -y --purge python-dev gcc cpp libpq-dev libyaml-dev git && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
