:title: Adding/Removing Hosts
:description: Considerations for adding or removing Deis hosts.

.. _add_remove_host:

Adding/Removing Hosts
=====================

Most Deis components handle new machines just fine. Care has to be taken when removing machines from
the cluster, however, since the deis-store components act as the backing store for all the
stateful data Deis needs to function properly.

Note that these instructions follow the Ceph documentation for `removing monitors`_ and `removing OSDs`_.
Should these instructions differ significantly from the Ceph documentation, the Ceph documentation
should be followed, and a PR to update this documentation would be much appreciated.

Since Ceph uses the Paxos algorithm, it is important to always have enough monitors in the cluster
to be able to achieve a majority: 1:1, 2:3, 3:4, 3:5, 4:6, etc. It is always preferable to add
a new node to the cluster before removing an old one, if possible.

This documentation will assume a running three-node Deis cluster.
We will add a fourth machine to the cluster, then remove the first machine.

Inspecting health
-----------------

Before we begin, we should check the state of the Ceph cluster to be sure it's healthy.
To do this, we use ``deis-store-admin`` - see :ref:`using-store-admin`.

.. code-block:: console

    core@deis-1 ~ $ nse deis-store-admin
    root@deis-1:/# ceph -s
        cluster 20038e38-4108-4e79-95d4-291d0eef2949
         health HEALTH_OK
         monmap e3: 3 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0}, election epoch 16, quorum 0,1,2 deis-1,deis-2,deis-3
         mdsmap e10: 1/1/1 up {0=deis-2=up:active}, 2 up:standby
         osdmap e36: 3 osds: 3 up, 3 in
          pgmap v2096: 1344 pgs, 12 pools, 369 MB data, 448 objects
                24198 MB used, 23659 MB / 49206 MB avail
                1344 active+clean

We see from the ``pgmap`` that we have 1344 placement groups, all of which are ``active+clean``. This is good!

Adding a node
-------------

To add a new node to your Deis cluster, simply provision a new CoreOS machine with the same
etcd discovery URL specified in the cloud-config file. When the new machine comes up, it will join the etcd cluster.
You can confirm this with ``fleetctl list-machines``.

Since the store components are global units, they will be automatically started on the new node.

Once the new machine is running, we can inspect the Ceph cluster health again:

.. code-block:: console

    root@deis-1:/# ceph -s
        cluster 20038e38-4108-4e79-95d4-291d0eef2949
         health HEALTH_WARN 4 pgs recovering; 7 pgs recovery_wait; 31 pgs stuck unclean; recovery 325/1353 objects degraded (24.021%); clock skew detected on mon.deis-4
         monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 20, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
         mdsmap e11: 1/1/1 up {0=deis-2=up:active}, 3 up:standby
         osdmap e40: 4 osds: 4 up, 4 in
          pgmap v2172: 1344 pgs, 12 pools, 370 MB data, 451 objects
                29751 MB used, 34319 MB / 65608 MB avail
                325/1353 objects degraded (24.021%)
                  88 active
                   7 active+recovery_wait
                1245 active+clean
                   4 active+recovering
      recovery io 2302 kB/s, 2 objects/s
      client io 204 B/s wr, 0 op/s

Note that we are in a ``HEALTH_WARN`` state, and we have placement groups recovering. Ceph is
copying data to our new node. We can query the status of this until it completes. Then, we should
we something like:

.. code-block:: console

    root@deis-1:/# ceph -s
        cluster 20038e38-4108-4e79-95d4-291d0eef2949
         health HEALTH_OK
         monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 20, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
         mdsmap e11: 1/1/1 up {0=deis-2=up:active}, 3 up:standby
         osdmap e40: 4 osds: 4 up, 4 in
          pgmap v2216: 1344 pgs, 12 pools, 372 MB data, 453 objects
                29749 MB used, 34324 MB / 65608 MB avail
                    1344 active+clean
      client io 409 B/s wr, 0 op/s

We're back in a ``HEALTH_OK``, and note the following:

.. code-block:: console

    monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 20, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
    mdsmap e11: 1/1/1 up {0=deis-2=up:active}, 3 up:standby
    osdmap e40: 4 osds: 4 up, 4 in

We have 4 monitors, OSDs, and metadata servers. Hooray!

.. note::

    If you have applied the `custom firewall script`_ to your cluster, you will have to run this
    script again and reboot your nodes for iptables to remove the duplicate rules.

Removing a node
---------------

When removing a node from the cluster that runs a deis-store component, you'll need to tell Ceph
that the store services on this host will be leaving the cluster.
In this example we're going to remove the first node in our cluster, deis-1.
That machine has an IP address of ``172.17.8.100``.

.. _removing_an_osd:

Removing an OSD
~~~~~~~~~~~~~~~

Before we can tell Ceph to remove an OSD, we need the OSD ID. We can get this from etcd:

.. code-block:: console

    core@deis-2 ~ $ etcdctl get /deis/store/osds/172.17.8.100
    2

Note: In some cases, we may not know the IP or hostname or the machine we want to remove.
In these cases, we can use ``ceph osd tree`` to see the current state of the cluster.
This will list all the OSDs in the cluster, and report which ones are down.

Now that we have the OSD's ID, let's remove it. We'll need a shell in any store container
on any host in the cluster (except the one we're removing). In this example, I am on ``deis-2``.

.. code-block:: console

    core@deis-2 ~ $ nse deis-store-admin
    root@deis-2:/# ceph osd out 2
    marked out osd.2.

This instructs Ceph to start relocating placement groups on that OSD to another host. We can watch this with ``ceph -w``:

.. code-block:: console

    root@deis-2:/# ceph -w
        cluster 20038e38-4108-4e79-95d4-291d0eef2949
         health HEALTH_WARN 4 pgs recovery_wait; 151 pgs stuck unclean; recovery 654/1365 objects degraded (47.912%); clock skew detected on mon.deis-4
         monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 20, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
         mdsmap e11: 1/1/1 up {0=deis-2=up:active}, 3 up:standby
         osdmap e42: 4 osds: 4 up, 3 in
         pgmap v2259: 1344 pgs, 12 pools, 373 MB data, 455 objects
                23295 MB used, 24762 MB / 49206 MB avail
                654/1365 objects degraded (47.912%)
                 151 active
                   4 active+recovery_wait
                1189 active+clean
      recovery io 1417 kB/s, 1 objects/s
      client io 113 B/s wr, 0 op/s

    2014-11-04 06:45:07.940731 mon.0 [INF] pgmap v2260: 1344 pgs: 142 active, 3 active+recovery_wait, 1199 active+clean; 373 MB data, 23301 MB used, 24757 MB / 49206 MB avail; 619/1365 objects degraded (45.348%); 1724 kB/s, 0 keys/s, 1 objects/s recovering
    2014-11-04 06:45:17.948788 mon.0 [INF] pgmap v2261: 1344 pgs: 141 active, 4 active+recovery_wait, 1199 active+clean; 373 MB data, 23301 MB used, 24757 MB / 49206 MB avail; 82 B/s rd, 0 op/s; 619/1365 objects degraded (45.348%); 843 kB/s, 0 keys/s, 0 objects/s recovering
    2014-11-04 06:45:18.962420 mon.0 [INF] pgmap v2262: 1344 pgs: 140 active, 5 active+recovery_wait, 1199 active+clean; 373 MB data, 23318 MB used, 24740 MB / 49206 MB avail; 371 B/s rd, 0 B/s wr, 0 op/s; 618/1365 objects degraded (45.275%); 0 B/s, 0 keys/s, 0 objects/s recovering
    2014-11-04 06:45:23.347089 mon.0 [INF] pgmap v2263: 1344 pgs: 130 active, 5 active+recovery_wait, 1209 active+clean; 373 MB data, 23331 MB used, 24727 MB / 49206 MB avail; 379 B/s rd, 0 B/s wr, 0 op/s; 572/1365 objects degraded (41.905%); 2323 kB/s, 0 keys/s, 4 objects/s recovering
    2014-11-04 06:45:37.970125 mon.0 [INF] pgmap v2264: 1344 pgs: 129 active, 4 active+recovery_wait, 1211 active+clean; 373 MB data, 23336 MB used, 24722 MB / 49206 MB avail; 568/1365 objects degraded (41.612%); 659 kB/s, 2 keys/s, 1 objects/s recovering
    2014-11-04 06:45:40.006110 mon.0 [INF] pgmap v2265: 1344 pgs: 129 active, 4 active+recovery_wait, 1211 active+clean; 373 MB data, 23336 MB used, 24722 MB / 49206 MB avail; 568/1365 objects degraded (41.612%); 11 B/s, 3 keys/s, 0 objects/s recovering
    2014-11-04 06:45:43.034215 mon.0 [INF] pgmap v2266: 1344 pgs: 129 active, 4 active+recovery_wait, 1211 active+clean; 373 MB data, 23344 MB used, 24714 MB / 49206 MB avail; 1010 B/s wr, 0 op/s; 568/1365 objects degraded (41.612%)
    2014-11-04 06:45:44.048059 mon.0 [INF] pgmap v2267: 1344 pgs: 129 active, 4 active+recovery_wait, 1211 active+clean; 373 MB data, 23344 MB used, 24714 MB / 49206 MB avail; 1766 B/s wr, 0 op/s; 568/1365 objects degraded (41.612%)
    2014-11-04 06:45:48.366555 mon.0 [INF] pgmap v2268: 1344 pgs: 129 active, 4 active+recovery_wait, 1211 active+clean; 373 MB data, 23345 MB used, 24713 MB / 49206 MB avail; 576 B/s wr, 0 op/s; 568/1365 objects degraded (41.612%)

Eventually, the cluster will return to a clean state and will once again report ``HEALTH_OK``.
Then, we can stop the daemon. Since the store units are global units, we can't target a specific
one to stop. Instead, we log into the host machine and instruct Docker to stop the container.

Reminder: make sure you're logged into the machine you're removing from the cluster!

.. code-block:: console

    core@deis-1 ~ $ docker stop deis-store-daemon
    deis-store-daemon

Back inside a store container on ``deis-2``, we can finally remove the OSD:

.. code-block:: console

    core@deis-2 ~ $ nse deis-store-admin
    root@deis-2:/# ceph osd crush remove osd.2
    removed item id 2 name 'osd.2' from crush map
    root@deis-2:/# ceph auth del osd.2
    updated
    root@deis-2:/# ceph osd rm 2
    removed osd.2

For cleanup, we should remove the OSD entry from etcd:

.. code-block:: console

    core@deis-2 ~ $ etcdctl rm /deis/store/osds/172.17.8.100

That's it! If we inspect the health, we see that there are now 3 osds again, and all of our placement groups are ``active+clean``.

.. code-block:: console

    core@deis-2 ~ $ nse deis-store-admin
    root@deis-2:/# ceph -s
        cluster 20038e38-4108-4e79-95d4-291d0eef2949
         health HEALTH_OK
         monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 20, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
         mdsmap e11: 1/1/1 up {0=deis-2=up:active}, 3 up:standby
         osdmap e46: 3 osds: 3 up, 3 in
          pgmap v2338: 1344 pgs, 12 pools, 375 MB data, 458 objects
                23596 MB used, 24465 MB / 49206 MB avail
                    1344 active+clean
      client io 326 B/s wr, 0 op/s

Removing a monitor
~~~~~~~~~~~~~~~~~~

Removing a monitor is much easier. First, we remove the etcd entry so any clients that are using Ceph won't use the monitor for connecting:

.. code-block:: console

    $ etcdctl rm /deis/store/hosts/172.17.8.100

Within 5 seconds, confd will run on all store clients and remove the monitor from the ``ceph.conf`` configuration file.

Next, we stop the container:

.. code-block:: console

    core@deis-1 ~ $ docker stop deis-store-monitor
    deis-store-monitor


Back on another host, we can again enter a store container and then remove this monitor:

.. code-block:: console

    core@deis-2 ~ $ nse deis-store-admin
    root@deis-2:/# ceph mon remove deis-1
    removed mon.deis-1 at 172.17.8.100:6789/0, there are now 3 monitors
    2014-11-04 06:57:59.712934 7f04bc942700  0 monclient: hunting for new mon
    2014-11-04 06:57:59.712934 7f04bc942700  0 monclient: hunting for new mon

Note that there may be faults that follow - this is normal to see when a Ceph client is
unable to communicate with a monitor. The important line is that we see ``removed mon.deis-1 at 172.17.8.100:6789/0, there are now 3 monitors``.

Finally, let's check the health of the cluster:

.. code-block:: console

    root@deis-2:/# ceph -s
        cluster 20038e38-4108-4e79-95d4-291d0eef2949
         health HEALTH_OK
         monmap e5: 3 mons at {deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 26, quorum 0,1,2 deis-2,deis-3,deis-4
         mdsmap e17: 1/1/1 up {0=deis-4=up:active}, 3 up:standby
         osdmap e47: 3 osds: 3 up, 3 in
          pgmap v2359: 1344 pgs, 12 pools, 375 MB data, 458 objects
                23605 MB used, 24455 MB / 49206 MB avail
                    1344 active+clean
      client io 816 B/s wr, 0 op/s

We're done!

Removing a metadata server
~~~~~~~~~~~~~~~~~~~~~~~~~~

Like the daemon, we'll just stop the Docker container for the metadata service.

Reminder: make sure you're logged into the machine you're removing from the cluster!

.. code-block:: console

    core@deis-1 ~ $ docker stop deis-store-metadata
    deis-store-metadata

This is actually all that's necessary. Ceph provides a ``ceph mds rm`` command, but has no
documentation for it. See: http://docs.ceph.com/docs/hammer/rados/operations/control/#mds-subsystem

Removing the host from etcd
~~~~~~~~~~~~~~~~~~~~~~~~~~~

The etcd cluster still has an entry for the host we've removed, so we'll need to remove this entry.
This can be achieved by making a request to the etcd API. See `remove machines`_ for details.

.. _`custom firewall script`: https://github.com/deis/deis/blob/master/contrib/util/custom-firewall.sh
.. _`remove machines`: https://coreos.com/docs/distributed-configuration/etcd-api/#remove-machines
.. _`removing monitors`: http://ceph.com/docs/hammer/rados/operations/add-or-rm-mons/#removing-monitors
.. _`removing OSDs`: http://docs.ceph.com/docs/hammer/rados/operations/add-or-rm-osds/#removing-osds-manual

Automatic Host Removal
======================

The ``contrib/coreos/user-data.example`` provides 2 units, ``graceful-etcd-shutdown.service`` and
``graceful-ceph-shutdown.service``, that contain some experimental logic to clean-up a Deis node's
cluster membership before reboot, shutdown or halt events. The units can be used independently or
together.

The ``graceful-etcd-shutdown`` unit is useful for any Deis node running its own etcd. To be used, it
must be enabled and started.

.. code-block:: console

    root@deis-1:/# systemctl enable graceful-etcd-shutdown
    root@deis-1:/# systemctl start graceful-etcd-shutdown

The ``graceful-ceph-shutdown`` script is only useful for nodes running deis-store components. To be used,
the unit requires that the optional ``deis-store-admin`` component is installed.

.. code-block:: console

    root@deis-1:/# deisctl install store-admin
    root@deis-1:/# deisctl start store-admin

Then the unit should be enabled and started.

.. code-block:: console

    root@deis-1:/# systemctl enable graceful-ceph-shutdown
    root@deis-1:/# systemctl start graceful-ceph-shutdown

At this point your node is ready to be gracefully removed whenever a halt, shutdown or reboot event occurs.
The graceful shutdown units insert themselves ahead of the etcd and Ceph units in the shutdown order. This
allows them to perform preemptive actions on etcd and Ceph while they are still healthy and in the cluster.

The units make use of the script ``/opt/bin/graceful-shutdown.sh`` to remove the node from the cluster. For
Ceph, this means determining if the Ceph cluster is healthy and has enough nodes to return to health - if it
does, it will remove its OSD and wait for the Ceph cluster to return to health. Once it is healthy, it will
remove its monitor and continue to shut down Ceph components. The end result should be a Ceph cluster that
returns its status as ``health_ok``.

For etcd, the script remove its etcd member and delete itself from the CoreOS discovery url.
