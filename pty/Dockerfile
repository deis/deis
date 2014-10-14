FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json \
    && chmod +x /usr/local/bin/confd

# install docker-in-docker
RUN echo "deb http://get.docker.io/ubuntu docker main" > /etc/apt/sources.list.d/docker.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9

# install pty and docker dependencies
RUN apt-get update && apt-get install -yq \
    git mercurial \
    aufs-tools iptables lxc \
    lxc-docker-1.2.0 \
    build-essential

# install go runtime
RUN wget -qO- https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz | tar -C /usr/local -xz

# prepare go environment
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:/go/bin

# prepare execution environment
WORKDIR /go/src/github.com/deis/deis/pty
EXPOSE 3333

ADD . /go/src/github.com/deis/deis/pty

# compile the binary
RUN go build pty.go
CMD ./pty
