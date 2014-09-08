#!/bin/sh

set -e

# determine from whence to download the installer
DEIS_INSTALLER=${DEIS_INSTALLER:-deisctl-0.11.0-linux-amd64.run}
DEIS_BASE_URL=${DEIS_BASE_URL:-https://s3-us-west-2.amazonaws.com/opdemand}
INSTALLER_URL=$DEIS_BASE_URL/$DEIS_INSTALLER

# download the installer archive to /tmp
curl -f -o /tmp/$DEIS_INSTALLER $INSTALLER_URL

# run the installer
sh /tmp/$DEIS_INSTALLER

# clean up after ourselves
rm -f /tmp/$DEIS_INSTALLER
