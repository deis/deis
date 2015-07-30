#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build mesos locally"
  echo
  exit 1
fi

# shellcheck disable=SC2034
DEBIAN_FRONTEND=noninteractive

echo "deb http://repos.mesosphere.io/ubuntu/ trusty main" > /etc/apt/sources.list.d/mesosphere.list

apt-key adv --keyserver keyserver.ubuntu.com --recv E56151BF

apt-get update && \
  apt-get -y install mesos="$MESOS"

apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
