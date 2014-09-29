#!/bin/bash
#
# Prepares the process environment to run a test

echo
echo "Preparing test environment..."
echo

# use GOPATH to determine project root
export DEIS_ROOT=${GOPATH?}/src/github.com/deis/deis
echo "DEIS_ROOT=$DEIS_ROOT"

# the application under test
export DEIS_TEST_APP=${DEIS_TEST_APP:-example-go}
echo "DEIS_TEST_APP=$DEIS_TEST_APP"

# SSH key name used for testing
export DEIS_TEST_AUTH_KEY=${DEIS_TEST_AUTH_KEY:-deis-test}
echo "DEIS_TEST_AUTH_KEY=$DEIS_TEST_AUTH_KEY"

# SSH key used for deisctl tunneling
export DEIS_TEST_SSH_KEY=${DEIS_TEST_SSH_KEY:-~/.vagrant.d/insecure_private_key}
echo "DEIS_TEST_SSH_KEY=$DEIS_TEST_SSH_KEY"

# domain used for wildcard DNS
export DEIS_TEST_DOMAIN=${DEIS_TEST_DOMAIN:-local3.deisapp.com}
echo "DEIS_TEST_DOMAIN=$DEIS_TEST_DOMAIN"

# SSH tunnel used by deisctl
export DEISCTL_TUNNEL=${DEISCTL_TUNNEL:-127.0.0.1:2222}
echo "DEISCTL_TUNNEL=$DEISCTL_TUNNEL"

# the registry used to host dev-release images
# must be accessible to local Docker engine and Deis cluster
export DEV_REGISTRY=${DEV_REGISTRY?}
echo "DEV_REGISTRY=$DEV_REGISTRY"

# bail if registry is not accessible
if ! curl -s $DEV_REGISTRY; then
  echo "DEV_REGISTRY is not accessible, exiting..."
  exit 1
fi
echo ; echo

# disable git+ssh host key checking
export GIT_SSH=$DEIS_ROOT/tests/bin/git-ssh-nokeycheck.sh

# cleanup any stale example applications
rm -rf $DEIS_ROOT/tests/example-*

# generate ssh key if it doesn't already exist
test -e ~/.ssh/$DEIS_TEST_AUTH_KEY || ssh-keygen -t rsa -f ~/.ssh/$DEIS_TEST_AUTH_KEY -N ''

# prepare the SSH agent
ssh-add -D || eval $(ssh-agent) && ssh-add -D
ssh-add $DEIS_TEST_SSH_KEY

# clean out deis session data
rm -rf ~/.deis

# clean out vagrant environment
$THIS_DIR/halt-all-vagrants.sh
vagrant destroy --force

# wipe out all vagrants & deis virtualboxen
function cleanup {
    echo
    echo "Cleaning up..."
    echo
    set +e
    ${GOPATH}/src/github.com/deis/deis/tests/bin/destroy-all-vagrants.sh
    VBoxManage list vms | grep deis | sed -n -e 's/^.* {\(.*\)}/\1/p' | xargs -L1 -I {} VBoxManage unregistervm {} --delete
    vagrant global-status --prune
    echo
    echo "Test run complete."
    echo
}

function dump_logs {
  echo
  echo "Error detected, dumping logs..."
  echo
  set +e
  export FLEETCTL_TUNNEL=$DEISCTL_TUNNEL
  set -x
  fleetctl -strict-host-key-checking=false list-units
  fleetctl -strict-host-key-checking=false ssh deis-controller@1 docker logs deis-controller
  fleetctl -strict-host-key-checking=false ssh deis-builder@1 docker logs deis-builder
  fleetctl -strict-host-key-checking=false ssh deis-logger@1 docker logs deis-logger
  set +x
  exit 1
}
