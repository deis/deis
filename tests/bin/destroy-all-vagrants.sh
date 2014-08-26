#!/bin/sh
#
# WARNING: Don't run this script unless you understand that it will destroy all Deis vagrant VMs.
#
# Tells vagrant to destroy all VMs with "deis" in the name.

VMS=$(vagrant global-status | grep deis | awk '{ print $5 }')
for dir in $VMS; do
    cd $dir && vagrant destroy --force
done
