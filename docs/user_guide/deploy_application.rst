:title: Deploy an Application on Deis
:description: First steps for developers using Deis to deploy and manage applications

.. _deploy_application:

Deploy an Application
=====================

:ref:`Applications <Application>` are typically deployed to Deis by pushing source code
and configuration to the system's API endpoint using the Deis client or other clients that
use the Deis API.

Supported Applications
----------------------

Deis can deploy any application or service that can run inside a Docker container.  In
order to be scaled horizontally, applications must follow Heroku's
`twelve-factor methodology`_ and store state in external backing services.

For example, if your application persists state to the local filesystem -- common with
content management systems like Wordpress and Drupal -- it cannot be scaled horizonally
using :ref:`deis scale <deis_ps>`.

Fortunately, most modern applications feature a stateless application tier that can scale
horizontally inside Deis.

Login to the Controller
-----------------------

Before deploying an application, users must first authenticate against the Deis
:ref:`Controller`.

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

Heroku buildpacks are useful if you want to follow Heroku's best practices for building
applications or if you are porting an application from Heroku.

Learn how to use deploy applications on Deis :ref:`using_buildpacks`.

Dockerfiles
^^^^^^^^^^^

Dockerfiles are a powerful way to define a portable execution environment built on a base
OS of your choosing.

Learn how to use deploy applications on Deis :ref:`using_dockerfiles`.

Push Application Code
---------------------

Change directories to your application's source code, and use :code:`git push deis master`
to deploy your application. If your application lives on a separate branch, you can use
:code:`git push deis <branch>` and it will deploy your app from the specified branch.

The output of the push command will look something like the following:

.. code-block:: console

    $ git push deis master
    Counting objects: 95, done.
    Delta compression using up to 8 threads.
    Compressing objects: 100% (52/52), done.
    Writing objects: 100% (95/95), 20.24 KiB | 0 bytes/s, done.
    Total 95 (delta 41), reused 85 (delta 37)
    -----> Ruby app detected
    -----> Compiling Ruby/Rack
    -----> Using Ruby version: ruby-1.9.3
    -----> Installing dependencies using 1.5.2
           Running: bundle install --without development:test --path vendor/bundle --binstubs vendor/bundle/bin -j4 --deployment
           Fetching gem metadata from http://rubygems.org/..........
           Fetching additional metadata from http://rubygems.org/..
           Using bundler (1.5.2)
           Installing tilt (1.3.6)
           Installing rack (1.5.2)
           Installing rack-protection (1.5.0)
           Installing sinatra (1.4.2)
           Your bundle is complete!
           Gems in the groups development and test were not installed.
           It was installed into ./vendor/bundle
           Bundle completed (8.81s)
           Cleaning up the bundler cache.
    -----> Discovering process types
           Procfile declares types -> web
           Default process types for Ruby -> rake, console, web
    -----> Compiled slug size is 12M
    -----> Building Docker image
    Uploading context 11.81 MB
    Uploading context
    Step 0 : FROM deis/slugrunner
     ---> 5567a808891d
    Step 1 : RUN mkdir -p /app
     ---> Running in a4f8e66a79c1
     ---> 5c07e1778b9e
    Removing intermediate container a4f8e66a79c1
    Step 2 : ADD slug.tgz /app
     ---> 52d48b1692e5
    Removing intermediate container e9dfce920e26
    Step 3 : ENTRYPOINT ["/runner/init"]
     ---> Running in 7a8416bce1f2
     ---> 4a18f93f1779
    Removing intermediate container 7a8416bce1f2
    Successfully built 4a18f93f1779
    -----> Pushing image to private registry

           Launching... done, v2

    -----> unisex-huntress deployed to Deis
           http://unisex-huntress.local.deisapp.com

           To learn more, use `deis help` or visit http://deis.io

    To ssh://git@local.deisapp.com:2222/unisex-huntress.git
     * [new branch]      master -> master

The Deis client will show logs for your application after it has deployed. To inspect
these logs, use the :ref:`deis logs <deis_apps>` command.


.. _`twelve-factor methodology`: http://12factor.net/
.. _`Heroku Buildpacks`: https://devcenter.heroku.com/articles/buildpacks
.. _`Dockerfiles`: http://docs.docker.io/en/latest/use/builder/
.. _`Docker Image`: http://docs.docker.io/introduction/understanding-docker/
