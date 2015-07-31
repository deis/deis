#!/usr/bin/env bash

# fail on any command exiting non-zero
#set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build zoopeeper locally"
  echo
  exit 1
fi

apt-get update && apt-get install -y curl

cd /tmp

mkdir -p /opt

curl -sSL "https://github.com/mesosphere/marathon/archive/v$MARATHON_VERSION.tar.gz" | tar -xzf - -C /opt
ln -s "/opt/marathon-$MARATHON_VERSION" /app
ln -s "/opt/marathon-$MARATHON_VERSION" /marathon

mkdir -p "/opt/marathon-$MARATHON_VERSION/target"
ln -s "/marathon-assembly.jar /opt/marathon-$MARATHON_VERSION/target/marathon-assembly-$MARATHON_VERSION.jar"

apt-get autoremove -y --purge && \
  apt-get clean -y && \
  rm -Rf /usr/share/man /usr/share/doc && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
