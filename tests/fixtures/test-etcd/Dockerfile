FROM alpine:3.1

# install common packages
RUN apk add --update-cache curl tar && rm -rf /var/cache/apk/*

ENV ETCD_VERSION v0.4.9

# install etcd and etcdctl
RUN curl -sSL https://github.com/coreos/etcd/releases/download/$ETCD_VERSION/etcd-$ETCD_VERSION-linux-amd64.tar.gz \
  | tar -vxz -C /usr/local/bin --strip=1

EXPOSE 4001 7001
CMD ["/usr/local/bin/etcd"]
