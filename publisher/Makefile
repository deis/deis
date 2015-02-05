include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/publisher

GO_FILES = main.go
GO_PACKAGES = server
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
DOCKER_IMAGE := $(IMAGE_PREFIX)$(COMPONENT)
BUILD_IMAGE := $(DOCKER_IMAGE)-build
RELEASE_IMAGE := $(DOCKER_IMAGE):$(BUILD_TAG)
REMOTE_IMAGE := $(REGISTRY)/$(RELEASE_IMAGE)

build: check-docker
	cp -pR ../Godeps .
	docker build -t $(BUILD_IMAGE) .
	docker cp `docker run -d $(BUILD_IMAGE)`:/go/bin/publisher image/
	docker build -t $(RELEASE_IMAGE) image
	rm -rf image/publisher

clean: check-docker check-registry
	rm -rf Godeps
	docker rmi $(RELEASE_IMAGE) $(REMOTE_IMAGE)

full-clean: check-docker check-registry
	docker images -q $(DOCKER_IMAGE) | xargs docker rmi -f
	docker images -q $(REGISTRY)/$(DOCKER_IMAGE) | xargs docker rmi -f

install: check-deisctl
	deisctl install publisher

dev-release: push set-image

push: check-registry
	docker tag -f $(RELEASE_IMAGE) $(REMOTE_IMAGE)
	docker push $(REMOTE_IMAGE)

set-image: check-deisctl
	deisctl config publisher set image=$(REMOTE_IMAGE)

release: check-docker
	docker push $(DOCKER_IMAGE)

deploy: build dev-release restart

restart: stop start

run: install start

start: check-deisctl
	deisctl start publisher

stop: check-deisctl
	deisctl stop publisher

test: test-style test-unit test-functional

test-functional:
	@echo no functional tests

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES) $(GO_FILES)
	@$(GOFMT) $(GO_PACKAGES) $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
# FIXME: make this mandatory
	-$(GOVET) $(repo_path) $(GO_PACKAGES_REPO_PATH)
# FIXME: make this mandatory
	-$(GOLINT) ./...

test-unit:
	godep go test -v ./...

uninstall: check-deisctl
	deisctl uninstall publisher
