:title: Release
:description: A Deis release is a combination of a Build with a Config. Each Deis application is associated with one release at a time.
:keywords: release, deis

.. _release:

Release
=======
A Deis release is a combination of a :ref:`Build` with a :ref:`Config`.
Each :ref:`Application` is associated with one release at a time.
Deis releases are numbered and increment by one (e.g. v1, v2, v3).

:ref:`Containers <container>` in the runtime :ref:`Layer` of a :ref:`formation`
use the release version to pull the correct code and configuration as
part of their `Chef run`_.

.. _`Chef run`: http://docs.opscode.com/essentials_nodes_chef_run.html
