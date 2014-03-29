#!/bin/bash
set -e

SSH_OPTIONS="-i $FLEETW_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"

# set debug if provided as an envvar
[[ $DEBUG ]] && set -x

# if fleet unit is defined, scp it to the remote host
if [[ $FLEETW_UNIT ]]; then
  unitfile=$(mktemp)
  echo $FLEETW_UNIT_DATA | base64 -d > $unitfile
  scp $SSH_OPTIONS $unitfile core@$FLEETW_HOST:$FLEETW_UNIT
fi

# run the fleetctl command remotely
ssh $SSH_OPTIONS core@$FLEETW_HOST fleetctl $@