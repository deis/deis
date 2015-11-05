include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/terminator

SHELL_SCRIPTS = $(shell find "." -name '*.sh') $(wildcard rootfs/bin/*)

COMPONENT = $(notdir $(repo_path))
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(REGISTRY)$(IMAGE)
BUILD_IMAGE = $(COMPONENT)-build

build: check-docker
	docker build -t $(IMAGE) .

clean: check-docker check-registry
	docker rmi $(IMAGE)

full-clean: check-docker check-registry
	docker images -q $(IMAGE_PREFIX)$(COMPONENT) | xargs docker rmi -f

dev-release: push set-image

push: check-registry
	docker tag -f $(IMAGE) $(DEV_IMAGE)
	docker push $(DEV_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_IMAGE)

release:
	docker push $(IMAGE)

deploy: build dev-release

test: test-style test-unit test-functional

test-functional:
	@echo no functional tests

test-style:
	shellcheck $(SHELL_SCRIPTS)

test-unit:
	@echo no unit tests

uninstall:
	@echo nothing to do

stop:
	@echo nothing to do