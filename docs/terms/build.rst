:title: Build
:description: A Deis build refers to the output of a specific application build. Deis builds are created automatically using git push.
:keywords: build, release, git, git push, deis

.. _build:

Build
=====
A Deis build refers to the output of a specific application build.
Deis builds are created automatically on the controller when a 
developer uses ``git push deis master``.

When a new build is created, a new :ref:`release` is created automatically.
