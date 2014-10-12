FROM golang:1.3

WORKDIR /go/src/github.com/deis/deis/publisher
CMD /go/bin/publisher

ADD . /go/src/github.com/deis/deis/publisher
RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/deis/deis/publisher
