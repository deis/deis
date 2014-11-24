FROM golang:1.3

RUN go get github.com/tools/godep

WORKDIR /go/src/github.com/deis/deis/controller

ADD . /go/src/github.com/deis/deis/controller

RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' github.com/deis/deis/controller

RUN cp controller /go/bin/procfile-introspection