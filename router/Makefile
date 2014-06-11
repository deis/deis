include ../includes.mk

build:
	$(call ssh_all,'cd share/router && sudo docker build -t deis/router .')

install: check-fleet
	@$(foreach R, $(ROUTER_UNITS), \
		cp systemd/deis-router.service ./$(R) ; \
		$(FLEETCTL) load ./$(R) ; \
		rm -f ./$(R) ; \
	)

uninstall: check-fleet stop
	@$(foreach R, $(ROUTER_UNITS), \
		cp systemd/deis-router.service ./$(R) ; \
		$(FLEETCTL) unload ./$(R) ; \
		$(FLEETCTL) destroy ./$(R) ; \
		rm -f ./$(R) ; \
	)

start: check-fleet
	@$(foreach R, $(ROUTER_UNITS), \
		cp systemd/deis-router.service ./$(R) ; \
		$(FLEETCTL) start -no-block ./$(R) ; \
		rm -f ./$(R) ; \
	)

stop: check-fleet
	@$(foreach R, $(ROUTER_UNITS), \
		cp systemd/deis-router.service ./$(R) ; \
		$(FLEETCTL) stop -block-attempts=600 ./$(R) ; \
		rm -f ./$(R) ; \
	)

restart: stop start

run: install start

clean: uninstall
	$(call ssh_all,'sudo docker rm -f deis-router')

full-clean: clean
	$(call ssh_all,'sudo docker rmi deis/router')

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@echo no functional tests
