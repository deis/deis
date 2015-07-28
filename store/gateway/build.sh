#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the store dashboard locally"
  echo
  exit 1
fi

# shellcheck disable=SC2034
DEBIAN_FRONTEND=noninteractive

apt-get update && apt-get install -yq radosgw radosgw-agent

# cleanup. indicate that python is a required package.
apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
