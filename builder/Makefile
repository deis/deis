build:
	docker build -t deis/builder .

config:
	-etcdctl -C $${ETCD:-127.0.0.1:4001} setdir /deis
	-etcdctl -C $${ETCD:-127.0.0.1:4001} setdir /deis/builder
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/builder/port $${PORT:-22}

run:
	docker run -privileged -e ETCD=$${ETCD:-127.0.0.1:4001} -p $${PORT:-2222}:$${PORT:-22} -rm deis/builder ; exit 0

shell:
	docker run -privileged -e $${ETCD:-127.0.0.1:4001} -t -i -rm deis/builder /bin/bash

clean:
	-docker rmi deis/builder