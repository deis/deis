:title: Addding/Removing Hosts
:description: Considerations for adding or removing Deis hosts.

.. _add_remove_host:

Adding/Removing Hosts
=====================

Most Deis components handle new machines just fine. Care has to be taken when removing machines from the cluster, however, since the deis-store components act as the backing store for all the stateful data Deis needs to function properly.

Note that these instructions follow the Ceph documentation for `removing monitors`_ and `removing OSDs`_. Should these instructions differ significantly from the Ceph documentation, the Ceph documentation should be followed, and a PR to update this documentation would be much appreciated.

Since Ceph uses the Paxos algorithm, it is important to always have enough monitors in the cluster to be able to achieve a majority: 1:1, 2:3, 3:4, 3:5, 4:6, etc. It is always preferable to add a new node to the cluster before removing an old one, if possible.

This documentation will assume a running three-node Deis cluster. We will add a fourth machine to the cluster, then remove the first machine.

Inspecting health
-----------------

Before we begin, we should check the state of the Ceph cluster to be sure it's healthy. We can do this by logging into any machine in the cluster, entering a store container, and then querying Ceph:

.. code-block:: console

    core@deis-1 ~ $ nse deis-store-monitor
    groups: cannot find name for group ID 11
    root@deis-1:/# ceph -s
        cluster c3ff2017-b0a8-4c5a-be00-636560ca567d
         health HEALTH_OK
         monmap e3: 3 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0}, election epoch 8, quorum 0,1,2 deis-1,deis-2,deis-3
         osdmap e18: 3 osds: 3 up, 3 in
          pgmap v31: 960 pgs, 9 pools, 1158 bytes data, 45 objects
                16951 MB used, 31753 MB / 49200 MB avail
                     960 active+clean

We see from the ``pgmap`` that we have 960 placement groups, all of which are ``active+clean``. This is good!

Adding a node
-------------

To add a new node to your Deis cluster, simply provision a new CoreOS machine with the same etcd discovery URL specified in the cloud-config file. When the new machine comes up, it will join the etcd cluster. You can confirm this with ``fleetctl list-machines``.

Since logspout, publisher, store-monitor, and store-daemon are global units, they will be automatically started on the new node.

Once the new machine is running, we can inspect the Ceph cluster health again:

.. code-block:: console

    core@deis-1 ~ $ nse deis-store-monitor
    groups: cannot find name for group ID 11
    root@deis-1:/# ceph -s
        cluster c3ff2017-b0a8-4c5a-be00-636560ca567d
         health HEALTH_WARN clock skew detected on mon.deis-4
         monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 12, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
         osdmap e22: 4 osds: 4 up, 4 in
          pgmap v43: 960 pgs, 9 pools, 1158 bytes data, 45 objects
                22584 MB used, 42352 MB / 65600 MB avail
                     960 active+clean

Note that we have:

.. code-block:: console

     monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 12, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
     osdmap e22: 4 osds: 4 up, 4 in

We have 4 monitors and OSDs. Hooray!

Removing a node
---------------

When removing a node from the cluster that runs a deis-store component, you'll need to tell Ceph that both the store-daemon and store-monitor running on this host will be leaving the cluster. We're going to remove the first node in our cluster, deis-1. That machine has an IP address of ``172.17.8.100``.

Removing an OSD
~~~~~~~~~~~~~~~

Before we can tell Ceph to remove an OSD, we need the OSD ID. We can get this from etcd:

.. code-block:: console

    core@deis-2 ~ $ etcdctl get /deis/store/osds/172.17.8.100
    1

Note: In some cases, we may not know the IP or hostname or the machine we want to remove. In these cases, we can use ``ceph osd tree`` to see the current state of the cluster. This will list all the OSDs in the cluster, and report which ones are down.

Now that we have the OSD's ID, let's remove it. We'll need a shell in any store-monitor or store-daemon container on any host in the cluster (except the one we're removing). In this example, I am on ``deis-2``.

.. code-block:: console

    core@deis-2 ~ $ nse deis-store-monitor
    groups: cannot find name for group ID 11
    root@deis-2:/# ceph osd out 1
    marked out osd.1.


This instructs Ceph to start relocating placement groups on that OSD to another host. We can watch this with ``ceph -w``:

.. code-block:: console

    root@deis-2:/# ceph -w
        cluster c3ff2017-b0a8-4c5a-be00-636560ca567d
         health HEALTH_WARN clock skew detected on mon.deis-4
         monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 12, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
         osdmap e24: 4 osds: 4 up, 3 in
          pgmap v58: 960 pgs, 9 pools, 1158 bytes data, 45 objects
                16900 MB used, 31793 MB / 49200 MB avail
                     960 active+clean

    2014-10-07 17:55:11.900151 mon.0 [INF] pgmap v58: 960 pgs: 960 active+clean; 1158 bytes data, 16900 MB used, 31793 MB / 49200 MB avail; 29 B/s, 3 objects/s recovering
    2014-10-07 17:56:38.860305 mon.0 [INF] pgmap v59: 960 pgs: 960 active+clean; 1158 bytes data, 16900 MB used, 31793 MB / 49200 MB avail

We can see that the placement groups are back in a clean state. We can now stop the daemon. Since the store units are global units, we can't target a specific one to stop. Instead, we log into the host machine and instruct Docker to stop the container:

.. code-block:: console

    core@deis-1 ~ $ docker stop deis-store-daemon
    deis-store-daemon

Back inside a store container on ``deis-2``, we can finally remove the OSD:

.. code-block:: console

    core@deis-2 ~ $ nse deis-store-monitor
    groups: cannot find name for group ID 11
    root@deis-2:/# ceph osd crush remove osd.1
    removed item id 1 name 'osd.1' from crush map
    root@deis-2:/# ceph auth del osd.1
    updated
    root@deis-2:/# ceph osd rm 1
    removed osd.1

For cleanup, we should remove the OSD entry from etcd:

.. code-block:: console

    core@deis-2 ~ $ etcdctl rm /deis/store/osds/172.17.8.100

That's it! If we inspect the health, we see that there are now 3 osds again, and all of our placement groups are ``active+clean``.

.. code-block:: console

    core@deis-2 ~ $ nse deis-store-monitor
    groups: cannot find name for group ID 11
    root@deis-2:/# ceph -s
        cluster c3ff2017-b0a8-4c5a-be00-636560ca567d
         health HEALTH_WARN clock skew detected on mon.deis-4
         monmap e4: 4 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 12, quorum 0,1,2,3 deis-1,deis-2,deis-3,deis-4
         osdmap e28: 3 osds: 3 up, 3 in
          pgmap v81: 960 pgs, 9 pools, 1158 bytes data, 45 objects
                16915 MB used, 31779 MB / 49200 MB avail
                     960 active+clean

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

    root@deis-2:/# ceph mon remove deis-1
    2014-10-07 18:14:38.055584 7fab0d6e7700  0 monclient: hunting for new mon
    2014-10-07 18:14:38.055584 7fab0d6e7700  0 monclient: hunting for new mon
    removed mon.deis-1 at 172.17.8.100:6789/0, there are now 3 monitors
    2014-10-07 18:14:38.072885 7fab0c5e4700  0 -- 172.17.8.101:0/1000361 >> 172.17.8.100:6789/0 pipe(0x7faafc007c90 sd=4 :0 s=1 pgs=0 cs=0 l=1 c=0x7faafc007f00).fault
    2014-10-07 18:14:38.072885 7fab0c5e4700  0 -- 172.17.8.101:0/1000361 >> 172.17.8.100:6789/0 pipe(0x7faafc007c90 sd=4 :0 s=1 pgs=0 cs=0 l=1 c=0x7faafc007f00).fault

Note the faults that follow - this is normal to see when a Ceph client is unable to communicate with a certain monitor. The important line is that we see ``removed mon.deis-1 at 172.17.8.100:6789/0, there are now 3 monitors``.

Finally, let's check the health of the cluster:

.. code-block:: console

    root@deis-2:/# ceph -s
        cluster c3ff2017-b0a8-4c5a-be00-636560ca567d
         health HEALTH_OK
         monmap e5: 3 mons at {deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0,deis-4=172.17.8.103:6789/0}, election epoch 16, quorum 0,1,2 deis-2,deis-3,deis-4
         osdmap e28: 3 osds: 3 up, 3 in
          pgmap v91: 960 pgs, 9 pools, 1158 bytes data, 45 objects
                16927 MB used, 31766 MB / 49200 MB avail
                     960 active+clean

We're done!

Removing the host from etcd
~~~~~~~~~~~~~~~~~~~~~~~~~~~

The etcd cluster still has an entry for the host we've removed, so we'll need to remove this entry.
This can be achieved by making a request to the etcd API. See `remove machines`_ for details.

.. _`remove machines`: https://coreos.com/docs/distributed-configuration/etcd-api/#remove-machines
.. _`removing monitors`: http://ceph.com/docs/v0.80.5/rados/operations/add-or-rm-mons/#removing-monitors
.. _`removing OSDs`: http://docs.ceph.com/docs/v0.80.5/rados/operations/add-or-rm-osds/#removing-osds-manual
