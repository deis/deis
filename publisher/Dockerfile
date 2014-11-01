FROM golang:1.3

WORKDIR /go/src/github.com/deis/deis/publisher
CMD /go/bin/publisher

RUN go get github.com/tools/godep

ADD . /go/src/github.com/deis/deis/publisher

RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' github.com/deis/deis/publisher 

RUN cp publisher /go/bin/publisher
