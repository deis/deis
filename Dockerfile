FROM deis/base
ADD https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz /tmp/
RUN tar -C /usr/local -xzf /tmp/go1.3.linux-amd64.tar.gz
RUN apt-get update && apt-get install -yq git mercurial
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/go/bin
ENV GOPATH /go
ADD . /go/src/github.com/deis/deisctl
ADD systemd /tmp/package/etc/systemd/system
ADD units /tmp/package/var/lib/deis/units
ADD hooks /tmp/package/var/lib/deis/hooks
ADD deis-version /tmp/package/etc/deis-version
WORKDIR /go/src/github.com/deis/deisctl
RUN go get -v ./... && go install -v ./...
RUN mkdir -p /tmp/package/opt/bin && cp /go/bin/deisctl /tmp/package/opt/bin/deisctl
RUN tar -C /tmp/package -czf /tmp/deisctl.tar.gz .
ENTRYPOINT ["/go/bin/deisctl"]
