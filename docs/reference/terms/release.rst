:title: Release
:description: A Deis application release combines a build with a config.

.. _release:

Release
=======
A Deis release is a combination of a :ref:`Build` with a :ref:`Config`.
Each :ref:`Application` is associated with one release at a time.
Deis releases are numbered and new releases always increment by
one (e.g. v1, v2, v3).

:ref:`Containers <container>` that host an application use these
release versions to pull the correct code and configuration.
