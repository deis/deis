:title: Container
:description: Deis containers are instances of Docker containers used to run applications.
:keywords: lxc container, lxc, container, docker, deis

.. _container:

Container
=========
Deis containers are instances of `Docker`_ containers used to run :ref:`Applications <application>`.
Containers perform the actual work of an :ref:`application` by servicing requests or running
background tasks as part of the :ref:`Formation's <formation>` runtime :ref:`Layer`.

Ephemeral Filesystem
--------------------

Each container gets its own ephemeral filesystem, with a fresh copy of the most recently
deployed code. During the container’s lifetime, its running processes can use the
filesystem as a temporary scratchpad, but no files that are written are visible to
processes in any other container. Any files written to the ephemeral filesystem will be
discarded the moment the container is either stopped or restarted.

Container States
----------------

There are several states that a container can be in at any time. The states are:

1) initialized - the initial state of the container before it is created
2) created - the container is built and ready for operation
3) up - the container is running
4) down - the container crashed or is stopped
5) destroyed - the container has been destroyed


.. _`Docker`: http://docker.io/
