include ../includes.mk

build: check-docker
	docker build -t deis/router:$(BUILD_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/router:$(BUILD_TAG) $(REGISTRY)/deis/router:$(BUILD_TAG)
	docker push $(REGISTRY)/deis/router:$(BUILD_TAG)
	deisctl config router set image=$(REGISTRY)/deis/router:$(BUILD_TAG)

clean: check-docker check-registry
	docker rmi deis/router:$(BUILD_TAG)
	docker rmi $(REGISTRY)/deis/router:$(BUILD_TAG)

full-clean: check-docker check-registry
	docker images -q deis/router | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/router | xargs docker rmi -f

install: check-deisctl
	deisctl scale router=1

uninstall: check-deisctl
	deisctl scale router=0

start: check-deisctl
	deisctl start router

stop: check-deisctl
	deisctl stop router

restart: stop start

run: install start

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 20m ./tests/...
