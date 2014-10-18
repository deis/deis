FROM golang:1.3

ADD . /go/src/github.com/deis/deis/boot

RUN DOCKER_BUILD=true /go/src/github.com/deis/deis/boot/build.sh

WORKDIR /go/src/github.com/deis/deis/boot

RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/deis/deis/boot
