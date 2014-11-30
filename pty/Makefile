include ../includes.mk

DOCKER_IMAGE := deis/pty
BUILD_IMAGE := $(DOCKER_IMAGE)-build
RELEASE_IMAGE := $(DOCKER_IMAGE):$(BUILD_TAG)
REMOTE_IMAGE := $(REGISTRY)/$(RELEASE_IMAGE)

build: check-docker
	docker build -t $(BUILD_IMAGE) .
	docker cp `docker run -d $(BUILD_IMAGE)`:/go/bin/pty image/bin
	docker build -t $(RELEASE_IMAGE) image
	rm -rf image/bin/pty

clean: check-docker check-registry
	docker rmi $(RELEASE_IMAGE) $(REMOTE_IMAGE)

full-clean: check-docker check-registry
	docker images -q $(DOCKER_IMAGE) | xargs docker rmi -f
	docker images -q $(REGISTRY)/$(DOCKER_IMAGE) | xargs docker rmi -f

install: check-deisctl
	echo "Installed"

dev-release: push set-image

push: check-registry
	docker tag $(RELEASE_IMAGE) $(REMOTE_IMAGE)
	docker push $(REMOTE_IMAGE)

set-image: check-deisctl
	deisctl config publisher set image=$(REMOTE_IMAGE)

release: check-docker
	docker push $(IMAGE)

deploy: build dev-release restart

restart: stop start

run: install start

uninstall: check-deisctl
	@echo "Nothing to do"

start: check-deisctl
	@echo "Nothing to do"

stop: check-deisctl
	@echo "Nothing to do"

test: test-unit test-functional

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

test-style:
	godep go vet -x ./...
	-golint .

test-unit: test-style
	godep go test -v -cover .

test-functional:
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v ./tests/...

coverage:
	godep go test -coverprofile coverage.out ./syslog
	godep go tool cover -html=coverage.out
