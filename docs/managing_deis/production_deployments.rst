:title: Production deployments
:description: Considerations for deploying Deis in production.

.. _production_deployments:

Production deployments
======================

Many Deis users are running Deis quite successfully in production. When readying a Deis deployment
for production workloads, there are some additional (but optional) recommendations.

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

Configure logging and monitoring
--------------------------------

Many users already have external monitoring or logging systems, and connecting Deis to these
platforms is quite simple. Review :ref:`platform_logging` and :ref:`platform_monitoring`.

Enable TLS
----------

Using TLS to encrypt traffic (including Deis client traffic, such as login credentials) is crucial.
See :ref:`platform_ssl`.

Large-scale cluster architecture
--------------------------------

Deis Control Plane, Router Mesh, and Data Plane components all depend on etcd
for service discovery.  For larger deployments, it is recommended that etcd
not run on greater than nine servers as doing so can greatly increase the time
required for etcd nodes to build consensus, which can have a deleterious effect
on a cluster's overall performance.

Those building larger clusters should reference `CoreOS cluster architecture
documentation`_.  This is considered an advanced use case, and as such Deis'
provisioning scripts do not address this concern directly.  The remainder
of this section addresses the means by which an operator might construct such
a large-scale cluster.

To realize such a topology, it is necessary to provide different cloud-config
for those servers that will run etcd and for those that will not.  By extension,
that process will involve a customized provisioning process.

For a small, fixed number of servers running etcd and satisfying the "central
services" role (as described in the CoreOS documentation), the cloud-config
provided with Deis is sufficient.

For servers not running etcd, satisfying the "worker" role (as described in the
CoreOS documentation), cloud-config must be tweaked to account for the following
considerations:

1. Etcd must be disabled.
2. Fleet must be configured to use remote etcd servers.
3. The ``ETCDCTL_PEERS`` environment variable must be set with the addresses of the remote etcd servers.
4. The ``ETCDCTL_PEERS`` environment variable must be added to ``/etc/environment``.  This will enable the use of etcdctl and derivative utilities such as Deis' ``get_image`` script from within the .service files that are used to launch each Deis component.

For example:

.. code-block:: yaml

  coreos:
    fleet:
      # allow etcd to slow down at times
      etcd_request_timeout: 3.0
      etcd_servers: "http://10.0.0.101:4001,http://10.0.0.102:4001,http://10.0.0.103:4001,http://10.0.0.104:4001,http://10.0.0.105:4001"
    # ...
    units:
      - name: etcd.service
        mask: true
      - name: setup-etcdctl-env.service
        command: start
        content: |
          [Unit]
          Description=Setup etcdctl env vars
          Wants=coreos-setup-environment.service
          After=coreos-setup-environment.service

          [Service]
          Type=oneshot
          RemainAfterExit=true
          ExecStart=/bin/bash -c "source /etc/profile.d/etcdctl.sh; env | grep ETCDCTL_PEERS >> /etc/environment"

          [Install]
          WantedBy=multi-user.target
    # ...
    write_files:
      - path: /etc/profile.d/etcdctl.sh
        permissions: 0644
        owner: core
        content: |
          export ETCDCTL_PEERS="http://10.0.0.101:4001,http://10.0.0.102:4001,http://10.0.0.103:4001,http://10.0.0.104:4001,http://10.0.0.105:4001"
      - path: /etc/profile.d/fleetctl.sh
        permissions: 0644
        owner: core
        content: |
          export FLEETCTL_ENDPOINT=unix:///var/run/fleet.sock
          export FLEETCTL_EXPERIMENTAL_API=true

In the example above, the sample IPs obviously must be replaced with the actual
IPs of the "central services" nodes.  The means of accomplishing that
programmatically at the time of provisioning will vary from one cloud provider
to the next.

.. _`CoreOS cluster architecture documentation`: https://coreos.com/docs/cluster-management/setup/cluster-architectures/#production-cluster-with-central-services
