:title: Operational tasks
:description: Common operational tasks for your Deis cluster.

.. _operational_tasks:

Operational tasks
~~~~~~~~~~~~~~~~~

Inspecting store
================
It is sometimes helpful to query the :Ref:`Store` component to ask about the health of the Ceph cluster.
To do this, log into any machine running a ``store-monitor`` or ``store-daemon`` service. Then,
``nse deis-store-monitor`` or ``nse deis-store-daemon`` and issue a ``ceph -s``. This should output the
health of the cluster like:

.. code-block:: console

    cluster 6506db0c-9eae-4bb6-a40a-95954dd3c4c3
    health HEALTH_OK
    monmap e3: 3 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0}, election epoch 8, quorum 0,1,2 deis-1,deis-2,deis-3
    osdmap e7: 3 osds: 3 up, 3 in
    pgmap v14: 192 pgs, 3 pools, 0 bytes data, 0 objects
    19378 MB used, 28944 MB / 49200 MB avail
    192 active+clean

If you see ``HEALTH_OK``, this means everything is working as it should.
Note also ``monmap e3: 3 mons at...`` which means all three monitor containers are up and responding,
and ``osdmap e7: 3 osds: 3 up, 3 in`` which means all three daemon containers are up and running.

We can also see from the ``pgmap`` that we have 192 placement groups, all of which are ``active+clean``.

For additional information on troubleshooting Ceph, see `troubleshooting`_.

Managing users
==============

There are two classes of Deis users: normal users and administrators.

* Users can use most of the features of Deis - creating and deploying applications, adding/removing domains, etc.
* Administrators can perform all the actions that users can, but they can also create, edit, and destroy clusters.

The first user created on a Deis installation is automatically an administrator.

Promoting users to administrators
---------------------------------

You can use the ``deis perms`` command to promote a user to an administrator:

.. code-block:: console

    $ deis perms:create john --admin

.. _`troubleshooting`: http://docs.ceph.com/docs/firefly/rados/troubleshooting/
