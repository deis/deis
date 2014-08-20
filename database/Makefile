include ../includes.mk

build:
	$(call rsync_all)
	$(call ssh_all,'cd share/database && sudo docker build -t deis/database .')

install: check-fleet
	$(FLEETCTL) load systemd/*.service

uninstall: check-fleet stop
	$(FLEETCTL) unload systemd/*.service
	$(FLEETCTL) destroy systemd/*.service

start: check-fleet
	$(FLEETCTL) start -no-block systemd/*.service

stop: check-fleet
	$(FLEETCTL) stop -block-attempts=600 systemd/*.service

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
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 20m ./tests/...
