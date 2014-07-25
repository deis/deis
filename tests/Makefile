test: test-unit test-integration

test-integration:
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestSmoke

test-unit:
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
	go test ./unit

test-smoke:
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -short -v -run TestSmoke

test-full:
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestGlobal
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestApps
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestAuth
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestBuilds
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestClusters
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestConfig
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestKeys
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestPerms
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestPs
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestReleases

nuke_from_orbit:
	-docker kill `docker ps -q`
	-docker rm `docker ps -a -q`
	-docker rmi -f `docker images -q`
