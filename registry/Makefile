include ../includes.mk

build: check-docker
	docker build -t deis/registry:$(GIT_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/registry:$(GIT_TAG) $(REGISTRY)/deis/registry:$(GIT_TAG)
	docker push $(REGISTRY)/deis/registry:$(GIT_TAG)
	deisctl config registry set image=$$DEIS_REGISTRY/deis/registry:$(GIT_TAG)

clean: check-docker check-registry
	docker rmi deis/registry:$(GIT_TAG)
	docker rmi $(REGISTRY)/deis/registry:$(GIT_TAG)

full-clean: check-docker check-registry
	docker images -q deis/registry | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/registry | xargs docker rmi -f

install: check-deisctl
	deisctl install registry-data
	deisctl scale registry=1

uninstall: check-deisctl
	deisctl scale registry=0
	deisctl uninstall registry-data

start: check-deisctl
	deisctl start registry

stop: check-deisctl
	deisctl stop registry

restart: stop start

run: install start

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 20m ./tests/...
