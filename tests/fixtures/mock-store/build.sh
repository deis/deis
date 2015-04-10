#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

# install required packages to build
apt-get update \
  && apt-get install -y build-essential git python-dev curl net-tools

# install etcdctl
curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
  && chmod +x /usr/local/bin/etcdctl


git clone https://github.com/jserver/mock-s3 /app/mock-s3 --depth 1
cd /app/mock-s3
#FIXME: This is a gisted patch to enable pseudo-handling of POST requests, otherwise wal-e crashes attempting to delete old wal segments
curl https://gist.githubusercontent.com/anonymous/c565f11a8d90d6e2d92b/raw/c5815f6c83aa5c2cfb7b0a34cfab4a075c97be16/mock-s3-post.diff|git apply

# install pip
curl -sSL https://raw.githubusercontent.com/pypa/pip/6.0.8/contrib/get-pip.py | python -

python setup.py install

# cleanup. indicate that python, libpq and libyanl are required packages.
apt-mark unmarkauto python && \
  apt-get remove -y --purge build-essential python-dev gcc cpp git && \
  apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*

