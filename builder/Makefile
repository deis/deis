include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/builder

GO_FILES = types.go utils.go utils_test.go
GO_PACKAGES = bin tests
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(DEV_REGISTRY)/$(IMAGE)
BINARIES := extract-domain extract-types extract-version generate-buildhook get-app-config get-app-values publish-release-controller yaml2json-procfile
BINARY_DEST_DIR := image/bin

compile:
	for i in $(BINARIES); do \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -v -ldflags '-s' -o $(BINARY_DEST_DIR)/$$i bin/$$i.go || exit 1; \
	done

build: check-docker
	docker build -t $(IMAGE) image

clean: check-docker check-registry
	for i in $(BINARIES); do \
		rm -rf $(BINARY_DEST_DIR)/$$i ; \
	done
	docker rmi $(IMAGE)

full-clean: check-docker check-registry clean
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

deploy: compile build dev-release restart

test: test-style test-unit test-functional

test-unit:
	godep go test -v .

test-functional:
	@docker history deis/test-etcd >/dev/null 2>&1 || docker pull deis/test-etcd:latest
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES) $(GO_FILES)
	@$(GOFMT) $(GO_PACKAGES) $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(repo_path) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

cedarish/build:
	mkdir -p build
	docker pull progrium/cedarish:latest
	docker save progrium/cedarish:latest | gzip -9 > build/progrium_cedarish_$(shell date "+%Y_%m_%d").tar.gz
