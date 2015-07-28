#!/usr/bin/env bash
#
# Usage: ./provision-do-cluster.sh <REGION_ID> <SSH_ID> <SIZE>
#

set -e

listcontains() {
  for i in $1; do
    [[ $i = "$2" ]] && return 0
  done
  return 1
}

THIS_DIR=$(cd "$(dirname "$0")"; pwd) # absolute path
CONTRIB_DIR=$(dirname "$THIS_DIR")

REGION_SLUG=$1
SSH_ID=$2
SIZE=$3
PREFIX=$4

if [ -z "$PREFIX" ]; then
  PREFIX="deis"
fi

source "$CONTRIB_DIR/utils.sh"

if [ -z "$3" ]; then
  echo_red 'Usage: provision-do-cluster.sh <REGION_SLUG> <SSH_ID> <SIZE> [PREFIX]'
  exit 1
fi

# check for DO tools in $PATH
if ! which docl > /dev/null; then
  echo_red 'Please install the docl gem and ensure it is in your PATH.'
  exit 1
fi

if [ -z "$DEIS_NUM_INSTANCES" ]; then
    DEIS_NUM_INSTANCES=3
fi

regions_without_private_networking_or_metadata="nyc1 nyc2 ams1"
if listcontains "$regions_without_private_networking_or_metadata" "$REGION_SLUG";
then
    echo_red "Invalid region. Please supply a region with private networking & metadata support."
    echo_red "Valid regions are (use the name in brackets):"
    docl regions --private_networking --metadata
    exit 1
fi

# check that the CoreOS user-data file is valid
"$CONTRIB_DIR/util/check-user-data.sh"

# TODO: Make it follow a specific ID once circumstances allow us to do so.
BASE_IMAGE_ID='coreos-stable'

if [ -z "$BASE_IMAGE_ID" ]; then
	echo_red "DigitalOcean Image not found..."
	exit 1
fi

# launch the Deis cluster on DigitalOcean
i=1 ; while [[ $i -le $DEIS_NUM_INSTANCES ]] ; do \
    NAME="$PREFIX-$i"
    echo_yellow "Provisioning ${NAME}..."
    docl create "$NAME" "$BASE_IMAGE_ID" "$SIZE" "$REGION_SLUG" --key="$SSH_ID" --private-networking --user-data="$CONTRIB_DIR/coreos/user-data" --wait
    ((i = i + 1)) ; \
done

echo_green "Your Deis cluster has successfully deployed to DigitalOcean."
echo_green "Please continue to follow the instructions in the README."
