include ../includes.mk

build: check-docker
	docker build -t deis/database:$(BUILD_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/database:$(BUILD_TAG) $(REGISTRY)/deis/database:$(BUILD_TAG)
	docker push $(REGISTRY)/deis/database:$(BUILD_TAG)
	deisctl config database set image=$(REGISTRY)/deis/database:$(BUILD_TAG)

clean: check-docker check-registry
	docker rmi deis/database:$(BUILD_TAG)
	docker rmi $(REGISTRY)/deis/database:$(BUILD_TAG)

full-clean: check-docker check-registry
	docker images -q deis/database | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/database | xargs docker rmi -f

install: check-deisctl
	deisctl install database-data
	deisctl scale database=1

uninstall: check-deisctl
	deisctl scale database=0
	deisctl uninstall database-data

start: check-deisctl
	deisctl start database

stop: check-deisctl
	deisctl stop database

restart: stop start

run: install start

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 20m ./tests/...
