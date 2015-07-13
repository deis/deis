etcd-sync
=========

An etcd-based sync module, aiming at implementing the Go sync pkg over etcd for cluster-wide synchronization.

Installation
============

~~~shell
go get github.com/leeor/etcd-sync
~~~

Usage
=====

At this time, only a simple mutex has been implemented.

## EtcdMutex

~~~go
mutex := NewMutexFromServers([]string{"http://127.0.0.1:4001"}, key, 0)
mutex.Lock()
// do some critical stuff
mutex.Unlock()
~~~
