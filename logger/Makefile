build:
	docker build -t deis/logger .

run:
	docker run -rm -p :514 -e ETCD=$${ETCD:-127.0.0.1:4001} deis/logger
	exit 0

shell:
	docker run -t -i -rm deis/logger /bin/bash

clean:
	-docker rmi deis/logger
