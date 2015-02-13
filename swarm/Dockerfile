FROM golang:1.3
WORKDIR /go/src/github.com/deis
RUN git clone https://github.com/smothiki/deis
WORKDIR /go/src/github.com/deis/deis
RUN git fetch origin swarm
RUN git checkout swarm
WORKDIR /go/src/github.com/deis/deis/swarm
ENV GOPATH /go/src/github.com/deis/deis/Godeps/_workspace:$GOPATH
RUN CGO_ENABLED=0 go install -v -a "github.com/deis/deis/swarm"
ENV SWARM_HOST :2375
EXPOSE 2375
RUN wget https://github.com/docker/swarm/releases/download/v0.1.0-rc2/docker-swarm_linux-amd64
RUN mv docker-swarm_linux-amd64 deis-swarm
RUN chmod 777 deis-swarm
VOLUME $HOME/.deis-swarm
ENTRYPOINT ["swarm"]
CMD ["--help"]
