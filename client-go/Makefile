include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/client-go

GO_FILES = $(wildcard *.go)
GO_PACKAGES = parser cmd controller/api controller/client $(wildcard controller/models/*)
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
IMAGE =  $(IMAGE_PREFIX)/$(COMPONENT):$(BUILD_TAG)

build:
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o deis .
	@$(call check-static-binary,deis) 

install:
	godep go install -v .

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v golang.org/x/tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v golang.org/x/tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/cover
	go get -u golang.org/x/tools/cmd/vet

test: test-style test-unit

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES) $(GO_FILES)
	@$(GOFMT) $(GO_PACKAGES) $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(repo_path) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

test-unit:
	$(GOTEST) ./...
