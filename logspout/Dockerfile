FROM alpine:3.1
ENTRYPOINT ["/bin/logspout"]
VOLUME /mnt/routes
EXPOSE 8000

ENV DEIS_RELEASE 1.10.0-dev

RUN mkdir -p /go/src/github.com/deis/deis/logspout
COPY . /go/src/github.com/deis/deis/logspout

RUN cd /go/src/github.com/deis/deis/logspout && ./build.sh $DEIS_RELEASE
