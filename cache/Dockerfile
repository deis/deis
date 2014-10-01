FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# build Redis from source
RUN buildDeps='gcc libc6-dev make'; \
    set -x; \
    apt-get update && apt-get install -y $buildDeps net-tools --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /usr/src/redis \
    && curl -sSL http://download.redis.io/releases/redis-2.8.15.tar.gz -o redis.tar.gz \
    && echo "afc0d753cea68a26038775df2dea75a76e3d0e1d *redis.tar.gz" | sha1sum -c - \
    && tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
    && make -C /usr/src/redis \
    && make -C /usr/src/redis install \
    && rm -r redis.tar.gz /usr/src/redis \
    && apt-get purge -y $buildDeps \
    && apt-get autoremove -y \
    && apt-get clean

WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 6379
ADD . /app
ADD redis.conf /etc/redis/redis.conf
