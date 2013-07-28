#!/bin/bash -ex

#
# Prepare a Deis-optimized AMI from a vanilla Ubuntu 12.04
#
# Instructions:
#
#   1. Launch a vanilla Ubuntu 12.04 instance (64-bit with an EBS root volume)
#   2. SSH in and install the 3.8 kernel with:
#      apt-get update && apt-get install -yq linux-generic-lts-raring && reboot
#   3. After reboot is complete, SSH in and `uname -r` to confirm kernel is 3.8
#   4. Run this script (as root!) to optimize the image for fast boot times
#   5. Create a new AMI from the root volume
#   6. Distribute the AMI to other regions using `ec2-copy-image`
#   7. Create/update your Deis flavors to use your new AMIs
#

# add docker ppa
apt-add-repository ppa:dotcloud/lxc-docker -y

# upgrade to latest packages
apt-get update
apt-get dist-upgrade -yq

# install required packages
apt-get install lxc-docker curl git python-setuptools python-pip -yq

# create buildstep docker image
git clone https://github.com/opdemand/buildstep.git
cd buildstep
./build.sh ./stack deis/buildstep
cd ..
rm -rf buildstep

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

# ssh host keys are automatically regenerated
# on system boot by ubuntu cloud init

# purge /var/log
find /var/log -type f | xargs rm

# flush writes to block storage
sync
