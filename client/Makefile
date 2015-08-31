include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/client

GO_FILES = $(wildcard *.go)
GO_PACKAGES = parser cmd controller/api controller/client $(wildcard controller/models/*) $(wildcard pkg/*)
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
IMAGE =  $(IMAGE_PREFIX)/$(COMPONENT):$(BUILD_TAG)

build:
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o deis .
	@$(call check-static-binary,deis)

install: build
	cp deis $$GOPATH/bin

installer: build
	@if [ ! -d makeself ]; then git clone -b single-binary https://github.com/deis/makeself.git; fi
	PATH=./makeself:$$PATH BINARY=deis makeself.sh --bzip2 --current --nox11 . \
		deis-cli-`cat deis-version`-`go env GOOS`-`go env GOARCH`.run \
		"Deis CLI" "echo \
		&& echo 'deis is in the current directory. Please' \
		&& echo 'move deis to a directory in your search PATH.' \
		&& echo \
		&& echo 'See http://docs.deis.io/ for documentation.' \
		&& echo"

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
