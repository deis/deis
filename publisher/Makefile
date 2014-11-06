include ../includes.mk

DOCKER_IMAGE := deis/publisher
BUILD_IMAGE := $(DOCKER_IMAGE)-build
RELEASE_IMAGE := $(DOCKER_IMAGE):$(BUILD_TAG)
REMOTE_IMAGE := $(REGISTRY)/$(RELEASE_IMAGE)

build: check-docker
	docker build -t $(BUILD_IMAGE) .
	docker cp `docker run -d $(BUILD_IMAGE)`:/go/bin/publisher image/
	docker build -t $(RELEASE_IMAGE) image
	rm -rf image/publisher

clean: check-docker check-registry
	docker rmi $(RELEASE_IMAGE) $(REMOTE_IMAGE)

full-clean: check-docker check-registry
	docker images -q $(DOCKER_IMAGE) | xargs docker rmi -f
	docker images -q $(REGISTRY)/$(DOCKER_IMAGE) | xargs docker rmi -f

install: check-deisctl
	deisctl install publisher

dev-release: push set-image

push: check-registry
	docker tag $(RELEASE_IMAGE) $(REMOTE_IMAGE)
	docker push $(REMOTE_IMAGE)

set-image: check-deisctl
	deisctl config publisher set image=$(REMOTE_IMAGE)

release: check-docker
	docker push $(DOCKER_IMAGE)
	
deploy: build dev-release restart

restart: stop start

run: install start

start: check-deisctl
	deisctl start publisher

stop: check-deisctl
	deisctl stop publisher

test: test-unit

test-unit:
	godep go test -v ./...

uninstall: check-deisctl
	deisctl uninstall publisher
