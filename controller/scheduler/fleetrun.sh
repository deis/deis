#!/bin/bash
set -e

SSH_OPTIONS="-i $FLEETW_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR"

# set debug if provided as an envvar
[[ $DEBUG ]] && set -x

# run the fleetctl command remotely
ssh $SSH_OPTIONS core@$FLEETW_HOST "$@"