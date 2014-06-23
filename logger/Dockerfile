FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install go runtime
RUN wget -O /tmp/go1.2.1.linux-amd64.tar.gz -q https://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf /tmp/go1.2.1.linux-amd64.tar.gz

# prepare go environment
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH /usr/local/bin:/usr/bin:/bin:/sbin:/usr/local/go/bin

# add the current build context to /app
ADD . /app
ADD . /go/src/github.com/deis/deis/logger

# compile the binary
RUN cd /go/src/github.com/deis/deis/logger/syslogd && go install -v .

# create /var/log/deis for holding logs (access via bind mount)
RUN mkdir -p /var/log/deis

# prepare execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 514
