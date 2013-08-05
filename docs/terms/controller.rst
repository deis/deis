:title: Controller
:description: What is a Deis Controller?
:keywords: deis, controller

.. _controller:

Controller
==========
The controller is the "brains" of the Deis platform.
Each controller is tied to a single Chef organization.
The controller is in charge of:

* Processing Client API calls
* Managing Chef Nodes
* Managing Docker Containers
* Configuring Nginx proxies

The controller stack includes:

* Django API Server for handling API calls
* PostgreSQL database as a backing store for Django
* Celery / RabbitMQ for dispatching tasks
* Gitosis to handle access control for Git Push over SSH
* Docker and Buildstep to process Heroku Buildpacks

Follow the :ref:`Installation` process to setup your own private
Deis controller.