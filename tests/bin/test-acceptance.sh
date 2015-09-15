#!/usr/bin/env bash
#
# Preps a test environment and runs `make test-integration`
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

log_phase "Provisioning Deis on old release"

function set_release {
  deisctl config $1 set image=deis/$1:$2
}
set_release logger ${OLD_TAG}
set_release router ${OLD_TAG}
set_release database ${OLD_TAG}
set_release controller ${OLD_TAG}
set_release registry ${OLD_TAG}

deisctl install platform
time deisctl start platform

log_phase "Running smoke tests"

time make test-smoke

log_phase "Publishing new release"

time make release

log_phase "Updating channel with new release"

updateservicectl channel update --app-id=${APP_ID} --channel=${CHANNEL} --version=${BUILD_TAG} --publish=true

log_phase "Waiting for upgrade to complete"

# configure platform settings
deisctl config platform set domain=$DEIS_TEST_DOMAIN
deisctl config platform set sshPrivateKey=$DEIS_TEST_SSH_KEY
deisctl config platform channel=${CHANNEL} autoupdate=true

function wait_for_update {
  set -x
  vagrant ssh $1 -c "journalctl -n 500 -u deis-updater.service -f" &
  pid_$1=$!
  vagrant ssh $1 -c "/bin/sh -c \"while [[ \"\$(cat /etc/deis-version)\" != \"${BUILD_TAG}\" ]]; do echo waiting for update to complete...; sleep 5; done\""
  kill $pid_$1
  set +x
}

wait_for_update deis-1 &
update1=$!
wait_for_update deis-2 &
update2=$!
wait_for_update deis-3 &
update3=$!
wait update1 update2 update3

log_phase "Running end-to-end integration test with Python client"

time make test-integration
