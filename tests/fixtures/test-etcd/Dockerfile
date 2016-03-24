FROM alpine:3.3

# install common packages
RUN apk add --no-cache curl tar

# ETCD_VERSION is actually used by the etcd daemon, and causes an issue if we
# format it for our use here. So, we call this something else.
ENV INSTALL_ETCD_VERSION v2.1.2

# install etcd and etcdctl
RUN curl -sSL https://github.com/coreos/etcd/releases/download/$INSTALL_ETCD_VERSION/etcd-$INSTALL_ETCD_VERSION-linux-amd64.tar.gz \
  | tar -vxz -C /usr/local/bin --strip=1

EXPOSE 4001 7001 2379 2380
CMD ["/usr/local/bin/etcd"]
