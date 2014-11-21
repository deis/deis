FROM golang:1.3

RUN go get github.com/tools/godep

WORKDIR /go/src/github.com/deis/deis/boot

ADD . /go/src/github.com/deis/deis/boot

RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' github.com/deis/deis/boot

RUN cp /go/src/github.com/deis/deis/boot/boot /tmp/boot