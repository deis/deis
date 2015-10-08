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

Deis v2
-------

Much of the Deis team's efforts are focused on Deis v2, which will see the Deis
platform running on Kubernetes.

For details on Deis v2, see the `Deis v2 Design Document`_
and issues with the `v2`_ tag.

Etcd Hardening
--------------

We have recently seen various issues with etcd performance, which can cause
issues with Deis. We will focus on re-evaluating our implementation and use
of etcd in Deis, with a focus on performance and reliability.

This feature is tracked as GitHub issue `#4404`_.

Rigger
------

We want the provisioning process of CoreOS hosts to be as smooth as the
``deisctl`` experience when users install and start Deis. We've started working
on a tool called `Rigger`_ which provisions hosts and deploys Deis on
any supported provider.

We're still defining the interface between Rigger and the provider-specific
provision logic. We'd love input from the community on this discussion in
`#4345`_.

New Default Scheduler
---------------------
Deis now has support for Docker Swarm, Apache Mesos, and Google Kubernetes as
application schedulers. With the known limitations of fleet (primarily, not being
a resource-aware scheduler), we should investigate using a different scheduler
as our default.

This feature is tracked as GitHub issue `#4222`_.

Permissions and Teams
---------------------
Deis deployments in larger organizations require more fine-grained control
over users and permissions. Implementation of teams and improved user permissions
are tracked in separate issues:

 - [ ] `Permissions`_
 - [ ] `Teams`_

Monitoring & Telemetry
----------------------
Deis installations today use custom solutions for monitoring, alerting and operational visibility.
Deis will standardize the monitoring interfaces and provide open source agent(s) that can be used to ship telemetry to arbitrary endpoints.

 - [ ] Host Telemetry (cpu, memory, network, disk)
 - [ ] Container Telemetry (cpu, memory, network, disk)
 - [ ] Platform Telemetry (control plane, data plane)
 - [ ] Controller Telemetry (app created, build created, containers scaled)

This feature is tracked as GitHub issue `#3699`_.

Internal Service Discovery
--------------------------
To provide a better container networking experience, Deis must provide
internal service discovery for components to coordinate.

This feature is tracked as GitHub issue `#3072`_.

Update Service
--------------
Deis must support 100% automated, zero-downtime updates of the control plane.
Like CoreOS, Deis clusters should be attached to an alpha, beta or stable channel and rely on an automatic update mechanism.
To accomplish this, Deis plans to use the `Google Omaha Protocol`_ as implemented by `CoreUpdate`_.

 - [X] `Update client/agent`_
 - [ ] Update server
 - [ ] `Automatic CoreOS upgrades`_
 - [ ] CI Integration

This feature is tracked as GitHub issue `#2106`_.

Deis Push
---------
End-users should be able to push Docker-based applications into Deis from their local machine using ``deis push user/app``.
This works around a number of authentication issues with private registries and ``deis pull``.

 - [ ] `Docker Registry v2`_
 - [ ] `Deis Push`_

TTY Broker
----------
Today Deis cannot provide bi-directional streams needed for log tailing and interactive batch processes.
By having the :ref:`Controller` drive a TTY Broker component, Deis can securely open WebSockets
through the routing mesh.

 - [ ] `TTY Broker component`_
 - [ ] `Interactive Deis Run`_ (``deis run bash``)
 - [ ] `Log Tailing`_ (``deis logs -f``)

Service Broker
--------------
In Deis, connections to :ref:`concepts_backing_services` are meant to be explicit and modeled as a series of environment variables.
Deis believes the Cloud Foundry `Service Broker API`_ is the best embodiment of this today.

 - [ ] Deis Addons CLI (deis addons)
 - [ ] PostgreSQL Service Broker
 - [ ] Redis Service Broker

This feature is tracked as GitHub issue `#231`_.

.. _`#231`: https://github.com/deis/deis/issues/231
.. _`#2106`: https://github.com/deis/deis/issues/2106
.. _`#3072`: https://github.com/deis/deis/issues/3072
.. _`#3699`: https://github.com/deis/deis/issues/3699
.. _`#4222`: https://github.com/deis/deis/issues/4222
.. _`#4345`: https://github.com/deis/deis/issues/4345
.. _`#4404`: https://github.com/deis/deis/issues/4404
.. _`Automatic CoreOS upgrades`: https://github.com/deis/deis/issues/1043
.. _`CoreUpdate`: https://coreos.com/docs/coreupdate/custom-apps/coreupdate-protocol/
.. _`Deis Push`: https://github.com/deis/deis/issues/2680
.. _`Deis v2 Design Document`: https://github.com/deis/deis/issues/4582
.. _`Docker Registry v2`: https://github.com/deis/deis/issues/3814
.. _`Google Omaha Protocol`: https://code.google.com/p/omaha/wiki/ServerProtocol
.. _`Interactive Deis Run`: https://github.com/deis/deis/issues/117
.. _`like CoreOS`: https://coreos.com/releases/
.. _`Log Tailing`: https://github.com/deis/deis/issues/465
.. _`Permissions`: https://github.com/deis/deis/issues/4150
.. _`Rigger`: https://github.com/deis/rigger
.. _`Service Broker API`: http://docs.cloudfoundry.org/services/api.html
.. _`Teams`: https://github.com/deis/deis/issues/4173
.. _`TTY Broker component`: https://github.com/deis/deis/issues/3808
.. _`Update client/agent`: https://github.com/deis/deis/issues/3811
.. _`v2`: https://github.com/deis/deis/labels/v2
