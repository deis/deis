#
# Deis Makefile
#

ifndef DEIS_NUM_INSTANCES
    DEIS_NUM_INSTANCES = 3
endif

define ssh_all
  i=1 ; while [ $$i -le $(DEIS_NUM_INSTANCES) ] ; do \
      vagrant ssh deis-$$i -c $(1) ; \
      i=`expr $$i + 1` ; \
  done
endef

# ordered list of deis components
COMPONENTS=registry logger database cache controller builder router

all: build run

pull:
	$(call ssh_all,'for c in $(COMPONENTS); do docker pull deis/$$c; done')

build:
	$(call ssh_all,'cd share && for c in $(COMPONENTS); do cd $$c && docker build -t deis/$$c . && cd ..; done')

install:
	for c in $(COMPONENTS); do fleetctl --strict-host-key-checking=false submit ../../$$c/systemd/*; done

uninstall: stop
	for c in $(COMPONENTS); do fleetctl --strict-host-key-checking=false destroy ../../$$c/systemd/*; done

start:
	for c in $(COMPONENTS); do fleetctl --strict-host-key-checking=false start ../../$$c/systemd/*; done

stop:
	for c in $(COMPONENTS); do fleetctl --strict-host-key-checking=false stop ../../$$c/systemd/*; done

restart: stop start

run: install start

clean: uninstall
	$(call ssh_all,'for c in $(COMPONENTS); do docker rm -f deis-$$c; done')

full-clean: clean
	$(call ssh_all,'for c in $(COMPONENTS); do docker rmi deis-$$c; done')
