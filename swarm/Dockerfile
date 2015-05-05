FROM golang:1.4

WORKDIR /go/src/github.com/docker
RUN git clone https://github.com/deis/swarm
WORKDIR /go/src/github.com/docker/swarm
RUN git fetch origin nodefailover
RUN git checkout nodefailover

ENV GOPATH /go/src/github.com/docker/swarm/Godeps/_workspace:$GOPATH
RUN CGO_ENABLED=0 go build -a -tags netgo -installsuffix cgo \
    -ldflags "-s -w -X github.com/docker/swarm/version.GITCOMMIT `git rev-parse --short HEAD`" \
    -o deis-swarm
