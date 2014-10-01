FROM ubuntu:14.04

ADD . /tmp

RUN /tmp/build.sh

RUN tar -C /usr/local/bin -czf /redis.tgz .