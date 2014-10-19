:title: Container
:description: Deis uses Docker containers to run applications.

.. _container:

Container
=========
Deis containers are instances of `Docker`_ containers used to run
:ref:`Applications <application>`. Containers perform the actual work
of an :ref:`application` by servicing requests or by running background
tasks as part of the cluster.

Ephemeral Filesystem
--------------------

Each container gets its own ephemeral filesystem, with a fresh copy of the most recently
deployed code. During the container’s lifetime, its running processes can use the
filesystem as a temporary scratchpad, but no files that are written are visible to
processes in any other container. Any files written to the ephemeral filesystem will be
discarded the moment the container is either stopped or restarted.

Container States
----------------
There are several states that a container can be in at any time. The
states are:

#. initialized - the state of the container before it is created
#. created - the container is built and ready for operation
#. up - the container is running
#. down - the container crashed or is stopped
#. destroyed - the container has been destroyed


.. _`Docker`: http://docker.io/
