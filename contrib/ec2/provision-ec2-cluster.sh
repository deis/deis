#!/usr/bin/env bash
#
# Usage: ./provision-ec2-cluster.sh
#

set -e

###############
#### Variables
###############
THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

source $CONTRIB_DIR/utils.sh

##################
#### Check prereqs
##################

# check for EC2 API tools in $PATH
if ! which aws > /dev/null; then
  echo_red 'Please install the AWS command-line tool and ensure it is in your $PATH.'
  exit 1
fi

if [ -z "$DEIS_NUM_INSTANCES" ]; then
    DEIS_NUM_INSTANCES=3
fi

###############
#### User-Data
###############

DISCOVERY_URL=$(curl -s https://discovery.etcd.io/new)
USER_DATA=$CONTRIB_DIR/coreos/user-data
# if its darwin (osx), (install and) use gnu sed
if [[ `uname` == 'Darwin' ]]; then
  sed  -u "8i\    discovery: ${DISCOVERY_URL}" ${USER_DATA}
else
  sed  -i "8i\    discovery: ${DISCOVERY_URL}" ${USER_DATA}
fi
echo_green "discovery URL configured"

# check that the CoreOS user-data file is valid
$CONTRIB_DIR/util/check-user-data.sh

###############
#### Create AWS
###############

# create an EC2 cloudformation stack based on CoreOS's default template
aws cloudformation create-stack \
    --template-body "$(./gen-json.py)" \
    --stack-name deis \
    --parameters "$(<cloudformation.json)"

#################
#### Confirmation
#################

echo_green "Your Deis cluster has successfully deployed to AWS CloudFormation."
echo_green "Please continue to follow the instructions in the README."
