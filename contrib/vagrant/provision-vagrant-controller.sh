#!/usr/bin/env bash

function echo_color {
  echo -e "\033[1m$1\033[0m"
}

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

# check for Deis' general dependencies
if ! $CONTRIB_DIR/check-deis-deps.sh; then
  echo 'Deis is missing some dependencies.'
  exit 1
fi

#################
# chef settings #
#################
node_name=deis-controller
run_list="recipe[deis::controller]"
chef_version=11.4.4

################
# SSH settings #
################
ssh_key_path=~/.vagrant.d/insecure_private_key
ssh_user="vagrant"
ssh_port="22"

# create data bags
knife data bag create deis-users 2>/dev/null
knife data bag create deis-formations 2>/dev/null
knife data bag create deis-apps 2>/dev/null

# start a controller and two nodes through vagrant
vagrant up --provider virtualbox --parallel

# trigger vagrant instance bootstrap
echo_color "Provisioning $node_name with knife..."
set -x
knife bootstrap 192.168.61.100 \
 --bootstrap-version $chef_version \
 --ssh-user $ssh_user \
 --ssh-port $ssh_port \
 --identity-file $ssh_key_path \
 --node-name $node_name \
 --run-list $run_list \
 --sudo
set +x
