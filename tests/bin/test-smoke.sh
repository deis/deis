#!/bin/bash
#
# Preps a test environment and runs `make test-smoke`
# against artifacts produced from the current source tree
#

# fail on any command exiting non-zero
set -eo pipefail

# absolute path to current directory
export THIS_DIR=$(cd $(dirname $0); pwd)

# setup the test environment
source $THIS_DIR/test-setup.sh

# setup callbacks on process exit and error
trap cleanup EXIT
trap dump_logs ERR

echo
echo "Running test-smoke..."
echo

# test building documentation
make -C docs/ test

echo
echo "Building from current source tree..."
echo

# build all docker images and client binaries
make build

# use the built client binaries
export PATH=$DEIS_ROOT/deisctl:$DEIS_ROOT/client/dist:$PATH

echo
echo "Running unit and functional tests..."
echo

make test-components

echo
echo "Provisioning 3-node CoreOS..."
echo

export DEIS_NUM_INSTANCES=3
git checkout $DEIS_ROOT/contrib/coreos/user-data
make discovery-url
vagrant up --provider virtualbox

echo
echo "Waiting for etcd/fleet..."
echo

until deisctl list >/dev/null 2>&1; do
    sleep 1
done

echo
echo "Publishing release from source tree..."
echo

make dev-release

echo
echo "Provisioning Deis..."
echo

time deisctl install platform
deisctl scale router=3
deisctl start router@1 router@2 router@3
time deisctl start platform

echo
echo "Starting smoke tests..."
echo

time make test-smoke
