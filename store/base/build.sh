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
curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/get-deis/etcdctl-v0.4.9 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
CONFD_VERSION=0.10.0
curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v$CONFD_VERSION/confd-$CONFD_VERSION-linux-amd64 \
	&& chmod +x /usr/local/bin/confd

curl -sSL 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | apt-key add -
echo "deb http://ceph.com/debian-hammer trusty main" > /etc/apt/sources.list.d/ceph.list

apt-get update && apt-get install -yq ceph lsb-release

apt-get clean -y

rm -Rf /usr/share/man /usr/share/doc
rm -rf /tmp/* /var/tmp/*
rm -rf /var/lib/apt/lists/*
