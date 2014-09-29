#!/bin/bash
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

echo
echo "Running test-acceptance..."
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

until deisctl list >/dev/null 2>&1; do
    sleep 1
done

echo
echo "Provisioning Deis on old release..."
echo

function set_release {
  deisctl config $1 set image=deis/$1:$2
}
set_release logger ${OLD_TAG}
set_release cache ${OLD_TAG}
set_release router ${OLD_TAG}
set_release database ${OLD_TAG}
set_release controller ${OLD_TAG}
set_release registry ${OLD_TAG}

deisctl install platform
deisctl scale router=3
deisctl start router@1 router@2 router@3
time deisctl start platform

echo
echo "Running smoke tests..."
echo

time make test-smoke

echo
echo "Publishing new release..."
echo

time make release

echo
echo "Updating channel with new release..."
echo

updateservicectl channel update --app-id=${APP_ID} --channel=${CHANNEL} --version=${BUILD_TAG} --publish=true

echo
echo "Waiting for upgrade to complete..."
echo

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

echo
echo "Running end-to-end integration test..."
echo

time make test-integration
