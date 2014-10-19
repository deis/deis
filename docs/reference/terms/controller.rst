:title: Controller
:description: The controller is the brain of the Deis platform.

.. _controller:

Controller
==========
The controller is the "brain" of the Deis platform. A controller
manages :ref:`Applications <application>` and their :ref:`Containers <container>`.

The controller is in charge of:

* Authenticating and authorizing clients
* Processing client API calls
* Managing containers that perform work for applications
* Managing proxies that route traffic to containers
* Managing users, keys and other base configuration

The controller stack includes:

* Django API Server for handling API calls

.. * PostgreSQL database as a backing store for Django
.. * A lightweight *gitreceive* hook for ``git push`` access control
.. * Docker and Slugbuilder to process Heroku Buildpacks and Dockerfiles

Follow the :ref:`Installing Deis <installing_deis>` guide to create your own
private Deis controller.
