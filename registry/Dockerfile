FROM alpine:3.2

# install common packages
RUN apk add --update-cache curl bash sudo && rm -rf /var/cache/apk/*

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/get-deis/etcdctl-v0.4.9 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.10.0/confd-0.10.0-linux-amd64 \
    && chmod +x /usr/local/bin/confd

ENV DOCKER_REGISTRY_TAG=radosgw DOCKER_REGISTRY_CONFIG=/etc/config.yml

# install registry binaries
RUN export DISTRIBUTION_DIR=/go/src/github.com/docker/distribution \
    && apk add --update-cache git go make py-boto \
    && git clone -b $DOCKER_REGISTRY_TAG --single-branch https://github.com/deis/distribution.git $DISTRIBUTION_DIR \
    && cd $DISTRIBUTION_DIR \
    && export GOPATH=/go:$DISTRIBUTION_DIR/Godeps/_workspace \
    && make binaries \
    && cp bin/* /bin/ \
    && rm -rf /go \
    && apk del --purge git go make \
    && rm -rf /var/cache/apk/*

# define the execution environment
CMD ["/bin/boot"]
EXPOSE 5000

COPY rootfs/ /

ENV DEIS_RELEASE 1.13.0-dev
