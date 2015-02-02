FROM ubuntu-debootstrap:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl net-tools sudo

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json \
    && chmod +x /usr/local/bin/confd

# install docker-in-docker
RUN echo "deb http://get.docker.io/ubuntu docker main" > /etc/apt/sources.list.d/docker.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9

# install builder, docker, and hook dependencies
RUN apt-get update && apt-get install -yq \
    openssh-server git \
    aufs-tools iptables lxc \
    lxc-docker-1.3.3

# configure ssh server
RUN rm /etc/ssh/ssh_host_*
RUN dpkg-reconfigure openssh-server
RUN mkdir -p /var/run/sshd

# configure locale
RUN echo LANG="en_US.UTF-8" > /etc/default/locale && dpkg-reconfigure locales

# install git and configure gituser
ENV GITHOME /home/git
ENV GITUSER git
RUN useradd -d $GITHOME $GITUSER
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
RUN addgroup --quiet --gid 2000 slug && useradd slug --uid=2000 --gid=2000

# $GITUSER is added to docker group to use docker without sudo and to slug
# group in order to share resources with the slug user
RUN usermod -a -G docker $GITUSER
RUN usermod -a -G slug $GITUSER

ADD templates/shim.dockerfile /home/git/
ADD etc /etc
ADD . /app
RUN chown -R root:root /app

ENV DEIS_RELEASE 1.4.0-dev
