#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

DEBIAN_FRONTEND=noninteractive

# install common packages
apt-get update && apt-get install -y curl net-tools sudo

# install etcdctl
curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# Use modified confd with a fix for /etc/hosts - see https://github.com/kelseyhightower/confd/pull/123
curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.8.0/confd-0.8.0-linux-amd64 \
	&& chmod +x /usr/local/bin/confd

curl -sSL 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | apt-key add -
echo "deb http://ceph.com/debian-giant trusty main" > /etc/apt/sources.list.d/ceph.list

apt-get update && apt-get install -yq ceph

apt-get clean -y

rm -Rf /usr/share/man /usr/share/doc
rm -rf /tmp/* /var/tmp/*
rm -rf /var/lib/apt/lists/*
