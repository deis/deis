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

log_phase "Running test-nightly on $DEIS_TEST_APP"

log_phase "Installing clients"

# FIXME: switch to deis CLI install from website
cd $DEIS_ROOT/client
sudo python setup.py install
cd $THIS_DIR

# install latest deisctl from the website
curl -sSL http://deis.io/deisctl/install.sh | sudo sh

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

# run the full integration suite
time make test-integration
