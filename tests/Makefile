test: test-smoke

test-smoke: test-style
	godep go test -tags integration -short -timeout 20m -v -run TestSmoke

test-full: test-style
	godep go test -tags integration -v -run TestGlobal
	godep go test -tags integration -v -run TestApps
	godep go test -tags integration -v -run TestAuth
	godep go test -tags integration -v -run TestBuilds
	godep go test -tags integration -v -run TestConfig
	godep go test -tags integration -v -run TestKeys
	godep go test -tags integration -v -run TestPerms
	godep go test -tags integration -v -run TestPs
	godep go test -tags integration -v -run TestReleases

setup-root-gotools:
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get -u -v code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH

setup-gotools:
	go get -v github.com/golang/lint/golint

test-style:
	go vet -x ./...
	-golint ./...

nuke_from_orbit:
	-docker kill `docker ps -q`
	-docker rm `docker ps -a -q`
	-docker rmi -f `docker images -q`
