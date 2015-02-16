FROM gliderlabs/alpine:3.1

# install common packages
RUN apk-install curl bash

ADD firewall /tmp/firewall

ADD build.sh /tmp/build.sh

WORKDIR /tmp

RUN mkdir -p /opt/nginx && DOCKER_BUILD=true /tmp/build.sh

RUN tar -C /opt/nginx -czf /nginx.tgz .
