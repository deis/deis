FROM busybox

ENV DOCKER_HOST unix:///tmp/docker.sock
ENV ROUTESPATH /tmp
CMD ["/bin/logspout"]

ADD logspout /bin/logspout

ENV DEIS_RELEASE 1.6.1
