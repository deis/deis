include ../includes.mk

COMPONENT = builder
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)
BUILD_IMAGE := $(COMPONENT)-build
BINARIES := extract-domain extract-types extract-version generate-buildhook get-app-config get-app-values publish-release-controller yaml2json-procfile
BINARY_DEST_DIR := image/bin

build: check-docker
	for i in $(BINARIES); do \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -v -ldflags '-s' -o $(BINARY_DEST_DIR)/$$i bin/$$i.go ; \
	done
	docker build -t $(IMAGE) image

clean: check-docker check-registry
	for i in $(BINARIES); do \
		rm -rf $(BINARY_DEST_DIR)/$$i ; \
	done
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

cedarish/build:
	mkdir -p build
	docker pull progrium/cedarish:latest
	docker save progrium/cedarish:latest | gzip -9 > build/progrium_cedarish_$(shell date "+%Y_%m_%d").tar.gz
