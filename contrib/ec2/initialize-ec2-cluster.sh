#!/usr/bin/env bash
#
# Usage: ./initialize-ec2-cluster.sh
#

set -e

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)
ROOT_DIR=$(dirname $CONTRIB_DIR)

# check for fleetctl in $PATH
if ! which fleetctl > /dev/null; then
  echo 'Please install fleetctl and ensure it is in your $PATH.'
  echo 'See https://github.com/coreos/fleet for more information'
  exit 1
fi

if [ -z "$FLEETCTL_TUNNEL" ]
then
    echo 'Please set $FLEETCTL_TUNNEL.'
    echo 'See https://github.com/coreos/fleet/blob/master/Documentation/remote-access.md'
    exit 1
fi

cd $ROOT_DIR

# upload each component's systemd unit to the fleet cluster
for component in registry logger database cache controller builder router
do
  pushd $component/systemd > /dev/null
  fleetctl submit deis-$component.service
  fleetctl start deis-$component.service
  popd > /dev/null
done

echo "done!"
