FROM alpine:3.1

# install common packages
RUN apk add --update-cache curl bash sudo && rm -rf /var/cache/apk/*

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/get-deis/etcdctl-v0.4.9 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.10.0/confd-0.10.0-linux-amd64 \
    && chmod +x /usr/local/bin/confd

RUN apk add --update-cache \
    coreutils \
    device-mapper \
    e2fsprogs \
    git \
    iptables \
    libudev \
    lxc \
    openssh \
    udev \
    util-linux \
    xz \
    && rm -rf /var/cache/apk/*

# the docker package in alpine disables aufs and devicemapper
RUN curl -sSL https://get.docker.com/builds/Linux/x86_64/docker-1.5.0 -o /usr/bin/docker && \
  chmod +x /usr/bin/docker

# configure ssh server
RUN mkdir -p /var/run/sshd && rm -rf /etc/ssh/ssh_host*

# install git and configure gituser
ENV GITHOME /home/git
ENV GITUSER git
RUN adduser -D -h $GITHOME $GITUSER
RUN mkdir -p $GITHOME/.ssh && chown git:git $GITHOME/.ssh
RUN chown -R $GITUSER:$GITUSER $GITHOME

# define the execution environment
# use VOLUME to remove /var/lib/docker from copy-on-write for performance
# we don't want to stack overlay filesystems
VOLUME /var/lib/docker

ENTRYPOINT ["/bin/entry"]
CMD ["/bin/boot"]
EXPOSE 2223
RUN addgroup -g 2000 slug && adduser -D -u 2000 -G slug slug

# $GITUSER is added to docker group to use docker without sudo and to slug
# group in order to share resources with the slug user
RUN addgroup -S docker
RUN addgroup $GITUSER docker
RUN addgroup $GITUSER slug
RUN passwd -u git

COPY . /

ENV DEIS_RELEASE 1.10.0-dev
