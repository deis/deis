build/container: build/logspout Dockerfile
	docker build --no-cache -t logspout .
	touch build/container

build/logspout: *.go
	go build -o build/logspout

release:
	docker tag logspout progrium/logspout
	docker push progrium/logspout

.PHONY: clean
clean:
	rm -rf build