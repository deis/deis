FROM golang:1.3

WORKDIR /go/src/github.com/deis/deis/pty
CMD /go/bin/pty

RUN go get github.com/tools/godep

ADD . /go/src/github.com/deis/deis/pty

RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' github.com/deis/deis/pty

RUN cp pty /go/bin/pty
