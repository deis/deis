FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl && apt-get clean

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json \
    && chmod +x /usr/local/bin/confd

# create /var/log/deis for holding logs (access via bind mount)
RUN mkdir -p /var/log/deis

# prepare execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 514

ADD . /app
