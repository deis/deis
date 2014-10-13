FROM golang:1.3

ENV CGO_ENABLED 0

ADD . /go/src/github.com/deis/deis/deisctl

ADD systemd /tmp/package/etc/systemd/system
ADD units /tmp/package/var/lib/deis/units
ADD hooks /tmp/package/var/lib/deis/hooks
ADD deis-version /tmp/package/etc/deis-version

WORKDIR /go/src/github.com/deis/deis/deisctl

RUN go get github.com/tools/godep

RUN godep go install -v -a -ldflags '-s' ./...

RUN mkdir -p /tmp/package/opt/bin && cp /go/bin/deisctl /tmp/package/opt/bin/deisctl

RUN tar -C /tmp/package -czf /tmp/deisctl.tar.gz .

ENTRYPOINT ["/go/bin/deisctl"]
