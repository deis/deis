#!/usr/bin/env bash
#
# Usage: ./provision-do-cluster.sh <REGION_ID> <IMAGE_ID> <SSH_ID> <SIZE>
#

set -e

listcontains() {
  for i in $1; do
    [[ $i = $2 ]] && return 0
  done
  return 1
}

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

source $CONTRIB_DIR/utils.sh

if [ -z "$4" ]; then
  echo_red 'Usage: provision-do-cluster.sh <REGION_ID> <IMAGE_ID> <SSH_ID> <SIZE>'
  exit 1
fi

# check for DO tools in $PATH
if ! which tugboat > /dev/null; then
  echo_red 'Please install the tugboat gem and ensure it is in your $PATH.'
  exit 1
fi

if [ -z "$DEIS_NUM_INSTANCES" ]; then
    DEIS_NUM_INSTANCES=3
fi

regions_with_private_networking="4 5 6 7"
if ! listcontains "$regions_with_private_networking" "$1";
then
    echo_red "Invalid region. Please supply a region with private networking support."
    echo_red "Valid regions are:"
    echo_red "4: New York 2"
    echo_red "5: Amsterdam 2"
    echo_red "6: Singapore 1"
    echo_red "7: London 1"
    exit 1
fi

# check that the CoreOS user-data file is valid
$CONTRIB_DIR/util/check-user-data.sh

# launch the Deis cluster on DigitalOcean
i=1 ; while [[ $i -le $DEIS_NUM_INSTANCES ]] ; do \
    NAME=deis-$i
    echo_yellow "Provisioning ${NAME}..."
    tugboat create $NAME -r $1 -i $2 -p true -k $3 -s $4
    ((i = i + 1)) ; \
done

echo_green "Your Deis cluster has successfully deployed to DigitalOcean."
echo_green "Please continue to follow the instructions in the README."
