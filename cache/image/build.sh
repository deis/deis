#!/usr/bin/env bash

if [[ -z $DOCKER_BUILD ]]; then
  echo 
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the cache locally"
  echo 
  exit 1
fi

# redis to use
REDIS=2.8.17
SHA=913479f9d2a283bfaadd1444e17e7bab560e5d1e

cd /tmp

curl -sSL http://download.redis.io/releases/redis-$REDIS.tar.gz -o redis.tar.gz

echo "$SHA *redis.tar.gz" | sha1sum -c -
mkdir /usr/src/redis
tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
make -C /usr/src/redis
make -C /usr/src/redis install
