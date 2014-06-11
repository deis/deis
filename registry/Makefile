include ../includes.mk

build:
	$(call ssh_all,'cd share/registry && sudo docker build -t deis/registry .')

install: check-fleet
	$(FLEETCTL) load systemd/*

uninstall: check-fleet stop
	$(FLEETCTL) unload systemd/*
	$(FLEETCTL) destroy systemd/*

start: check-fleet
	$(FLEETCTL) start -no-block systemd/*

stop: check-fleet
	$(FLEETCTL) stop -block-attempts=600 systemd/*

restart: stop start

run: install start

clean: uninstall
	$(call ssh_all,'sudo docker rm -f deis-registry')

full-clean: clean
	$(call ssh_all,'sudo docker rmi deis/registry')

test: test-unit test-functional

test-unit: .PHONY
	@if [ ! -d test-unit ]; then \
	    git clone -b 0.6.8 https://github.com/dotcloud/docker-registry test-unit; \
	fi
	@if [ ! -d venv ]; then virtualenv venv; fi
	venv/bin/pip install -q -r test-unit/requirements.txt
	venv/bin/pip install -q -r test-unit/test-requirements.txt
	cd test-unit && \
		DOCKER_REGISTRY_CONFIG=config_test.yml SETTINGS_FLAVOR=test PYTHONPATH=test \
		../venv/bin/python -m unittest discover -s test

test-functional:
	@echo no functional tests

.PHONY:
