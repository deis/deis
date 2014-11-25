include ../includes.mk

COMPONENT = router
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)
BUILD_IMAGE = $(COMPONENT)-build

build: check-docker
	docker build -t $(BUILD_IMAGE) parent
	docker cp `docker run -d $(BUILD_IMAGE)`:/nginx.tgz `pwd`
	docker cp `docker run -d $(BUILD_IMAGE)`:/go/bin/boot bin/
	docker build -t $(IMAGE) .
	rm nginx.tgz
	rm bin/boot

clean: check-docker check-registry
	docker rmi $(IMAGE)

full-clean: check-docker check-registry
	docker images -q $(IMAGE_PREFIX)$(COMPONENT) | xargs docker rmi -f

install: check-deisctl
	deisctl scale $(COMPONENT)=1

uninstall: check-deisctl
	deisctl scale $(COMPONENT)=0

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
	@echo no unit tests

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v ./tests/...
