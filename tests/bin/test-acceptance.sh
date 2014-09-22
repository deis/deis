#!/bin/bash
#
# Preps a test environment and runs `make test-integration` with single-node vagrant.

echo Testing ${DEIS_TEST_APP?}...
THIS_DIR=$(cd $(dirname $0); pwd)  # absolute path

cd ${GOPATH?}/src/github.com/deis/deis
echo HOST_IPADDR=${HOST_IPADDR?}
echo DEISCTL_TUNNEL=${DEISCTL_TUNNEL?}

# Environment reset and configuration
rm -rf ~/.deis
ssh-add -D || eval $(ssh-agent) && ssh-add -D
ssh-add ~/.vagrant.d/insecure_private_key
ssh-add ~/.ssh/deis
$THIS_DIR/halt-all-vagrants.sh
vagrant destroy --force
rm -rf tests/example-*

set -e

make -C docs/ test
make build
make test-components

if ! [[ -x deisctl ]]; then
    curl -sSL http://deis.io/deisctl/install.sh | sudo sh
fi

if [[ -z "$DEIS_REGISTRY" ]]; then
    docker ps | grep -q registry && curl -s http://$HOST_IPADDR:5000 2>&1 >/dev/null
    [[ $? -eq 0 ]] || make dev-registry
    export DEIS_REGISTRY=$HOST_IPADDR:5000
fi

vagrant up --provider virtualbox
make push
deisctl install platform
deisctl start platform
make test-integration
deisctl uninstall platform
vagrant halt
