FROM progrium/cedarish:latest

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl net-tools sudo

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json \
    && chmod +x /usr/local/bin/confd

# install herokuish
RUN curl -sSL https://github.com/aledbf/herokuish/releases/download/v0.1.2/herokuish_0.1.2_linux_x86_64.tgz \
    | tar -xzC /usr/local/bin && chmod 755 /usr/local/bin/herokuish

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

# install git and configure gituser
ENV GITHOME /home/git
ENV GITUSER git
RUN useradd -d $GITHOME $GITUSER
RUN mkdir -p $GITHOME/.ssh && chown git:git $GITHOME/.ssh
RUN chown -R $GITUSER:$GITUSER $GITHOME

# let the git user run `sudo /home/git/builder` (not writeable)
RUN echo "%git    ALL=(ALL:ALL) NOPASSWD:/home/git/builder" >> /etc/sudoers

WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 22
RUN addgroup --quiet --gid 2000 slug && useradd slug --uid=2000 --gid=2000

ADD templates/shim.dockerfile /home/git/
ADD etc /etc
ADD . /app
RUN chown -R root:root /app

RUN /app/bin/install-buildpacks.sh

ENV DEIS_RELEASE 1.3.0-dev

