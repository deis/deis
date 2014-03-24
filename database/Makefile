build:
	docker build -t deis/database .

env:
	export PG_ADMIN_USER=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/database/admin-user`
	export PG_ADMIN_PASS=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/database/admin-pass`
	export PG_USER_NAME=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/database/user`
	export PG_USER_PASS=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/database/password`
	export PG_USER_DB=`etcdctl -C $${ETCD:-127.0.0.1:4001} get /deis/database/name`
	env

config:
	-etcdctl -C $${ETCD:-127.0.0.1:4001} setdir /deis
	-etcdctl -C $${ETCD:-127.0.0.1:4001} setdir /deis/database
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/database/engine $${DATABASE_ENGINE:-databaseql_psycopg2}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/database/admin-user $${DATABASE_ADMIN_USER:-database}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/database/admin-pass $${DATABASE_ADMIN_PASS:-changeme123}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/database/host $${DATABASE_HOST:-127.0.0.1}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/database/user $${DATABASE_USER:-docker}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/database/password $${DATABASE_PASS:-changeme123}
	etcdctl -C $${ETCD:-127.0.0.1:4001} set /deis/database/name $${DATABASE_NAME:-docker}

run:
	docker run -p :5432 -e ETCD=$${ETCD:-127.0.0.1:4001} deis/database
	exit 0

shell:
	docker run -t -i -rm deis/database /bin/bash

clean:
	-docker rmi deis/database
