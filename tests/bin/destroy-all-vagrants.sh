#!/bin/sh
#
# WARNING: Don't run this script unless you understand that it will destroy all Deis vagrant VMs.
#
# Tells vagrant to destroy all VMs with "deis" in the name.

VMS=$(vagrant global-status | grep deis | awk '{ print $5 }')
for dir in $VMS; do
    cd "$dir" && vagrant destroy --force
done

# optional commands to remove all VirtualBox vms, since sometimes they are orphaned
#VBoxManage list vms | sed -n -e 's/^.* {\(.*\)}/\1/p' | xargs -L1 -I {} VBoxManage unregistervm {} --delete
#rm -rf $HOME/VirtualBox\ VMs/deis*
