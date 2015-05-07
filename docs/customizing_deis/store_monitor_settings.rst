:title: Customizing store-monitor
:description: Learn how to tune custom Deis settings.

.. _store_monitor_settings:

Customizing store-monitor
=========================
The following settings are tunable for the :ref:`store` component's monitor service.

Dependencies
------------
Requires: none

Required by: :ref:`store-daemon <store_daemon_settings>`, :ref:`store-gateway <store_gateway_settings>`

Considerations: none

Settings set by store-monitor
-----------------------------
The following etcd keys are set by the store-monitor component, typically in its /bin/boot script.

===============================          ==================================================================================
setting                                  description
===============================          ==================================================================================
/deis/store/adminKeyring                 keyring for an admin user to access the Ceph cluster
/deis/store/fsid                         Ceph filesystem ID
/deis/store/hosts/$HOST                  hostname (not IP) of the host running this store-monitor instance
/deis/store/maxPGsPerOSDWarning          threshold for warning on number of placement groups per OSD (set by store-monitor)
/deis/store/monKeyring                   keyring for the monitor to access the Ceph cluster
/deis/store/monSetupComplete             set when the Ceph cluster setup is complete
/deis/store/monSetupLock                 IP address of the monitor instance that is or has set up the Ceph cluster
/deis/store/minSize                      minimum number of store-daemons necessary for the cluster to accept writes
/deis/store/pgNum                        number of Ceph placement groups for the storage pools
/deis/store/size                         number of replicas for data stored in Ceph
===============================          ==================================================================================

Settings used by store-monitor
------------------------------
The store-monitor component uses no keys from etcd other than the ones it sets.

Using a custom store-monitor image
----------------------------------
You can use a custom Docker image for the store-monitor component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config store-monitor set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config store-monitor set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock store-monitor image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock store-monitor image`: https://github.com/deis/deis/tree/master/store/monitor
