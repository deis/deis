FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install go runtime
RUN wget -qO- https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz | tar -C /usr/local -xz

# prepare go environment
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH /usr/local/bin:/usr/bin:/bin:/sbin:/usr/local/go/bin

# install dependencies
RUN apt-get update -y && apt-get install -yq \
	build-essential \
    git \
	mercurial \
  lxc-docker-1.1.2 \
	--no-install-recommends

# prepare execution environment
WORKDIR /go/src/github.com/deis/deis/controller/pty
ENV PORT 5309 
EXPOSE 5309
CMD ["/bin/boot"]

ADD bin/boot /bin/boot
ADD . /go/src/github.com/deis/deis/pty

# compile the binary
RUN go get .
