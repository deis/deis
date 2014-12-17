include ../includes.mk

COMPONENT := logspout
DOCKER_IMAGE := deis/$(COMPONENT)
BUILD_IMAGE := $(DOCKER_IMAGE)-build
RELEASE_IMAGE := $(DOCKER_IMAGE):$(BUILD_TAG)
DEV_DOCKER_IMAGE := $(DEV_REGISTRY)/$(RELEASE_IMAGE)

build: check-docker
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -v -ldflags '-s' -o image/logspout
	docker build -t $(RELEASE_IMAGE) image

clean: check-docker check-registry
	rm -rf image/logspout
	docker rmi $(RELEASE_IMAGE) $(BUILD_IMAGE)

full-clean: check-docker check-registry
	docker images -q $(DOCKER_IMAGE) | xargs docker rmi -f
	docker images -q $(BUILD_IMAGE) | xargs docker rmi -f
	docker images -q $(DEV_DOCKER_IMAGE) | xargs docker rmi -f

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
	docker tag $(RELEASE_IMAGE) $(DEV_DOCKER_IMAGE)
	docker push $(DEV_DOCKER_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_DOCKER_IMAGE)

release:
	docker push $(RELEASE_IMAGE)

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
	@echo no tests

test-functional:
	@echo no functional tests
