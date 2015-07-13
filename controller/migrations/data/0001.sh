#!/usr/bin/env bash

if [[ "$(etcdctl get /deis/migrations/data/0001 2> /dev/null)" != "done" ]];
then
    for i in $(etcdctl ls /deis/domains 2> /dev/null);
    do
        for j in $(etcdctl get $i);
        do
            etcdctl set /deis/domains/$j "$(basename $i)" 1> /dev/null;
            echo "migrated $j"
        done;
        etcdctl rm $i;
    done
    etcdctl set /deis/migrations/data/0001 "done"
fi
