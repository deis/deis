include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/tests

GO_FILES = $(wildcard *.go)
GO_PACKAGES = dockercli etcdutils mock utils
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

MOCK_STORE_IMAGE = $(IMAGE_PREFIX)mock-store:$(BUILD_TAG)
TEST_ETCD_IMAGE = $(IMAGE_PREFIX)test-etcd:$(BUILD_TAG)
TEST_POSTGRESQL_IMAGE = $(IMAGE_PREFIX)test-postgresql:$(BUILD_TAG)

test: test-smoke

test-smoke: test-style
# TODO: remove GOMAXPROCS=1 when concurrent stdout issues are fixed
	GOMAXPROCS=1 godep go test -tags integration -short -timeout 20m -v -run TestSmoke

test-full: test-style
# TODO: remove GOMAXPROCS=1 when concurrent stdout issues are fixed
	GOMAXPROCS=1 godep go test -tags integration -v -run TestGlobal
	GOMAXPROCS=1 godep go test -tags integration -v -run TestApps
	GOMAXPROCS=1 godep go test -tags integration -v -run TestAuth
	GOMAXPROCS=1 godep go test -tags integration -v -run TestBuilds
	GOMAXPROCS=1 godep go test -tags integration -v -run TestConfig
	GOMAXPROCS=1 godep go test -tags integration -v -run TestDomains
	GOMAXPROCS=1 godep go test -tags integration -v -run TestHealthcheck
	GOMAXPROCS=1 godep go test -tags integration -v -run TestKeys
	GOMAXPROCS=1 godep go test -tags integration -v -run TestPerms
	GOMAXPROCS=1 godep go test -tags integration -v -run TestPs
	GOMAXPROCS=1 godep go test -tags integration -v -run TestReleases
	GOMAXPROCS=1 godep go test -tags integration -v -run TestUsers

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

nuke_from_orbit:
	-docker kill `docker ps -q`
	-docker rm -v `docker ps -a -q`
	-docker rmi -f `docker images -q`
