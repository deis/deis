:title: Backing up Data
:description: Backing up stateful data on Deis.

.. _backing_up_data:

Backing up Data
========================

While applications deployed on Deis follow the Twelve-Factor methodology and are thus stateless,
Deis maintains platform state in two places: the :ref:`Store` component, and in etcd.

Store component
---------------
The store component runs `Ceph`_, and is used by the :ref:`Database` and :ref:`Registry` components
as a data store. This enables the components themselves to freely move around the cluster while
their state is backed by store.

The store component is configured to still operate in a degraded state, and will automatically
recover should a host fail and then rejoin the cluster. Total data loss of Ceph is only possible
if all of the store containers are removed. However, backup of Ceph is fairly straightforward.

Data in Ceph is stored on the filesystem in ``/var/lib/ceph``, and metadata information is stored
within Ceph. Ceph provides the ability to take snapshots of storage pools with the `rados`_ command.

Using pg_dump
-------------
Since the database component runs PostgreSQL, ``pg_dumpall`` can also be used to generate a text
dump of the database.

.. code-block:: console

    dev $ fleetctl ssh deis-database.service
    coreos $ nse deis-database
    coreos $ sudo -u postgres pg_dumpall > pg_dump.sql

etcd
----
Service state and fleet scheduling data is stored in etcd. Unfortunately, there is currently no
recommended backup solution for etcd. However, there is a third-party tool called `etcd-dump`_ which
can be used to dump the data stored in etcd.

Official backup recommendations for etcd are forthcoming. The CoreOS team is tracking etcd update
documentation in `#683`_.

.. _`#683`: https://github.com/coreos/etcd/issues/683
.. _`etcd-dump`: https://github.com/AaronO/etcd-dump
.. _`Ceph`: http://ceph.com
.. _`rados`: http://ceph.com/docs/master/man/8/rados
