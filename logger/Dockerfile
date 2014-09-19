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

# install go runtime
RUN curl -sSL https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz | tar -C /usr/local -xz

# prepare go environment
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH /usr/local/bin:/usr/bin:/bin:/sbin:/usr/local/go/bin

# create /var/log/deis for holding logs (access via bind mount)
RUN mkdir -p /var/log/deis

# prepare execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 514

ADD . /app
ADD . /go/src/github.com/deis/deis/logger

# compile the binary
RUN cd /go/src/github.com/deis/deis/logger/syslogd && go install -v .
