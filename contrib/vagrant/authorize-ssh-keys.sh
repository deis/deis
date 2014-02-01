#!/bin/bash

function echo_color {
  echo -e "\033[1m$1\033[0m"
}

################
# SSH settings #
################
ssh_key_path=~/.vagrant.d/insecure_private_key

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
