include ../includes.mk

COMPONENT = deisctl
IMAGE =  $(IMAGE_PREFIX)/$(COMPONENT):$(BUILD_TAG)

build:
	CGO_ENABLED=0 godep go build -a -ldflags '-s' .

installer:
	rm -rf dist && mkdir -p dist
	CGO_ENABLED=0 godep go build -a -ldflags '-s' -o dist/deisctl .
	@if [ ! -d makeself ]; then git clone -b single-binary https://github.com/deis/makeself.git; fi
	PATH=./makeself:$$PATH BINARY=deisctl makeself.sh --bzip2 --current --nox11 dist \
		dist/deisctl-`cat deis-version`-`go env GOOS`-`go env GOARCH`.run "Deis Control Utility" \
		"./deisctl refresh-units \
		&& echo \
		&& echo '\033[0;36mdeisctl\033[0m is in the current directory and unit files are' \
		&& echo 'in \$$HOME/.deis/units. Please move \033[0;36mdeisctl\033[0m to a' \
		&& echo 'directory in your search PATH.' \
		&& echo \
		&& echo 'See http://docs.deis.io/ for documentation.' \
		&& echo"

install:
	godep go install -v .

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

test-style:
	go vet ./...
	-golint ./...

test-unit:
	godep go test -v -cover ./...

test: test-style test-unit
