include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/logspout

GO_FILES = $(wildcard *.go)
GO_PACKAGES = adapters/deis
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
DOCKER_IMAGE := $(IMAGE_PREFIX)$(COMPONENT)
RELEASE_IMAGE := $(DOCKER_IMAGE):$(BUILD_TAG)
DEV_DOCKER_IMAGE := $(REGISTRY)$(RELEASE_IMAGE)

DEV_LOGSPOUT_IMAGE = $(DOCKER_IMAGE):dev
DEV_ETCD_IMAGE = quay.io/coreos/etcd:v2.1.2
DEV_LOG_SERVER_IMAGE = $(IMAGE_PREFIX)logger:$(BUILD_TAG)

build: check-docker
	docker build --no-cache -t $(RELEASE_IMAGE) .

clean:

full-clean:
	docker images -q $(DOCKER_IMAGE) | xargs docker rmi -f
	docker images -q $(DEV_DOCKER_IMAGE) | xargs docker rmi -f
	docker images -q $(RELEASE_IMAGE) | xargs docker rmi -f

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
	docker tag -f $(RELEASE_IMAGE) $(DEV_DOCKER_IMAGE)
	docker push $(DEV_DOCKER_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_DOCKER_IMAGE)

release:
	docker push $(RELEASE_IMAGE)

deploy: build dev-release restart

test: test-style test-unit test-functional

test-functional:
	@echo no functional tests

test-style:
# Run gofmt first to display its output
	gofmt -l .
# Run it again and assess whether it failed; if so, exit
	@gofmt -l . | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi

test-unit:
	@echo no unit tests

start-dev-etcd:
	@docker history $(DEV_ETCD_IMAGE) &> /dev/null || docker pull $(DEV_ETCD_IMAGE) &> /dev/null
	@docker run --name dev-etcd -d --net=host $(DEV_ETCD_IMAGE) &> /dev/null
	@echo "dev-etcd container is started."

stop-dev-etcd:
	@docker kill dev-etcd &> /dev/null
	@docker rm dev-etcd &> /dev/null
	@echo "dev-etcd container is stopped."

run-dev-log-server:
	@docker history $(DEV_LOG_SERVER_IMAGE) &> /dev/null || $(MAKE) -C ../logger/ build
	@docker run --name deis-logger --rm --net=host -e EXTERNAL_PORT=514 -e HOST=127.0.0.1 -e LOGSPOUT=ignore -v /var/lib/deis/store:/data $(DEV_LOG_SERVER_IMAGE)

dev-build:
# Build go code BEFORE trying to build the container-- the point is to fail fast
	@go build deis.go config.go adapter.go transport.go && docker build -f Dockerfile.dev -t $(DEV_LOGSPOUT_IMAGE) .

dev-run: dev-build
	@docker run --name logspout-dev --rm -e ETCD_HOST=127.0.0.1 -e ETCD_PORT=2379 -e DEBUG=true -e LOGSPOUT=ignore --net=host -v /var/run/docker.sock:/var/run/docker.sock $(DEV_LOGSPOUT_IMAGE)
