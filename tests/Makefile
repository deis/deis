test: test-smoke

test-smoke: test-style
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -short -timeout 20m -v -run TestSmoke

test-full: test-style
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestGlobal
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestApps
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestAuth
	GOPATH=$(CURDIR)/_vendor:$(GOPATH) \
		go test -tags integration -v -run TestBuilds
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

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

test-style:
	go vet -x ./...
	-golint; golint dockercli; golint etcdutils; golint mock; golint utils

nuke_from_orbit:
	-docker kill `docker ps -q`
	-docker rm `docker ps -a -q`
	-docker rmi -f `docker images -q`
