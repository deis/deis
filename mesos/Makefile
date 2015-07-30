include ../includes.mk

REPO = deis
MESOS = 0.22.1-1.0.ubuntu1404
ZOOKEEPER_VERSION = 3.5.0
MARATHON_VERSION = 0.8.2-RC3

repo_path = github.com/deis/deis/mesos

GO = godep go
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test -cover -race -v
GOVET = $(GO) vet

COMPONENT = $(notdir $(repo_path))
GO_PACKAGES = pkg/boot pkg/confd pkg/etcd pkg/fleet pkg/log pkg/net pkg/os pkg/types
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

SHELL_SCRIPTS = $(shell find "." -name '*.sh')

MASTER_IMAGE = $(IMAGE_PREFIX)mesos-master:$(BUILD_TAG)
MASTER_DEV_IMAGE = $(REGISTRY)$(MASTER_IMAGE)
MARATHON_IMAGE = $(IMAGE_PREFIX)mesos-marathon:$(BUILD_TAG)
MARATHON_DEV_IMAGE = $(REGISTRY)$(MARATHON_IMAGE)
SLAVE_IMAGE = $(IMAGE_PREFIX)mesos-slave:$(BUILD_TAG)
SLAVE_DEV_IMAGE = $(REGISTRY)$(SLAVE_IMAGE)
ZOOKEEPER_IMAGE = $(IMAGE_PREFIX)zookeeper:$(BUILD_TAG)
ZOOKEEPER_DEV_IMAGE = $(REGISTRY)$(ZOOKEEPER_IMAGE)

build: test-style zookeeper-go mesos-template mesos-master mesos-slave mesos-marathon zookeeper

install: check-deisctl
	deisctl install $(COMPONENT)

uninstall: check-deisctl
	deisctl uninstall $(COMPONENT)

start: check-deisctl
	deisctl start $(COMPONENT)

stop: check-deisctl
	deisctl stop $(COMPONENT)

restart: stop start

run: install start

dev-release: push set-image

push: check-registry
	docker tag -f $(MASTER_IMAGE) $(MASTER_DEV_IMAGE)
	docker push $(MASTER_DEV_IMAGE)
	docker tag -f $(SLAVE_IMAGE) $(SLAVE_DEV_IMAGE)
	docker push $(SLAVE_DEV_IMAGE)
	docker tag -f $(MARATHON_IMAGE) $(MARATHON_DEV_IMAGE)
	docker push $(MARATHON_DEV_IMAGE)
	docker tag -f $(ZOOKEEPER_IMAGE) $(ZOOKEEPER_DEV_IMAGE)
	docker push $(ZOOKEEPER_DEV_IMAGE)

set-image: check-deisctl
	deisctl config mesos-master set image=$(MASTER_DEV_IMAGE)
	deisctl config mesos-slave set image=$(SLAVE_DEV_IMAGE)
	deisctl config mesos-marathon set image=$(MARATHON_DEV_IMAGE)
	deisctl config zookeeper set image=$(ZOOKEEPER_DEV_IMAGE)

release:
	docker push $(MASTER_IMAGE)
	docker push $(SLAVE_IMAGE)
	docker push $(MARATHON_IMAGE)
	docker push $(ZOOKEEPER_IMAGE)

deploy: build dev-release restart

setup-gotools:
	go get -u github.com/jteeuwen/go-bindata/...

mesos-go: setup-gotools
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o bin/master-boot pkg/boot/mesos/master/main.go
	@$(call check-static-binary,bin/master-boot)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o bin/slave-boot pkg/boot/mesos/slave/main.go
	@$(call check-static-binary,bin/slave-boot)
	go-bindata -pkg bindata -o bindata/marathon/bindata.go pkg/boot/mesos/marathon/bash/; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o bin/marathon-boot pkg/boot/mesos/marathon/main.go
	@$(call check-static-binary,bin/marathon-boot)

mesos-template:
	sed "s/#VERSION#/$(MESOS)/g" template > Dockerfile
	docker build -t $(IMAGE_PREFIX)$@:$(BUILD_TAG) .
	rm -f Dockerfile

mesos-master: mesos-go mesos-template
	sed "s@#PREFIX#@$(IMAGE_PREFIX)@;s/#VERSION#/$(BUILD_TAG)/g" master > Dockerfile
	docker build -t $(IMAGE_PREFIX)$@:$(BUILD_TAG) .
	rm -f Dockerfile

mesos-slave: mesos-go mesos-template
	sed "s@#PREFIX#@$(IMAGE_PREFIX)@;s/#VERSION#/$(BUILD_TAG)/g" slave > Dockerfile
	docker build -t $(IMAGE_PREFIX)$@:$(BUILD_TAG) .
	rm -f Dockerfile

build-mesos-marathon: mesos-template
	sed "s/#MARATHON_VERSION#/$(MARATHON_VERSION)/;s@#PREFIX#@$(IMAGE_PREFIX)@;s/#VERSION#/$(BUILD_TAG)/g" build-marathon > Dockerfile
	docker build -t $(IMAGE_PREFIX)$@:$(BUILD_TAG) .
	rm -f Dockerfile

mesos-marathon: mesos-go build-mesos-marathon
	cp marathon Dockerfile
	docker cp `docker create $(IMAGE_PREFIX)build-mesos-marathon:$(BUILD_TAG) /bin/bash`:/marathon/target/marathon-assembly-$(MARATHON_VERSION).jar .
	mv marathon-assembly-$(MARATHON_VERSION).jar marathon-assembly.jar
	sed "s/#MARATHON_VERSION#/$(MARATHON_VERSION)/;s@#PREFIX#@$(IMAGE_PREFIX)@;s/#VERSION#/$(BUILD_TAG)/" marathon > Dockerfile
	docker build -t $(IMAGE_PREFIX)$@:$(BUILD_TAG) .
	rm -f Dockerfile

zookeeper: zookeeper-go
	docker build -t $(IMAGE_PREFIX)$@:$(BUILD_TAG) zookeeper/.

zookeeper-go:
	echo "Building..."
	go-bindata -pkg bindata -o bindata/zookeeper/bindata.go pkg/boot/zookeeper/bash/; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o zookeeper/bin/boot pkg/boot/zookeeper/main/boot.go
	@$(call check-static-binary,zookeeper/bin/boot)

test: mesos-go zookeeper-go
	@$(GOFMT) -timeout 10s $(GO_PACKAGES)

test-style:
	@$(GOFMT) $(GO_PACKAGES)
	@$(GOFMT) $(GO_PACKAGES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(GO_PACKAGES_REPO_PATH)
	@for i in $(addsuffix /...,$(GO_PACKAGES)); do \
		$(GOLINT) $$i; \
	done
	shellcheck $(SHELL_SCRIPTS)
