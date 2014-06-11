include ../includes.mk

build:
	$(call ssh_all,'cd share/database && sudo docker build -t deis/database .')

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
	$(call ssh_all,'sudo docker rm -f deis-database')

full-clean: clean
	$(call ssh_all,'sudo docker rmi deis/database')

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@echo no functional tests
