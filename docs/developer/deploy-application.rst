:title: Deploy an Application on Deis
:description: First steps for developers using Deis to deploy and scale applications.
:keywords: tutorial, guide, walkthrough, howto, deis, developer, dev

Deploy an Application
=====================
Deis allows you to deploy and scale your :ref:`Application` in seconds
using Docker's industry-standard `Linux container engine`_.

Create an Application
---------------------
Change directory into a git repository for the app you'd like to deploy,
then use the ``deis create`` command to create a new Deis application.

.. code-block:: console

    $ cd example-java-jetty    # change into your application's git root
    $ deis create
    Creating application... done, created peachy-waxworks
    Git remote deis added

Deploy the Application
----------------------
With the application created and associated with the SSH :ref:`Key` on your account,
deploy it with ``git push deis master``. If you don't have an application to test with,
you can use `our Dockerfile example`_.

.. code-block:: console

    ><> deis create --formation=dev
    Creating application... done, created owlish-huntress
    Git remote deis added
    ><> git push deis master
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
    
    ><> curl -s http://owlish-huntress.example.com
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

 * Java
 * Python
 * Ruby
 * Node.js
 * Clojure
 * Scala
 * Play Framework
 * PHP
 * Perl
 * Dart
 * Go

Support for many other languages and frameworks is possible through
use of custom `Heroku Buildpacks`_ and `Dockerfiles`_.

Example Applications
--------------------

 * Clojure: https://github.com/opdemand/example-clojure-ring
 * Dart: https://github.com/opdemand/example-dart
 * Dockerfile: https://github.com/opdemand/example-dockerfile-python
 * Golang: https://github.com/opdemand/example-go
 * Java: https://github.com/opdemand/example-java-jetty
 * Node.js: https://github.com/opdemand/example-nodejs-express
 * Perl: https://github.com/opdemand/example-perl
 * PHP: https://github.com/opdemand/example-php
 * Play: https://github.com/opdemand/example-play
 * Python/Django: https://github.com/opdemand/example-python-django
 * Python/Flask: https://github.com/opdemand/example-python-flask
 * Ruby: https://github.com/opdemand/example-ruby-sinatra
 * Scala: https://github.com/opdemand/example-scala

.. _`Linux container engine`: http://docker.io/
.. _`twelve-factor methodology`: http://12factor.net/
.. _`Heroku Buildpacks`: https://devcenter.heroku.com/articles/buildpacks
.. _`Dockerfiles`: http://docs.docker.io/en/latest/use/builder/
.. _`our Dockerfile example`: https://github.com/opdemand/example-dockerfile-python