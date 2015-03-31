include ../includes.mk

COMPONENT = database
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)
BINARY_DEST_DIR = bin

build: check-docker
	go-bindata -pkg bindata -o bindata/bindata.go bash/
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o $(BINARY_DEST_DIR)/boot boot.go || exit 1
	$(call check-static-binary,$(BINARY_DEST_DIR)/boot)
	docker build -t $(IMAGE) .

clean: check-docker check-registry
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
	docker tag -f $(IMAGE) $(DEV_IMAGE)
	docker push $(DEV_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_IMAGE)

release:
	docker push $(IMAGE)

deploy: build dev-release restart

test: test-style test-unit test-functional

test-unit:
	@echo no unit tests

test-style:
	@echo no style tests

test-functional:
	@docker history deis/mock-store >/dev/null 2>&1 || $(MAKE) -C ../tests/ mock-store
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...
