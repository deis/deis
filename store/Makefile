include ../includes.mk

TEMPLATE_IMAGES=daemon monitor gateway
BUILT_IMAGES=base $(TEMPLATE_IMAGES)

build: check-docker
	@# Build base as normal
	docker build -t deis/store-base:$(BUILD_TAG) base/
	$(foreach I, $(TEMPLATE_IMAGES), \
		cp $(I)/Dockerfile.template $(I)/Dockerfile ; \
		sed -i '' "s/#FROM is generated dynamically by the Makefile/FROM deis\/store-base:${BUILD_TAG}/" $(I)/Dockerfile ; \
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
	deisctl scale store-monitor=3
	deisctl scale store-daemon=3
	deisctl scale store-gateway=1

uninstall: check-deisctl
	deisctl scale store-gateway=0
	deisctl scale store-daemon=0
	deisctl scale store-monitor=0

start: check-deisctl
	deisctl start store-monitor@1
	deisctl start store-monitor@2
	deisctl start store-monitor@3
	deisctl start store-daemon@1
	deisctl start store-daemon@2
	deisctl start store-daemon@3
	deisctl start store-gateway@1

stop: check-deisctl
	deisctl stop store-gateway@1
	deisctl stop store-daemon@1
	deisctl stop store-daemon@2
	deisctl stop store-daemon@3
	deisctl stop store-monitor@1
	deisctl stop store-monitor@2
	deisctl stop store-monitor@3

restart: stop start

run: install start

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@echo no functional tests
