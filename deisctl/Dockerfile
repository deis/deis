FROM deis/base
ADD https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz /tmp/
RUN tar -C /usr/local -xzf /tmp/go1.3.linux-amd64.tar.gz
RUN apt-get update && apt-get install -yq git mercurial make
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/go/bin
ENV GOPATH /go
RUN go get github.com/tools/godep
ADD . /go/src/github.com/deis/deis/deisctl
ADD units /tmp/package/var/lib/deis/units
WORKDIR /go/src/github.com/deis/deis/deisctl
RUN godep go install -v ./...
RUN mkdir -p /tmp/package/opt/bin && cp /go/bin/deisctl /tmp/package/opt/bin/deisctl
RUN tar -C /tmp/package -czf /tmp/deisctl.tar.gz .
ENTRYPOINT ["/go/bin/deisctl"]
