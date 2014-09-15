FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# Envars for git user
ENV GITHOME /home/git
ENV GITUSER git

# install docker-in-docker
# install builder, docker, and hook dependencies
# Regen ssh keys
# install pip
# install deps
# build locales
# setup gituser
# give gituser god for builder
# import cedarish as a tarball
RUN echo "deb http://get.docker.io/ubuntu docker main" > /etc/apt/sources.list.d/docker.list && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9 && \
    apt-get update && apt-get install -yq \
    openssh-server git \
    aufs-tools iptables lxc \
    curl \
    lxc-docker-1.1.2 && \
    rm /etc/ssh/ssh_host_* && \
    dpkg-reconfigure openssh-server && \
    mkdir -p /var/run/sshd && \
    wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python - && \
    pip install pyyaml requests && \
    ln -s /usr/share/i18n/SUPPORTED /var/lib/locales/supported.d/all && locale-gen && \
    useradd -d $GITHOME $GITUSER && \
    mkdir -p $GITHOME/.ssh && chown git:git $GITHOME/.ssh && \
    chown -R $GITUSER:$GITUSER $GITHOME && \
    echo "%git    ALL=(ALL:ALL) NOPASSWD:/home/git/builder" >> /etc/sudoers && \
    wget -O /progrium_cedarish.tar --progress=dot:giga \
        https://s3-us-west-2.amazonaws.com/opdemand/progrium_cedarish_cedar.tar

# define the execution environment
VOLUME /var/lib/docker
WORKDIR /app
ENTRYPOINT ["/app/bin/entry"]
CMD ["/app/bin/boot"]
EXPOSE 22

ADD . /app
ADD sshd_config /etc/ssh/sshd_config
RUN chown -R root:root /app
