#
# Deis Makefile
#

include includes.mk

COMPONENTS=builder cache controller database logger registry router
START_ORDER=logger database cache registry controller builder router

all: build run

dev-registry: check-docker
	@docker run -d -p 5000:5000 --name registry registry:0.8.1
	@echo
	@echo "To configure the registry for local Deis development:"
	@echo "    export DEIS_REGISTRY=`boot2docker ip 2>/dev/null`:5000"

discovery-url:
	@mv -f contrib/coreos/user-data contrib/coreos/user-data.orig
	@cat contrib/coreos/user-data.orig | \
	sed -e "s,# discovery: https://discovery.etcd.io/,discovery: https://discovery.etcd.io/)," | \
	sed -e "s,discovery: https://discovery.etcd.io/.*,discovery: $$(curl -s -w '\n' https://discovery.etcd.io/new)," > \
	contrib/coreos/user-data

build: check-docker
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) build || exit 1;)

push: check-docker check-registry
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) push || exit 1;)
	
full-clean:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) full-clean || exit 1;)

install:
	@$(foreach C, $(START_ORDER), $(MAKE) -C $(C) install || exit 1;)

uninstall:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) uninstall || exit 1;)

start:
	@$(foreach C, $(START_ORDER), $(MAKE) -C $(C) start || exit 1;)

stop:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) stop || exit 1;)

restart: stop start

run: install start

test: test-components push test-integration

test-components:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) test || exit 1;)

test-integration:
	$(MAKE) -C tests/ test-full

test-smoke:
	$(MAKE) -C tests/ test-smoke
