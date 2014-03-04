#!/bin/bash -ex

#
# Prepare a Deis Controller image for Amazon EC2.
#
# Instructions:
#
#   1. Launch a vanilla Ubuntu 12.04 instance (64-bit with an EBS root volume)
#   2. SSH in and install the 3.11 kernel as root with:
#      apt-get update && apt-get install -yq linux-image-generic-lts-saucy linux-headers-generic-lts-saucy && reboot
#   3. After reboot is complete, SSH in and `uname -r` to confirm kernel is 3.11
#   4. Run this script (as root!) to optimize the image for fast boot times
#   5. Create a new AMI from the root volume
#   6. Distribute the AMI to other regions using `ec2-copy-image`
#   7. Update `provision-ec2-controller.sh` script with new AMIs
#

# Remove old kernel(s)
dpkg -l 'linux-*' | sed '/^ii/!d;/'"$(uname -r | sed "s/\(.*\)-\([^0-9]\+\)/\1/")"'/d;s/^[^ ]* [^ ]* \([^ ]*\).*/\1/;/[0-9]/!d' | xargs sudo apt-get -y purge

# Add the Docker repository key to your local keychain
# using apt-key finger you can check the fingerprint matches 36A1 D786 9245 C895 0F96 6E92 D857 6A8B A88D 21E9
curl https://get.docker.io/gpg | apt-key add -

# Add the Docker repository to your apt sources list.
echo deb https://get.docker.io/ubuntu docker main > /etc/apt/sources.list.d/docker.list

# upgrade to latest packages
apt-get update
apt-get dist-upgrade -yq

# install required packages
apt-get install lxc-docker-0.8.0 fail2ban curl git inotify-tools make python-setuptools python-pip -yq

# wait for docker to start
while [ ! -e /var/run/docker.sock ] ; do
  inotifywait -t 2 -e create $(dirname /var/run/docker.sock)
done

# pull current docker images
docker pull deis/data:latest
docker pull deis/discovery:latest
docker pull deis/registry:latest
docker pull deis/cache:latest
docker pull deis/logger:latest
docker pull deis/database:latest
docker pull deis/server:latest
docker pull deis/worker:latest
docker pull deis/builder:latest

# install chef 11.x deps
apt-get install -yq ruby1.9.1 ruby1.9.1-dev make
update-alternatives --set ruby /usr/bin/ruby1.9.1
update-alternatives --set gem /usr/bin/gem1.9.1

# clean and remove old packages
apt-get clean
apt-get autoremove -yq

# reset cloud-init
rm -rf /var/lib/cloud

# purge SSH authorized keys
rm -f /home/ubuntu/.ssh/authorized_keys
rm -f /root/.ssh/authorized_keys

# remove /etc/chef so contents can't intefere with
# node being converged (i.e. old keys)
rm -f /etc/chef/*

# purge /var/log
find /var/log -type f | xargs rm

# flush writes to block storage
sync
