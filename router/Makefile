build:
	docker build -t deis/router .

env:
	export ROUTER_HOST=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/router/host`
	export ROUTER_PORT=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/router/port`
	export ROUTER_WORKER_CONNECTIONS=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/router/workerConnections`
	export ROUTER_WORKER_PROCESSES=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/router/workerProcesses`
	env

config:
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/router/host $${ROUTER_HOST:-127.0.0.1}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/router/port $${ROUTER_PORT:-80}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/router/workerConnections $${ROUTER_WORKER_CONNECTIONS:-1024}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/router/workerProcesses $${ROUTER_WORKER_PROCESSES:-2}

run:
	docker run -d -p 80:80 -e ETCD=$${ETCD:-172.17.42.1:4001} -e HOST=$${HOST:-172.17.42.1} --name deis-router deis/router

shell:
	docker run -t -i -rm deis/router /bin/bash

kill:
	docker kill deis-router
	docker rm deis-router

clean:
	docker rmi deis/router

test:
	docker build .
