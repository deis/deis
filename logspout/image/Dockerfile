FROM alpine:3.2

ENV DOCKER_HOST unix:///tmp/docker.sock
ENV ROUTESPATH /tmp
CMD ["/bin/logspout"]

ADD logspout /bin/logspout

ENV DEIS_RELEASE 1.12.0
