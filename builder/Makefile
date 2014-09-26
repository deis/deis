include ../includes.mk

build: check-docker
	docker build -t deis/builder:$(BUILD_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/builder:$(BUILD_TAG) $(REGISTRY)/deis/builder:$(BUILD_TAG)
	docker push $(REGISTRY)/deis/builder:$(BUILD_TAG)
	deisctl config builder set image=$(REGISTRY)/deis/builder:$(BUILD_TAG)

clean: check-docker check-registry
	docker rmi deis/builder:$(BUILD_TAG)
	docker rmi $(REGISTRY)/deis/builder:$(BUILD_TAG)

full-clean: check-docker check-registry
	docker images -q deis/builder | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/builder | xargs docker rmi -f

install: check-deisctl
	deisctl scale builder=1

uninstall: check-deisctl
	deisctl scale builder=0

start: check-deisctl
	deisctl start builder

stop: check-deisctl
	deisctl stop builder

restart: stop start

run: install start

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v ./tests/...
