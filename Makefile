#
# Deis Makefile
#

include includes.mk

SHELL := /bin/bash

define check_for_errors
	@if $(FLEETCTL) list-units -no-legend | awk '(($$2 == "launched") && ($$5 == "failed"))' | egrep -q "deis-.+service"; then \
		echo "\033[0;31mOne or more services failed! Check which services by running 'make status'\033[0m" ; \
		echo "\033[0;31mYou can get detailed output with 'fleetctl status deis-servicename.service'\033[0m" ; \
		echo "\033[0;31mThis usually indicates an error with Deis - please open an issue on GitHub or ask for help in IRC\033[0m" ; \
		exit 1 ; \
	fi
endef

define deis_units
	$(shell $(FLEETCTL) list-units -no-legend=true | \
	  awk '($$2 ~ "$(1)" && ($$4 ~ "$(2)"))' | \
	  sed -n 's/\(deis-.*\.service\).*/\1/p' | tr '\n' ' ')
endef

# TODO: re-evaluate the fragile start order
COMPONENTS=builder cache controller database logger registry
ALL_COMPONENTS=$(COMPONENTS) router
START_COMPONENTS=registry logger cache database

ALL_UNITS = $(foreach C,$(COMPONENTS),$(wildcard $(C)/systemd/*.service))
START_UNITS = $(foreach C,$(START_COMPONENTS),$(wildcard $(C)/systemd/*.service))

DATA_CONTAINER_TEMPLATES=builder/systemd/deis-builder-data.service database/systemd/deis-database-data.service logger/systemd/deis-logger-data.service registry/systemd/deis-registry-data.service

all: build run

build: rsync
	$(call ssh_all,'cd share && for c in $(ALL_COMPONENTS); do cd $$c && docker build -t deis/$$c . && cd ..; done')

clean: uninstall
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker rm -f deis-$$c; done')

full-clean: clean
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker rmi deis-$$c; done')

install: check-fleet install-routers install-data-containers
	$(FLEETCTL) load $(START_UNITS)
	$(FLEETCTL) load controller/systemd/*.service
	$(FLEETCTL) load builder/systemd/*.service

install-data-containers: check-fleet
	@$(foreach T, $(DATA_CONTAINER_TEMPLATES), \
		UNIT=`basename $(T)` ; \
		if [[ `$(FLEETCTL) list-units | grep $$UNIT` ]]; then \
			echo $$UNIT already loaded. Skipping... ; \
		else \
			cp $(T).template . ; \
			NEW_FILENAME=`ls *.template | sed 's/\.template//g'`; \
			mv *.template $$NEW_FILENAME ; \
			MACHINE_ID=`$(FLEETCTL) list-machines --no-legend --full list-machines | awk 'BEGIN { OFS="\t"; srand() } { print rand(), $$1 }' | sort -n | cut -f2- | head -1` ; \
			sed -e "s/CHANGEME/$$MACHINE_ID/" $$NEW_FILENAME > $$NEW_FILENAME.bak ; \
			rm -f $$NEW_FILENAME ; \
			mv $$NEW_FILENAME.bak $$NEW_FILENAME ; \
			$(FLEETCTL) load $$NEW_FILENAME ; \
			rm -f $$NEW_FILENAME ; \
		fi ; \
	)

install-routers: check-fleet
	@$(foreach R, $(ROUTER_UNITS), \
		cp router/systemd/deis-router.service ./$(R) ; \
		$(FLEETCTL) load ./$(R) ; \
		rm -f ./$(R) ; \
	)

pull:
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker pull deis/$$c:latest; done')
	$(call ssh_all,'docker pull deis/slugrunner:latest')

restart: stop start

rsync:
	$(call rsync_all)

run: install start

start: check-fleet start-warning start-routers
	@# registry logger cache database
	$(call echo_yellow,"Waiting for deis-registry to start...")
	$(FLEETCTL) start -no-block $(START_UNITS)
	@until $(FLEETCTL) list-units | egrep -q "deis-registry.+(running|failed)"; \
		do sleep 2; \
			printf "\033[0;33mStatus:\033[0m "; $(FLEETCTL) list-units | \
			grep "deis-registry" | awk '{printf "%-10s (%s)    \r", $$4, $$5}'; \
			sleep 8; \
		done
	$(call check_for_errors)

	@# controller
	$(call echo_yellow,"Waiting for deis-controller to start...")
	$(FLEETCTL) start -no-block controller/systemd/*
	@until $(FLEETCTL) list-units | egrep -q "deis-controller.+(running|failed)"; \
		do sleep 2; \
			printf "\033[0;33mStatus:\033[0m "; $(FLEETCTL) list-units | \
			grep "deis-controller" | awk '{printf "%-10s (%s)    \r", $$4, $$5}'; \
			sleep 8; \
		done
	$(call check_for_errors)

	@# builder
	$(call echo_yellow,"Waiting for deis-builder to start...")
	$(FLEETCTL) start -no-block builder/systemd/*.service
	@until $(FLEETCTL) list-units | egrep -q "deis-builder.+(running|failed)"; \
		do sleep 2; \
			printf "\033[0;33mStatus:\033[0m "; $(FLEETCTL) list-units | \
			grep "deis-builder" | awk '{printf "%-10s (%s)    \r", $$4, $$5}'; \
			sleep 8; \
		done
	$(call check_for_errors)

	$(call echo_yellow,"Your Deis cluster is ready to go! Continue following the README to login and use Deis.")

start-routers: check-fleet start-warning
	$(call echo_yellow,"Waiting for 1 of $(DEIS_NUM_ROUTERS) deis-routers to start...")
	$(foreach R,$(ROUTER_UNITS),$(FLEETCTL) start -no-block $(R);)
	@until $(FLEETCTL) list-units | egrep -q "deis-router.+(running)"; \
		do sleep 2; \
			printf "\033[0;33mStatus:\033[0m "; $(FLEETCTL) list-units | \
			grep "deis-router" | head -n 1 | \
			awk '{printf "%-10s (%s)    \r", $$4, $$5}'; \
			sleep 8; \
		done
	$(call check_for_errors)

start-warning:
	$(call echo_cyan,"Deis components may take a long time to start the first time they are initialized.")

status: check-fleet
	$(FLEETCTL) list-units

stop: check-fleet
	$(FLEETCTL) stop -block-attempts=600 $(strip $(call deis_units,launched,active))

test: test-components test-integration

test-components:
	@$(foreach C,$(ALL_COMPONENTS), \
		echo \\nTesting deis/$(C) ; \
		$(MAKE) -C $(C) test ; \
	)

test-integration:
	$(MAKE) -C tests/ test

test-smoke:
	$(MAKE) -C tests/ test-smoke

uninstall: check-fleet stop
	$(FLEETCTL) unload -block-attempts=600 $(call deis_units,launched,.)
	$(FLEETCTL) destroy $(strip $(call deis_units,.,.))
