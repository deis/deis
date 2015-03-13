include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/cache

GO_PACKAGES = . tests
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)
BUILD_IMAGE = $(COMPONENT)-build
BINARY_DEST_DIR = image/bin

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -v -ldflags '-s' -o $(BINARY_DEST_DIR)/boot main.go || exit 1

build: check-docker
	docker build -t $(BUILD_IMAGE) .
	docker cp `docker run -d $(BUILD_IMAGE)`:/usr/local/bin/redis-server $(BINARY_DEST_DIR)/
	docker build -t $(IMAGE) image

clean: check-docker check-registry
	rm -f $(BINARY_DEST_DIR)/boot
	rm -f $(BINARY_DEST_DIR)/redis-server
	docker rmi -f $(IMAGE)

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
	docker tag -f $(IMAGE) $(DEV_IMAGE)
	docker push $(DEV_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_IMAGE)

release:
	docker push $(IMAGE)

deploy: compile build dev-release restart

test: test-style test-unit test-functional

test-unit:
	@echo no unit tests

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES)
	@$(GOFMT) $(GO_PACKAGES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...
