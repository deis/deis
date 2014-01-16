#!/usr/bin/env bash
#
# Usage: ./provision-vagrant-controller.sh
#

function echo_color {
  echo -e "\033[1m$1\033[0m"
}

THIS_DIR="$(cd $(dirname $0); pwd)" # absolute path
CONTRIB_DIR=$(dirname "$THIS_DIR")

# check for Deis' general dependencies
if ! "$CONTRIB_DIR/check-deis-deps.sh"; then
  echo 'Deis is missing some dependencies.'
  exit 1
fi

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

#################
# chef settings #
#################
node_name=deis-controller
run_list="recipe[deis::controller]"
chef_version=11.6.2

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

# Boot the deis-controller VM
echo_color "Booting $node_name with 'vagrant up'"
pushd $THIS_DIR/../../
vagrant up --provision
if [ $? -gt 0 ]; then
  echo_color "Canceling provision because 'vagrant up' failed"
  exit 1
fi

# Add the Controller's public SSH key to user's machine. This allows the Controller to
# issue vagrant commands on the host machine.
read -p "Add the Deis Controller's SSH key to your authorized_keys file? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then #TODO: Might be nice to have flag to make manual confirmation optional?

  # Create an SSH key pair for the deis user
  vagrant ssh -c "
  if [ ! -f ~/.ssh/id_rsa ]; then
    ssh-keygen -t rsa -N \"\" -f ~/.ssh/id_rsa
    chmod a+r ~/.ssh/id_rsa # Not strictly best practice, but the deis user needs to be able to read it
  fi"

  # Copy the created key over to your local machine
  scp \
    -P22 \
    -o IdentityFile=$ssh_key_path \
    'vagrant@deis-controller.local:/home/vagrant/.ssh/id_rsa.pub' \
    '/tmp/vagrant_key'
  KEY=$(cat /tmp/vagrant_key)

  if [ ! -n "$KEY" ]; then
    echo_color "Aborting. No SSH key copied from the Deis Controller"
    exit 1
  fi

  if [ -z "$(grep "$KEY" ~/.ssh/authorized_keys )" ]; then
    echo $KEY >> ~/.ssh/authorized_keys;
    echo_color "Key added."
  else
    echo_color "Key already added."
  fi

fi

echo_color "Provisioning $node_name with knife vagrant..."
set -x
knife bootstrap "$node_name.local" \
  --bootstrap-version $chef_version \
  --ssh-user $ssh_user \
  --ssh-port $ssh_port \
  --identity-file $ssh_key_path \
  --node-name $node_name \
  --run-list $run_list \
  --sudo
set +x

echo_color "Updating Django site object from 'example.com' to 'deis-controller'..."
vagrant ssh -c "sudo su deis -c \"psql deis -c \\\" \
  UPDATE django_site \
  SET domain = 'deis-controller.local', \
      name = 'deis-controller.local' \
  WHERE id = 1 \\\"\"" >/dev/null

if [ $? -eq 0 ]; then
  echo_color "Site object updated."
fi
popd

echo_color "Setting devmode flag on 'deis-controller'..."
knife exec -E 'nodes.transform("name:deis-controller") {|n| n.normal_attrs["deis"]["devmode"] = true; n.save }'

# Need Chef admin permission in order to add and remove nodes and clients
echo -e "\033[35mPlease ensure that \"deis-controller\" is added to the Chef \"admins\" group.\033[0m"
