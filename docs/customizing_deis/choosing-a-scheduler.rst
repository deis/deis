:title: Choosing a Scheduler
:description: How to choose a scheduler backend for Deis.


.. _choosing_a_scheduler:

Choosing a Scheduler
====================

The :ref:`scheduler` creates, starts, stops, and destroys each :ref:`container`
of your app. For example, a command such as ``deis scale web=3`` tells the
scheduler to run three containers from the Docker image for your app.

Deis defaults to using the `Fleet Scheduler`_. The Tech previews of the `Swarm Scheduler`_,
`Mesos`_ and `Kubernetes`_ are available for testing. Keep watching `dynamic metadata fleet PR 1077`_
for support of runtime changes to metadata.

.. note::

    If you are using a scheduler other than fleet app containers will not be rescheduled if deis-registry is unavailable.
    For more information look at `deis-registry issue 3619`.

Settings set by scheduler
-------------------------

The following etcd keys are set by the scheduler module of the controller component.

Some keys will exist only if a particular ``schedulerModule`` backend is enabled.

===================================            ==========================================================
setting                                        description
===================================            ==========================================================
/deis/scheduler/swarm/host                     the swarm manager's host IP address
/deis/scheduler/swarm/node                     used to identify other nodes in the cluster
/deis/scheduler/mesos/marathon                 used to identify Marathon framework's host IP address
/deis/scheduler/k8s/master                     used to identify host IP address of kubernetes ApiService
===================================            ==========================================================


Settings used by scheduler
--------------------------

The following etcd keys are used by the scheduler module of the controller component.

====================================      ===============================================
setting                                   description
====================================      ===============================================
/deis/controller/schedulerModule          scheduler backend, either "fleet" or "swarm" or
                                          "mesos_marathon" or "k8s" (default: "fleet")
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

To use the Fleet Scheduler backend explicitly, set the controller's
``schedulerModule`` to "fleet":

.. code-block:: console

    $ deisctl config controller set schedulerModule=fleet


Swarm Scheduler
---------------

.. important::

    The Swarm Scheduler is a technology preview and is not recommended for
    production use.

`swarm`_ is a scheduling backend for Docker:

    Docker Swarm is native clustering for Docker. It turns a pool of Docker hosts
    into a single, virtual host.

..

    Swarm serves the standard Docker API, so any tool which already communicates
    with a Docker daemon can use Swarm to transparently scale to multiple hosts...

Deis includes an enhanced version of swarm v0.2.0 with node failover and optimized
locking on container creation. The Swarm Scheduler uses a `soft affinity`_ filter
to spread app containers out among available machines.

Swarm requires the Docker Remote API to be available at TCP port 2375. If you are
upgrading an earlier installation of Deis, please refer to the CoreOS documentation
to `enable the remote API`_.

.. note::

    **Known Issues**

    - It is not yet possible to change the default affinity filter.

To test the Swarm Scheduler backend, first install and start the swarm components:

.. code-block:: console

    $ deisctl install swarm && deisctl start swarm

Then set the controller's ``schedulerModule`` to "swarm":

.. code-block:: console

    $ deisctl config controller set schedulerModule=swarm

The Swarm Scheduler is now active. Commands such as ``deis destroy`` or
``deis scale web=9`` will use `swarm`_ to manage app containers.

To monitor Swarm Scheduler operations, watch the logs of the swarm-manager
component, or spy on Docker events directly on the swarm-manager machine:

Kubernetes Scheduler
--------------------

.. important::

    The Kubernetes Scheduler is a technology preview and is not recommended for
    production use.

`Kubernetes`_ orchestration system for docker containers:

    Kubernetes provides APIs to manage, deploy and scale docker containers. Kubernetes deploys docker containers as `pods`_,
    unique entity across a cluster. Containers inside a pod share the same namespaces. More information about `pods`_.

Kubernetes requires overlay network so that each pod can get a unique IP across the cluster which is achieved by using `flannel`_.
Deis repository provides a new user-data file to bootstrap CoreOs machines on the cluster with flannel.

To test the Kubernetes Scheduler, first install and start the kubernetes components:

.. code-block:: console

    $deisctl install k8s && deisctl start k8s

Then set the controller's ``schedulerModule`` to "k8s":

.. code-block:: console

    $ deisctl config controller set schedulerModule=k8s

The Kubernetes scheduler is now active. Commands such as ``deis destroy`` or
``deis scale web=9`` will use kubernetes ApiServer to manage app pods.

Deis creates a `replication controller`_ to manage pods and a `service`_ which acts as a proxy and routes traffic to the pods associated with the App.
A new release uses rolling deploy mechanism where in we create and replace the pod with the new release one by one until all the pods are replaced properly or rollback to the older release
if any error occurs which differs from the existing schedulers where we create all the containers at a time and replace them with the old ones.

.. note::

    **Known Issues**

    - Kubernetes ApiServer is not HA. If ApiServer is rescheduled it will reschedule all the kubernetes units.
    - Kubernetes does resource based scheduling. Specifying limits will create reservation of the specified resource on the node.
    - Installing flannel is not backwards compatible as it changes Docker bridge which requires restarting
      Docker Engine.
    - Docker Engine is started with arguments ``--iptables=false --ip-masq=false`` for the docker not to create IP masquerading rule and use the one created by flannel.


Mesos with Marathon framework
-----------------------------

.. important::

    The Mesos with Marathon framework Scheduler is a technology preview and is not recommended for
    production use.

`Mesos`_ is a distributed system kernel:

    Mesos provides APIs for resource management and scheduling. A framework interacts with Mesos master
    and schedules and task. A Zookeeper cluster elects Mesos master node. Mesos slaves are installed on
    each node and they communicate to master with available resources.

`Marathon`_ is a Mesos_ framework for long running applications:

    Marathon provides a Paas like feel for long running applications and features like high-availablilty, host constraints,
    service discovery, load balancing and REST API to control your Apps.

Deis uses the Marathon framework to schedule containers. Since Marathon is a framework for long-running
jobs, Deis uses the `Fleet Scheduler`_ to run batch processing jobs. ``deisctl`` installs a standalone Mesos
cluster. To install an HA Mesos cluster, follow the directions at `aledbf-mesos`_, and set the etcd key
``/deis/scheduler/mesos/marathon`` to any Marathon node IP address. If a request is received by a regular
Marathon node, it is proxied to the master Marathon node.

To test the Marathon Scheduler backend, first install and start the mesos components:

.. code-block:: console

    $ deisctl install mesos && deisctl start mesos

Then set the controller's ``schedulerModule`` to "mesos_marathon":

.. code-block:: console

    $ deisctl config controller set schedulerModule=mesos_marathon

The Marathon framework is now active. Commands such as ``deis destroy`` or
``deis scale web=9`` will use `Marathon`_ to manage app containers.

Deis starts Marathon on port 8180. You can manage apps through the Marathon UI, which is accessible at http://<Marathon-node-IP>:8180

.. note::

    **Known Issues**

    - deisctl installs a standalone mesos cluster as fleet doesn't support runtime change to metadata.
      You can specify this in cloud-init during the deployment of the node. keep watching `dynamic metadata fleet PR 1077`_.
    - If you want to access Marathon UI, you'll have to expose port 8180 in the security group settings.
      This is blocked off by default for security purposes.
    - Deis does not yet use Marathon's docker container API to create containers.
    - CPU shares are integers representing the number of CPUs. Memory limits should be specified in MB.


.. _Kubernetes: http://kubernetes.io/
.. _Mesos: http://mesos.apache.org
.. _Marathon: https://github.com/mesosphere/marathon
.. _pods: https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/user-guide/pods.md
.. _replication controller: https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/user-guide/replication-controller.md
.. _service: https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/user-guide/services.md
.. _flannel: https://github.com/coreos/flannel
.. _fleet: https://github.com/coreos/fleet#fleet---a-distributed-init-system
.. _swarm: https://github.com/docker/swarm#swarm-a-docker-native-clustering-system
.. _`soft affinity`: https://docs.docker.com/swarm/scheduler/filter/#soft-affinitiesconstraints
.. _`enable the remote API`: https://coreos.com/docs/launching-containers/building/customizing-docker/
.. _`deis-kubernetes issue 3850`: https://github.com/deis/deis/issues/3850
.. _`dynamic metadata fleet PR 1077`: https://github.com/coreos/fleet/pull/1077
.. _`aledbf-mesos`: https://github.com/aledbf/coreos-mesos-zookeeper
.. _`deis-registry issue 3619`: https://github.com/deis/deis/issues/3619
