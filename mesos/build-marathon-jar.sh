#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build marathon locally"
  echo
  exit 1
fi

# shellcheck disable=SC2034
DEBIAN_FRONTEND=noninteractive

apt-get update && apt-get install --no-install-recommends -y \
  openjdk-7-jdk \
  scala \
  curl

curl -SsL -O http://dl.bintray.com/sbt/debian/sbt-0.13.5.deb && \
  dpkg -i sbt-0.13.5.deb

curl -sSL "https://github.com/mesosphere/marathon/archive/v$MARATHON_VERSION.tar.gz" | tar -xzf - -C /opt
ln -s "/opt/marathon-$MARATHON_VERSION" /app
ln -s "/opt/marathon-$MARATHON_VERSION" /marathon

cd /app

# Word splitting wanted in the situation.
# shellcheck disable=SC2046
sbt assembly && \
  mv $(find target -name 'marathon-assembly-*.jar' | sort | tail -1) ./ && \
  rm -rf target/* ~/.sbt ~/.ivy2 && \
  mv marathon-assembly-*.jar target

# cleanup. indicate that python, libpq and libyanl are required packages.
apt-get clean -y && \
  rm -rf /tmp/* /var/tmp/* && \
  rm -rf /var/lib/apt/lists/*
