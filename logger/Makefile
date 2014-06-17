include ../includes.mk

build:
	$(call ssh_all,'cd share/logger && sudo docker build -t deis/logger .')

install: check-fleet
	$(FLEETCTL) load systemd/*

uninstall: check-fleet stop
	$(FLEETCTL) unload systemd/*
	$(FLEETCTL) destroy systemd/*

start: check-fleet
	$(FLEETCTL) start -no-block systemd/*

stop: check-fleet
	$(FLEETCTL) stop -block-attempts=600 systemd/*

restart: stop start

run: install start

clean: uninstall
	$(call ssh_all,'sudo docker rm -f deis-logger')

full-clean: clean
	$(call ssh_all,'sudo docker rmi deis/logger')

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v ./test/...
