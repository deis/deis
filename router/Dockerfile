FROM alpine:3.2

# install common packages
RUN apk add --update-cache curl bash sudo && rm -rf /var/cache/apk/*

ADD firewall /tmp/firewall

ADD build.sh /tmp/build.sh

WORKDIR /tmp

RUN mkdir -p /opt/nginx && DOCKER_BUILD=true /tmp/build.sh

RUN tar -C /opt/nginx -czf /nginx.tgz .
