:title: Deploy an Application on Deis
:description: First steps for developers using Deis to deploy and manage applications

.. _deploy-application:

Deploy an Application
=====================
An :ref:`Application` is deployed to Deis using ``git push`` or the ``deis`` client.

Supported Applications
----------------------
Deis can deploy any application or service that can run inside a Docker container.  In order to be scaled horizontally, applications must follow Heroku's `twelve-factor methodology`_ and store state in external backing services.

For example, if your application persists state to the local filesystem -- common with content management systems like Wordpress and Drupal -- it cannot be scaled horizonally using ``deis scale``.

Fortunately, most modern applications feature a stateless application tier that can scale horizontally inside Deis.

Login to the Controller
-----------------------
Before deploying an application, users must first authenticate against the Deis :ref:`Controller`.

.. code-block:: console

    $ deis login http://deis.example.com
    username: deis
    password:
    Logged in as deis

Select a Build Process
----------------------
Deis supports three different ways of building applications:

 1. `Heroku Buildpacks`_
 2. `Dockerfiles`_
 3. `Docker Image`_ (coming soon)

Buildpacks
^^^^^^^^^^
Heroku buildpacks are useful if you want to follow Heroku's best practices for building applications or if you are porting an application from Heroku.

Learn how to use deploy applications on Deis :ref:`using-buildpacks`.

Dockerfiles
^^^^^^^^^^^
Dockerfiles are a powerful way to define a portable execution environment built on a base OS of your choosing.

Learn how to use deploy applications on Deis :ref:`using-dockerfiles`.


.. _`twelve-factor methodology`: http://12factor.net/
.. _`Heroku Buildpacks`: https://devcenter.heroku.com/articles/buildpacks
.. _`Dockerfiles`: http://docs.docker.io/en/latest/use/builder/
.. _`Docker Image`: http://docs.docker.io/introduction/understanding-docker/
