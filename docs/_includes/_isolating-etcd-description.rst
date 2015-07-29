The Deis Control Plane, Data Plane, and Router Mesh components all depend on an
etcd cluster for service discovery and configuration.

Whether built for evaluation or to host production applications, when managing a
small Deis cluster (three to five nodes), it is reasonable to accept the
platform's default behavior wherein etcd runs on every node within the cluster.

In larger Deis clusters however, running etcd on every node can have a
deleterious effect on overall cluster performance since it increases the time
required for nodes to reach consensus on writes and leader elections. In such
cases, it is beneficial to isolate etcd to a small, fixed number of nodes.  All
other nodes in the Deis cluster may run an etcd proxy.  Proxies will forward
read and write requests to active participants in the etcd cluster (leader or
followers) without affecting the time required for etcd nodes to reach consensus
on writes or leader elections.

.. note::

    The benefit of running an etcd proxy on any node not running a full etcd
    process is that any container or service depending on etcd can connect to
    etcd easily via ``localhost`` from any node in the Deis cluster.

    Also see `CoreOS cluster architecture documentation`_ for further details.

.. _`CoreOS cluster architecture documentation`: https://coreos.com/os/docs/latest/cluster-architectures.html#production-cluster-with-central-services