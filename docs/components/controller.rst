:title: Controller
:description: The Deis controller is the brains of the Deis platform. Details on what the Deis controller is in charge of and what the Deis controller stack includes.
:keywords: controller, deis

.. _controller:

Controller
==========
The controller is the "brains" of the Deis platform.
The controller manages container :ref:`Formations <formation>`,
comprised of clusters of nodes providing proxy and runtime services for
the application platform.  A single controller manages multiple 
container formations.

Controllers are tied to a configuration management backend (typically a 
Chef Server) where data about users, applications and formations is stored.

The controller is in charge of:

* Processing client API calls
* Managing nodes that provide services to a formation
* Managing containers that perform work for applications
* Managing proxies that route traffic to containers
* Managing users, providers, flavors, keys and other base configuration

The controller stack includes:

* Django API Server for handling API calls
* PostgreSQL database as a backing store for Django
* Celery / RabbitMQ for dispatching tasks
* Gitosis to handle access control for Git Push over SSH
* Docker and Buildstep to process Heroku Buildpacks

Follow the :ref:`Operations Guide <operations>` to setup your own private
Deis controller.