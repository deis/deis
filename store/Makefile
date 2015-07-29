include ../includes.mk

TEMPLATE_IMAGES=admin daemon gateway metadata monitor
BUILT_IMAGES=base $(TEMPLATE_IMAGES)

ADMIN_IMAGE = $(IMAGE_PREFIX)store-admin:$(BUILD_TAG)
ADMIN_DEV_IMAGE = $(REGISTRY)$(ADMIN_IMAGE)
DAEMON_IMAGE = $(IMAGE_PREFIX)store-daemon:$(BUILD_TAG)
DAEMON_DEV_IMAGE = $(REGISTRY)$(DAEMON_IMAGE)
GATEWAY_IMAGE = $(IMAGE_PREFIX)store-gateway:$(BUILD_TAG)
GATEWAY_DEV_IMAGE = $(REGISTRY)$(GATEWAY_IMAGE)
METADATA_IMAGE = $(IMAGE_PREFIX)store-metadata:$(BUILD_TAG)
METADATA_DEV_IMAGE = $(REGISTRY)$(METADATA_IMAGE)
MONITOR_IMAGE = $(IMAGE_PREFIX)store-monitor:$(BUILD_TAG)
MONITOR_DEV_IMAGE = $(REGISTRY)$(MONITOR_IMAGE)

build: check-docker
	@# Build base as normal
	docker build -t $(IMAGE_PREFIX)store-base:$(BUILD_TAG) base/
	$(foreach I, $(TEMPLATE_IMAGES), \
		sed -e "s,#FROM is generated dynamically by the Makefile,FROM $(IMAGE_PREFIX)store-base:${BUILD_TAG}," $(I)/Dockerfile.template > $(I)/Dockerfile ; \
		docker build -t $(IMAGE_PREFIX)store-$(I):$(BUILD_TAG) $(I)/ || exit 1; \
		rm $(I)/Dockerfile ; \
	)

clean: check-docker check-registry
	$(foreach I, $(BUILT_IMAGES), \
		docker rmi $(IMAGE_PREFIX)store-$(I):$(BUILD_TAG) ; \
		docker rmi $(REGISTRY)/$(IMAGE_PREFIX)store-$(I):$(BUILD_TAG) ; \
	)

full-clean: check-docker check-registry
	$(foreach I, $(BUILT_IMAGES), \
		docker images -q $(IMAGE_PREFIX)store-$(I) | xargs docker rmi -f ; \
		docker images -q $(REGISTRY)/$(IMAGE_PREFIX)store-$(I) | xargs docker rmi -f ; \
	)

install: check-deisctl
	deisctl install store-monitor
	deisctl install store-daemon
	deisctl install store-metadata
	deisctl install store-volume
	deisctl scale store-gateway=1

uninstall: check-deisctl
	deisctl scale store-gateway=0
	deisctl uninstall store-volume
	deisctl uninstall store-metadata
	deisctl uninstall store-daemon
	deisctl uninstall store-monitor

start: check-deisctl
	deisctl start store-monitor
	deisctl start store-daemon
	deisctl start store-metadata
	deisctl start store-volume
	deisctl start store-gateway@*

stop: check-deisctl
	deisctl stop store-gateway@*
	deisctl stop store-volume
	deisctl stop store-metadata
	deisctl stop store-daemon
	deisctl stop store-monitor

restart: stop start

run: install start

dev-release: push set-image

push: check-registry
	docker tag -f $(ADMIN_IMAGE) $(ADMIN_DEV_IMAGE)
	docker push $(ADMIN_DEV_IMAGE)
	docker tag -f $(DAEMON_IMAGE) $(DAEMON_DEV_IMAGE)
	docker push $(DAEMON_DEV_IMAGE)
	docker tag -f $(GATEWAY_IMAGE) $(GATEWAY_DEV_IMAGE)
	docker push $(GATEWAY_DEV_IMAGE)
	docker tag -f $(METADATA_IMAGE) $(METADATA_DEV_IMAGE)
	docker push $(METADATA_DEV_IMAGE)
	docker tag -f $(MONITOR_IMAGE) $(MONITOR_DEV_IMAGE)
	docker push $(MONITOR_DEV_IMAGE)

set-image: check-deisctl
	deisctl config store-admin set image=$(ADMIN_DEV_IMAGE)
	deisctl config store-daemon set image=$(DAEMON_DEV_IMAGE)
	deisctl config store-gateway set image=$(GATEWAY_DEV_IMAGE)
	deisctl config store-metadata set image=$(METADATA_DEV_IMAGE)
	deisctl config store-monitor set image=$(MONITOR_DEV_IMAGE)

release:
	docker push $(ADMIN_IMAGE)
	docker push $(DAEMON_IMAGE)
	docker push $(GATEWAY_IMAGE)
	docker push $(METADATA_IMAGE)
	docker push $(MONITOR_IMAGE)

deploy: build dev-release restart

test: test-style test-unit test-functional

test-functional:
	@$(MAKE) -C ../tests/ test-etcd
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...

test-style:
	@echo no style tests

test-unit:
	@echo no unit tests
