#!/usr/bin/env bash

function echo_color {
  echo -e "\033[1m$1\033[0m"
}

# Require a command-line arg for the formation name
if [ -z $1 ]; then
  echo usage: $0 [formation]
  exit 1
fi
formation=$1

# Create a "static" formation, where nodes must be added manually
deis formations:create $formation --flavor=static

# Update the layer to SSH as "vagrant", and retrieve its SSH public key
echo_color "Updating the runtime layer to SSH as \"vagrant\"..."
ssh_key=$(deis layers:update $formation runtime --ssh_username=vagrant |
	      grep -Eo '\"ssh_public_key\"\: \"(.*)\"' | cut -d\" -f4)
authfile=.ssh/authorized_keys
tmpfile=/tmp/authorized_keys.tmp

# SSH into deis-node-1 and authorize the SSH public key
echo_color "Adding the layer's public key to deis-node-1..."
vagrant ssh deis-node-1 -c "echo $ssh_key|cat - $authfile > $tmpfile && mv $tmpfile $authfile"

# Add deis-node-1 to the formation
echo_color "Adding node deis-node-1 to formation and provisioning..."
deis nodes:create $formation deis-node-1.local --layer=runtime

# SSH into deis-node-2 and authorize the SSH public key
echo_color "Adding the layer's public key to deis-node-2..."
vagrant ssh deis-node-2 -c "echo $ssh_key|cat - $authfile > $tmpfile && mv $tmpfile $authfile"

# Add deis-node-1 to the formation
echo_color "Adding node deis-node-2 to formation and provisioning..."
deis nodes:create $formation deis-node-2.local --layer=runtime

echo_color "Done. Now run \"deis create --formation=$formation\" in your app repository."
