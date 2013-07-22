#!/bin/sh -e

# upgrade to latest packages
apt-get update
apt-get upgrade -yq

# install 3.8 kernel
apt-get install -yq linux-generic-lts-raring

# install chef 11.x deps
apt-get install -yq ruby1.9.1 ruby1.9.1-dev make

# cleanup for bundle
rm -rf /var/lib/cloud
rm -f /root/.ssh/authorized_keys
find /var/log -type f | xargs rm

