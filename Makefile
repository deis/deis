#
# Deis Makefile
#

ifndef DEIS_NUM_INSTANCES
    DEIS_NUM_INSTANCES = 1
endif

define ssh_all
  i=1 ; while [ $$i -le $(DEIS_NUM_INSTANCES) ] ; do \
      vagrant ssh deis-$$i -c $(1) ; \
      i=`expr $$i + 1` ; \
  done
endef

# ordered list of deis components
# we don't manage the router if we're setting up a local cluster
ifeq ($(DEIS_NUM_INSTANCES),1)
	COMPONENTS=builder cache controller database logger registry router
else
	COMPONENTS=builder cache controller database logger registry
endif

UNIT_FILES = $(foreach C, $(COMPONENTS), $(wildcard $(C)/systemd/*))

all: build run

pull:
	$(call ssh_all,'for c in $(COMPONENTS); do docker pull deis/$$c; done')

build:
	$(call ssh_all,'cd share && for c in $(COMPONENTS); do cd $$c && docker build -t deis/$$c . && cd ..; done')

install:
	fleetctl --strict-host-key-checking=false submit $(UNIT_FILES)

uninstall:  stop
	fleetctl --strict-host-key-checking=false destroy $(UNIT_FILES)

start:
	echo "\033[0;33mStarting services can take some time... grab some coffee!\033[0m"
	fleetctl --strict-host-key-checking=false start $(UNIT_FILES)

stop:
	fleetctl --strict-host-key-checking=false stop $(UNIT_FILES)

restart: stop start

run: install start

clean: uninstall
	$(call ssh_all,'for c in $(COMPONENTS); do docker rm -f deis-$$c; done')

full-clean: clean
	$(call ssh_all,'for c in $(COMPONENTS); do docker rmi deis-$$c; done')
