include ../includes.mk

build:
	$(call rsync_all)
	$(call ssh_all,'cd share/store/base && sudo docker build -t deis/store-base .')
	$(call ssh_all,'cd share/store/daemon && sudo docker build -t deis/store-daemon .')
	$(call ssh_all,'cd share/store/monitor && sudo docker build -t deis/store-monitor .')
	$(call ssh_all,'cd share/store/gateway && sudo docker build -t deis/store-gateway .')

stop: check-fleet
	$(foreach S,$(STORE_DAEMON_UNITS),$(FLEETCTL) stop $(S);)
	$(foreach S,$(STORE_MONITOR_UNITS),$(FLEETCTL) stop $(S);)
	$(FLEETCTL) stop deis-store-gateway
