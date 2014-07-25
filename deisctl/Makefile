COMPONENTS=builder cache controller database logger registry router

copy-units:
	@$(foreach C, $(COMPONENTS), \
	mkdir -p units/$(C) ; \
	cp ../$(C)/systemd/* units/$(C) ;\
	)

test:
	godep go test -v -cover ./...

build: copy-units
	cd deisctl && godep go build .
	@echo "binary placed in ./deisctl/deisctl"

install: copy-units
	godep go install ./...
	@echo "binaries placed in \$$GOPATH/bin"

package: copy-units
	docker build -t deis/deisctl .
	CONTAINER=`docker run -d deis/deisctl` ; \
	docker cp $$CONTAINER:/tmp/deisctl.tar.gz package/ ; \
	docker rm -f $$CONTAINER
	@echo "package placed in ./package/deisctl.tar.gz"
