#!/usr/bin/env bash
#
# Usage: ./provision-vagrant-controller.sh
#

set -e

function echo_color {
  echo -e "\033[1m$1\033[0m"
}

THIS_DIR="$(cd $(dirname $0); pwd)" # absolute path
CONTRIB_DIR=$(dirname "$THIS_DIR")
CODE_BASE_DIR=$(dirname "$CONTRIB_DIR/../../")

# For those upgrading from the pre-containerize branch, they may still have some redundant files
# in their code base.
if [ -h $CODE_BASE_DIR/deis/local_settings.py ]; then
  echo "Removing old local_settings symlink"
  rm $CODE_BASE_DIR/deis/local_settings.py
  rm $CODE_BASE_DIR/deis/local_settings.pyc
fi

echo_color "Checking for Deis dependecnies..."

# check for Deis' general dependencies
# if ! "$CONTRIB_DIR/check-deis-deps.sh"; then
#   echo 'Deis is missing some dependencies.'
#   exit 1
# fi

# Make sure SSHD is installed
# TODO: Better SSH server detection
if [ ! -f /etc/ssh/sshd_config ] && [ ! -f /etc/sshd_config ]; then
  echo 'Please install an SSH server'
  exit 1
fi

# Make sure avahi-daemon is installed and running
if [[ `uname -s` =~ Linux ]]; then
  if ! pgrep avahi-daemon >/dev/null; then
    echo 'Please install avahi-daemon to broadcast your hostname to the local network.'
    exit 1
  fi
fi

echo_color "Ensuring submodules are cloned and up to date..."

# Ensure all submodules are cloned and up to date
git submodule init && git submodule update

#################
# chef settings #
#################
node_name=deis-controller
run_list="recipe[deis::controller]"
chef_version=11.6.2

################
# SSH settings #
################
ssh-keygen -R "$node_name.local"
ssh_key_path=~/.vagrant.d/insecure_private_key
ssh_user="vagrant"
ssh_port="22"

echo_color "Creating knife data bags..."

# create data bags
# knife data bag create deis-formations 2>/dev/null
# knife data bag create deis-apps 2>/dev/null

# Boot the deis-controller VM
echo_color "Booting $node_name with 'vagrant up'"
pushd "$THIS_DIR"
vagrant up --provision
if [ $? -gt 0 ]; then
  echo_color "Canceling provision because 'vagrant up' failed"
  exit 1
fi

# 'deis provider:discover' detects the host machine's user and IP address, however, that command cannot
# be guareteed to run inside the deis codebase. Therefore we can't use that opportunity to discover
# the path of the codebase on the host machine. Therefore we do it now as this script has to exist
# inside the codebase.
nodes_dir="$CONTRIB_DIR/vagrant/nodes"
nodes_path_file="$CONTRIB_DIR/vagrant/.host_nodes_dir"
echo $nodes_dir > $nodes_path_file

# Add the Controller's public SSH key to user's machine. This allows the Controller to
# issue vagrant commands on the host machine.
echo_color "Ensuring presence of Controller's public key in your ~/.ssh/authorized_keys file..."
KEY=$(cat util/ssh_keys/id_rsa_vagrant-deis-controller.pub)
if [ -z "$(grep "$KEY" ~/.ssh/authorized_keys )" ]; then
  echo $KEY >> ~/.ssh/authorized_keys;
  echo_color "Key added."
else
  echo_color "Key already added."
fi

chef_json=$(echo '
{
  "deis": {
      "public_ip": "192.168.61.100",
      "dev": {
        "mode": true,
        "source": "/vagrant"
      }
  }
}' | tr '\n' ' ')

echo_color "Provisioning $node_name with knife vagrant..."
set -x
knife bootstrap "$node_name.local" \
  --bootstrap-version $chef_version \
  --ssh-user $ssh_user \
  --ssh-port $ssh_port \
  --identity-file $ssh_key_path \
  --node-name $node_name \
  --run-list $run_list \
  --json-attributes "$chef_json" \
  --sudo
set +x

echo_color "Running post chef run setup..."
vagrant ssh -c '/vagrant/contrib/vagrant/util/_post_chef_run.sh'
popd

# Need Chef admin permission in order to add and remove nodes and clients
echo -e "\033[35mPlease ensure that \"deis-controller\" is added to the Chef \"admins\" group.\033[0m"

set +e
