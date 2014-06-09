#!/bin/bash

set -e

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

source $CONTRIB_DIR/utils.sh

if ! which tugboat &>/dev/null; then
	echo_red 'Digital Ocean command line client tugboat not found.'
	exit 1
fi

function wait_ssh () {
	echo -n "Trying to ssh into $1@$2..."
	while ! $SSH -o ConnectTimeout=10 $1@$2 hostname &>/dev/null; do
		echo -n "."
		sleep 1
	done
	echo_green "done"
}

# parse parameters
if [ -z "$1" ]; then
	echo_red 'Usage: $0 SSH_ID [REGION_ID]'
	echo
	echo 'Use `tugboat keys` to list available SSH_IDs.'
	exit 1
fi

SSH_ID="$1"
REGION="${2:-5}" # Amsterdam 1 by default
SIZE="66" # 512 MB
NAME="deis-controller-image-$(date +%Y%m%d%H%M%S)"

BASE_IMAGE='Ubuntu 14.04 x64'
BASE_IMAGE_ID=$(tugboat images -g | grep "$BASE_IMAGE" | sed 's/.*id: \([0-9]*\).*/\1/')
SSH_OPTIONS="-o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o PasswordAuthentication=no -o ConnectTimeout=1 -o ConnectionAttempts=10"
SSH="ssh $SSH_OPTIONS"
SCP="scp $SSH_OPTIONS -q"

# create droplet
tugboat create "$NAME" -i $BASE_IMAGE_ID -p true -k $SSH_ID -s $SIZE -r $REGION

# destroy droplets on error
function cleanup () {
	set +e
	tugboat destroy -c -n "$NAME"
	exit 1
}
trap cleanup ERR
trap cleanup SIGINT

# wait for droplet to come up with ssh login
tugboat wait -n "$NAME" --state active
IP=$(tugboat info -n "$NAME" | egrep '^IP' | awk '{ print $2; }')
wait_ssh root $IP

# bootstrap
echo "Deploying CoreOS on top of $BASE_IMAGE..."
(
	set -ex
	$SCP update-coreos root@$IP:/usr/sbin
	$SSH root@$IP mkdir -p /usr/share/oem/bin
	$SCP cloud-config.yml root@$IP:/usr/share/oem/cloud-config.yml.template
	$SCP ../coreos/user-data root@$IP:/usr/share/oem/user-data.yml
	$SCP create-coreos-docker-store coreos-setup-environment coreos-apply-user-data root@$IP:/usr/share/oem/bin
	$SCP rc.local root@$IP:/etc
	$SSH root@$IP <<EOF
set -xe

DEBIAN_FRONTEND=noninteractive
apt-get install debconf-utils -y
echo "kexec-tools kexec-tools/load_kexec boolean false" | debconf-set-selections
apt-get install squashfs-tools kexec-tools -y
shutdown -h now
EOF
) 2>&1 | sed 's/^/    /'; test ${PIPESTATUS[0]} -eq 0

# switch off and make snapshot
tugboat wait -n "$NAME" --state off
tugboat snapshot $NAME -n $NAME

# wait for snapshot to finish
echo -n "Waiting for snapsnot to appear..."
while ! tugboat images | grep $NAME &>/dev/null; do
	sleep 5
	echo -n "."
done
IMAGE_ID=$(tugboat images | grep $NAME | awk '{print $3;}' | sed 's/,//')
echo_green "done"

echo -n "Trying to destroy droplet..."
while ! tugboat destroy -c -n $NAME &>/dev/null; do
	sleep 5
	echo -n "."
done
echo_green "done"

echo
echo_green "Congratulations: The Deis controller image $NAME was created with ID $IMAGE_ID."
echo
echo "Spawn controllers with:"
echo
echo "  tugboat create deis1 -i $IMAGE_ID -r $REGION -p true -k $SSH_ID -s 65"