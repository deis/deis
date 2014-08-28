include ../includes.mk

build:
	$(call rsync_all)
	$(call ssh_all,'cd share/store && sudo docker build -t deis/store .')

stop: check-fleet
	$(foreach S,$(STORE_UNITS),$(FLEETCTL) stop $(S);)

clean: uninstall
	$(call ssh_all,'sudo docker rm -f deis-store')

full-clean: clean
	$(call ssh_all,'sudo docker rmi deis/store')
