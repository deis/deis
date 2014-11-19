FROM golang:1.3

# install dependencies
RUN go get github.com/tools/godep

RUN mkdir -p /go/bin

COPY . /go/src/github.com/deis/deis/builder

WORKDIR /go/src/github.com/deis/deis/builder

RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/extract-domain.go && mv extract-domain /go/bin/
RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/extract-types.go && mv extract-types /go/bin/
RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/extract-version.go && mv extract-version /go/bin/
RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/generate-buildhook.go && mv generate-buildhook /go/bin/
RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/get-app-config.go && mv get-app-config /go/bin/
RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/get-app-values.go && mv get-app-values /go/bin/
RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/yaml2json-procfile.go && mv yaml2json-procfile /go/bin/
RUN CGO_ENABLED=0 godep go build -a -ldflags '-s' bin/publish-release-controller.go && mv publish-release-controller /go/bin/
