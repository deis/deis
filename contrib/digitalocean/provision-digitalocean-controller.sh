#!/usr/bin/env bash
#
# Usage: ./provision-digitalocean-controller.sh <region-id>
#
# Retrieve the region-id by using `knife digital_ocean region list`
#

if [[ -z $1 ]]; then
  echo usage: $0 [region]
  exit 1
fi

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

# connection details for using digital ocean's API
client_id=$(knife exec -E"puts Chef::Config[:knife][:digital_ocean_client_id]")
api_key=$(knife exec -E"puts Chef::Config[:knife][:digital_ocean_api_key]")

# Check that client ID and API key was set
if test -z $client_id; then
  echo "Please add your client id to ${knife_file}."
fi

if test -z $api_key; then
  echo "Please add your api key to ${knife_file}."
fi

#################
# chef settings #
#################
node_name=deis-controller
run_list="recipe[deis::controller]"
chef_version=11.6.2

#########################
# digitalocean settings #
#########################

# the name of the location we want to work with
region_id=$1
# The snapshot that we want to use (deis-base)
image_id=$(knife digital_ocean image list | grep "deis-base" | awk '{print $1}')
# the ID of the size (1GB)
size_id=$(knife digital_ocean size list | grep "2GB" | awk '{print $1}')

if [[ -z $image_id ]]; then
  echo "Can't find saved image \"deis-base\" in region $region_id. Please follow the"
  echo "instructions in prepare-digitalocean-snapshot.sh before provisioning a Deis controller."
  exit 1
fi

if [[ -z $size_id ]]; then
  echo "Cannot find a droplet with the size '2GB' in region $region_id."
  exit 1
fi

################
# SSH settings #
################
key_name=deis-controller
ssh_key_path=~/.ssh/$key_name

# create ssh keypair and store it
if ! test -e $ssh_key_path; then
  echo_color "Creating new SSH key: $key_name"
  set -x
  ssh-keygen -f $ssh_key_path -t rsa -N '' -C "deis-controller" >/dev/null
  curl -X GET \
    --data-urlencode "name=$node_name" \
    --data-urlencode "ssh_pub_key=$(cat $ssh_key_path.pub)" \
    --data-urlencode "client_id=$client_id" \
    --data-urlencode "api_key=$api_key" \
    https://api.digitalocean.com/ssh_keys/new
  ssh-add $ssh_key_path
  set +x
  echo_color "Saved to $ssh_key_path"
else
  echo_color "WARNING: SSH key $ssh_key_path exists, skipping upload"
fi

# get the id of the SSH key that we just uploaded
ssh_key_id=$(knife digital_ocean sshkey list | grep "$key_name" | awk '{print $1}')

# create data bags
knife data bag create deis-users 2>/dev/null
knife data bag create deis-formations 2>/dev/null
knife data bag create deis-apps 2>/dev/null

# trigger digital ocean instance bootstrap
echo_color "Provisioning $node_name with knife digital_ocean..."

set -x
knife digital_ocean droplet create \
  --bootstrap-version $chef_version \
  --server-name $node_name \
  --image $image_id \
  --location $region_id \
  --size $size_id \
  --ssh-keys $ssh_key_id \
  --identity-file $ssh_key_path \
  --bootstrap \
  --run-list $run_list
result=$?
set +x

if [ $result -ne 0 ]; then
  echo_color "Knife botstrap failed."
  # Destroy droplet
  droplet_id=$(knife digital_ocean droplet list | grep deis-controller | awk '{print $1}')
  echo_color "Destroying Droplet $droplet_id..."
  knife digital_ocean droplet destroy -S $droplet_id
  # Remove node and client from Chef Server
  echo_color "Deleting Chef client..."
  knife client delete deis-controller
  echo_color "Deleting Chef node..."
  knife node delete deis-controller
else
  echo_color "Knife bootstrap successful."
  # Need Chef admin permission in order to add and remove nodes and clients
  echo -e "\033[35mPlease ensure that \"deis-controller\" is added to the Chef \"admins\" group.\033[0m"
fi

