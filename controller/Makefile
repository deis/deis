include ../includes.mk

.PHONY: all test logs

all: build run

build:
	$(call ssh_all,'cd share/controller && sudo docker build -t deis/controller .')

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
	$(call ssh_all,'sudo docker rm -f deis-controller')

full-clean: clean
	$(call ssh_all,'sudo docker rmi deis/controller')

runserver:
	python manage.py runserver

db:
	python manage.py syncdb --migrate --noinput

coverage:
	coverage run manage.py test --noinput api
	coverage html

flake8:
	flake8

test: test-unit test-functional

test-unit:
	@if [ ! -d venv ]; then virtualenv venv; fi
	venv/bin/pip install -q -r requirements.txt
	venv/bin/python manage.py test --noinput api

test-functional:
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -timeout 30m -v ./test/...
