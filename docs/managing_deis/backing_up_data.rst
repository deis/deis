:title: Backing up Data
:description: Backing up stateful data on Deis.

.. _backing_up_data:

Backing up Data
========================

While applications deployed on Deis follow the Twelve-Factor methodology and are thus stateless,
Deis maintains platform state in two places: data containers and etcd.

Data containers
---------------
Data containers are simply Docker containers that expose a volume which is shared with another container.
The components with data containers are builder, logger, and registry. Since these are just
Docker containers, they can be exported with ordinary Docker commands:

.. code-block:: console

    dev $ fleetctl ssh deis-builder.service
    coreos $ sudo docker export deis-builder-data > /home/coreos/deis-builder-data-backup.tar
    dev $ fleetctl ssh deis-logger.service
    coreos $ sudo docker export deis-logger-data > /home/coreos/deis-logger-data-backup.tar
    dev $ fleetctl ssh deis-registry.service
    coreos $ sudo docker export deis-registry-data > /home/coreos/deis-registry-data-backup.tar

Importing looks very similar:

.. code-block:: console

    dev $ fleetctl ssh deis-builder.service
    coreos $ cat /home/coreos/deis-builder-data-backup.tar | sudo docker import - deis-builder-data
    dev $ fleetctl ssh deis-logger.service
    coreos $ cat /home/coreos/deis-logger-data-backup.tar | sudo docker import - deis-logger-data
    dev $ fleetctl ssh deis-registry.service
    coreos $ cat /home/coreos/deis-registry-data-backup.tar | sudo docker import - deis-registry-data

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
