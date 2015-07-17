include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/publisher

GO_FILES = $(wildcard *.go)
GO_PACKAGES = server
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
DOCKER_IMAGE := $(IMAGE_PREFIX)$(COMPONENT)
RELEASE_IMAGE := $(DOCKER_IMAGE):$(BUILD_TAG)
REMOTE_IMAGE := $(REGISTRY)$(RELEASE_IMAGE)
BINARY_DEST_DIR = image/bin

build: check-docker
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o $(BINARY_DEST_DIR)/publisher github.com/deis/deis/publisher || exit 1
	@$(call check-static-binary,$(BINARY_DEST_DIR)/publisher)
	docker build -t $(RELEASE_IMAGE) image

clean: check-docker check-registry
	rm -rf $(BINARY_DEST_DIR)
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
	$(GOVET) $(repo_path) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

test-unit:
	$(GOTEST) ./...

uninstall: check-deisctl
	deisctl uninstall publisher
