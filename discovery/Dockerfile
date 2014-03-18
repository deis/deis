FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# download and extract latest stable etcd
ADD https://github.com/coreos/etcd/releases/download/v0.3.0/etcd-v0.3.0-linux-amd64.tar.gz /tmp/etcd.tar.gz
RUN find /tmp
RUN cd /tmp && \
    tar xfz etcd.tar.gz && \
    mv etcd-v0.3.0-linux-amd64/etcd etcd-v0.3.0-linux-amd64/etcdctl /usr/local/bin && \
    rm -rf etcd-v0.3.0-linux-amd64 etcd.tar.gz

# create etcd daemon user
RUN useradd etcd --home-dir /app --create-home --shell /bin/bash

# add the current build context to /app
ADD . /app
RUN chown -R etcd /app

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 4001 7001
