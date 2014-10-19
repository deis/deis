:title: Concepts
:description: Deis scales Twelve-Factor apps as containers over a cluster of machines.

.. _concepts:

Concepts
========
Deis is a lightweight, flexible and powerful application platform that
deploys and scales :ref:`concepts_twelve_factor` apps as
:ref:`concepts_docker` containers across a cluster of
:ref:`concepts_coreos` machines.

.. _concepts_twelve_factor:

Twelve-Factor
-------------
The `Twelve-Factor App`_ is a DevOps manifesto for building and
deploying scalable, modern applications and services.

We consider it an invaluable synthesis of much experience with
software-as-a-service apps in the wild, especially on the
Heroku platform.

Deis works best with applications in a `Twelve-Factor App`_ style.
Following the twelve-factor model, Deis enforces a strict separation of
the :ref:`Build and Run <concepts_build_release_run>` stages.

.. _concepts_docker:

Docker
------
`Docker`_ is an open source project to pack, ship and run any
application as a lightweight, portable, self-sufficient container.

When you deploy an app with ``git push deis master``, Deis builds and
packages it as a Docker image, then distributes it as Docker containers
across your cluster.

(Deis itself is also a set of coordinated Docker containers.)

.. _concepts_coreos:

CoreOS
------
`CoreOS`_ is a lean new Linux distribution, rearchitected for features
needed by modern infrastructure stacks and targeted at massive
server deployments.

Deis applications are processes running on CoreOS machines, which can
be private or public cloud instances, or bare metal. CoreOS clusters
allow Deis to host applications and services at scale with
high resilience.

Yet Deis and CoreOS run identically in a Vagrant virtual machine on
your laptop, for convenient testing and rapid development.

.. _concepts_applications:

Applications
------------
An :ref:`application`, or app, lives on a cluster where it uses
:ref:`Containers <container>` to process requests and run tasks for a
deployed git repository.

Developers use :ref:`Applications <application>` to push code, change
configuration, scale processes, view logs, or run admin commands --
regardless of the cluster's underlying infrastructure.

.. _concepts_build_release_run:

Build, Release, Run
-------------------

Build Stage
^^^^^^^^^^^
The :ref:`Controller` includes a *gitreceive* hook that receives incoming git push requests over
SSH and builds applications inside ephemeral Docker containers. Tarballs of the /app directory are
extracted into a slug and is injected into another container, which will create the app image. The
image is then pushed to a private registry for later execution.

Release Stage
^^^^^^^^^^^^^
During the release stage, a :ref:`build` is combined with :ref:`config` to create a new numbered
:ref:`release`. The release stage is triggered any time a new build is created or config is
changed, making it easy to rollback code and configuration.

Run Stage
^^^^^^^^^
The run stage dispatches containers to a scheduler and updates the router accordingly.
The scheduler is in control of placing containers on hosts and balancing them evenly across the cluster.
Containers are published to the router once they are healthy.  Old containers are only collected
after the new containers are live and serving traffic -- providing zero-downtime deploys.

.. _concepts_backing_services:

Backing Services
----------------
Deis treats databases, caches, storage, messaging systems, and other
`backing services`_ as attached resources, in keeping with Twelve-Factor
best practices.

Applications can be decoupled this way, using simple
`environment variables`_ to configure and attach to any services needed.
Apps are then free to scale up independently, to use services provided
by other apps, or to switch easily to external or third-party vendor
services.

See Also
--------
* :ref:`Using Deis <using_deis>`
* :ref:`Managing Deis <managing_deis>`
* The `Twelve-Factor App`_


.. _`Twelve-Factor App`: http://12factor.net/
.. _`Docker`: http://docker.io/
.. _`CoreOS`: https://coreos.com/
.. _`Build and Run`: http://12factor.net/build-release-run
.. _`backing services`: http://12factor.net/backing-services
.. _`environment variables`: http://12factor.net/config
