#!/usr/bin/env bash
#
# Usage: ./provision-aws-cluster.sh [name]
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

# check for AWS API tools in $PATH
if ! which aws > /dev/null; then
  echo_red 'Please install the AWS command-line tool and ensure it is in your $PATH.'
  exit 1
fi

if [ ! -z "$AWS_CLI_PROFILE" ]; then
    EXTRA_AWS_CLI_ARGS+="--profile $AWS_CLI_PROFILE"
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

# Prepare bailout function to prevent us polluting the namespace
bailout() {
  aws cloudformation delete-stack --stack-name $STACK_NAME
}

# create an AWS cloudformation stack based on CoreOS's default template
aws cloudformation create-stack \
    --template-body "$($THIS_DIR/gen-json.py --channel $COREOS_CHANNEL --version $COREOS_VERSION)" \
    --stack-name $STACK_NAME \
    --parameters "$(<$THIS_DIR/cloudformation.json)" \
    --stack-policy-body "$(<$THIS_DIR/stack_policy.json)" \
    $EXTRA_AWS_CLI_ARGS

# loop until the instances are created
ATTEMPTS=60
SLEEPTIME=10
COUNTER=1
INSTANCE_IDS=""
until [ $(wc -w <<< $INSTANCE_IDS) -eq $DEIS_NUM_INSTANCES -a "$STACK_STATUS" = "CREATE_COMPLETE" ]; do
    if [ $COUNTER -gt $ATTEMPTS ]; then
        echo "Provisioning instances failed (timeout, $(wc -w <<< $INSTANCE_IDS) of $DEIS_NUM_INSTANCES provisioned after $(expr $ATTEMPTS * $SLEEPTIME)s)"
        echo "Destroying stack $STACK_NAME"
        bailout
        exit 1
    fi

    STACK_STATUS=$(aws --output text cloudformation describe-stacks --stack-name $STACK_NAME --query 'Stacks[].StackStatus' $EXTRA_AWS_CLI_ARGS)
    if [ $STACK_STATUS != "CREATE_IN_PROGRESS" -a $STACK_STATUS != "CREATE_COMPLETE" ] ; then
      echo "error creating stack: "
      aws --output text cloudformation describe-stack-events \
          --stack-name $STACK_NAME \
          --query 'StackEvents[?ResourceStatus==`CREATE_FAILED`].[LogicalResourceId,ResourceStatusReason]' \
          $EXTRA_AWS_CLI_ARGS
      bailout
      exit 1
    fi

    INSTANCE_IDS=$(aws ec2 describe-instances \
        --filters Name=tag:aws:cloudformation:stack-name,Values=$STACK_NAME Name=instance-state-name,Values=running \
        --query 'Reservations[].Instances[].[ InstanceId ]' \
        --output text \
        $EXTRA_AWS_CLI_ARGS)

    echo "Waiting for instances to be provisioned ($STACK_STATUS, $(expr $ATTEMPTS + 1 - $COUNTER)0s) ..."
    sleep $SLEEPTIME

    let COUNTER=COUNTER+1
done

# loop until the instances pass health checks
COUNTER=1
INSTANCE_STATUSES=""
until [ `wc -w <<< $INSTANCE_STATUSES` -eq $DEIS_NUM_INSTANCES ]; do
    if [ $COUNTER -gt $ATTEMPTS ];
        then echo "Health checks not passed after $(expr $ATTEMPTS * $SLEEPTIME)s, giving up"
        echo "Destroying stack $STACK_NAME"
        bailout
        exit 1
    fi

    if [ $COUNTER -ne 1 ]; then sleep $SLEEPTIME; fi
    echo "Waiting for instances to pass initial health checks ($(expr $ATTEMPTS + 1 - $COUNTER)0s) ..."
    INSTANCE_STATUSES=$(aws ec2 describe-instance-status \
        --filters Name=instance-status.reachability,Values=passed \
        --instance-ids $INSTANCE_IDS \
        --query 'InstanceStatuses[].[ InstanceId ]' \
        --output text \
        $EXTRA_AWS_CLI_ARGS)
    let COUNTER=COUNTER+1
done

# print instance info
echo "Instances are available:"
aws ec2 describe-instances \
    --filters Name=tag:aws:cloudformation:stack-name,Values=$STACK_NAME Name=instance-state-name,Values=running \
    --query 'Reservations[].Instances[].[InstanceId,PublicIpAddress,InstanceType,Placement.AvailabilityZone,State.Name]' \
    --output text \
    $EXTRA_AWS_CLI_ARGS

# get ELB public DNS name through cloudformation
export ELB_DNS_NAME=$(aws cloudformation describe-stacks \
    --stack-name $STACK_NAME \
    --max-items 1 \
    --query 'Stacks[].[ Outputs[?OutputKey==`DNSName`].OutputValue ]' \
    --output=text \
    $EXTRA_AWS_CLI_ARGS)

# get ELB friendly name through aws elb
ELB_NAME=$(aws elb describe-load-balancers \
    --query "LoadBalancerDescriptions[?DNSName=='$ELB_DNS_NAME'].[ LoadBalancerName ]" \
    --output=text \
    $EXTRA_AWS_CLI_ARGS)
echo "Using ELB $ELB_NAME at $ELB_DNS_NAME"

# Instance launched into a VPC may not have a PublicIPAddress
for ip_type in PublicIpAddress PrivateIpAddress; do
    FIRST_INSTANCE=$(aws ec2 describe-instances \
        --filters Name=tag:aws:cloudformation:stack-name,Values=$STACK_NAME Name=instance-state-name,Values=running \
        --query "Reservations[].Instances[].[$ip_type][0]" \
        --output text \
        $EXTRA_AWS_CLI_ARGS)
    if [[ ! $FIRST_INSTANCE == "None" ]]; then
        break
    fi
done

if [[ $FIRST_INSTANCE == "None" ]]; then
    echo_red "The IP of any of your instances could not be found."
    echo_red "Ensure that the cloudformation stack was created successfully."
    exit 1
fi

export DEISCTL_TUNNEL=$FIRST_INSTANCE

# loop until etcd2 / fleet are up and running
COUNTER=1
until deisctl list >/dev/null; do
    if [ $COUNTER -gt $ATTEMPTS ]; then
        echo_red "Timed out waiting for fleet, giving up"
        echo_red "Ensure that the private key in cloudformation.json"
        echo_red "is added to your ssh-agent."
        break
    fi
    echo "Waiting until fleet is up and running ..."
    sleep 5
    let COUNTER=COUNTER+1
done

echo_green "Enabling proxy protocol"
if ! deisctl config router set proxyProtocol=1; then
    echo_red "# WARNING: Enabling proxy protocol failed."
    echo_red "# Ensure that the private key in cloudformation.json is added to"
    echo_red "# your ssh-agent, then enable proxy protocol before continuing:"
    echo_red "#"
    echo_red "# deisctl config router set proxyProtocol=1\n"
fi

echo_green "Your Deis cluster was deployed to AWS CloudFormation as stack "$STACK_NAME".\n"
echo_green "Now run this command in your shell:"
echo_green "export DEISCTL_TUNNEL=$FIRST_INSTANCE"
echo_green "and continue to follow the documentation for \"Installing the Deis Platform.\""
