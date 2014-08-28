build/container: stage/logspout Dockerfile
	docker build --no-cache -t logspout .
	touch build/container

build/logspout: *.go
	GOOS=linux GOARCH=amd64 go build -o build/logspout

stage/logspout: build/logspout
	mkdir -p stage
	cp build/logspout stage/logspout

release:
	docker tag logspout progrium/logspout
	docker push progrium/logspout

.PHONY: clean
clean:
	rm -rf build