test:
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
	go test -tags integration -timeout 30m -v ./...

test-smoke:
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
	go test -tags integration -timeout 10m -short -v ./...

nuke_from_orbit:
	-docker kill `docker ps -q`
	-docker rm `docker ps -a -q`
	-docker rmi -f `docker images -q`
