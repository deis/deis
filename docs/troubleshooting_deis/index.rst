:title: Troubleshooting Deis
:description: Resolutions for common issues encountered when running Deis.

.. _troubleshooting_deis:

Troubleshooting Deis
====================

Common issues that users have run into when provisioning Deis are detailed below.

A deis-store component fails to start
-------------------------------------

The store component is the most complex component of Deis. As such, there are many ways for it to fail.
Recall that the store components represent Ceph services as follows:

* ``store-monitor``: http://ceph.com/docs/giant/man/8/ceph-mon/
* ``store-daemon``: http://ceph.com/docs/giant/man/8/ceph-osd/
* ``store-gateway``: http://ceph.com/docs/giant/radosgw/
* ``store-metadata``: http://ceph.com/docs/giant/man/8/ceph-mds/
* ``store-volume``: a system service which mounts a `Ceph FS`_ volume to be used by the controller and logger components

Log output for store components can be viewed with ``deisctl status store-<component>`` (such as
``deisctl status store-volume``). Additionally, the Ceph health can be queried by entering
a store container with ``nse deis-store-monitor`` and then issuing a ``ceph -s``. This should output the
health of the cluster like:

.. code-block:: console

    core@deis-1 ~ $ nse deis-store-monitor
    root@deis-1:/# ceph -s
        cluster 20038e38-4108-4e79-95d4-291d0eef2949
         health HEALTH_OK
         monmap e3: 3 mons at {deis-1=172.17.8.100:6789/0,deis-2=172.17.8.101:6789/0,deis-3=172.17.8.102:6789/0}, election epoch 16, quorum 0,1,2 deis-1,deis-2,deis-3
         mdsmap e10: 1/1/1 up {0=deis-2=up:active}, 2 up:standby
         osdmap e36: 3 osds: 3 up, 3 in
          pgmap v2096: 1344 pgs, 12 pools, 369 MB data, 448 objects
                24198 MB used, 23659 MB / 49206 MB avail
                1344 active+clean

If you see ``HEALTH_OK``, this means everything is working as it should.
Note also ``monmap e3: 3 mons at...`` which means all three monitor containers are up and responding,
``mdsmap e10: 1/1/1 up...`` which means all three metadata containers are up and responding,
and ``osdmap e7: 3 osds: 3 up, 3 in`` which means all three daemon containers are up and running.

We can also see from the ``pgmap`` that we have 1344 placement groups, all of which are ``active+clean``.

For additional information on troubleshooting Ceph, see `troubleshooting`_. Common issues with
specific store components are detailed below.

store-monitor
~~~~~~~~~~~~~

The monitor is the first store component to start, and is required for any of the other store
components to function properly. If a ``deisctl list`` indicates that any of the monitors are failing,
it is likely due to a host issue. Common failure scenarios include not
having adequate free storage on the host node - in that case, monitors will fail with errors similar to:

.. code-block:: console

  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.053693 7fd0586a6700  0 mon.deis-staging-node1@0(leader).data_health(6) update_stats avail 1% total 5960684 used 56655
  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.053770 7fd0586a6700 -1 mon.deis-staging-node1@0(leader).data_health(6) reached critical levels of available space on
  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.053772 7fd0586a6700  0 ** Shutdown via Data Health Service **
  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.053821 7fd056ea3700 -1 mon.deis-staging-node1@0(leader) e3 *** Got Signal Interrupt ***
  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.053834 7fd056ea3700  1 mon.deis-staging-node1@0(leader) e3 shutdown
  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.054000 7fd056ea3700  0 quorum service shutdown
  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.054002 7fd056ea3700  0 mon.deis-staging-node1@0(shutdown).health(6) HealthMonitor::service_shutdown 1 services
  Oct 29 20:04:00 deis-staging-node1 sh[1158]: 2014-10-29 20:04:00.054065 7fd056ea3700  0 quorum service shutdown

This is typically only an issue when deploying Deis on bare metal, as most cloud providers have adequately
large volumes.

store-daemon
~~~~~~~~~~~~

The daemons are responsible for actually storing the data on the filesystem. The cluster is configured
to allow writes with just one daemon running, but the cluster will be running in a degraded state, so
restoring all daemons to a running state as quickly as possible is paramount.

Daemons can be safely restarted with ``deisctl restart store-daemon``, but this will restart all daemons,
resulting in downtime of the storage cluster until the daemons recover. Alternatively, issuing a
``sudo systemctl restart deis-store-daemon`` on the host of the failing daemon will restart just
that daemon.

store-gateway
~~~~~~~~~~~~~

The gateway runs Apache and a FastCGI server to communicate with the cluster. Restarting the gateway
will result in a short downtime for the registry component (and will prevent the database from
backing up), but those components should recover as soon as the gateway comes back up.

store-metadata
~~~~~~~~~~~~~~

The metadata servers are required for the **volume** to function properly. Only one is active at
any one time, and the rest operate as hot standbys. The monitors will promote a standby metadata
server should the active one fail.

store-volume
~~~~~~~~~~~~

Without functioning monitors, daemons, and metadata servers, the volume service will likely hang
indefinitely (or restart constantly). If the controller or logger happen to be running on a host with a
failing store-volume, application logs will be lost until the volume recovers.

Note that store-volume requires CoreOS >= 471.1.0 for the CephFS kernel module.

Any component fails to start
----------------------------

Use ``deisctl status <component>`` to view the status of the component.
You can also use ``deisctl journal <component>`` to tail logs for a component, or ``deisctl list``
to list all components.

Failed initializing SSH client
------------------------------

A ``deisctl`` command fails with: 'Failed initializing SSH client: ssh: handshake failed: ssh: unable to authenticate'.
Did you remember to add your SSH key to the ssh-agent? ``ssh-add -L`` should list the key you used
to provision the servers. If it's not there, ``ssh-add -K /path/to/your/key``.

All the given peers are not reachable
-------------------------------------

A ``deisctl`` command fails with: 'All the given peers are not reachable (Tried to connect to each peer twice and failed)'.
The most common cause of this issue is that a new discovery URL wasn't generated and updated in
``contrib/coreos/user-data`` before the cluster was launched. Each Deis cluster must have a unique
discovery URL, or else ``etcd`` will try and fail to connect to old hosts. Try destroying the cluster
and relaunching the cluster with a fresh discovery URL.

You can use ``make discovery-url`` to automatically fetch a new discovery URL.

Could not find unit template...
-------------------------------

If you built ``deisctl`` locally or didn't use its installer, you may see an error like this:

    .. code-block:: console

        $ deisctl install platform

        Storage subsystem...
        Could not find unit template for store-daemon

This is because ``deisctl`` could not find unit files for Deis locally. Run
``deisctl help refresh-units`` to see where ``deisctl`` searches, and then run a command such as
``deisctl refresh-units --tag=v1.0.2``, or set the ``$DEISCTL_UNITS`` environment variable to a directory
containing the unit files.

Other issues
------------

Running into something not detailed here? Please `open an issue`_ or hop into #deis on Freenode IRC and we'll help!

.. _`Ceph FS`: https://ceph.com/docs/giant/cephfs/
.. _`open an issue`: https://github.com/deis/deis/issues/new
.. _`troubleshooting`: http://docs.ceph.com/docs/giant/rados/troubleshooting/
