include ../includes.mk

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/deis/deisctl

GO_FILES = $(wildcard *.go)
GO_PACKAGES = backend client cmd config units utils
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

COMPONENT = $(notdir $(repo_path))
IMAGE =  $(IMAGE_PREFIX)/$(COMPONENT):$(BUILD_TAG)

build:
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' .

installer:
	rm -rf dist && mkdir -p dist
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o dist/deisctl .
	@if [ ! -d makeself ]; then git clone -b single-binary https://github.com/deis/makeself.git; fi
	PATH=./makeself:$$PATH BINARY=deisctl makeself.sh --bzip2 --current --nox11 dist \
		dist/deisctl-`cat deis-version`-`go env GOOS`-`go env GOARCH`.run "Deis Control Utility" \
		"./deisctl refresh-units \
		&& echo \
		&& echo 'deisctl is in the current directory and unit files are' \
		&& echo 'in \$$HOME/.deis/units. Please move deisctl to a' \
		&& echo 'directory in your search PATH.' \
		&& echo \
		&& echo 'See http://docs.deis.io/ for documentation.' \
		&& echo"

install:
	godep go install -v .

test: test-style test-unit

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES) $(GO_FILES)
	@$(GOFMT) $(GO_PACKAGES) $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(repo_path) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

test-unit:
# TODO: remove GOMAXPROCS=1 when concurrent stdout issues are fixed
	GOMAXPROCS=1 $(GOTEST) ./...
