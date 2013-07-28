#!/bin/bash
#
# Distribute AMIs across regions using EC2 CLI tools
#

if [ -z $2 ] ; then
  echo usage: $0 [src-region] [src-ami]
  exit 1
fi

set -ex

src_region=$1
src_ami=$2

# copy the ami to every other region
for region in ap-northeast-1 ap-southeast-1 ap-southeast-2 eu-west-1 sa-east-1 us-east-1 us-west-1 us-west-2; do
  [ $region = $src_region ] && continue
  ec2-copy-image -r $src_region -s $src_ami --region $region
done
