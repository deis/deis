:title: Choosing a Scheduler
:description: How to choose a scheduler backend for Deis.


.. _choosing_a_scheduler:

Choosing a Scheduler
====================

The scheduler creates, starts, stops, and destroys your app containers. For example,
a command such as ``deis scale web=3`` tells the scheduler to run three containers
from the Docker image for your app.

The scheduler must decide which machines are eligible to run these container jobs.
Scheduler backends vary in the details of their job allocation policies and whether
or not they are resource-aware, among other features.

Deis defaults to using the `Fleet Scheduler`_. A tech preview of the `Swarm Scheduler`_
is available for testing. And work is ongoing on `Kubernetes`_ and `Mesos`_-based
schedulers with the intent to test those alternatives in future releases of Deis.


Settings set by scheduler
-------------------------

The following etcd keys are set by the scheduler module of the controller component.

Some keys will exist only if a particular ``schedulerModule`` backend is enabled.

=============================            ================================================
setting                                  description
=============================            ================================================
/deis/scheduler/swarm/host               the swarm manager's host IP address
/deis/scheduler/swarm/node               used to identify other nodes in the cluster
=============================            ================================================


Settings used by scheduler
--------------------------

The following etcd keys are used by the scheduler module of the controller component.

====================================      ===============================================
setting                                   description
====================================      ===============================================
/deis/controller/schedulerModule          scheduler backend, either "fleet" or "swarm"
                                          (default: "fleet")
====================================      ===============================================


Fleet Scheduler
---------------

`fleet`_ is a scheduling backend included with CoreOS:

    fleet ties together systemd and etcd into a distributed init system. Think of
    it as an extension of systemd that operates at the cluster level instead of the
    machine level. This project is very low level and is designed as a foundation
    for higher order orchestration.

``fleetd`` is already running on the machines provisioned for Deis: no additional
configuration is needed. Commands such as ``deis ps:restart web.1`` or
``deis scale cmd=10`` will use `fleet`_ by default to manage app containers.

To use the Fleet Scheduler backend explicitly, set ``schedulerModule`` to "fleet":

.. code-block:: console

    $ deisctl config scheduler set schedulerModule=fleet


Swarm Scheduler
---------------

.. important::

    The Swarm Scheduler is a tech preview and is not recommended for production use.

`swarm`_ is a scheduling backend for Docker:

    Docker Swarm is native clustering for Docker. It turns a pool of Docker hosts
    into a single, virtual host.

..

    Swarm serves the standard Docker API, so any tool which already communicates
    with a Docker daemon can use Swarm to transparently scale to multiple hosts...

Deis includes an enhanced version of swarm v0.2.0 with node failover and optimized
locking for container creation. The Swarm Scheduler uses a `soft affinity`_ filter
to spread app containers out among available machines. It is not yet possible to
change this default affinity filter.

Swarm requires the Docker Remote API to be available at TCP port 2375. If you are
upgrading an earlier installation of Deis, please refer to the CoreOS documentation
to `enable the remote API`_.

To test the Swarm Scheduler backend, it is necessary to install swarm components:

.. code-block:: console

    $ deisctl install swarm
    $ deisctl start swarm

Then set ``schedulerModule`` to "swarm":

.. code-block:: console

    $ deisctl config controller set schedulerModule=swarm

The Swarm Scheduler is now active. Commands such as ``deis destroy`` or
``deis scale web=9`` will use `swarm`_ to manage app containers.

To monitor Swarm Scheduler operations, watch the logs of the swarm-manager
component, or spy on Docker events directly on the swarm-manager machine:

.. code-block:: console

    $ deisctl journal swarm-manager
    $ docker -H 172.17.8.102:2395 events
    2015-04-30T17:31 172.17.8.100:5000/hungry-variable:v5: (from  node:deis-01) pull
    2015-04-30T17:31 172.17.8.100:5000/hungry-variable:v5: (from  node:deis-02) pull
    2015-04-30T17:31 02a570: (from 172.17.8.100:5000/hungry-variable:v5 node:deis-01) create
    2015-04-30T17:31 02a570: (from 172.17.8.100:5000/hungry-variable:v5 node:deis-01) start
    2015-04-30T17:31 61e59c: (from 172.17.8.100:5000/hungry-variable:v5 node:deis-02) create
    2015-04-30T17:31 61e59c: (from 172.17.8.100:5000/hungry-variable:v5 node:deis-02) start


.. _Kubernetes: http://kubernetes.io/
.. _Mesos: http://mesos.apache.org/
.. _fleet: https://github.com/coreos/fleet#fleet---a-distributed-init-system
.. _swarm: https://github.com/docker/swarm#swarm-a-docker-native-clustering-system
.. _`soft affinity`: https://docs.docker.com/swarm/scheduler/filter/#soft-affinitiesconstraints
.. _`enable the remote API`: https://coreos.com/docs/launching-containers/building/customizing-docker/
