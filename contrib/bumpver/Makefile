build:
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' .

upx:
	upx -9 --brute bumpver && upx -t bumpver

clean:
	rm -f bumpver

test-style:
	go vet -x ./...
	-golint .

test: test-style
	godep go test -v -cover ./...
