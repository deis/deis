#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the controller locally"
  echo
  exit 1
fi

DEBIAN_FRONTEND=noninteractive

# install common packages
apt-get update && apt-get install -y curl net-tools sudo

# install etcdctl
curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# Use modified confd with a fix for /etc/hosts - see https://github.com/kelseyhightower/confd/pull/123
curl -sSL https://s3-us-west-2.amazonaws.com/opdemand/confd-git-0e563e5 -o /usr/local/bin/confd
chmod +x /usr/local/bin/confd

curl -sSL 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | apt-key add -
echo "deb http://ceph.com/debian-giant trusty main" > /etc/apt/sources.list.d/ceph.list

apt-get update && apt-get install -yq ceph

apt-get clean -y

rm -Rf /usr/share/man /usr/share/doc
rm -rf /tmp/* /var/tmp/*
rm -rf /var/lib/apt/lists/*
