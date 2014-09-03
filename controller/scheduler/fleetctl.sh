#!/bin/bash
set -e

read FLEETW_HOST PORT <<<$(IFS=":"; echo $FLEETW_HOST)
SCP_PORT="-P ${PORT:=22}"
SSH_PORT="-p ${PORT:=22}"
SSH_OPTIONS="-i $FLEETW_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR"

# set debug if provided as an envvar
[[ $DEBUG ]] && set -x

# if fleet unit is defined, scp it to the remote host
if [[ $FLEETW_UNIT ]]; then
  unitfile=$(mktemp)
  echo $FLEETW_UNIT_DATA | base64 -d > $unitfile
  scp $SCP_PORT $SSH_OPTIONS $unitfile core@$FLEETW_HOST:$FLEETW_UNIT
fi

# run the fleetctl command remotely
ssh $SSH_PORT $SSH_OPTIONS core@$FLEETW_HOST fleetctl $@

# clean up
if [[ $FLEETW_UNIT ]]; then
  ssh $SSH_PORT $SSH_OPTIONS core@$FLEETW_HOST "rm -f /home/core/$FLEETW_UNIT"
fi
