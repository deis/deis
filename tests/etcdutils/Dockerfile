FROM deis/common:0.12.0-dev

# build etcd and then clean up
RUN buildDeps='curl git-core'; \
    set -x; \
    apt-get update && apt-get install -y $buildDeps --no-install-recommends \
    && curl -sSL https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz | tar -C /usr/local -xz \
    && git clone -q https://github.com/coreos/etcd.git /opt/etcd \
    && cd /opt/etcd && git checkout -q v0.4.6 && PATH=/usr/local/go/bin:$PATH ./build \
    && cp /opt/etcd/bin/etcd /usr/local/bin \
    && cd && rm -rf /usr/local/go /opt/etcd \
    && apt-get purge -y $buildDeps \
    && apt-get autoremove -y && apt-get clean

EXPOSE 4001 7001
CMD ["/usr/local/bin/etcd"]
