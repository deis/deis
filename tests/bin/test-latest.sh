#!/bin/bash
#
# Preps a test environment and runs `make test-integration`
# using the latest published artifacts available on Docker Hub
# and the deis.io website.
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

log_phase "Running test-latest on $DEIS_TEST_APP"

log_phase "Installing clients"

# install deis CLI from http://deis.io/ website
pushd $DEIS_ROOT/deisctl
curl -sSL http://deis.io/deis-cli/install.sh | sh -s 0.13.0
popd

# install deisctl from http://deis.io/ website
# installs latest unit files to $HOME/.deis/units
pushd $DEIS_ROOT/client
curl -sSL http://deis.io/deisctl/install.sh | sh -s 0.13.0
popd

# ensure we use distributed unit files
unset DEISCTL_UNITS

# use the built client binaries
export PATH=$DEIS_ROOT/deisctl:$DEIS_ROOT/client:$PATH

log_phase "Provisioning 3-node CoreOS"

export DEIS_NUM_INSTANCES=3
git checkout contrib/coreos/user-data
make discovery-url
vagrant up --provider virtualbox

log_phase "Waiting for etcd/fleet"

until deisctl list >/dev/null 2>&1; do
    sleep 1
done

log_phase "Provisioning Deis"

# provision deis from master using :latest
deisctl install platform
deisctl scale router=3
deisctl start router@1 router@2 router@3
time deisctl start platform

log_phase "Running integration tests"

# run the smoke tests unless another target is specified
make ${TEST_TYPE:-test-smoke}
