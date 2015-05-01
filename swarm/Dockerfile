FROM golang:1.4
WORKDIR /go/src/github.com/docker
RUN git clone https://github.com/smothiki/swarm
WORKDIR /go/src/github.com/docker/swarm
RUN git fetch origin fail
RUN git checkout fail
ENV GOPATH /go/src/github.com/docker/swarm/Godeps/_workspace:$GOPATH
RUN CGO_ENABLED=0 go install -v -a -tags netgo -ldflags "-w -X github.com/docker/swarm/version.GITCOMMIT `git rev-parse --short HEAD`"
WORKDIR /go/src/github.com/deis
RUN git clone https://github.com/smothiki/deis
WORKDIR /go/src/github.com/deis/deis
RUN git fetch origin swarm
RUN git checkout swarm
WORKDIR /go/src/github.com/deis/deis/swarm
RUN mv /go/bin/swarm deis-swarm
ENV GOPATH /go/src/github.com/deis/deis/Godeps/_workspace:$GOPATH
RUN CGO_ENABLED=0 go install -v -a "github.com/deis/deis/swarm"
ENV SWARM_HOST :2375
EXPOSE 2375
VOLUME $HOME/.deis-swarm
ENTRYPOINT ["swarm"]
CMD ["--help"]
