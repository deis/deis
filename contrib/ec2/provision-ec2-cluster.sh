#!/usr/bin/env bash
#
# Usage: ./provision-ec2-cluster.sh
#

set -e

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

source $CONTRIB_DIR/utils.sh

# check for EC2 API tools in $PATH
if ! which aws > /dev/null; then
  echo_red 'Please install the AWS command-line tool and ensure it is in your $PATH.'
  exit 1
fi

# create an EC2 cloudformation stack based on CoreOS's default template
aws cloudformation create-stack \
    --template-url https://s3.amazonaws.com/coreos.com/dist/aws/coreos-alpha.template \
    --stack-name deis \
    --parameters "$(<cloudformation.json)"

echo_green "Your Deis cluster has successfully deployed to AWS CloudFormation."
echo_green "Please continue to follow the instructions in the README."
