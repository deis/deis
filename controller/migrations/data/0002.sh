#!/usr/bin/env bash

ETCD_PORT=${ETCD_PORT:-4001}
ETCD="$HOST:$ETCD_PORT"
ETCDCTL="etcdctl -C $ETCD"

# April 8, 2015: If registrationEnabled key exists, migrate it to registrationMode and delete it.

if [[ "$($ETCDCTL get /deis/migrations/data/0002 2> /dev/null)" != "done" ]];
then
    if $ETCDCTL ls /deis/controller | grep -q '/deis/controller/registrationEnabled'
    then
        if [[ "$($ETCDCTL get /deis/controller/registrationEnabled 2> /dev/null)" == "false" ]]
        then
        $ETCDCTL set /deis/controller/registrationMode "disabled"
        else
          $ETCDCTL set /deis/controller/registrationMode "enabled"
        fi

        $ETCDCTL rm /deis/controller/registrationEnabled
    else
        echo "registrationEnabled key doesn't exist, skipping migration"
    fi

    $ETCDCTL set /deis/migrations/data/0002 "done"
fi
