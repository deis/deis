FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y libgeoip1 curl && apt-get clean

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-git-b8e693c \
    && chmod +x /usr/local/bin/confd

WORKDIR /app

EXPOSE 80 2222

CMD ["/app/bin/boot"]

ADD . /app

ADD nginx.tgz /opt/nginx

RUN rm nginx.tgz
