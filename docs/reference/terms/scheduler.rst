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

The Deis scheduler is implemented in the :ref:`controller` component. The
scheduler implementation (or "backend") can be changed dynamically to support
different strategies or cluster types.

Follow the :ref:`choosing_a_scheduler` guide to learn about available
options for the Deis scheduler.
