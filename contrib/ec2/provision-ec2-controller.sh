#!/usr/bin/env bash
#
# Usage: ./provision-ec2-controller.sh <region> <knife params>
#

if [ -z $1 ]; then
  echo usage: $0 [region] [knife params]
  exit 1
fi

function echo_color {
  echo -e "\033[1m$1\033[0m"
}

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname "$THIS_DIR")

# check for Deis' general dependencies
if ! "$CONTRIB_DIR/check-deis-deps.sh"; then
  echo 'Deis is missing some dependencies.'
  exit 1
fi

# check for EC2 API tools in $PATH
if ! which ec2-describe-group > /dev/null; then
  echo 'Please install the EC2 API command-line tools and ensure they are in your $PATH.'
  exit 1
fi

# check for AWS environment variables
: ${AWS_ACCESS_KEY:?'Please set AWS_ACCESS_KEY in your environment for EC2 API access.'}
: ${AWS_SECRET_KEY:?'Please set AWS_SECRET_KEY in your environment for EC2 API access.'}

#################
# chef settings #
#################
node_name=deis-controller-ec2
run_list="recipe[deis::controller]"
chef_version=11.8.2

#######################
# Amazon EC2 settings #
#######################
region=$1
# see contrib/prepare-ubuntu-ami.sh for instructions
# on creating your own deis-optmized AMIs
if [ "$region" == "ap-northeast-1" ]; then
  image=ami-d3225fd2
elif [ "$region" == "ap-southeast-1" ]; then
  image=ami-9e4919cc
elif [ "$region" == "ap-southeast-2" ]; then
  image=ami-afdb4395
elif [ "$region" == "eu-west-1" ]; then
  image=ami-ebeb159c
elif [ "$region" == "sa-east-1" ]; then
  image=ami-cdcc6ed0
elif [ "$region" == "us-east-1" ]; then
  image=ami-15382c7c
elif [ "$region" == "us-west-1" ]; then
  image=ami-f60d33b3
elif [ "$region" == "us-west-2" ]; then
  image=ami-ae85ec9e
else
  echo "Cannot find AMI for region: $region"
  exit 1
fi
flavor="m1.large"
ebs_size=100
sg_name=deis-controller
sg_src=0.0.0.0/0
export EC2_URL=https://ec2.$region.amazonaws.com/

################
# SSH settings #
################
key_name=deis-controller
ssh_key_path=~/.ssh/$key_name
ssh_user="ubuntu"

# create security group and authorize ingress
if ! ec2-describe-group | grep -q "$sg_name"; then
  echo_color "Creating security group: $sg_name"
  set -x
  ec2-create-group $sg_name -d "Created by Deis"
  set +x
  echo_color "Authorizing TCP ports 22,80,443,514,2222,5000,8000 from $sg_src..."
  set -x
  ec2-authorize deis-controller -P tcp -p 22 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 80 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 443 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 514 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 2222 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 5000 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 8000 -s $sg_src >/dev/null
  set +x
else
  echo_color "Security group $sg_name exists"
fi

# create ssh keypair and store it
if ! test -e $ssh_key_path; then
  echo_color "Creating new SSH key: $key_name"
  set -x
  ec2-create-keypair $key_name > $ssh_key_path
  chmod 600 $ssh_key_path
  set +x
  echo_color "Saved to $ssh_key_path"
else
  echo_color "WARNING: SSH key $ssh_key_path exists"
fi

# create data bags
knife data bag create deis-formations 2>/dev/null
knife data bag create deis-apps 2>/dev/null

# trigger ec2 instance bootstrap
echo_color "Provisioning $node_name with knife ec2..."
set -x
knife ec2 server create \
 --bootstrap-version $chef_version \
 --region $region \
 --image $image \
 --flavor $flavor \
 --groups $sg_name \
 --tags Name=$node_name \
 --ssh-key $key_name \
 --ssh-user $ssh_user \
 --identity-file $ssh_key_path \
 --node-name $node_name \
 --ebs-size $ebs_size \
 --run-list $run_list \
 ${@//$region/}
set +x

# Need Chef admin permission in order to add and remove nodes and clients
echo -e "\033[35mPlease ensure that the \"deis-controller-ec2\" client object is added to the Chef \"admins\" group.\033[0m"
