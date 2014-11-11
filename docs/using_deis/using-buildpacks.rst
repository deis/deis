:title: Deploying with Heroku Buildpacks on Deis
:description: How to deploy applications on Deis using Heroku Buildpacks

.. _using-buildpacks:

Using Buildpacks
================
Deis supports deploying applications via `Heroku Buildpacks`_.
Buildpacks are useful if you're interested in following Heroku's best practices for building applications or if you are deploying an application that already runs on Heroku.

Prepare an Application
----------------------
If you do not have an existing application, you can clone an example application that demonstrates the Heroku Buildpack workflow.

.. code-block:: console

    $ git clone https://github.com/deis/example-ruby-sinatra.git
    $ cd example-ruby-sinatra

Create an Application
---------------------
Use ``deis create`` to create an application on the :ref:`controller`.

.. code-block:: console

    $ deis create
    Creating application... done, created unisex-huntress
    Git remote deis added

Push to Deploy
--------------
Use ``git push deis master`` to deploy your application.

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

    $ curl -s http://unisex-huntress.local.deisapp.com
    Powered by Deis!

Because a Heroku-style application is detected, the ``web`` process type is automatically scaled to 1 on first deploy.

Included Buildpacks
-------------------
For convenience, a number of buildpacks come bundled with Deis:

 * `Ruby Buildpack`_
 * `Nodejs Buildpack`_
 * `Java Buildpack`_
 * `Gradle Buildpack`_
 * `Grails Buildpack`_
 * `Play Buildpack`_
 * `Python Buildpack`_
 * `PHP Buildpack`_
 * `Clojure Buildpack`_
 * `Scala Buildpack`_
 * `Go Buildpack`_
 * `Multi Buildpack`_

Deis will cycle through the ``bin/detect`` script of each buildpack to match the code you
are pushing.

Using a Custom Buildpack
------------------------
To use a custom buildpack, set the ``BUILDPACK_URL`` environment variable.

.. code-block:: console

    $ deis config:set BUILDPACK_URL=https://github.com/dpiddy/heroku-buildpack-ruby-minimal
    Creating config... done, v2

    === humble-autoharp
    BUILDPACK_URL: https://github.com/dpiddy/heroku-buildpack-ruby-minimal

On your next ``git push``, the custom buildpack will be used.


.. _`Ruby Buildpack`: https://github.com/heroku/heroku-buildpack-ruby
.. _`Nodejs Buildpack`: https://github.com/heroku/heroku-buildpack-nodejs
.. _`Java Buildpack`: https://github.com/heroku/heroku-buildpack-java
.. _`Gradle Buildpack`: https://github.com/heroku/heroku-buildpack-gradle
.. _`Grails Buildpack`: https://github.com/heroku/heroku-buildpack-grails
.. _`Play Buildpack`: https://github.com/heroku/heroku-buildpack-play
.. _`Python Buildpack`: https://github.com/heroku/heroku-buildpack-python
.. _`PHP Buildpack`: https://github.com/deis/heroku-buildpack-php
.. _`Clojure Buildpack`: https://github.com/heroku/heroku-buildpack-clojure
.. _`Scala Buildpack`: https://github.com/heroku/heroku-buildpack-scala
.. _`Go Buildpack`: https://github.com/kr/heroku-buildpack-go
.. _`Multi Buildpack`: https://github.com/heroku/heroku-buildpack-multi
.. _`Heroku Buildpacks`: https://devcenter.heroku.com/articles/buildpacks
