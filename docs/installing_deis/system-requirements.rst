:title: System Requirements
:description: System requirements for provisioning Deis.

.. _system-requirements:

System Requirements
===================

When deploying Deis, it's important to provision machines with adequate resources. Deis is a
highly-available distributed system, which means that Deis components and your deployed applications
will move around the cluster onto healthy hosts as hosts leave the cluster for various reasons
(failures, reboots, autoscalers, etc.). Because of this, you should have ample spare resources on
any machine in your cluster to withstand the additional load of running services for failed machines.

Resources
---------

Deis components consume approximately 2 - 2.5GB of memory across the cluster, and approximately
30GB of hard disk space. Because each machine should be able to absorb additional load should a
machine fail, each machine must have:

* At least 4GB of RAM (more is better)
* At least 40GB of hard disk space

Note that these estimates are for Deis and CoreOS only, and there should be ample room for deployed
applications.

Running smaller machines will likely result in increased system load and has been known to result
in component failures, issues with etcd/fleet, and other problems.

.. _cluster-size:

Cluster size
------------

For the :ref:`deis-store <Store>` component to work properly, clusters must have at least
three nodes. The ``etcd`` service must always be able to obtain a quorum, and the Ceph
data store must maintain at least three replicas of persistent data.

See `optimal etcd cluster size`_ and `etcd disaster recovery`_ for further information.

.. important::

    Deis clusters of less than three nodes are unsupported.

If running multiple (at least three) machines of an adequate size is unreasonable, it is recommended to
investigate the `Dokku`_ project instead. Dokku is `sponsored`_ by Deis and is ideal for environments
where a highly-available distributed system is not necessary (i.e. local development, testing, etc.).

Network
-------

.. include:: ../_includes/_private-network.rst

.. _`dokku`: https://github.com/progrium/dokku
.. _`optimal etcd cluster size`: https://coreos.com/docs/cluster-management/scaling/etcd-optimal-cluster-size/
.. _`etcd disaster recovery`: https://github.com/coreos/etcd/blob/master/Documentation/admin_guide.md#disaster-recovery
.. _`sponsored`: http://deis.io/deis-sponsors-dokku/
