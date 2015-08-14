:title: Production deployments
:description: Considerations for deploying Deis in production.

.. _production_deployments:

Production deployments
======================

Many Deis users are running Deis quite successfully in production. When readying a Deis deployment
for production workloads, there are some additional (but optional) recommendations.

Isolating the Planes
--------------------

.. include:: ../_includes/_isolating-planes-description.rst

See :ref:`isolating-planes` for further details.

Isolating etcd
--------------

.. include:: ../_includes/_isolating-etcd-description.rst

See :ref:`isolating-etcd` for further details.

Running Deis without Ceph
-------------------------

.. include:: ../_includes/_ceph-dependency-description.rst

See :ref:`running-deis-without-ceph` for details on removing this operational
complexity.

.. _preseeding_continers:

Preseeding containers
---------------------

When a host in your CoreOS cluster fails or becomes unresponsive, the CoreOS scheduler will relocate
any cluster services on that machine to another host. These services come up on the new host just fine,
but a component's first task is to pull the corresponding Docker image from Docker Hub. Depending
on factors such as available bandwidth, network latency, and performance of the Docker Hub platform,
this can take some time. Failover is not finished until the pull completes and the component starts.

To minimize component downtime should failover occur, it is recommended to preseed the Docker images
for Deis on all hosts in a cluster. This will pull all the images to the host's local Docker graph,
so if failover should occur, a component can start quickly.

A preseed script is provided as a script already loaded on CoreOS hosts.

On all hosts in the cluster, run:

.. code-block:: console

    $ /run/deis/bin/preseed

This will pull all component images for the installed version of Deis.

Review security considerations
------------------------------

There are some additional security-related considerations when running Deis in production, and users
can consider enabling a firewall on the CoreOS hosts as well as the router component.

See :ref:`security_considerations` for details.

Back up data
------------

Backing up data regularly is recommended. See :ref:`backing_up_data` for steps.

Change Registration Mode
------------------------

Changing the registration process is highly recommended in production.  By default, registrations for a
new cluster are open to anyone with the proper URL.  Once the admin user has registered with a new cluster,
it is recommended that you either turn off registrations entirely or enable the admin only registration feature.

Please see the following documentation:  :ref:`controller_settings`

Configure logging and monitoring
--------------------------------

Many users already have external monitoring or logging systems, and connecting Deis to these
platforms is quite simple. Review :ref:`platform_logging` and :ref:`platform_monitoring`.

Enable TLS
----------

Using TLS to encrypt traffic (including Deis client traffic, such as login credentials) is crucial.
See :ref:`platform_ssl`.
