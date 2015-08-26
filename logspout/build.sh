#!/bin/sh
set -e

apk add --update go git mercurial

mkdir -p /go/src/github.com/gliderlabs
cd /go/src/github.com/gliderlabs
git clone https://github.com/gliderlabs/logspout.git
cd logspout
mv /go/src/github.com/deis/deis/logspout/modules.go .

export GOPATH=/go
go get
go build -ldflags "-X main.Version $1" -o /bin/logspout

apk del go git mercurial
rm -rf /go
rm -rf /var/cache/apk/*

# backwards compatibility
ln -fs /tmp/docker.sock /var/run/docker.sock
