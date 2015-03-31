FROM golang:1.4

ADD firewall /tmp/firewall

ADD build.sh /go/src/github.com/deis/deis/router/build.sh

WORKDIR /go/src/github.com/deis/deis/router

RUN DOCKER_BUILD=true /go/src/github.com/deis/deis/router/build.sh

RUN tar -C /opt/nginx -czf /nginx.tgz .

RUN go get github.com/tools/godep

ADD . /go/src/github.com/deis/deis/router/

