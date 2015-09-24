#!/usr/bin/env bash
#
# Usage: ./update-aws-cluster.sh [name]
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

# check for AWS API tools in $PATH
if ! which aws > /dev/null; then
  echo_red 'Please install the AWS command-line tool and ensure it is in your $PATH.'
  exit 1
fi

if [ ! -z "$AWS_CLI_PROFILE" ]; then
    EXTRA_AWS_CLI_ARGS+="--profile $AWS_CLI_PROFILE"
fi

# check that the CoreOS user-data file is valid
$CONTRIB_DIR/util/check-user-data.sh

# update the AWS CloudFormation stack
aws cloudformation update-stack \
    --template-body "$($THIS_DIR/gen-json.py --channel $COREOS_CHANNEL --version $COREOS_VERSION)" \
    --stack-name $NAME \
    --parameters "$(<$THIS_DIR/cloudformation.json)" \
    --stack-policy-body "$(<$THIS_DIR/stack_policy.json)" \
    $EXTRA_AWS_CLI_ARGS

echo_green "Your Deis cluster on AWS CloudFormation has been successfully updated."
