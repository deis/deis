FROM alpine:3.1

# install common packages
RUN apk add --update-cache curl tar && rm -rf /var/cache/apk/*

# install etcdctl
ENV ETCD_VERSION=2.0.13
RUN curl -sSL https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz | \
    tar zxv -C /usr/local/bin && ln -s /usr/local/bin/etcd-v${ETCD_VERSION}-linux-amd64/etcd* /usr/local/bin

EXPOSE 2379 2380
CMD ["/usr/local/bin/etcd"]
