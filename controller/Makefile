include ../includes.mk

.PHONY: all test logs

all: build run

build: check-docker
	docker build -t deis/controller:$(GIT_TAG) .

push: check-docker check-registry check-deisctl
	docker tag deis/controller:$(GIT_TAG) $(REGISTRY)/deis/controller:$(GIT_TAG)
	docker push $(REGISTRY)/deis/controller:$(GIT_TAG)
	deisctl config controller set image=$$DEIS_REGISTRY/deis/controller:$(GIT_TAG)

clean: check-docker check-registry
	docker rmi deis/controller:$(GIT_TAG)
	docker rmi $(REGISTRY)/deis/controller:$(GIT_TAG)

full-clean: check-docker check-registry
	docker images -q deis/controller | xargs docker rmi -f
	docker images -q $(REGISTRY)/deis/controller | xargs docker rmi -f

install: check-deisctl
	deisctl scale controller=1

uninstall: check-deisctl
	deisctl scale controller=0

start: check-deisctl
	deisctl start controller

stop: check-deisctl
	deisctl stop controller

restart: stop start

run: install start

runserver:
	python manage.py runserver

db:
	python manage.py syncdb --migrate --noinput

coverage:
	coverage run manage.py test --noinput api
	coverage html

test: test-unit test-functional

setup-venv:
	@if [ ! -d venv ]; then virtualenv venv; fi
	venv/bin/pip install -q -r requirements.txt -r dev_requirements.txt

test-style: setup-venv
	venv/bin/flake8

test-unit: setup-venv test-style
	venv/bin/python manage.py test --noinput api

test-functional:
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v -timeout 20m ./tests/...
