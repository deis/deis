REPO = deis
MESOS = 0.22.1-1.0.ubuntu1404
MESOS_VERSION = v1.9.0
ZOOKEEPER_VERSION = 3.5.0

MARATHON_VERSION = 0.8.2-RC3

repo_path = github.com/deis/deis/mesos

GO = godep go
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test -cover -race -v
GOVET = $(GO) vet

GO_PACKAGES = pkg/boot pkg/confd pkg/etcd pkg/fleet pkg/log pkg/net
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

build: test-style zookeeper-go build-tools mesos-template mesos-master mesos-slave mesos-marathon zookeeper

mesos-go:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o bin/master-boot pkg/boot/mesos/master/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o bin/slave-boot pkg/boot/mesos/slave/main.go
	go-bindata -pkg bindata -o pkg/boot/mesos/marathon/bindata/bindata.go pkg/boot/mesos/marathon/bash/; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o bin/marathon-boot pkg/boot/mesos/marathon/main.go

mesos-template:
	sed "s/#VERSION#/$(MESOS)/g" template > Dockerfile
	docker build -t $(REPO)/$@:$(MESOS_VERSION) .
	rm -f Dockerfile

mesos-master: mesos-go mesos-template
	sed "s/#VERSION#/$(MESOS_VERSION)/g" master > Dockerfile
	docker build -t $(REPO)/$@:$(MESOS_VERSION) .
	rm -f Dockerfile

mesos-slave: mesos-go mesos-template
	sed "s/#VERSION#/$(MESOS_VERSION)/g" slave > Dockerfile
	docker build -t $(REPO)/$@:$(MESOS_VERSION) .
	rm -f Dockerfile

build-mesos-marathon: mesos-template
	sed "s/#MARATHON_VERSION#/$(MARATHON_VERSION)/;s/#VERSION#/$(MESOS_VERSION)/" build-marathon > Dockerfile
	docker build -t $(REPO)/$@:$(MESOS_VERSION) .
	rm -f Dockerfile

mesos-marathon: build-mesos-marathon mesos-go
	cp marathon Dockerfile
	docker cp `docker create $(REPO)/build-mesos-marathon:$(MESOS_VERSION) /bin/bash`:/marathon/target/marathon-assembly-$(MARATHON_VERSION).jar .
	mv marathon-assembly-$(MARATHON_VERSION).jar marathon-assembly.jar
	sed "s/#MARATHON_VERSION#/$(MARATHON_VERSION)/;s/#VERSION#/$(MESOS_VERSION)/" marathon > Dockerfile
	docker build -t $(REPO)/$@:$(MESOS_VERSION) .
	rm -f Dockerfile

zookeeper: zookeeper-go
	docker build -t $(REPO)/$@:$(ZOOKEEPER_VERSION) zookeeper/.

zookeeper-go:
	echo "Building..."
	go-bindata -pkg bindata -o pkg/boot/zookeeper/bindata/bindata.go pkg/boot/zookeeper/bash/; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o zookeeper/bin/boot pkg/boot/zookeeper/main/boot.go

build-tools:
	echo "Building tools..."

test: mesos-go zookeeper-go
	@$(GOFMT) -timeout 10s $(GO_PACKAGES)

test-style:
	@$(GOFMT) $(GO_PACKAGES)
	@$(GOFMT) $(GO_PACKAGES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(GO_PACKAGES_REPO_PATH)
	@for i in $(addsuffix /...,$(GO_PACKAGES)); do \
		$(GOLINT) $$i; \
	done
