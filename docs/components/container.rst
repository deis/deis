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

.. _`Docker`: http://docker.io/
