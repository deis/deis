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

DEIS_LAST_ROUTER = $(shell echo $(DEIS_FIRST_ROUTER)\+$(DEIS_NUM_ROUTERS)\-1 | bc)

define ssh_all
  for host in $(DEIS_HOSTS); do ssh -o LogLevel=FATAL -o Compression=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o PasswordAuthentication=no core@$$host -t $(1); done
endef

define echo_cyan
  @echo "\033[0;36m$(subst ",,$(1))\033[0m"
endef

define echo_yellow
  @echo "\033[0;33m$(subst ",,$(1))\033[0m"
endef

ROUTER_UNITS = $(shell seq -f "deis-router.%g.service" -s " " $(DEIS_FIRST_ROUTER) 1 $(DEIS_LAST_ROUTER))

check-fleet:
  @LOCAL_VERSION=`$(FLEETCTL) -version`; \
  REMOTE_VERSION=`ssh -o StrictHostKeyChecking=no core@$(subst :, -p ,$(FLEETCTL_TUNNEL)) fleetctl -version`; \
  if [ "$$LOCAL_VERSION" != "$$REMOTE_VERSION" ]; then \
      echo "Your fleetctl client version should match the server. Local version: $$LOCAL_VERSION, server version: $$REMOTE_VERSION. Uninstall your local version and install the latest build from https://github.com/coreos/fleet/releases"; exit 1; \
  fi
