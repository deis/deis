FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y curl && apt-get clean

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 6379
ADD . /app
ADD redis.conf /etc/redis/redis.conf

ADD redis.tgz /usr/local/bin
