FROM ubuntu:12.04
MAINTAINER Matt Boersma <matt@opdemand.com>

# Install prerequisites: wget, git, and go.
RUN apt-get update -q && DEBIAN_FRONTEND=noninteractive apt-get install -qy wget git-core
RUN wget -qO- https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH /usr/local/go/bin:$PATH

# Build etcd v0.4.5. Keep this in sync with Deis' version.
RUN git clone -q https://github.com/coreos/etcd.git /opt/etcd
RUN cd /opt/etcd && git checkout -q v0.4.5 && ./build

# Download latest stable etcdctl. Keep this in sync with Deis' version.
ADD https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.5 /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

EXPOSE 4001 7001
CMD ["/opt/etcd/bin/etcd"]
