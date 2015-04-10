#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

DEBIAN_FRONTEND=noninteractive

apt-get update && apt-get install -yq radosgw radosgw-agent

# cleanup. indicate that python is a required package.
apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*

