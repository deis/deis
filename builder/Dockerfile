FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update

# install ssh server
RUN apt-get install -yq openssh-server
RUN rm /etc/ssh/ssh_host_*
RUN dpkg-reconfigure openssh-server
RUN mkdir -p /var/run/sshd

# install docker in docker deps
RUN apt-get install -yq aufs-tools iptables ca-certificates lxc
RUN echo "deb http://get.docker.io/ubuntu docker main" > /etc/apt/sources.list.d/docker.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
RUN apt-get update -qy
RUN apt-get install -yq lxc-docker-1.0.0

# 45 seconds ssh client keepalive, for intermediary firewalls/loadbalancers.
RUN echo ClientAliveInterval 45 >> /etc/ssh/sshd_config

# install recent pip
RUN wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python -

# install hook dependencies
RUN pip install pyyaml requests

# install hook utilities
RUN apt-get install -yq curl vim

# install all i18n locales
RUN ln -s /usr/share/i18n/SUPPORTED /var/lib/locales/supported.d/all && locale-gen

# install git and configure gituser
ENV GITHOME /home/git
ENV GITUSER git
RUN apt-get install -yq git
RUN useradd -d $GITHOME $GITUSER
RUN mkdir -p $GITHOME/.ssh && chown git:git $GITHOME/.ssh
RUN chown -R $GITUSER:$GITUSER $GITHOME

# let the git user run `sudo /home/git/builder` (not writeable)
RUN apt-get install -yq sudo
RUN echo "%git    ALL=(ALL:ALL) NOPASSWD:/home/git/builder" >> /etc/sudoers

# add the current build context to /app
ADD . /app
RUN chown -R root:root /app

# define the execution environment
VOLUME /var/lib/docker
WORKDIR /app
ENTRYPOINT ["/app/bin/entry"]
CMD ["/app/bin/boot"]
EXPOSE 22
