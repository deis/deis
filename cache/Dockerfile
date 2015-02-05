FROM ubuntu-debootstrap:14.04

# install deps
RUN apt-get update && apt-get install -yq curl make gcc

ENV REDIS 2.8.19
ENV SHA 3e362f4770ac2fdbdce58a5aa951c1967e0facc8

# install redis-server to /usr/local/bin
RUN mkdir -p /usr/src/redis
RUN curl -sSL http://download.redis.io/releases/redis-$REDIS.tar.gz | tar -xzC /usr/src/redis --strip-components=1 && \
    cd /usr/src/redis && \
    make && \
    make install
