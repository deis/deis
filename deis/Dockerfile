FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

ENV DEBIAN_FRONTEND noninteractive

# install docker in docker deps
RUN apt-get install -yq aufs-tools iptables ca-certificates
RUN echo "deb http://get.docker.io/ubuntu docker main" > /etc/apt/sources.list.d/docker.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
RUN apt-get update -q
RUN apt-get install -yq lxc-docker-0.9.0

# add scripts to /app/bin
RUN mkdir -p /app
ADD bin /app/bin
RUN chown -R root:root /app

# define the execution environment
ENTRYPOINT ["/app/bin/entry"]
CMD ["/app/bin/boot"]