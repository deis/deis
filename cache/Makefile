include ../includes.mk

build: check-docker
	docker build -t deis/cache:$(GIT_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/cache:$(GIT_TAG) $(REGISTRY)/deis/cache:$(GIT_TAG)
	docker push $(REGISTRY)/deis/cache:$(GIT_TAG)
	deisctl config cache set image=$$DEIS_REGISTRY/deis/cache:$(GIT_TAG)

clean: check-docker check-registry
	docker rmi deis/cache:$(GIT_TAG)
	docker rmi $(REGISTRY)/deis/cache:$(GIT_TAG)

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
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 20m ./tests/...
