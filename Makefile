COMPONENTS=builder cache controller database logger registry router

build:
	godep go build ./...

installer:
	rm -rf dist && mkdir -p dist
	godep go build -a -o dist/deisctl .
	@if [ ! -d makeself ]; then git clone -b deisctl-hack https://github.com/deis/makeself.git; fi
	PATH=./makeself:$$PATH makeself.sh --bzip2 --nox11 --target /usr/local/bin dist \
		dist/deisctl-`cat deis-version`-`go env GOOS`-`go env GOARCH`.run \
		"Deis Control Utility" "deisctl refresh-units"

install:
	godep go install -v ./...

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

test-style:
	go vet ./...
	-golint *.go client/*.go cmd/*.go config/*.go constant/*.go lock/*.go update/*.go utils/*.go

test: test-style
	godep go test -v -cover ./...

package:
	rm -f package
	docker build -t deis/deisctl .
	mkdir -p package
	-docker cp `docker run -d deis/deisctl`:/tmp/deisctl.tar.gz package/
	mv package/deisctl.tar.gz package/deisctl-v`cat deis-version`.tar.gz
