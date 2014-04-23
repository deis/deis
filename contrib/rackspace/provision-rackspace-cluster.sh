#!/usr/bin/env bash
#
# Usage: ./provision-rackspace-cluster.sh <key pair name> [flavor]
#

set -e

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

source $CONTRIB_DIR/utils.sh

if [ -z "$1" ]; then
  echo_red 'Usage: provision-rackspace-cluster.sh <key pair name> [flavor]'
  exit 1
fi

if [ -z "$2" ]; then
  FLAVOR="performance1-2"
else
  FLAVOR=$2
fi

if ! which supernova > /dev/null; then
  echo_red 'Please install the dependencies listed in the README and ensure they are in your $PATH.'
  exit 1
fi

if [ -z "$DEIS_NUM_INSTANCES" ]; then
    DEIS_NUM_INSTANCES=3
fi

i=1 ; while [[ $i -le $DEIS_NUM_INSTANCES ]] ; do \
    echo_yellow "Provisioning deis-$i..."
    supernova production boot --image 6bdbd558-e66c-49cc-9ff3-126e7411f602 --flavor $FLAVOR --key-name $1 --user-data ../coreos/user-data --config-drive true deis-$i ; \
    ((i = i + 1)) ; \
done

echo_green "Your Deis cluster has successfully deployed to Rackspace."
echo_green "Please continue to follow the instructions in the README."
