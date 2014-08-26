include ../includes.mk

build:
	$(call rsync_all)
	$(call ssh_all,'cd share/store && sudo docker build -t deis/store .')

start: check-fleet
	$(FLEETCTL) start -no-block deis-store*

stop: check-fleet
	$(FLEETCTL) stop -block-attempts=600 deis-store*

restart: stop start

clean: uninstall
	$(call ssh_all,'sudo docker rm -f deis-store')

full-clean: clean
	$(call ssh_all,'sudo docker rmi deis/store')
