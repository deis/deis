ifndef FLEETCTL
  FLEETCTL = fleetctl --strict-host-key-checking=false
endif

ifndef FLEETCTL_TUNNEL
$(error You need to set FLEETCTL_TUNNEL to the IP address of a server in the cluster.)
endif

ifndef DEIS_NUM_INSTANCES
  DEIS_NUM_INSTANCES = 1
endif

ifndef DEIS_HOSTS
  DEIS_HOSTS = $(shell seq -f "172.17.8.%g" -s " " 100 1 `expr $(DEIS_NUM_INSTANCES) + 99` )
endif

ifndef DEIS_NUM_ROUTERS
  DEIS_NUM_ROUTERS = 1
endif

ifndef DEIS_FIRST_ROUTER
  DEIS_FIRST_ROUTER = 1
endif

ifndef DEIS_NUM_DATABASES
  DEIS_NUM_DATABASES = 2
endif

ifndef DEIS_FIRST_DATABASE
  DEIS_FIRST_DATABASE = 1
endif

define ssh_all
  for host in $(DEIS_HOSTS); do ssh -o LogLevel=FATAL -o Compression=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o PasswordAuthentication=no core@$$host -t $(1); done
endef

define rsync_all
  for host in $(DEIS_HOSTS); do rsync -Pave "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" --exclude=venv/ --exclude=.git/ --exclude='*.pyc' $(SELF_DIR)/* core@$$host:/home/core/share; done
endef

define echo_cyan
  @echo "\033[0;36m$(subst ",,$(1))\033[0m"
endef

define echo_yellow
  @echo "\033[0;33m$(subst ",,$(1))\033[0m"
endef


SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

DEIS_LAST_ROUTER = $(shell echo $(DEIS_FIRST_ROUTER)\+$(DEIS_NUM_ROUTERS)\-1 | bc)
DEIS_LAST_DATABASE = $(shell echo $(DEIS_FIRST_DATABASE)\+$(DEIS_NUM_DATABASES)\-1 | bc)

ROUTER_UNITS = $(shell seq -f "deis-router.%g.service" -s " " $(DEIS_FIRST_ROUTER) 1 $(DEIS_LAST_ROUTER))
DATABASE_UNITS = $(shell seq -f "deis-database.%g.service" -s " " $(DEIS_FIRST_DATABASE) 1 $(DEIS_LAST_DATABASE))

check-fleet:
  @LOCAL_VERSION=`$(FLEETCTL) -version`; \
  REMOTE_VERSION=`ssh -o StrictHostKeyChecking=no core@$(subst :, -p ,$(FLEETCTL_TUNNEL)) fleetctl -version`; \
  if [ "$$LOCAL_VERSION" != "$$REMOTE_VERSION" ]; then \
      echo "Your fleetctl client version should match the server. Local version: $$LOCAL_VERSION, server version: $$REMOTE_VERSION. Uninstall your local version and install the latest build from https://github.com/coreos/fleet/releases"; exit 1; \
  fi
