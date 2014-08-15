FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install go runtime
RUN wget -qO- https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz | tar -C /usr/local -xz

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
