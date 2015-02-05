FROM golang:1.3

WORKDIR /go/src/github.com/skynetservices/skydns

ENV CGO_ENABLED 0

RUN go get github.com/tools/godep

RUN git clone https://github.com/skynetservices/skydns .

COPY Godeps /go/src/github.com/skynetservices/skydns/Godeps

RUN godep go build -a -ldflags '-s' && cp skydns /go/bin/skydns
