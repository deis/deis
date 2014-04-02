:title: Deploy an Application on Deis
:description: First steps for developers using Deis to deploy and scale applications.
:keywords: tutorial, guide, walkthrough, howto, deis, developer, dev

Deploy an Application
=====================

An :ref:`Application` is typically deployed to Deis by pushing source code using the deis
client or other clients that communicate with Deis' API endpoints. Deploying
applications will be different depending on the source code and its requirements.

Authenticating with the API
---------------------------

Before deploying an application, all users must first authenticate against the Deis
:ref:`Controller`. For example:

.. code-block:: console

    $ deis login http://example.com
    username: deis
    password:
    Logged in as deis

Create an Application
---------------------

Change to the root directory of your project you'd like to deploy, then use the ``deis
create`` command to create a remote repository for you to push your application to.

.. code-block:: console

    $ cd example-java-jetty    # change into your application's git root
    $ deis create --formation=dev
    Creating application... done, created peachy-waxworks
    Git remote deis added

Deploy the Application
----------------------

With the application created and associated with the SSH :ref:`Key` on your account,
deploy it with ``git push deis master``. If you don't have an application to test with,
you can use `our Dockerfile example`_.

.. code-block:: console

    $ deis create --formation=dev
    Creating application... done, created owlish-huntress
    Git remote deis added
    $ git push deis master
    Counting objects: 10, done.
    Delta compression using up to 8 threads.
    Compressing objects: 100% (9/9), done.
    Writing objects: 100% (10/10), 1.70 KiB | 0 bytes/s, done.
    Total 10 (delta 0), reused 0 (delta 0)
    -----> Building Docker image
    Uploading context 5.632 kB
    Uploading context
    Step 0 : FROM ubuntu:12.04
     ---> 9cd978db300e
    Step 1 : MAINTAINER OpDemand <info@opdemand.com>
     ---> Running in 9aefab8ad92c
     ---> da93d76703b7
    Step 2 : ENV PORT 8000
     ---> Running in 8ce25ddf4405
     ---> b6046ec54bb3
    Step 3 : ADD . /app
     ---> 5567f79d87fe
    Step 4 : WORKDIR /app
     ---> Running in 0b2c7906381c
     ---> 444006758e39
    Step 5 : CMD python -m SimpleHTTPServer $PORT
     ---> Running in b33074f3c0ea
     ---> 5a55b32b8da2
    Successfully built 5a55b32b8da2
    -----> Pushing image to private registry
    
           Launching... done, v2
    
    -----> owlish-huntress deployed to Deis
           http://owlish-huntress.example.com
    
           To learn more, use `deis help` or visit http://deis.io
    
    $ curl -s http://owlish-huntress.example.com
    <h1>Powered by Deis</h1>

Supported Applications
----------------------

As a Heroku-inspired Platform-as-a-Service, Deis is designed to deploy and scale
apps that adhere to `twelve-factor methodology`_.  

For example, if your application persists state to the local filesystem 
-- common with content management systems like Wordpress and Drupal -- 
it is not twelve-factor compatible and may not be suitable for Deis or other PaaSes.

Fortunately, most modern applications feature a stateless application tier that
can scale horizontally behind a load balancer.  These applications are a perfect
fit for Deis.  Deis currently suppports the following languages:

 * `Clojure`_
 * `Dart`_
 * `Dockerfile`_
 * `Golang`_
 * `Java`_
 * `Nodejs`_
 * `Perl`_
 * `PHP`_
 * `Play`_
 * `Python`_
 * `Ruby`_
 * `Scala`_

Support for many other languages and frameworks is possible through
use of custom `Heroku Buildpacks`_ and `Dockerfiles`_.

.. _`Clojure`: https://github.com/opdemand/example-clojure-ring
.. _`Dart`: https://github.com/opdemand/example-dart
.. _`Dockerfile`: https://github.com/opdemand/example-dockerfile-python
.. _`Golang`: https://github.com/opdemand/example-go
.. _`Java`: https://github.com/opdemand/example-java-jetty
.. _`Nodejs`: https://github.com/opdemand/example-nodejs-express
.. _`Perl`: https://github.com/opdemand/example-perl
.. _`PHP`: https://github.com/opdemand/example-php
.. _`Play`: https://github.com/opdemand/example-play
.. _`Python`: https://github.com/opdemand/example-python-flask
.. _`Ruby`: https://github.com/opdemand/example-ruby-sinatra
.. _`Scala`: https://github.com/opdemand/example-scala
.. _`Linux container engine`: http://docker.io/
.. _`twelve-factor methodology`: http://12factor.net/
.. _`Heroku Buildpacks`: https://devcenter.heroku.com/articles/buildpacks
.. _`Dockerfiles`: http://docs.docker.io/en/latest/use/builder/
.. _`our Dockerfile example`: https://github.com/opdemand/example-dockerfile-python
