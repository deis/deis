include ../includes.mk

COMPONENT = builder
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)
BUILD_IMAGE := $(COMPONENT)-build
BINARY_DEST_DIR := image/bin

build: check-docker
	cp -pR ../Godeps .
	docker build -t $(BUILD_IMAGE) .
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/extract-domain $(BINARY_DEST_DIR)/
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/extract-types $(BINARY_DEST_DIR)/
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/extract-version $(BINARY_DEST_DIR)/
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/generate-buildhook $(BINARY_DEST_DIR)/
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/get-app-config $(BINARY_DEST_DIR)/
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/get-app-values $(BINARY_DEST_DIR)/
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/publish-release-controller $(BINARY_DEST_DIR)/
	docker cp `docker create $(BUILD_IMAGE)`:/go/bin/yaml2json-procfile $(BINARY_DEST_DIR)/
	docker build -t $(IMAGE) image

clean: check-docker check-registry
	rm -rf Godeps
	rm -f $(BINARY_DEST_DIR)/extract-domain
	rm -f $(BINARY_DEST_DIR)/extract-types
	rm -f $(BINARY_DEST_DIR)/extract-version
	rm -f $(BINARY_DEST_DIR)/generate-buildhook
	rm -f $(BINARY_DEST_DIR)/get-app-config
	rm -f $(BINARY_DEST_DIR)/get-app-values
	rm -f $(BINARY_DEST_DIR)/publish-release-controller
	rm -f $(BINARY_DEST_DIR)/yaml2json-procfile
	docker rmi $(IMAGE)

full-clean: check-docker check-registry clean
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

test-unit:
	godep go test -v .

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...
