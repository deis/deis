#!/usr/bin/env bash
#
# Preps a test environment and runs `make test-integration`
# against artifacts produced from the current source tree
#

# fail on any command exiting non-zero
set -eo pipefail

# absolute path to current directory
export THIS_DIR=$(cd $(dirname $0); pwd)

# setup the test environment
source $THIS_DIR/test-setup.sh

# AWS credentials required for aws cli and boto
export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID?}
export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY?}

# install python requirements for this script
pip install --disable-pip-version-check awscli boto docopt

function cleanup_aws {
    if [ "$SKIP_CLEANUP" != true ]; then
        log_phase "Cleaning up"
        aws cloudformation delete-stack --stack-name $STACK_NAME
        python $DEIS_ROOT/contrib/aws/route53-wildcard.py delete $DEIS_TEST_DOMAIN $ELB_DNS_NAME
    fi
}

# setup callbacks on process exit and error
trap cleanup_aws EXIT
trap dump_logs ERR

log_phase "Running style tests"

make test-style

log_phase "Running documentation tests"

# test building documentation
make -C docs/ test

log_phase "Running unit tests"

make test-unit

log_phase "Building from current source tree"

# build all docker images and client binaries
make build

# use the built client binaries
export PATH=$DEIS_ROOT/deisctl:$DEIS_ROOT/client/dist:$PATH

log_phase "Running functional tests"

make test-functional

export DEIS_NUM_INSTANCES=3
log_phase "Provisioning $DEIS_NUM_INSTANCES-node CoreOS"

# TODO: don't hardcode --key-names
if ! aws ec2 describe-key-pairs --key-names "deis" >/dev/null ; then
    echo "Importing $DEIS_TEST_AUTH_KEY keypair to EC2"
    aws ec2 import-key-pair --key-name deis \
        --public-key-material file://~/.ssh/$DEIS_TEST_AUTH_KEY.pub \
        --output text
fi

make discovery-url

# customize cloudformation.json to use m3.medium instances
cat > $DEIS_ROOT/contrib/aws/cloudformation.json <<EOF
[
    {
        "ParameterKey":     "KeyPair",
        "ParameterValue":   "deis"
    },
    {
        "ParameterKey":     "InstanceType",
        "ParameterValue":   "m3.medium"
    }
]
EOF

# add random characters after STACK_TAG to avoid collisions
STACK_TAG=${STACK_TAG:-test}-$DEIS_TEST_ID
STACK_NAME=deis-$STACK_TAG
echo "Creating CloudFormation stack $STACK_NAME"
$DEIS_ROOT/contrib/aws/provision-aws-cluster.sh $STACK_NAME

# discard changes to cloudformation.json
git checkout -- $DEIS_ROOT/contrib/aws/cloudformation.json

# use the first cluster node for now
INSTANCE_IDS=$(aws ec2 describe-instances \
    --filters Name=tag:aws:cloudformation:stack-name,Values=$STACK_NAME Name=instance-state-name,Values=running \
    --query 'Reservations[].Instances[].[ InstanceId ]' \
    --output text)
export INSTANCE_ID=$(cut -d " " -f1 <<< $INSTANCE_IDS)
export DEISCTL_TUNNEL=$(aws ec2 describe-instances \
    --instance-ids=$INSTANCE_ID \
    --filters Name=tag:aws:cloudformation:stack-name,Values=$STACK_NAME Name=instance-state-name,Values=running \
    --query 'Reservations[].Instances[].[ PublicDnsName ]' \
    --output text)

log_phase "Waiting for etcd/fleet at $DEISCTL_TUNNEL"

# wait for etcd up to 5 minutes
WAIT_TIME=1
until deisctl --request-timeout=1 list >/dev/null 2>&1; do
   (( WAIT_TIME += 1 ))
   if [ $WAIT_TIME -gt 300 ]; then
    log_phase "Timeout waiting for etcd/fleet"
    # run deisctl one last time without eating the error, so we can see what's up
    deisctl --request-timeout=1 list
    exit 1;
  fi
done

log_phase "etcd available after $WAIT_TIME seconds"

log_phase "Publishing release from source tree"

set +e
trap - ERR

RETRY_COUNT=1

while [ $RETRY_COUNT -le 3 ]; do
  # TODO: detect where IMAGE_PREFIX=deis/ and DEV_REGISTRY=registry.hub.docker.com
  # and disallow it so we can't pollute the production account.
  make dev-release
  RESULT=$?

  if [ $RESULT -ne 0 ]; then
    echo "Docker Hub push failed. Attempt $RETRY_COUNT of 3."
  else
    break
  fi

  (( RETRY_COUNT += 1 ))
done

set -e
trap dump_logs ERR

if [ $RETRY_COUNT -gt 3 ]; then
  echo "Docker Hub push failed the maximum number of times, aborting."
  false
fi

log_phase "Provisioning Deis"

export DEIS_TEST_DOMAIN=$STACK_TAG.$DEIS_TEST_DOMAIN

# configure platform settings
deisctl config platform set domain=$DEIS_TEST_DOMAIN
deisctl config platform set sshPrivateKey=$DEIS_TEST_SSH_KEY

time deisctl install platform
time deisctl start platform

# get ELB public DNS name through cloudformation
# TODO: is "first output value" going to be reliable enough?
ELB_DNS_NAME=$(aws cloudformation describe-stacks \
    --stack-name $STACK_NAME \
    --max-items 1 \
    --query 'Stacks[].[ Outputs[0].[ OutputValue ] ]' \
    --output=text)

# get ELB friendly name through aws elb
ELB_NAME=$(aws elb describe-load-balancers \
    --query 'LoadBalancerDescriptions[].[ DNSName,LoadBalancerName ]' \
    --output=text | grep -F $ELB_DNS_NAME | head -n1 | cut -f2)
echo "Using ELB $ELB_NAME"

# add or update a route53 alias record set to route queries to the ELB
# this python script won't return until the wildcard domain is accessible
python $DEIS_ROOT/contrib/aws/route53-wildcard.py create $DEIS_TEST_DOMAIN $ELB_DNS_NAME

# loop until at least one instance is "in service" with the ELB
ATTEMPTS=45
SLEEPTIME=10
COUNTER=1
IN_SERVICE=0
until [ $IN_SERVICE -ge 1 ]; do
    if [ $COUNTER -gt $ATTEMPTS ]; then exit 1; fi  # timeout after 7 1/2 minutes
    if [ $COUNTER -ne 1 ]; then sleep $SLEEPTIME; fi
    echo "Waiting for ELB to see an instance in service..."
    IN_SERVICE=$(aws elb describe-instance-health \
        --load-balancer-name $ELB_NAME \
        --query 'InstanceStates[].State' \
        | grep InService | wc -l)
done

log_phase "Running integration test suite"

time make test-integration
