FROM golang:1.3
MAINTAINER OpDemand <info@opdemand.com>

WORKDIR /go/src/github.com/deis/deis/logspout

ENV CGO_ENABLED 0

RUN go get github.com/tools/godep

ADD . /go/src/github.com/deis/deis/logspout

RUN godep go build -a -ldflags '-s' && cp logspout /go/bin/logspout
