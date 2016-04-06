FROM ubuntu-debootstrap:14.04

ENV DEBIAN_FRONTEND noninteractive

ADD build.sh /tmp/build.sh

RUN DOCKER_BUILD=true /tmp/build.sh

# Add shared confd configuration
ADD . /app

ENV DEIS_RELEASE 1.14.0-dev

