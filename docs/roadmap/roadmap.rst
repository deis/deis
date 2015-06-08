:title: Roadmap
:description: The Deis project roadmap

.. _roadmap:

Deis Roadmap
============
The Deis Roadmap is a community document created as part of the open :ref:`planning`.
Each roadmap item describes a high-level capability or grouping of features that are deemed
important to the future of Deis.

Given the project's rapid :ref:`Release Schedule, <release_schedule>` roadmap items are designed to provide a sense of
direction over many releases.

Pluggable Storage Subsystem
---------------------------
Deis uses Ceph to provide a highly-available storage subsystem for stateful control plane components.
While Ceph is the right default storage subsystem, it is tricky to operate and can result in control plane instability.
Ceph should be optional, especially for users on AWS with direct access to services like S3.

 - [ ] Registry S3 configuration
 - [ ] Database S3 configuration
 - [ ] Stateless logger (in-memory ring buffer)

This feature is tracked as GitHub issue `#2812`_.

TTY Broker
----------
Today Deis cannot provide bi-directional streams needed for log tailing and interactive batch processes.
By having the :ref:`Controller` drive a TTY Broker component, Deis can securely open WebSockets
through the routing mesh.

 - [ ] `TTY Broker component`_
 - [ ] `Interactive Deis Run`_ (``deis run bash``)
 - [ ] `Log Tailing`_ (``deis logs -f``)

Scheduling and Orchestration
----------------------------
Today Deis uses `Fleet`_ for scheduling.  Unfortunately, Fleet does not support
resource-based scheduling, which results in poor cluster utilization at scale.

Fortunately, Deis is composable and can easily hot-swap orchestration APIs.
Because the most promising container orchestration solutions are under heavy development,
the Deis project is focused on releasing "technology previews".

These technology previews will help the community try different orchestration solutions easily,
report their findings and help guide the future direction of Deis.

 - [X] Swarm preview
 - [ ] `Mesos preview`_
 - [ ] `Kubernetes preview`_

Networking v2
-------------
To provide a better container networking experience, Deis must provide an overlay network
that can facilitate SDN and improved service discovery.

 - [ ] Overlay Network
 - [ ] `Internal Service Discovery`_
 - [ ] Migration Strategy

This feature is tracked as GitHub issue `#3812`_.

Update Service
--------------
Deis must support 100% automated, zero-downtime updates of the control plane.
Like CoreOS, Deis clusters should be attached to an alpha, beta or stable channel and rely on an automatic update mechanism.
To accomplish this, Deis plans to use the `Google Omaha Protocol`_ as implemented by `CoreUpdate`_.

 - [ ] `Update client/agent`_
 - [ ] Update server
 - [ ] `Automatic CoreOS upgrades`_
 - [ ] CI Integration

This feature is tracked as GitHub issue `#2106`_.

Etcd 2
------
A CP database like etcd is central to Deis, which requires a distributed lock service and key/value store.
As problems with etcd directly impact platform stability, Deis must move to the more stable etcd2.

 - [ ] Switch to etcd2
 - [ ] Migration strategy for etcd 0.4.x -> etcd2

This feature is tracked as GitHub issue `#3564`_.

User-defined Health Checks
--------------------------
Today Deis relies on TCP port checks as evidence of container readiness during deploys.
To facilitate a better zero-downtime deploy experience, Deis should allow user-defined
health checks that are respected during the rolling deploy process and to verify that containers
are healthy before publishing them to the routing mesh.

 - [ ] HealthCheck API and CLI
 - [ ] Publisher / HealthCheck integration

This feature is tracked as GitHub issue `#3813`_.

Deis Push
---------
End-users should be able to push Docker-based applications into Deis from their local machine using ``deis push user/app``.
This works around a number of authentication issues with private registries and ``deis pull``.

 - [ ] `Docker Registry v2`_
 - [ ] `Deis Push`_

Service Broker
--------------
In Deis, connections to :ref:`concepts_backing_services` are meant to be explicit and modeled as a series of environment variables.
Deis believes the Cloud Foundry `Service Broker API`_ is the best embodiment of this today.

 - [ ] Deis Addons CLI (deis addons)
 - [ ] PostgreSQL Service Broker
 - [ ] Redis Service Broker

This feature is tracked as GitHub issue `#231`_.

Monitoring & Telemetry
----------------------
Deis installations today use custom solutions for monitoring, alerting and operational visibility.
Deis will standardize the monitoring interfaces and provide open source agent(s) that can be used to ship telemetry to arbitrary endpoints.

 - [ ] Host Telemetry (cpu, memory, network, disk)
 - [ ] Container Telemetry (cpu, memory, network, disk)
 - [ ] Platform Telemetry (control plane, data plane)
 - [ ] Controller Telemetry (app created, build created, containers scaled)

This feature is tracked as GitHub issue `#3699`_.

.. _`#231`: https://github.com/deis/deis/issues/231
.. _`#2106`: https://github.com/deis/deis/issues/2106
.. _`#2812`: https://github.com/deis/deis/issues/2812
.. _`#3564`: https://github.com/deis/deis/issues/3564
.. _`#3699`: https://github.com/deis/deis/issues/3699
.. _`#3812`: https://github.com/deis/deis/issues/3812
.. _`#3813`: https://github.com/deis/deis/issues/3813
.. _`Automatic CoreOS upgrades`: https://github.com/deis/deis/issues/1043
.. _`CoreUpdate`: https://coreos.com/docs/coreupdate/custom-apps/coreupdate-protocol/
.. _`Deis Push`: https://github.com/deis/deis/issues/2680
.. _`Docker Registry v2`: https://github.com/deis/deis/issues/3814
.. _`Fleet`: https://github.com/coreos/fleet#readme
.. _`Google Omaha Protocol`: https://code.google.com/p/omaha/wiki/ServerProtocol
.. _`Interactive Deis Run`: https://github.com/deis/deis/issues/117
.. _`Internal Service Discovery`: https://github.com/deis/deis/issues/3072
.. _`Kubernetes preview`: https://github.com/deis/deis/issues/2744
.. _`like CoreOS`: https://coreos.com/releases/
.. _`Log Tailing`: https://github.com/deis/deis/issues/465
.. _`Mesos preview`: https://github.com/deis/deis/issues/3809
.. _`Service Broker API`: http://docs.cloudfoundry.org/services/api.html
.. _`TTY Broker component`: https://github.com/deis/deis/issues/3808
.. _`Update client/agent`: https://github.com/deis/deis/issues/3811
