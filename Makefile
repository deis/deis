#
# Deis Makefile
#

include includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis

GO_PACKAGES = version
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENTS=builder controller database logger logspout publisher registry router $(STORE_IF_STATEFUL)
START_ORDER=publisher $(STORE_IF_STATEFUL) logger logspout database registry controller builder router
CLIENTS=client deisctl

all: build run

dev-registry: check-docker
	@docker inspect registry >/dev/null 2>&1 && docker start registry || docker run --restart="always" -d -p 5000:5000 --name registry registry:0.9.1
	@echo
	@echo "To use a local registry for Deis development:"
	@echo "    export DEV_REGISTRY=`docker-machine ip $$(docker-machine active 2>/dev/null) 2>/dev/null || echo $(HOST_IPADDR) `:5000"

dev-cluster: discovery-url
	vagrant up
	ssh-add ~/.vagrant.d/insecure_private_key
	deisctl config platform set sshPrivateKey=$(HOME)/.vagrant.d/insecure_private_key
	deisctl config platform set domain=local3.deisapp.com
	deisctl config platform set enablePlacementOptions=true
	deisctl install platform

discovery-url:
	@for i in 1 2 3 4 5; do \
		URL=`curl -s -w '\n' https://discovery.etcd.io/new?size=$$DEIS_NUM_INSTANCES`; \
		if [ ! -z $$URL ]; then \
			sed -e "s,discovery: #DISCOVERY_URL,discovery: $$URL," contrib/coreos/user-data.example > contrib/coreos/user-data; \
			echo "Wrote $$URL to contrib/coreos/user-data"; \
		    break; \
		fi; \
		if [ $$i -eq 5 ]; then \
			echo "Failed to contact https://discovery.etcd.io after $$i tries"; \
		else \
			sleep 3; \
		fi \
	done

build: check-docker
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) build &&) echo done
	@$(foreach C, $(CLIENTS), $(MAKE) -C $(C) build &&) echo done

clean:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) clean &&) echo done
	@$(foreach C, $(CLIENTS), $(MAKE) -C $(C) clean &&) echo done

full-clean:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) full-clean &&) echo done

install:
	@$(foreach C, $(START_ORDER), $(MAKE) -C $(C) install &&) echo done

uninstall:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) uninstall &&) echo done

start:
	@$(foreach C, $(START_ORDER), $(MAKE) -C $(C) start &&) echo done

stop:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) stop &&) echo done

restart: stop start

run: install start

dev-release:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) dev-release &&) echo done

push:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) push &&) echo done

set-image:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) set-image &&) echo done

release: check-registry
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) release &&) echo done

deploy: build dev-release restart

setup-gotools:
	go get -u -v github.com/tools/godep
	go get -u -v github.com/golang/lint/golint
	go get -u -v golang.org/x/tools/cmd/cover
	go get -u -v golang.org/x/tools/cmd/vet

setup-root-gotools:
# "go vet" and "go cover" must be installed as root on some systems
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v golang.org/x/tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v golang.org/x/tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

test: test-style test-unit test-functional push test-integration

test-functional:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) test-functional &&) echo done

test-unit:
	@$(foreach C, $(COMPONENTS), $(MAKE) -C $(C) test-unit &&) echo done
	@$(foreach C, pkg $(CLIENTS), $(MAKE) -C $(C) test-unit &&) echo done

test-integration:
	$(MAKE) -C tests/ test-full

test-smoke:
	$(MAKE) -C tests/ test-smoke

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES)
	@$(GOFMT) $(GO_PACKAGES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(GO_PACKAGES_REPO_PATH)
	@for i in $(addsuffix /...,$(GO_PACKAGES)); do \
		$(GOLINT) $$i; \
	done
	@$(foreach C, tests pkg $(CLIENTS) $(COMPONENTS), $(MAKE) -C $(C) test-style &&) echo done

commit-hook:
	cp contrib/util/commit-msg .git/hooks/commit-msg
