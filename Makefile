#
# Deis Makefile
#

ifndef FLEETCTL
	FLEETCTL = fleetctl --strict-host-key-checking=false
endif

ifndef FLEETCTL_TUNNEL
$(error You need to set FLEETCTL_TUNNEL to the IP address of a server in the cluster.)
endif

ifndef DEIS_NUM_INSTANCES
	DEIS_NUM_INSTANCES = 1
endif

ifndef DEIS_NUM_ROUTERS
	DEIS_NUM_ROUTERS = 1
endif

ifndef DEIS_FIRST_ROUTER
	DEIS_FIRST_ROUTER = 1
endif

DEIS_LAST_ROUTER = $(shell echo $(DEIS_FIRST_ROUTER)\+$(DEIS_NUM_ROUTERS)\-1 | bc)

# TODO refactor to support non-vagrant installations, since this Makefile
# is now used by the various contrib/ scripts.
define ssh_all
	i=1 ; while [ $$i -le $(DEIS_NUM_INSTANCES) ] ; do \
			vagrant ssh deis-$$i -c $(1) ; \
			i=`expr $$i + 1` ; \
	done
endef

define check_for_errors
	@if $(FLEETCTL) list-units | egrep -q "(failed|dead)"; then \
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

define echo_yellow
	@echo "\033[0;33m$(subst ",,$(1))\033[0m"
endef

# TODO: re-evaluate the start order now that we're on fleet 0.3.2.
# due to scheduling problems with fleet 0.2.0, start order of components
# is fragile. hopefully this can be changed soon...
COMPONENTS=builder cache controller database logger registry
ALL_COMPONENTS=$(COMPONENTS) router
START_COMPONENTS=registry logger cache database

ALL_UNITS = $(foreach C, $(COMPONENTS), $(wildcard $(C)/systemd/*))
START_UNITS = $(foreach C, $(START_COMPONENTS), $(wildcard $(C)/systemd/*))
ROUTER_UNITS = $(shell seq -f "deis-router.%g.service" -s " " $(DEIS_FIRST_ROUTER) 1 $(DEIS_LAST_ROUTER))

all: build run

build:
	$(call ssh_all,'cd share && for c in $(ALL_COMPONENTS); do cd $$c && docker build -t deis/$$c . && cd ..; done')

check-fleet:
	@LOCAL_VERSION=`$(FLEETCTL) -version`; \
	REMOTE_VERSION=`ssh -o StrictHostKeyChecking=no core@$(subst :, -p ,$(FLEETCTL_TUNNEL)) fleetctl -version`; \
	if [ "$$LOCAL_VERSION" != "$$REMOTE_VERSION" ]; then \
			echo "Your fleetctl client version should match the server. Local version: $$LOCAL_VERSION, server version: $$REMOTE_VERSION. Uninstall your local version and install the latest build from https://github.com/coreos/fleet/releases"; exit 1; \
	fi

clean: uninstall
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker rm -f deis-$$c; done')

full-clean: clean
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker rmi deis-$$c; done')

install: check-fleet
	$(FLEETCTL) submit $(START_UNITS)

pull:
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker pull deis/$$c; done')
	$(call ssh_all,'docker pull deis/slugrunner')

restart: stop start

run: install start

start: check-fleet start-routers
	@# registry logger cache database
	$(call echo_yellow,"Starting Deis! Deis will be functional once all services are reported as running... ")
	$(FLEETCTL) start $(START_UNITS)
	$(call echo_yellow,"Waiting for deis-registry to start (this can take some time)... ")
	@until $(FLEETCTL) list-units | egrep -q "deis-registry.+(running|failed|dead)"; do printf "\033[0;33mStatus:\033[0m "; $(FLEETCTL) list-units | grep "registry" | awk '{printf $$3}'; printf "\r" ; sleep 10; done
	$(call check_for_errors)

	@# controller
	$(call echo_yellow,"Done! Waiting for deis-controller...")
	$(FLEETCTL) submit controller/systemd/*
	$(FLEETCTL) start controller/systemd/*
	@until $(FLEETCTL) list-units | egrep -q "deis-controller.+(running|failed|dead)"; do printf "\033[0;33mStatus:\033[0m "; $(FLEETCTL) list-units | grep "controller" | awk '{printf $$3}'; printf "\r" ; sleep 10; done
	$(call check_for_errors)

	@# builder
	$(call echo_yellow,"Done! Waiting for deis-builder to start (this can also take some time)... ")
	$(FLEETCTL) submit builder/systemd/*
	$(FLEETCTL) start builder/systemd/*
	@until $(FLEETCTL) list-units | egrep -q "deis-builder.+(running|failed|dead)"; do printf "\033[0;33mStatus:\033[0m "; $(FLEETCTL) list-units | grep "builder" | awk '{printf $$3}'; printf "\r" ; sleep 10; done
	$(call check_for_errors)

	$(call echo_yellow,"Your Deis cluster is ready to go! Continue following the README to login and use Deis.")

start-routers:
	$(call echo_yellow,"Starting $(DEIS_NUM_ROUTERS) router(s)...")
	@ $(foreach U, $(ROUTER_UNITS), \
		cp router/systemd/deis-router.service ./$(U) ; \
		$(FLEETCTL) submit ./$(U) ; \
		$(FLEETCTL) start ./$(U) ; \
		rm -f ./$(U) ; \
	)

status: check-fleet
	$(FLEETCTL) list-units

stop: check-fleet
	$(FLEETCTL) stop $(call deis_units,loaded,.)

tests:
	cd test && bundle install && bundle exec rake

uninstall: check-fleet stop
	$(FLEETCTL) destroy $(call deis_units,loaded,.)
