include ../includes.mk

build: check-docker
	docker build -t deis/database:$(GIT_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/database:$(GIT_TAG) $(REGISTRY)/deis/database:$(GIT_TAG)
	docker push $(REGISTRY)/deis/database:$(GIT_TAG)
	deisctl config database set image=$$DEIS_REGISTRY/deis/database:$(GIT_TAG)

clean: check-docker check-registry
	docker rmi deis/database:$(GIT_TAG)
	docker rmi $(REGISTRY)/deis/database:$(GIT_TAG)

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
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 20m ./tests/...
