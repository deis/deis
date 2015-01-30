include ../includes.mk

COMPONENT = logger
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
BUILD_IMAGE := $(COMPONENT)-build
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)
BINARY_DEST_DIR = image/bin

build: check-docker
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -v -ldflags '-s' -o $(BINARY_DEST_DIR)/logger github.com/deis/deis/logger || exit 1
	docker build -t $(IMAGE) image

clean: check-docker check-registry
	rm -f image/bin/logger
	docker rmi $(IMAGE)

full-clean: check-docker check-registry
	docker images -q $(IMAGE_PREFIX)$(COMPONENT) | xargs docker rmi -f

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
	docker tag $(IMAGE) $(DEV_IMAGE)
	docker push $(DEV_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_IMAGE)

release:
	docker push $(IMAGE)

deploy: build dev-release restart

test: test-unit test-functional

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

test-style:
	go vet -x ./...
	-golint ./...

test-unit: test-style
	go test -v -cover ./syslog

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...

coverage:
	go test -coverprofile coverage.out ./syslog
	go tool cover -html=coverage.out
