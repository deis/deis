#!/bin/sh

if [ -z $1 ]; then
  echo usage: $0 [region]
  exit 1
fi

region=$1

# see contrib/prepare-ubuntu-ami.sh for instructions
# on creating your own deis-optmized AMIs
if [ "$region" == "ap-northeast-1" ]; then
  image=ami-a57aeca4
elif [ "$region" == "ap-southeast-1" ]; then
  image=ami-e03a72b2
elif [ "$region" == "ap-southeast-2" ]; then
  image=ami-bd801287
elif [ "$region" == "eu-west-1" ]; then
  image=ami-d9d3cdad
elif [ "$region" == "sa-east-1" ]; then
  image=ami-a7df7bba
elif [ "$region" == "us-east-1" ]; then
  image=ami-e85a2081
elif [ "$region" == "us-west-1" ]; then
  image=ami-ac6942e9
elif [ "$region" == "us-west-2" ]; then
  image=ami-b55ac885
else
  echo "Cannot find AMI for region: $region"
  exit 1
fi

# ec2 settings
flavor="m1.large"
ebs_size=100
sg_name=deis-controller
sg_src=0.0.0.0/0
key_name=deis-controller
export EC2_URL=https://ec2.$region.amazonaws.com/

# ssh settings
ssh_key_path=~/.ssh/$key_name
ssh_user="ubuntu"

# chef settings
node_name="deis-controller"
run_list="recipe[deis],recipe[deis::postgresql],recipe[deis::server],recipe[deis::gitosis],recipe[deis::build]"
chef_version=11.4.4

function echo_color {
  echo "\033[1m$1\033[0m"
}

# create security group and authorize ingress
if ! ec2-describe-group | grep -q "$sg_name"; then
  echo_color "Creating security group: $sg_name"
  set -x
  ec2-create-group $sg_name -d "Created by Deis"
  set +x
  echo_color "Authorizing TCP ports 22,80,443 from $sg_src..."
  set -x
  ec2-authorize deis-controller -P tcp -p 22 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 80 -s $sg_src >/dev/null
  ec2-authorize deis-controller -P tcp -p 443 -s $sg_src >/dev/null
  set +x
else
  echo_color "Security group $sg_name exists"
fi

# create ssh keypair and store it
if ! test -e $ssh_key_path; then
  echo_color "Creating new SSH key: $key_name"
  set -x
  ec2-create-keypair $key_name > $ssh_key_path
  chmod 600 $ssh_key_path
  set +x
  echo "Saved to $ssh_key_path"
else
  echo_color "SSH key $ssh_key_path exists"
fi

# create data bags
knife data bag create deis-build 2>/dev/null
knife data bag create deis-formations 2>/dev/null

# create data bag item using a temp file
tempfile=$(mktemp -t deis)
mv $tempfile $tempfile.json
cat > $tempfile.json <<EOF
{ "id": "gitosis", "ssh_keys": {}, "formations": {} }
EOF
knife data bag from file deis-build $tempfile.json
rm -f $tempfile.json

# trigger ec2 instance bootstrap
echo_color "Provisioning $node_name with knife ec2..."
set -x
knife ec2 server create \
 --bootstrap-version $chef_version \
 --region $region \
 --image $image \
 --flavor $flavor \
 --groups $sg_name \
 --tags Name=$node_name \
 --ssh-key $key_name \
 --ssh-user $ssh_user \
 --identity-file $ssh_key_path \
 --node-name $node_name \
 --ebs-size $ebs_size \
 --run-list $run_list
set +x
