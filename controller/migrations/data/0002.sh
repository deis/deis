#!/usr/bin/env bash

# April 8, 2015: If registrationEnabled key exists, migrate it to registrationMode and delete it.

if [[ "$(etcdctl get /deis/migrations/data/0002 2> /dev/null)" != "done" ]];
then
    if etcdctl ls /deis/controller | grep -q '/deis/controller/registrationEnabled'
    then
        if [[ "$(etcdctl get /deis/controller/registrationEnabled 2> /dev/null)" == "false" ]]
        then
        etcdctl set /deis/controller/registrationMode "disabled"
        else
          etcdctl set /deis/controller/registrationMode "enabled"
        fi

        etcdctl rm /deis/controller/registrationEnabled
    else
        echo "registrationEnabled key doesn't exist, skipping migration"
    fi

    etcdctl set /deis/migrations/data/0002 "done"
fi
