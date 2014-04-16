#
# Deis Makefile
#

ifndef $(DEIS_NUM_INSTANCES)
    DEIS_NUM_INSTANCES = 3
endif

define ssh_all
  i=1 ; while [[ $$i -le $(DEIS_NUM_INSTANCES) ]] ; do \
      vagrant ssh deis-$$i -c $(1) ; \
      ((i = i + 1)) ; \
  done
endef

# ordered list of deis components
COMPONENTS=registry logger database cache controller builder router

all: build run

test_client:
	python -m unittest ../../discover client.tests

pull:
	$(call ssh_all,'for c in $(COMPONENTS); do docker pull deis/$$c; done')

build:
	$(call ssh_all,'cd share && for c in $(COMPONENTS); do cd $$c && docker build -t deis/$$c . && cd ..; done')

install:
	for c in $(COMPONENTS); do fleetctl enable ../../$$c/systemd/*; done

uninstall: stop
	for c in $(COMPONENTS); do fleetctl disable ../../$$c/systemd/*; done

start:
	for c in $(COMPONENTS); do fleetctl start ../../$$c/systemd/*; done

stop:
	for c in $(COMPONENTS); do fleetctl stop ../../$$c/systemd/*; done

restart:
	for c in $(COMPONENTS); do fleetctl restart ../../$$c/systemd/*; done

run: install restart

clean: uninstall
	$(call ssh_all,'cd share && for c in $(COMPONENTS); do docker rm -f deis-$$c; done')

full-clean: clean
	$(call ssh_all,'cd share && for c in $(COMPONENTS); do docker rmi deis-$$c; done')
