FROM flynn/busybox
MAINTAINER Jeff Lindsay <progrium@gmail.com>

ADD ./build/logspout /bin/logspout

ENV DOCKER unix:///tmp/docker.sock
ENV ROUTESPATH /mnt/routes
VOLUME /mnt/routes

EXPOSE 8000

ENTRYPOINT ["/bin/logspout"]
CMD []