:title: Customizing store-metadata
:description: Learn how to tune custom Deis settings.

.. _store_metadata_settings:

Customizing store-metadata
==========================
The following settings are tunable for the :ref:`store` component's metadata service.

Dependencies
------------
Requires: :ref:`store-daemon <store_daemon_settings>`, :ref:`store-monitor <store_monitor_settings>`

Required by: store-volume service (runs on all hosts - not a Deis component)

Considerations: none

Settings set by store-metadata
------------------------------
The following etcd keys are set by the store-metadata component, typically in its /bin/boot script.

===================================       ==============================================
setting                                   description
===================================       ==============================================
/deis/store/filesystemSetupComplete       Set when the Ceph filesystem setup is complete
===================================       ==============================================

Settings used by store-metadata
-------------------------------
The following etcd keys are used by the store-metadata component.

====================================      =================================================================================================
setting                                   description
====================================      =================================================================================================
/deis/store/adminKeyring                  keyring for an admin user to access the Ceph cluster (set by store-monitor)
/deis/store/fsid                          Ceph filesystem ID (set by store-monitor)
/deis/store/hosts/*                       deis-monitor hosts (set by store-monitor)
/deis/store/maxPGsPerOSDWarning           threshold for warning on number of placement groups per OSD (set by store-monitor)
/deis/store/monKeyring                    keyring for the monitor to access the Ceph cluster (set by store-monitor)
/deis/store/monSetupComplete              set when the Ceph cluster setup is complete (set by store-monitor)
/deis/store/monSetupLock                  host of store-monitor that completed setup (set by store-monitor)
/deis/store/minSize                       minimum number of store-daemons necessary for the cluster to accept writes (set by store-monitor)
/deis/store/pgNum                         number of Ceph placement groups for the storage pools (set by store-monitor)
/deis/store/size                          number of replicas for data stored in Ceph (set by store-monitor)
====================================      =================================================================================================

Using a custom store-metadata image
-----------------------------------
You can use a custom Docker image for the store-metadata component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config store-metadata set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config store-metadata set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock store-metadata image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock store-metadata image`: https://github.com/deis/deis/tree/master/store/metadata
