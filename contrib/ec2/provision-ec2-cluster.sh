#!/usr/bin/env bash
#
# Usage: ./provision-ec2-cluster.sh [name]
# The [name] is the CloudFormation stack name, and defaults to 'deis'

if [ -z "$1" ]
  then
    NAME=deis
  else
    NAME=$1
fi

set -e

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

source $CONTRIB_DIR/utils.sh

# check for EC2 API tools in $PATH
if ! which aws > /dev/null; then
  echo_red 'Please install the AWS command-line tool and ensure it is in your $PATH.'
  exit 1
fi

if [ -z "$DEIS_NUM_INSTANCES" ]; then
    DEIS_NUM_INSTANCES=3
fi

# make sure we have all VPC info
if [ -n "$VPC_ID" ]; then
  if [ -z "$VPC_SUBNETS" ] || [ -z "$VPC_ZONES" ]; then
    echo_red 'To provision Deis in a VPC, you must also specify VPC_SUBNETS and VPC_ZONES.'
    exit 1
  fi
fi

# check that the CoreOS user-data file is valid
$CONTRIB_DIR/util/check-user-data.sh

# create an EC2 cloudformation stack based on CoreOS's default template
aws cloudformation create-stack \
    --template-body "$(./gen-json.py)" \
    --stack-name $NAME \
    --parameters "$(<cloudformation.json)"

echo_green "Your Deis cluster has successfully deployed to AWS CloudFormation."
echo_green "Please continue to follow the instructions in the README."
