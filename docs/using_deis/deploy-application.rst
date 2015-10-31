:title: Deploy an Application on Deis
:description: First steps for developers using Deis to deploy and manage applications

.. _deploy-application:

Deploy an Application
=====================
An :ref:`Application` is deployed to Deis using ``git push`` or the ``deis`` client.

Supported Applications
----------------------
Deis can deploy any application or service that can run inside a Docker container.  In order to be scaled horizontally, applications must follow Heroku's `twelve-factor methodology`_ and store state in external backing services.

For example, if your application persists state to the local filesystem -- common with content management systems like Wordpress and Drupal -- it cannot be scaled horizontally using ``deis scale``.

Fortunately, most modern applications feature a stateless application tier that can scale horizontally inside Deis.

Login to the Controller
-----------------------
Before deploying an application, users must first authenticate against the Deis :ref:`Controller`.

.. code-block:: console

    $ deis login http://deis.example.com
    username: deis
    password:
    Logged in as deis

.. note::

    For Vagrant clusters: ``deis login http://deis.local3.deisapp.com``

Select a Build Process
----------------------
Deis supports three different ways of building applications:

 1. `Heroku Buildpacks`_
 2. `Dockerfiles`_
 3. `Docker Images`_

Buildpacks
^^^^^^^^^^
Heroku buildpacks are useful if you want to follow Heroku's best practices for building applications or if you are porting an application from Heroku.

Learn how to use deploy applications on Deis :ref:`using-buildpacks`.

Dockerfiles
^^^^^^^^^^^
Dockerfiles are a powerful way to define a portable execution environment built on a base OS of your choosing.

Learn how to use deploy applications on Deis :ref:`using-dockerfiles`.

Docker Image
^^^^^^^^^^^^
Deploying a Docker image onto Deis allows you to take a Docker image from either a public
or a private registry and copy it over bit-for-bit, ensuring that you are running the same
image in development or in your CI pipeline as you are in production.

Learn how to use deploy applications on Deis :ref:`using-docker-images`.


.. _`twelve-factor methodology`: http://12factor.net/
.. _`Heroku Buildpacks`: https://devcenter.heroku.com/articles/buildpacks
.. _`Dockerfiles`: https://docs.docker.com/reference/builder/
.. _`Docker Images`: https://docs.docker.com/introduction/understanding-docker/
