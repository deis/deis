FROM gliderlabs/alpine:3.1

# install common packages
RUN apk-install curl bash sudo

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-git-b8e693c \
    && chmod +x /usr/local/bin/confd

# install builder, docker, and hook dependencies
<<<<<<< HEAD
RUN apt-get update && apt-get install -yq \
    openssh-server git \
    aufs-tools iptables lxc \
    lxc-docker-1.5.0
=======
RUN apk-install openssh \
    git \
    coreutils \
    iptables \
    docker \
    lxc
>>>>>>> feat(*): use alpine linux.

# configure ssh server
RUN mkdir -p /var/run/sshd && rm -rf /etc/ssh/ssh_host*

# install git and configure gituser
ENV GITHOME /home/git
ENV GITUSER git
RUN adduser -D -h $GITHOME $GITUSER
RUN mkdir -p $GITHOME/.ssh && chown git:git $GITHOME/.ssh
RUN chown -R $GITUSER:$GITUSER $GITHOME

# HACK: import progrium/cedarish as a tarball
# see https://github.com/deis/deis/issues/1027
RUN curl -#SL -o /progrium_cedarish.tar.gz \
    https://s3-us-west-2.amazonaws.com/opdemand/progrium_cedarish_2014_12_15.tar.gz

# define the execution environment
# use VOLUME to remove /var/lib/docker from copy-on-write for performance
# we don't want to stack overlay filesystems
VOLUME /var/lib/docker
WORKDIR /app
ENTRYPOINT ["/app/bin/entry"]
CMD ["/app/bin/boot"]
EXPOSE 22
RUN addgroup -g 2000 slug && adduser -D -u 2000 -G slug slug

# $GITUSER is added to docker group to use docker without sudo and to slug
# group in order to share resources with the slug user
RUN addgroup $GITUSER docker
RUN addgroup $GITUSER slug
RUN passwd -u git

ADD templates/shim.dockerfile /home/git/
ADD etc /etc
ADD . /app
RUN chown -R root:root /app

ENV DEIS_RELEASE 1.5.0-dev

