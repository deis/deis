FROM deis/go:latest
MAINTAINER OpDemand <info@opdemand.com>

# install source and install
ADD . /go/src/github.com/deis/pty-proxy
RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/deis/pty-proxy
