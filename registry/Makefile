build:
	docker build -t deis/registry .

run:
	docker run -rm -p :5000 -e ETCD=$${ETCD:-127.0.0.1:4001} deis/registry
	exit 0

shell:
	docker run -t -i -rm deis/registry /bin/bash

clean:
	-docker rmi deis/registry
