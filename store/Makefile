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
	deisctl scale store=3

uninstall: check-deisctl
	deisctl scale store=0

start: check-deisctl
	deisctl start store

stop: check-deisctl
	deisctl stop store

restart: stop start

run: install start

test: test-unit test-functional

test-unit:
	@echo no unit tests

test-functional:
	@echo no functional tests
