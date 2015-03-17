#!/usr/bin/env bash
#
# Usage: ./provision-ec2-cluster.sh [name]
# The [name] is the CloudFormation stack name, and defaults to 'deis'

if [ -z "$1" ]
  then
    STACK_NAME=deis
  else
    STACK_NAME=$1
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
    --template-body "$($THIS_DIR/gen-json.py)" \
    --stack-name $STACK_NAME \
    --parameters "$(<$THIS_DIR/cloudformation.json)"

# loop until the instances are created
ATTEMPTS=45
SLEEPTIME=10
COUNTER=1
INSTANCE_IDS=""
until [ `wc -w <<< $INSTANCE_IDS` -eq $DEIS_NUM_INSTANCES ]; do
    if [ $COUNTER -gt $ATTEMPTS ]; then exit 1; fi  # timeout after 7 1/2 minutes
    if [ $COUNTER -ne 1 ]; then sleep $SLEEPTIME; fi
    echo "Waiting for instances to be created..."
    INSTANCE_IDS=$(aws ec2 describe-instances \
        --filters Name=tag:aws:cloudformation:stack-name,Values=$STACK_NAME Name=instance-state-name,Values=running \
        --query 'Reservations[].Instances[].[ InstanceId ]' \
        --output text)
    let COUNTER=COUNTER+1
done

# loop until the instances pass health checks
COUNTER=1
INSTANCE_STATUSES=""
until [ `wc -w <<< $INSTANCE_STATUSES` -eq $DEIS_NUM_INSTANCES ]; do
    if [ $COUNTER -gt $ATTEMPTS ]; then exit 1; fi  # timeout after 7 1/2 minutes
    if [ $COUNTER -ne 1 ]; then sleep $SLEEPTIME; fi
    echo "Waiting for instances to pass initial health checks..."
    INSTANCE_STATUSES=$(aws ec2 describe-instance-status \
        --filters Name=instance-status.reachability,Values=passed \
        --instance-ids $INSTANCE_IDS \
        --query 'InstanceStatuses[].[ InstanceId ]' \
        --output text)
    let COUNTER=COUNTER+1
done

# print instance info
echo "Instances are available:"
aws ec2 describe-instances \
    --filters Name=tag:aws:cloudformation:stack-name,Values=$STACK_NAME Name=instance-state-name,Values=running \
    --query 'Reservations[].Instances[].[InstanceId,PublicIpAddress,InstanceType,Placement.AvailabilityZone,State.Name]' \
    --output text

# get ELB public DNS name through cloudformation
# TODO: is "first output value" going to be reliable enough?
export ELB_DNS_NAME=$(aws cloudformation describe-stacks \
    --stack-name $STACK_NAME \
    --max-items 1 \
    --query 'Stacks[].[ Outputs[0].[ OutputValue ] ]' \
    --output=text)

# get ELB friendly name through aws elb
ELB_NAME=$(aws elb describe-load-balancers \
    --query 'LoadBalancerDescriptions[].[ DNSName,LoadBalancerName ]' \
    --output=text | grep -F $ELB_DNS_NAME | head -n1 | cut -f2)
echo "Using ELB $ELB_NAME at $ELB_DNS_NAME"

echo_green "Your Deis cluster has been successfully deployed to AWS CloudFormation and is started."
echo_green "Please continue to follow the instructions in the documentation."
