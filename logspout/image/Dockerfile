FROM flynn/busybox
MAINTAINER OpDemand <info@opdemand.com>

ENV DOCKER_HOST unix:///tmp/docker.sock
ENV ROUTESPATH /tmp
CMD ["/bin/logspout"]

ADD logspout /bin/logspout

ENV DEIS_RELEASE 1.4.0-dev

