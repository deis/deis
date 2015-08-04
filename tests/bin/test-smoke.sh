#!/usr/bin/env bash
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

log_phase "Running documentation tests"

# test building documentation
make -C docs/ test

log_phase "Running unit tests"

make test-unit

log_phase "Building from current source tree"

# build all docker images and client binaries
make build

# use the built client binaries
export PATH=$DEIS_ROOT/deisctl:$DEIS_ROOT/client/dist:$PATH

log_phase "Running functional tests"

make test-functional

log_phase "Provisioning 3-node CoreOS"

export DEIS_NUM_INSTANCES=3
make discovery-url
vagrant up --provider virtualbox

log_phase "Waiting for etcd/fleet"

until deisctl list >/dev/null 2>&1; do
    sleep 1
done

log_phase "Publishing release from source tree"

make dev-release

log_phase "Provisioning Deis"

# configure platform settings
deisctl config platform set domain=$DEIS_TEST_DOMAIN
deisctl config platform set sshPrivateKey=$DEIS_TEST_SSH_KEY

time deisctl install platform
time deisctl start platform

log_phase "Starting smoke tests"

time make test-smoke
