:title: HA database
:description: Highly-available database configuration.

.. _ha_database:

HA Database
=========================

Currently, a Deis installation contains one database server. Should that host go down, the cluster
will become unavailable. Application traffic will be unaffected, but no changes can be made.
A highly-available database cluster is necessary, and is planned for Deis 1.0.

The current plan is to use traditional PostgreSQL replication and use an etcd lock to determine
which host is the master. A slave would take over as master should that host go down.

Details on this approach can be found in GitHub issue `#923`_. Feedback is welcome.

.. note::

  Some Deis users connect their installation to an external HA database cluster. This can be
  accomplished by setting relevant options in etcd. For details on replacing the database, see
  :ref:`database_settings`.

.. _`#923`: https://github.com/deis/deis/issues/923
