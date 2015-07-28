#!/bin/sh
#
# Tells vagrant to halt all running VMs with "deis" in the name.

RUNNING_VMS=$(vagrant global-status | grep deis | grep running | awk '{ print $5 }')
for dir in $RUNNING_VMS; do
    cd "$dir" && vagrant halt
done
