#!/usr/bin/env bash

# redis to use
REDIS=2.8.17
SHA=913479f9d2a283bfaadd1444e17e7bab560e5d1e

cd /tmp

# build Redis from source
apt-get update && apt-get install -y gcc libc6-dev make net-tools curl --no-install-recommends

curl -sSL http://download.redis.io/releases/redis-$REDIS.tar.gz -o redis.tar.gz

echo "$SHA *redis.tar.gz" | sha1sum -c -
mkdir /usr/src/redis
tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
make -C /usr/src/redis
make -C /usr/src/redis install
