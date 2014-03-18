build:
	docker build -t deis/cache .

run:
	docker run -rm -p :6379 -e ETCD=$${ETCD:-127.0.0.1:4001} deis/cache
	exit 0

shell:
	docker run -t -i -rm deis/cache /bin/bash

clean:
	-docker rmi deis/cache
