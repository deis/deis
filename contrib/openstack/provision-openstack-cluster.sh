#!/usr/bin/env bash
#
# Usage: ./provision-openstack-cluster.sh <image id> <flavor> <key pair name> <net id>
#

set -e

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

source $CONTRIB_DIR/utils.sh
source $THIS_DIR/novaclient-env.sh

if [ $# != 4 ]; then
  echo_red 'Usage: provision-openstack-cluster.sh <image id> <flavor> <key pair name> <net id>'
  exit 1
fi

if [ -z "$DEIS_NUM_INSTANCES" ]; then
    DEIS_NUM_INSTANCES=3
fi

i=1 ; while [[ $i -le $DEIS_NUM_INSTANCES ]] ; do \
    echo_yellow "Provisioning deis-$i..."
    nova boot --image $1 --flavor $2 --key-name $3 --user-data \
    ../coreos/user-data --nic net-id=$4 --config-drive true deis-$i ; \
    ((i = i + 1)) ; \
done

echo_green "Your Deis cluster has successfully deployed to openstack."
echo_green "Please continue to follow the instructions in the README."
