FROM golang:1.3

ADD . /go/src/github.com/deis/deis/builder

WORKDIR /go/src/github.com/deis/deis/builder

RUN go get github.com/franela/goreq
RUN go get github.com/moraes/config

RUN CGO_ENABLED=0 go build -a -ldflags '-s' extract-domain.go
RUN CGO_ENABLED=0 go build -a -ldflags '-s' extract-types.go
RUN CGO_ENABLED=0 go build -a -ldflags '-s' extract-version.go
RUN CGO_ENABLED=0 go build -a -ldflags '-s' get-app-config.go
RUN CGO_ENABLED=0 go build -a -ldflags '-s' get-app-values.go
RUN CGO_ENABLED=0 go build -a -ldflags '-s' yaml2json-procfile.go
RUN CGO_ENABLED=0 go build -a -ldflags '-s' publish-release-controller.go

RUN mkdir -p /go/bin && cp * /go/bin
