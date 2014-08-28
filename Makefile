COMPONENTS=builder cache controller database logger registry router

build:
	godep go build ./...

installer:
	rm -rf dist && mkdir -p dist
	godep go build -a -o dist/deisctl .
	command -v upx >/dev/null 2>&1 && upx --best --ultra-brute -q dist/deisctl
	makeself.sh --current --nox11 dist \
		dist/deisctl-`cat deis-version`-`go env GOOS`-`go env GOARCH`.run \
		"Deis Control CLI" "./deisctl refresh-units"

install:
	godep go install -v ./...

test:
	godep go test -v ./...

package:
	rm -f package
	docker build -t deis/deisctl .
	mkdir -p package
	-docker cp `docker run -d deis/deisctl`:/tmp/deisctl.tar.gz package/
	mv package/deisctl.tar.gz package/deisctl-v`cat deis-version`.tar.gz
