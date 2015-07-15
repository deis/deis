FROM alpine:3.1

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apk add --update-cache curl bash sudo && rm -rf /var/cache/apk/*

# install etcdctl
ENV ETCD_VERSION=2.0.13
RUN curl -sSL https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz | \
    tar zxv -C /usr/local/bin etcd-v${ETCD_VERSION}-linux-amd64/etcdctl && ln -s /usr/local/bin/etcd-v${ETCD_VERSION}-linux-amd64/etcdctl /usr/local/bin && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.10.0/confd-0.10.0-linux-amd64 \
    && chmod +x /usr/local/bin/confd

ADD build.sh /tmp/build.sh
RUN DOCKER_BUILD=true /tmp/build.sh

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
ADD . /app

ENV DEIS_RELEASE 1.9.0-dev
