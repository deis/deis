:title: Scheduler
:description: The scheduler starts and manages containers.

.. _scheduler:

Scheduler
=========

The scheduler is responsible for creating, starting, stopping, and destroying
app :ref:`Containers <container>`. For example, a command such as
``deis scale cmd=10`` tells the scheduler to run ten containers from the
Docker image for your app.

The scheduler must decide which machines are eligible to run these container
jobs. Scheduler backends vary in the details of their job allocation policies
and whether or not they are resource-aware, among other features.

The Deis scheduler client is implemented in the :ref:`controller` component.
Deis uses `Fleet`_ to schedule the containers across the cluster.


.. _`Fleet`: https://github.com/coreos/fleet
