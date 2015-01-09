include ../includes.mk

build: check-docker
	cp -pR ../../Godeps .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -v -ldflags '-s' -o read-dockercfg main.go

clean: check-docker check-registry
	rm -rf Godeps

full-clean:
	@echo no available

install: check-deisctl
	@echo no available

dev-release:
	@echo no available

push:
	@echo no available

set-image:
	@echo no available

release:
	@echo no available

deploy:
	@echo no available

restart:
	@echo no available

run: install start

start: check-deisctl
	@echo no available

stop:
	@echo no available

test: test-unit

test-functional:
	@echo no functional tests

test-unit:
	godep go test -v ./...

uninstall: check-deisctl
	@echo no functional tests
