build/container: build/logspout Dockerfile
	docker build --no-cache -t logspout .
	touch build/container

build/logspout: *.go
	go build -o build/logspout

.PHONY: clean
clean:
	rm -rf build