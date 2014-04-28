#
# Deis Makefile
#

ifndef FLEETCTL_TUNNEL
	$(error You need to set FLEETCTL_TUNNEL to the IP address of a server in the cluster.)
endif

ifndef DEIS_NUM_INSTANCES
	DEIS_NUM_INSTANCES = 1
endif

# TODO refactor to support non-vagrant installations, since this Makefile
# is now used by the various contrib/ scripts.
define ssh_all
	i=1 ; while [ $$i -le $(DEIS_NUM_INSTANCES) ] ; do \
			vagrant ssh deis-$$i -c $(1) ; \
			i=`expr $$i + 1` ; \
	done
endef

define check_for_errors
	@if fleetctl --strict-host-key-checking=false list-units|grep -q "failed"; then \
		echo "\033[0;31mOne or more services failed! Check which services by running 'make status'\033[0m" ; \
		echo "\033[0;31mYou can get detailed output with 'fleetctl status deis-servicename.service'\033[0m" ; \
		echo "\033[0;31mThis usually indicates an error with Deis - please open an issue on GitHub or ask for help in IRC\033[0m" ; \
	fi
endef

define echo_yellow
	@echo "\033[0;33m$(subst ",,$(1))\033[0m"
endef

# due to scheduling problems with fleet 0.2.0, start order of components
# is fragile. hopefully this can be changed soon...
ALL_COMPONENTS=builder cache controller database logger registry router

ifeq ($(DEIS_NUM_INSTANCES),1)
	export SKIP_ROUTER=false
	START_COMPONENTS=registry logger cache database router
else
	export SKIP_ROUTER=true
	START_COMPONENTS=registry logger cache database
endif

ALL_UNITS = $(foreach C, $(ALL_COMPONENTS), $(wildcard $(C)/systemd/*))
START_UNITS = $(foreach C, $(START_COMPONENTS), $(wildcard $(C)/systemd/*))

all: build run

build:
	$(call ssh_all,'cd share && for c in $(ALL_COMPONENTS); do cd $$c && docker build -t deis/$$c . && cd ..; done')

check-fleet:
	@LOCAL_VERSION=`fleetctl -version`; \
	REMOTE_VERSION=`ssh -o StrictHostKeyChecking=no core@$(FLEETCTL_TUNNEL) fleetctl -version`; \
	if [ "$$LOCAL_VERSION" != "$$REMOTE_VERSION" ]; then \
			echo "Your fleetctl client version should match the server. Local version: $$LOCAL_VERSION, server version: $$REMOTE_VERSION. Uninstall your local version and install the latest build from https://github.com/coreos/fleet/releases"; exit 1; \
	fi

clean: uninstall
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker rm -f deis-$$c; done')

full-clean: clean
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker rmi deis-$$c; done')

install: check-fleet
	fleetctl --strict-host-key-checking=false submit $(START_UNITS)

pull:
	$(call ssh_all,'for c in $(ALL_COMPONENTS); do docker pull deis/$$c; done')
	$(call ssh_all,'docker pull deis/slugrunner')

restart: stop start

run: install start

start: check-fleet
	# registry logger cache database (router)
	fleetctl --strict-host-key-checking=false start $(START_UNITS)
	$(call echo_yellow,"Waiting for initial services to start (this can take some time)... ")
	until fleetctl --strict-host-key-checking=false list-units | egrep -q "deis-registry.+(running|failed)"; do sleep 10; done
	$(call check_for_errors)
	$(call echo_yellow,"Done! Waiting for deis-builder...")

	# builder
	fleetctl --strict-host-key-checking=false submit builder/systemd/*
	fleetctl --strict-host-key-checking=false start builder/systemd/*
	until fleetctl --strict-host-key-checking=false list-units | egrep -q "deis-builder.+(running|failed)"; do sleep 10; done
	$(call check_for_errors)
	$(call echo_yellow,"Done! Waiting for deis-controller...")

	# controller
	fleetctl --strict-host-key-checking=false submit controller/systemd/*
	fleetctl --strict-host-key-checking=false start controller/systemd/*
	until fleetctl --strict-host-key-checking=false list-units | egrep -q "deis-controller.+(running|failed)"; do sleep 10; done
	$(call check_for_errors)
	@if [ "$$SKIP_ROUTER" = true ]; then \
		echo "\033[0;33mYou'll need to configure DNS and start the router manually for multi-node clusters.\033[0m" ; \
		echo "\033[0;33mRun 'make start-router' to schedule and start deis-router.\033[0m" ; \
	else \
		echo "\033[0;33mYour Deis cluster is ready to go! Follow the README to login and use Deis.\033[0m" ; \
	fi

start-router: check-fleet
	fleetctl --strict-host-key-checking=false submit router/systemd/*
	fleetctl --strict-host-key-checking=false start router/systemd/*
	$(call echo_yellow,"Use 'make status' to monitor the service and note the IP it has been scheduled to.")
	$(call echo_yellow,"Create a wildcard DNS domain which resolves to this host and use that domain when creating clusters/apps in the README.")
	$(call echo_yellow,"Your Deis cluster is ready to go! Follow the README to login and use Deis.")

status: check-fleet
	fleetctl --strict-host-key-checking=false list-units

stop: check-fleet
	fleetctl --strict-host-key-checking=false stop $(ALL_UNITS)

uninstall: check-fleet stop
	fleetctl --strict-host-key-checking=false destroy $(ALL_UNITS)
