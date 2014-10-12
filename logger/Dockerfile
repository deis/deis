FROM golang:1.3

WORKDIR /go/src/github.com/deis/deis/logger
CMD /go/bin/logger

ADD . /go/src/github.com/deis/deis/logger
RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/deis/deis/logger
