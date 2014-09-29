FROM ubuntu:14.04

ADD . /tmp

RUN /tmp/build.sh

RUN tar -C /nginx -czf /nginx.tgz .
