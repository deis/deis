include ../includes.mk

build: check-docker
	docker build -t deis/cache:$(BUILD_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/cache:$(BUILD_TAG) $(REGISTRY)/deis/cache:$(BUILD_TAG)
	docker push $(REGISTRY)/deis/cache:$(BUILD_TAG)
	deisctl config cache set image=$(REGISTRY)/deis/cache:$(BUILD_TAG)

clean: check-docker check-registry
	docker rmi deis/cache:$(BUILD_TAG)
	docker rmi $(REGISTRY)/deis/cache:$(BUILD_TAG)

full-clean: check-docker check-registry
	docker images -q deis/cache | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/cache | xargs docker rmi -f

install: check-deisctl
	deisctl scale cache=1

uninstall: check-deisctl
	deisctl scale cache=0

start: check-deisctl
	deisctl start cache

stop: check-deisctl
	deisctl stop cache

restart: stop start

run: install start

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v ./tests/...
