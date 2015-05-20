:title: Roadmap
:description: The Deis project roadmap

.. _roadmap:

Deis Roadmap
============
The Deis Roadmap is a community document created as part of our open :ref:`planning`.
Each roadmap item describes a high-level capability or grouping of features we deem
important to the future of Deis.

Given our rapid :ref:`Release Schedule, <release_schedule>` roadmap items are designed to provide a sense of
direction over many releases.

Update Service
--------------
Deis must support 100% automated, zero-downtime updates of the control plane.
Like CoreOS, Deis clusters should be attached to an alpha, beta or stable channel and rely on an automatic update mechanism.
To accomplish this, we plan to use the `Google Omaha Protocol`_ as implemented by `CoreUpdate`_.

 - [ ] Update client/agent
 - [ ] Update server
 - [ ] CI Integration

Scheduling and Orchestration
----------------------------
Today Deis uses `Fleet`_ for scheduling.  Unfortunately, Fleet does not support
resource-based scheduling which results in poor cluster utilization at scale.

Fortunately, Deis is composable and can easily hot-swap orchestration APIs.
Because the most promising container orchestration solutions are under heavy development,
we are currently focused on releasing "technology previews".

These technology previews will help the community try different orchestration solutions easily,
report their findings and help guide the future direction of Deis.

 - [X] Swarm preview
 - [ ] Mesos preview
 - [ ] Kubernetes preview

TTY Broker
----------
Today Deis cannot provide bi-directional streams needed for log tailing and interactive batch processes.
By having the :ref:`Controller` drive a TTY Broker component, we can securely open WebSockets
through the routing mesh.

 - [ ] TTY Broker component
 - [ ] Interactive Deis Run (deis run bash)
 - [ ] Log Tailing (deis logs -f)

Deis Push
---------
End-users should be able to push Docker-based applications into Deis from their local machine using ``deis push user/app``.
This works around a number of authentication issues with private registries and ``deis pull``.

 - [ ] Docker Registry v2
 - [ ] Deis Push

Service Broker
--------------
In Deis, connections to :ref:`concepts_backing_services` are meant to be explicit and modeled as a series of environment variables.
We believe the Cloud Foundry `Service Broker API`_ is the best embodiment of this today.

 - [ ] Deis Addons CLI (deis addons)
 - [ ] PostgreSQL Service Broker
 - [ ] Redis Service Broker

Monitoring & Telemetry
----------------------
Deis installations today use custom solutions for monitoring, alerting and operational visibility.
We will standardize the monitoring interfaces and provide open source agent(s) that can be used to ship telemetry to arbitrary endpoints.

 - [ ] Host Telemetry (cpu, memory, network, disk)
 - [ ] Container Telemetry (cpu, memory, network, disk)
 - [ ] Platform Telemetry (control plane, data plane)
 - [ ] Controller Telemetry (app created, build created, containers scaled)

.. _`like CoreOS`: https://coreos.com/releases/
.. _`Google Omaha Protocol`: https://code.google.com/p/omaha/wiki/ServerProtocol
.. _`CoreUpdate`: https://coreos.com/docs/coreupdate/custom-apps/coreupdate-protocol/
.. _`Fleet`: https://github.com/coreos/fleet#readme
.. _`Service Broker API`: http://docs.cloudfoundry.org/services/api.html
