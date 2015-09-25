#!/usr/bin/env bash
#
# Preps a test environment and runs `make test-integration`
# using the latest published artifacts available on Docker Hub
# and the deis.io website.
#
# Usage: ./test-latest.sh [COREOS_CHANNEL]
#

# fail on any command exiting non-zero
set -eo pipefail

export COREOS_CHANNEL=${1:-stable}

# absolute path to current directory
export THIS_DIR=$(cd $(dirname $0); pwd)

# setup the test environment
source $THIS_DIR/test-setup.sh

# setup callbacks on process exit and error
trap cleanup EXIT
trap dump_logs ERR

log_phase "Running test-latest on $DEIS_TEST_APP"

log_phase "Installing clients"

# install deis CLI from http://deis.io/ website
pushd $DEIS_ROOT/client
curl -sSL http://deis.io/deis-cli/install.sh | sh
popd

# install deisctl from http://deis.io/ website
# installs latest unit files to $HOME/.deis/units
pushd $DEIS_ROOT/deisctl
curl -sSL http://deis.io/deisctl/install.sh | sh
popd

# ensure we use distributed unit files
unset DEISCTL_UNITS

# use the built client binaries
export PATH=$DEIS_ROOT/deisctl:$DEIS_ROOT/client:$PATH

log_phase "Provisioning 3-node CoreOS"

export DEIS_NUM_INSTANCES=3
make discovery-url
vagrant up --provider virtualbox

log_phase "Waiting for etcd/fleet"

until deisctl list >/dev/null 2>&1; do
    sleep 1
done

log_phase "Provisioning Deis"

# configure platform settings
deisctl config platform set domain=$DEIS_TEST_DOMAIN
deisctl config platform set sshPrivateKey=$DEIS_TEST_SSH_KEY

# provision deis from master using :latest
time deisctl install platform
time deisctl start platform

log_phase "Running integration tests"

# run the smoke tests unless another target is specified
make ${TEST_TYPE:-test-smoke}
