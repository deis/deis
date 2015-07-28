include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/tests

GO_FILES = $(wildcard *.go)
GO_PACKAGES = dockercli etcdutils mock utils
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

SHELL_SCRIPTS = $(wildcard bin/*) \
								$(wildcard fixtures/mock-store/bin/*) \
								$(wildcard fixtures/test-postgresql/rootfs/*/*)

MOCK_STORE_IMAGE = $(IMAGE_PREFIX)mock-store:latest
TEST_ETCD_IMAGE = $(IMAGE_PREFIX)test-etcd:latest
TEST_POSTGRESQL_IMAGE = $(IMAGE_PREFIX)test-postgresql:latest

test: test-smoke

test-smoke: test-style
	godep go test -tags integration -short -timeout 20m -v -run TestSmoke

test-full: test-style
	godep go test -tags integration -v -run TestGlobal
	godep go test -tags integration -v -run TestApps
	godep go test -tags integration -v -run TestAuth
	godep go test -tags integration -v -run TestBuilds
	godep go test -tags integration -v -run TestConfig
	godep go test -tags integration -v -run TestDomains
	godep go test -tags integration -v -run TestKeys
	godep go test -tags integration -v -run TestPerms
	godep go test -tags integration -v -run TestPs
	godep go test -tags integration -v -run TestReleases
	godep go test -tags integration -v -run TestUsers

mock-store:
	docker build -t $(MOCK_STORE_IMAGE) fixtures/mock-store/

test-etcd:
	docker build -t $(TEST_ETCD_IMAGE) fixtures/test-etcd/

test-postgresql:
	docker build -t $(TEST_POSTGRESQL_IMAGE) fixtures/test-postgresql/

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES) $(GO_FILES)
	@$(GOFMT) $(GO_PACKAGES) $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(repo_path) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...
	shellcheck $(SHELL_SCRIPTS)

nuke_from_orbit:
	-docker kill `docker ps -q`
	-docker rm -v `docker ps -a -q`
	-docker rmi -f `docker images -q`
