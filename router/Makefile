include ../includes.mk

COMPONENT = router
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)

build: check-docker
	cd parent && docker build -t deis/binary-router .
	docker cp `docker run -d deis/binary-router`:/nginx.tgz .
	docker build -t $(IMAGE) .	
	rm nginx.tgz

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

dev-release: check-registry check-deisctl
	docker tag $(IMAGE) $(DEV_IMAGE)
	docker push $(DEV_IMAGE)
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
