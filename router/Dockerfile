FROM golang:1.4

ADD firewall /tmp/firewall

WORKDIR /go/src/github.com/deis/deis/router

ADD build.sh /go/src/github.com/deis/deis/router/build.sh
RUN /go/src/github.com/deis/deis/router/build.sh

RUN tar -C /opt/nginx -czf /nginx.tgz .

RUN go get github.com/tools/godep

ADD . /go/src/github.com/deis/deis/router/

RUN CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' boot.go

RUN cp /go/src/github.com/deis/deis/router/boot /go/bin/boot
