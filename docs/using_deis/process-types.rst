:title: Process Types and the Procfile
:description: First steps for using the 12 factor process model with Deis.

.. _process-types:

Process Types and the Procfile
==============================

A Procfile is a mechanism for declaring what commands are run by your application’s containers on
the Deis platform. It follows the `process model`_. You can use a Procfile to declare various
process types, such as multiple types of workers, a singleton process like a clock, or a consumer
of the Twitter streaming API.

Process Types as Templates
--------------------------

A Procfile is a text file named ``Procfile`` placed in the root of your application that lists the
process types in an application. Each process type is a declaration of a command that is executed
when a container of that process type is started.

All the language and frameworks using :ref:`Heroku's Buildpacks <using-buildpacks>` declare a
``web`` process type, which starts the application server. Rails 3 has the following process type:

.. code-block:: console

    web: bundle exec rails server -p $PORT

All applications using :ref:`Dockerfile deployments <using-dockerfiles>` have an implied ``cmd``
process type, which spawns the default process of a Docker image:

.. code-block:: console

    $ cat Dockerfile
    FROM centos:latest
    COPY . /app
    WORKDIR /app
    CMD python -m SimpleHTTPServer 5000
    EXPOSE 5000

For applications using :ref:`Docker image deployments <using-docker-images>`, a ``cmd`` process
type is also implied and spawns the default process of the image.

Declaring Process Types
-----------------------

Process types are declared via a file named ``Procfile``, placed in the root of your app. Its
format is one process type per line, with each line containing:

.. code-block:: console

    <process type>: <command>

The syntax is defined as:

``<process type>`` – an alphanumeric string, is a name for your command, such as web, worker, urgentworker, clock, etc.

``<command>`` – a command line to launch the process, such as ``rake jobs:work``.

.. note::

    The web and cmd process types are special as they’re the only process types that will receive
    HTTP traffic from Deis’s routers. Other process types can be named arbitrarily.

Deploying to Deis
-----------------

A ``Procfile`` is not necessary to deploy most languages supported by Deis. The platform
automatically detects the language and supplies a default ``web`` process type to boot the server.

Creating an explicit Procfile is recommended for greater control and flexibility over your app.

For Deis to use your Procfile, add the Procfile to the root of your application, then push to Deis:

.. code-block:: console

    $ git add .
    $ git commit -m "Procfile"
    $ git push deis master
    ...
    -----> Procfile declares process types: web, worker
    Compiled slug size is 10.4MB

           Launching... done, v2

    -----> unisex-huntress deployed to Deis
           http://unisex-huntress.example.com

For Docker image deployments, a Procfile in the current directory or specified by
``deis pull --procfile`` will define the default process types for the application.

Use ``deis scale web=3`` to increase ``web`` processes to 3, for example. Scaling a
process type directly changes the number of :ref:`Containers <container>`
running that process.

Web vs Cmd Process Types
------------------------

When deploying to Deis using a Heroku Buildpack, Deis boots the ``web`` process type to boot the
application server. When you deploy an application that has a Dockerfile or uses :ref:`Docker
images <using-docker-images>`, Deis boots the ``cmd`` process type. Both act similarly in that they
are exposed to the router as web applications. However, The ``cmd`` process type is special because
it is equivalent to running the :ref:`container` without any additional arguments. Every other
process type is equivalent to running the relevant command that is provided in the Procfile.

When migrating from Heroku Buildpacks to a Docker-based deployment, Deis will not convert ``web``
process types to ``cmd``. To do this, you'll have to manually scale down the old process type and
scale the new process type up.


.. _`process model`: https://devcenter.heroku.com/articles/process-model
