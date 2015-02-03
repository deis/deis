#!/bin/bash -e
#
# Usage: setup-proxy-protocol [STACK]
# [STACK] -  should be the same as the name you ran provision-ec2-cluster.sh with. (default: deis)
# 
# The script can be executed at any time, as long as the CoreOS cluster is active.


if [[ -z "$1" ]]; then
	STACK_NAME=deis
else
	STACK_NAME=$1
fi

if ! aws cloudformation describe-stacks --stack-name ${STACK_NAME}; then
	echo "Could not find cloudformation stack: ${STACK_NAME}"
	exit 1
fi

LOAD_BALANCER_DNS=$(aws cloudformation describe-stacks --stack-name ${STACK_NAME} | grep -m 1 "OutputValue" | sed 's/ *//g' | cut -d ':' -f2 | cut -d '"' -f2)
LOAD_BALANCER_NAME=$(echo -e ${LOAD_BALANCER_DNS} | cut -d '.' -f1 | cut -d '-' -f1-4)

echo "Enabling Proxy Protocol"
deisctl config router set proxyProtocol=1

echo "Changing port 80 from HTTP -> TCP"
aws elb delete-load-balancer-listeners \
	--load-balancer-name ${LOAD_BALANCER_NAME} \
	--load-balancer-ports 80

aws elb create-load-balancer-listeners \
	--load-balancer-name ${LOAD_BALANCER_NAME} \
	--listeners Protocol=TCP,LoadBalancerPort=80,InstanceProtocol=TCP,InstancePort=80

SSL_CERTIFICATE=$(aws elb describe-load-balancers --load-balancer-names ${LOAD_BALANCER_NAME} | grep -m 1 "SSLCertificateId" | sed 's/ *//g' | cut -d ':' -f2-10 | cut -d '"' -f2)
if [[ -n "${SSL_CERTIFICATE}" ]]; then
	echo "Detected SSL Certificate. Changing port 443 from HTTPS -> SSL"
	aws elb delete-load-balancer-listeners \
		--load-balancer-name ${LOAD_BALANCER_NAME} \
		--load-balancer-ports 443

	aws elb create-load-balancer-listeners \
		--load-balancer-name ${LOAD_BALANCER_NAME} \
		--listeners Protocol=SSL,LoadBalancerPort=443,InstanceProtocol=TCP,InstancePort=80,SSLCertificateId=${SSL_CERTIFICATE}
fi

echo "Setup ELB Proxy Protocol policy"
aws elb create-load-balancer-policy \
	--load-balancer-name ${LOAD_BALANCER_NAME} \
	--policy-name EnableProxyProtocol \
	--policy-type-name ProxyProtocolPolicyType \
	--policy-attributes AttributeName=ProxyProtocol,AttributeValue=true

aws elb set-load-balancer-policies-for-backend-server \
	--load-balancer-name ${LOAD_BALANCER_NAME} \
	--instance-port 80 \
	--policy-names EnableProxyProtocol

echo "Setup TCP Health check"
aws elb configure-health-check \
	--load-balancer-name ${LOAD_BALANCER_NAME} \
	--health-check Target=TCP:80,Interval=15,Timeout=5,UnhealthyThreshold=2,HealthyThreshold=4

