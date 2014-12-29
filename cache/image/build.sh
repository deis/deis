#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the cache locally"
  echo
  exit 1
fi

# redis to use
REDIS=2.8.19
SHA=3e362f4770ac2fdbdce58a5aa951c1967e0facc8

cd /tmp

curl -sSL http://download.redis.io/releases/redis-$REDIS.tar.gz -o redis.tar.gz

echo "$SHA *redis.tar.gz" | sha1sum -c -
mkdir /usr/src/redis
tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
make -C /usr/src/redis
make -C /usr/src/redis install
