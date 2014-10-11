FROM deis/go:latest

WORKDIR /go/src/github.com/deis/deis/logger
CMD /go/bin/syslogd

ADD . /go/src/github.com/deis/deis/logger
RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/deis/deis/logger
