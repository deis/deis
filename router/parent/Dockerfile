FROM golang:1.3

ADD build.sh /go/src/github.com/deis/deis/router/build.sh

WORKDIR /go/src/github.com/deis/deis/router

RUN DOCKER_BUILD=true /go/src/github.com/deis/deis/router/build.sh

RUN tar -C /opt/nginx -czf /nginx.tgz .

ADD boot.go /go/src/github.com/deis/deis/router/boot.go
ADD logger/stdout_formatter.go /go/src/github.com/deis/deis/router/logger/stdout_formatter.go

RUN go get

RUN CGO_ENABLED=0 go build -a -ldflags '-s' boot.go

RUN cp /go/src/github.com/deis/deis/router/boot /go/bin/boot
