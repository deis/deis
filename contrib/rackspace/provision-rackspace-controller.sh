#!/usr/bin/env bash
#
# Usage: ./provision-rackspace-controller.sh <region>
#

if [ -z $1 ]; then
  echo usage: $0 [region]
  exit 1
fi

function echo_color {
  echo -e "\033[1m$1\033[0m"
}

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

# check for Deis' general dependencies
if ! $CONTRIB_DIR/check-deis-deps.sh; then
  echo 'Deis is missing some dependencies.'
  exit 1
fi

# check for knife-rackspace
if ! knife rackspace server list > /dev/null; then
  echo 'Please install the knife-rackspace Ruby gem and configure knife.rb.'
  exit 1
fi

#################
# chef settings #
#################
node_name=deis-controller
run_list="recipe[deis::controller]"
chef_version=11.6.2

######################
# Rackspace settings #
######################
region=$1
# see contrib/prepare-rackspace-image.sh for instructions
# on creating your own deis-optmized images
if ! [[ "ord dfw iad lon syd" =~ $region ]]; then
  echo "Unrecognized region: $region"
  exit 1
fi
flavor=$(knife rackspace flavor list | grep '2GB Standard Instance' | awk '{print $1}')
image=$(knife rackspace image list --rackspace-region $region | grep 'deis-base-image' | awk '{print $1}')
if [[ -z $image ]]; then
  echo "Can't find saved image \"deis-base-image\" in region $region. Please follow the"
  echo "instructions in prepare-rackspace-image.sh before provisioning a Deis controller."
  echo
  exit 1
fi

################
# SSH settings #
################
key_name=deis-controller
ssh_key_path=~/.ssh/$key_name
ssh_user="ubuntu"  # doesn't work?

# create ssh keypair and store it
if ! test -e $ssh_key_path; then
  echo_color "Creating new SSH key: $key_name"
  set -x
  ssh-keygen -f $ssh_key_path -t rsa -N '' -C "deis-controller" >/dev/null
  set +x
  echo_color "Saved to $ssh_key_path"
else
  echo_color "WARNING: SSH key $ssh_key_path exists"
fi

# create data bags
knife data bag create deis-users 2>/dev/null
knife data bag create deis-formations 2>/dev/null
knife data bag create deis-apps 2>/dev/null

# trigger Rackspace instance bootstrap
echo_color "Provisioning $node_name with knife rackspace..."
set -x
knife rackspace server create \
 --bootstrap-version $chef_version \
 --rackspace-region $region \
 --image $image \
 --flavor $flavor \
 --rackspace-metadata "{\"Name\": \"$node_name\"}" \
 --rackspace-disk-config MANUAL \
 --identity-file $ssh_key_path \
 --server-name $node_name \
 --node-name $node_name \
 --run-list $run_list
set +x

# Need Chef admin permission in order to add and remove nodes and clients
echo -e "\033[35mPlease ensure that \"deis-controller\" is added to the Chef \"admins\" group.\033[0m"
