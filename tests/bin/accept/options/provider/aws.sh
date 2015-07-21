# add random characters after STACK_TAG to avoid collisions
STACK_TAG="${STACK_TAG:-test}-${DEIS_TEST_ID}"

export DEIS_NUM_INSTANCES=${DEIS_NUM_INSTANCES:-3}
export DEIS_TEST_DOMAIN="${STACK_TAG}.${DEIS_TEST_DOMAIN}"
export STACK_NAME="${STACK_NAME:-deis-${STACK_TAG}}"

function _setup-provider-dependencies {
  # install python requirements for this script
  pip install --disable-pip-version-check awscli boto docopt
}

function aws-setup-keypair {
  local deis_auth_key="${1}"

  rerun_log "Importing ${deis_auth_key} keypair to EC2"

  # TODO: don't hardcode --key-names
  if ! aws ec2 describe-key-pairs --key-names "deis" >/dev/null ; then
    rerun_log "Importing ${deis_auth_key} keypair to EC2"
    aws ec2 import-key-pair --key-name deis \
        --public-key-material file://~/.ssh/${deis_auth_key}.pub \
        --output text
  fi
}

function aws-provision-cluster {
  local stack_name="${1}"

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

  rerun_log "Provisioning ${DEIS_NUM_INSTANCES}-node CoreOS"

  "${DEIS_ROOT}/contrib/aws/provision-aws-cluster.sh" "${stack_name}"

  # discard changes to cloudformation.json
  git checkout -- "${DEIS_ROOT}/contrib/aws/cloudformation.json"
}

function aws-get-elb-dns-name {
  local stack_name="${1}"

  aws cloudformation describe-stacks \
      --stack-name "${stack_name}" \
      --max-items 1 \
      --query 'Stacks[].[ Outputs[0].[ OutputValue ] ]' \
      --output=text
}

function aws-get-elb-name {
  local elb_dns_name="${1}"

  aws elb describe-load-balancers \
    --query 'LoadBalancerDescriptions[].[ DNSName,LoadBalancerName ]' \
    --output=text | grep -F ${elb_dns_name} | head -n1 | cut -f2
}

function aws-setup-route53 {
  local stack_name="${1}"
  local domain="${2}"

  rerun_log "Setting up Route53 zone..."

  python "${DEIS_ROOT}/contrib/aws/route53-wildcard.py" create "${domain}" "$(aws-get-elb-dns-name ${stack_name})"
}

function aws-destroy-route53 {
  local stack_name="${1}"
  local domain="${2}"

  local elb_dns_name="$(aws-get-elb-dns-name ${stack_name})"

  if [ -n "${elb_dns_name}" ]; then
    rerun_log "Removing Route53 zone..."
    python "${DEIS_ROOT}/contrib/aws/route53-wildcard.py" delete "${domain}" "${elb_dns_name}"
  fi
}

function aws-get-instance-id {
  local stack_name="${1}"

  local instance_ids=$(aws ec2 describe-instances \
      --filters Name=tag:aws:cloudformation:stack-name,Values=${stack_name} Name=instance-state-name,Values=running \
      --query 'Reservations[].Instances[].[ InstanceId ]' \
      --output text)

  cut -d " " -f1 <<< ${instance_ids}
}

function aws-deisctl-tunnel {
  local stack_name="${1}"

  aws ec2 describe-instances \
      --instance-ids=$(aws-get-instance-id ${stack_name}) \
      --filters Name=tag:aws:cloudformation:stack-name,Values=${stack_name} Name=instance-state-name,Values=running \
      --query 'Reservations[].Instances[].[ PublicDnsName ]' \
      --output text
}

function check-elb-service {
  local elb_name="${1}"

  ATTEMPTS=45
  SLEEPTIME=10
  COUNTER=1
  IN_SERVICE=0
  until [ $IN_SERVICE -ge 1 ]; do
      if [ $COUNTER -gt $ATTEMPTS ]; then exit 1; fi  # timeout after 7 1/2 minutes
      if [ $COUNTER -ne 1 ]; then sleep $SLEEPTIME; fi
      rerun_log "Waiting for ELB (${elb_name}) to see an instance in InService..."
      IN_SERVICE=$(aws elb describe-instance-health \
          --load-balancer-name "${elb_name}" \
          --query 'InstanceStates[].State' \
          | grep InService | wc -l)
  done
}

function _create {
  rerun_log "Creating CloudFormation stack ${STACK_NAME}"

  aws-setup-keypair "${DEIS_TEST_AUTH_KEY}"

  aws-provision-cluster "${STACK_NAME}"

  export ELB_DNS_NAME=$(aws-get-elb-dns-name "${STACK_NAME}")
  export ELB_NAME=$(aws-get-elb-name "${ELB_DNS_NAME}")

  aws-setup-route53 "${STACK_NAME}" "${DEIS_TEST_DOMAIN}"

  aws-get-instance-id "${STACK_NAME}"

  export DEISCTL_TUNNEL="$(aws-deisctl-tunnel ${STACK_NAME})"

  rerun_log "DEISCTL_TUNNEL=${DEISCTL_TUNNEL}"
}

function _destroy {
  rerun_log "Attempting to destroy ${STACK_NAME}..."

  aws cloudformation delete-stack --stack-name "${STACK_NAME}"

  aws-destroy-route53 "${STACK_NAME}" "${DEIS_TEST_DOMAIN}"
}

function _check-cluster {
  check-elb-service "${ELB_NAME}"
}
