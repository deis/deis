include ../includes.mk

build: check-docker
	docker build -t deis/logger:$(BUILD_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/logger:$(BUILD_TAG) $(REGISTRY)/deis/logger:$(BUILD_TAG)
	docker push $(REGISTRY)/deis/logger:$(BUILD_TAG)
	deisctl config logger set image=$(REGISTRY)/deis/logger:$(BUILD_TAG)

clean: check-docker check-registry
	docker rmi deis/logger:$(BUILD_TAG)
	docker rmi $(REGISTRY)/deis/logger:$(BUILD_TAG)

full-clean: check-docker check-registry
	docker images -q deis/logger | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/logger | xargs docker rmi -f

install: check-deisctl
	deisctl install logger-data
	deisctl scale logger=1

uninstall: check-deisctl
	deisctl scale logger=0
	deisctl uninstall logger-data

start: check-deisctl
	deisctl start logger

stop: check-deisctl
	deisctl stop logger

restart: stop start

run: install start

test: test-unit test-functional

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

test-style:
	go vet -x ./...
	-golint .

test-unit: test-style
	go test -v -cover ./syslog

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v ./tests/...

coverage:
	go test -coverprofile coverage.out ./syslog
	go tool cover -html=coverage.out
