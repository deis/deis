include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/tests

GO_FILES = apps_test.go auth_test.go builds_test.go config_test.go integration_test.go \
			keys_test.go limits_test.go main.go perms_test.go ps_test.go releases_test.go \
			smoke_test.go tags_test.go
GO_PACKAGES = dockercli etcdutils mock utils
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

MOCK_STORE_IMAGE = $(IMAGE_PREFIX)mock-store:latest
MOCK_STORE_SRC=fixtures/mock-store
MOCK_STORE_GIT_SHA = $(shell git log -n 1 --format="%h" -- $(MOCK_STORE_SRC))

check-envsubst:
ifneq ($(shell which envsubst &> /dev/null; echo $$?), 0)
	$(info Could not find GNU envsubst utility, which usually comes with gettext or gettext-base package)
	exit 1
endif

test: test-smoke

test-smoke: test-style
	godep go test -tags integration -short -timeout 20m -v -run TestSmoke

test-full: test-style
	godep go test -tags integration -v -run TestGlobal
	godep go test -tags integration -v -run TestApps
	godep go test -tags integration -v -run TestAuth
	godep go test -tags integration -v -run TestBuilds
	godep go test -tags integration -v -run TestConfig
	godep go test -tags integration -v -run TestKeys
	godep go test -tags integration -v -run TestPerms
	godep go test -tags integration -v -run TestPs
	godep go test -tags integration -v -run TestReleases

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

mock-store: check-envsubst
	cat $(MOCK_STORE_SRC)/Dockerfile | \
	  GIT_SHA=$(MOCK_STORE_GIT_SHA) envsubst | \
	  ./utils/dockerfeed -d - $(MOCK_STORE_SRC) | \
	  docker build -t $(MOCK_STORE_IMAGE) -

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
