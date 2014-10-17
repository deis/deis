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

log_phase "Building from current source tree"

# build all docker images and client binaries
make build

# use the built client binaries
export PATH=$DEIS_ROOT/deisctl:$DEIS_ROOT/client/dist:$PATH

log_phase "Running test-smoke"

make -C docs/ test

log_phase "Running unit and functional tests"

make test-components

log_phase "Provisioning 3-node CoreOS"

export DEIS_NUM_INSTANCES=3
git checkout $DEIS_ROOT/contrib/coreos/user-data
make discovery-url
vagrant up --provider virtualbox

log_phase "Waiting for etcd/fleet"

until deisctl list >/dev/null 2>&1; do
    sleep 1
done

log_phase "Publishing release from source tree"

make dev-release

log_phase "Provisioning Deis"

time deisctl install platform
time deisctl start platform

log_phase "Starting smoke tests"

time make test-smoke
