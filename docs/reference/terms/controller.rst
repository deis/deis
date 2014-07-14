:title: Controller
:description: The controller is the brain of the Deis platform.

.. _controller:

Controller
==========

The controller is the "brain" of the Deis platform. A controller manages
:ref:`Clusters <cluster>`, comprised of groups of nodes providing proxy and runtime
services for the application platform. A single controller manages multiple clusters and
applications.

The controller is in charge of:

* Authenticating and authorizing clients
* Processing client API calls
* Managing a cluster with nodes to host services
* Managing containers that perform work for applications
* Managing proxies that route traffic to containers
* Managing users, keys and other base configuration

The controller stack includes:

* Django API Server for handling API calls
* Celery, backed by Redis, for dispatching tasks

.. * PostgreSQL database as a backing store for Django
.. * A lightweight *gitreceive* hook for ``git push`` access control
.. * Docker and Slugbuilder to process Heroku Buildpacks and Dockerfiles

Follow the :ref:`installation_guide` to create your own private Deis controller.
