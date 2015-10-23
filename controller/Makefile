include ../includes.mk

.PHONY: all test logs

all: build run

SHELL_SCRIPTS = $(wildcard bin/*) $(shell find "." -name '*.sh')

COMPONENT = controller
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(REGISTRY)$(IMAGE)

build: check-docker
	docker build -t $(IMAGE) .

clean: check-docker check-registry
	docker rmi $(IMAGE)

full-clean: check-docker check-registry
	docker images -q $(IMAGE_PREFIX)$(COMPONENT) | xargs docker rmi -f

install: check-deisctl
	deisctl install $(COMPONENT)

uninstall: check-deisctl
	deisctl uninstall $(COMPONENT)

start: check-deisctl
	deisctl start $(COMPONENT)

stop: check-deisctl
	deisctl stop $(COMPONENT)

restart: stop start

run: install start

dev-release: push set-image

push: check-registry
	docker tag -f $(IMAGE) $(DEV_IMAGE)
	docker push $(DEV_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_IMAGE)

release:
	docker push $(IMAGE)

deploy: build dev-release restart

runserver:
	python manage.py runserver

postgres:
	docker start postgres || docker run --restart="always" -d -p 5432:5432 --name postgres postgres:9.4.1
	docker exec postgres createdb -U postgres deis 2>/dev/null || true
	@echo "To use postgres for local development:"
	@echo "    export PGHOST=`docker-machine ip $$(docker-machine active) 2>/dev/null || echo 127.0.0.1`"
	@echo "    export PGPORT=5432"
	@echo "    export PGUSER=postgres"

db:
	python manage.py syncdb --migrate --noinput

coverage:
	coverage run manage.py test --noinput api
	coverage html

test: test-unit test-functional

setup-venv:
	@if [ ! -d venv ]; then virtualenv venv; fi
	venv/bin/pip install --disable-pip-version-check -q -r requirements.txt -r dev_requirements.txt

test-style: setup-venv
	venv/bin/flake8 --show-pep8 --show-source
	shellcheck $(SHELL_SCRIPTS)

test-unit: setup-venv test-style
	venv/bin/coverage run manage.py test --noinput web registry api
	venv/bin/coverage report -m

test-functional:
	@$(MAKE) -C ../tests/ test-etcd
	@$(MAKE) -C ../tests/ test-postgresql
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...
