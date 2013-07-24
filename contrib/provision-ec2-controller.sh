#!/bin/sh

# ec2 settings
region="us-west-2"
image="ami-bf41d28f"
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
run_list="recipe[deis::default],recipe[deis::gitosis],recipe[deis::build],recipe[deis::postgresql],recipe[deis::server]"
chef_version=11.4.4

function echo_color {
  echo "\033[1m$1\033[0m"
}

# create security group and authorize ingress
if ! ec2-describe-group | grep -q "$sg_name"; then
  echo_color "Creating security group: $sg_name"
  set -x
  ec2-create-group $sg_name -d "Managed by Deis"
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
  echo "Creating new SSH key: $key_name"
  set -x
  ec2-create-keypair $key_name > $ssh_key_path
  chmod 600 $ssh_key_path
  set +x
  echo "Saved to $ssh_key_path"
else
  echo_color "SSH key $ssh_key_path exists"
fi

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
