include ../includes.mk

build: check-docker
	docker build -t deis/builder:$(GIT_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/builder:$(GIT_TAG) $(REGISTRY)/deis/builder:$(GIT_TAG)
	docker push $(REGISTRY)/deis/builder:$(GIT_TAG)
	deisctl config builder set image=$$DEIS_REGISTRY/deis/builder:$(GIT_TAG)

clean: check-docker check-registry
	docker rmi deis/builder:$(GIT_TAG)
	docker rmi $(REGISTRY)/deis/builder:$(GIT_TAG)

full-clean: check-docker check-registry
	docker images -q deis/builder | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/builder | xargs docker rmi -f

install: check-deisctl
	deisctl install builder-data
	deisctl scale builder=1

uninstall: check-deisctl
	deisctl scale builder=0
	deisctl uninstall builder-data

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
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 30m ./tests/...
