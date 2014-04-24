#!/usr/bin/env bash
#
# Usage: ./initialize-rackspace-cluster.sh
#

set -e

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)
ROOT_DIR=$(dirname $CONTRIB_DIR)

source $CONTRIB_DIR/utils.sh

# check for fleetctl in $PATH
if ! which fleetctl > /dev/null; then
  echo_red 'Please install fleetctl and ensure it is in your $PATH.'
  echo_red 'See https://github.com/coreos/fleet for more information'
  exit 1
fi

if [ -z "$FLEETCTL_TUNNEL" ]
then
    echo_red 'Please set $FLEETCTL_TUNNEL.'
    echo_red 'See https://github.com/coreos/fleet/blob/master/Documentation/remote-access.md'
    exit 1
fi

cd $ROOT_DIR

# upload all systemd unit to the fleet cluster
units=()
for component in builder cache controller database logger registry router; do
  units+=($component/systemd/*)
done

fleetctl submit ${units[@]}
fleetctl start ${units[@]}

echo_green "Done! Inspect the state of the services with: fleetctl list-units"
echo_green "Once all the services are running, you can register with your Deis cluster: deis register 1.2.3.4:8000"
