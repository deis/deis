build:
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' .

upx:
	upx -9 --brute bumpver && upx -t bumpver

clean:
	rm -f bumpver

setup-gotools-sudo:
	sudo GOPATH=/tmp/tmpGOPATH go get code.google.com/p/go.tools/cmd/cover
	sudo GOPATH=/tmp/tmpGOPATH go get code.google.com/p/go.tools/cmd/vet
	sudo rm -rf /tmp/tmpGOPATH
	@echo 'Please ensure that $$GOPATH/bin is in your $$PATH.'

setup-gotools:
	go get github.com/tools/godep
	go get github.com/golang/lint/golint
	@echo 'Please ensure that $$GOPATH/bin is in your $$PATH.'

test-style:
	go vet -x ./...
	-golint .

test: test-style
	godep go test -v -cover ./...
