include ../includes.mk

TEMPLATE_IMAGES=daemon monitor gateway
BUILT_IMAGES=base $(TEMPLATE_IMAGES)

DAEMON_IMAGE = $(IMAGE_PREFIX)store-daemon:$(BUILD_TAG)
DAEMON_DEV_IMAGE = $(DEV_REGISTRY)/$(DAEMON_IMAGE)
MONITOR_IMAGE = $(IMAGE_PREFIX)store-monitor:$(BUILD_TAG)
MONITOR_DEV_IMAGE = $(DEV_REGISTRY)/$(MONITOR_IMAGE)
GATEWAY_IMAGE = $(IMAGE_PREFIX)store-gateway:$(BUILD_TAG)
GATEWAY_DEV_IMAGE = $(DEV_REGISTRY)/$(GATEWAY_IMAGE)

build: check-docker
	@# Build base as normal
	docker build -t deis/store-base:$(BUILD_TAG) base/
	$(foreach I, $(TEMPLATE_IMAGES), \
		sed -e "s/#FROM is generated dynamically by the Makefile/FROM deis\/store-base:${BUILD_TAG}/" $(I)/Dockerfile.template > $(I)/Dockerfile ; \
		docker build -t deis/store-$(I):$(BUILD_TAG) $(I)/ ; \
		rm $(I)/Dockerfile ; \
	)

push: check-docker check-registry check-deisctl
	$(foreach I, $(BUILT_IMAGES), \
		docker tag deis/store-$(I):$(BUILD_TAG) $(REGISTRY)/deis/store-$(I):$(BUILD_TAG) ; \
		docker push $(REGISTRY)/deis/store-$(I):$(BUILD_TAG) ; \
	)

	$(foreach I, $(TEMPLATE_IMAGES), \
		deisctl config store-$(I) set image=$(REGISTRY)/deis/store-$(I):$(BUILD_TAG) ; \
	)

clean: check-docker check-registry
	$(foreach I, $(BUILT_IMAGES), \
		docker rmi deis/store-$(I):$(BUILD_TAG) ; \
		docker rmi $(REGISTRY)/deis/store-$(I):$(BUILD_TAG) ; \
	)

full-clean: check-docker check-registry
	$(foreach I, $(BUILT_IMAGES), \
		docker images -q deis/store-$(I) | xargs docker rmi -f ; \
		docker images -q $(REGISTRY)/deis/store-$(I) | xargs docker rmi -f ; \
	)

install: check-deisctl
	deisctl install store-monitor
	deisctl install store-daemon
	deisctl install store-gateway

uninstall: check-deisctl
	deisctl uninstall store-gateway
	deisctl uninstall store-daemon
	deisctl uninstall store-monitor

start: check-deisctl
	deisctl start store-monitor
	deisctl start store-daemon
	deisctl start store-gateway

stop: check-deisctl
	deisctl stop store-gateway
	deisctl stop store-daemon
	deisctl stop store-monitor

restart: stop start

run: install start

dev-release: push set-image

push: check-registry
	docker tag $(DAEMON_IMAGE) $(DAEMON_DEV_IMAGE)
	docker push $(DAEMON_DEV_IMAGE)
	docker tag $(MONITOR_IMAGE) $(MONITOR_DEV_IMAGE)
	docker push $(MONITOR_DEV_IMAGE)
	docker tag $(GATEWAY_IMAGE) $(GATEWAY_DEV_IMAGE)
	docker push $(GATEWAY_DEV_IMAGE)

set-image: check-deisctl
	deisctl config store-daemon set image=$(DAEMON_DEV_IMAGE)
	deisctl config store-monitor set image=$(MONITOR_DEV_IMAGE)
	deisctl config store-gateway set image=$(GATEWAY_DEV_IMAGE)

release:
	docker push $(DAEMON_IMAGE)
	docker push $(MONITOR_IMAGE)
	docker push $(GATEWAY_IMAGE)

deploy: build dev-release restart

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=$(CURDIR)/../tests/_vendor:$(GOPATH) go test -v ./tests/...
