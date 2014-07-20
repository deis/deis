include ../includes.mk

build:
	$(call rsync_all)
	$(call ssh_all,'cd share/database && sudo docker build -t deis/database .')

install: check-fleet
	@$(foreach U, $(DATABASE_UNITS), \
		cp systemd/deis-database.service ./$(U) ; \
		$(FLEETCTL) load ./$(U) ; \
		rm -f ./$(U) ; \
	)

uninstall: check-fleet stop
	@$(foreach U, $(DATABASE_UNITS), \
		cp systemd/deis-database.service ./$(U) ; \
		$(FLEETCTL) unload ./$(U) ; \
		$(FLEETCTL) destroy ./$(U) ; \
		rm -f ./$(U) ; \
	)

start: check-fleet
	@$(foreach U, $(DATABASE_UNITS), \
		cp systemd/deis-database.service ./$(U) ; \
		$(FLEETCTL) start -no-block ./$(U) ; \
		rm -f ./$(U) ; \
	)

stop: check-fleet
	@$(foreach U, $(DATABASE_UNITS), \
		cp systemd/deis-database.service ./$(U) ; \
		$(FLEETCTL) stop -block-attempts=600 ./$(U) ; \
		rm -f ./$(U) ; \
	)

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
