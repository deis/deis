#!/bin/bash
#
# Preps a test environment and runs `make test-integration` with single-node vagrant.

echo Testing ${DEIS_TEST_APP?}...
THIS_DIR=$(cd $(dirname $0); pwd)  # absolute path

# Environment reset and configuration
rm -rf ~/.deis ~/.fleetctl ~/.ssh/known_hosts
ssh-add -D || eval $(ssh-agent) && ssh-add -D
ssh-add ~/.vagrant.d/insecure_private_key
ssh-add ~/.ssh/deis
cd ${GOPATH?}/src/github.com/deis/deis
rm -rf tests/example-*

# Vagrant provisioning
$THIS_DIR/halt-all-vagrants.sh
vagrant destroy --force
vagrant up --provider virtualbox --provision

# Trap exit signal to halt vagrant
function cleanup {
    set +e
    make stop
    vagrant halt
}
trap cleanup EXIT

set -e

# Build updated Deis CLI and use it for testing
virtualenv --system-site-packages venv
. venv/bin/activate
pip install docopt==0.6.2 python-dateutil==2.2 PyYAML==3.11 requests==2.3.0 pyinstaller==2.1 termcolor==1.1.0
make -C client/ client
chmod +x client/dist/deis
export PATH=`pwd`/client/dist:$PATH

# Install Deis and run tests
make build
make run
make test-integration
