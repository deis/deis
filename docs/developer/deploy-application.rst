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
deploy it with ``git push deis master``.

.. code-block:: console

    $ git push deis master
    Counting objects: 17, done.
    Delta compression using up to 8 threads.
    Compressing objects: 100% (10/10), done.
    Writing objects: 100% (17/17), 2.35 KiB, done.
    Total 17 (delta 2), reused 0 (delta 0)
           Java app detected
    -----> Installing OpenJDK 1.6... done
    -----> Installing Maven 3.0.3... done
    -----> Installing settings.xml... done
    -----> executing /cache/.maven/bin/mvn -B -Duser.home=/build/app -Dmaven.repo.local=/cache/.m2/repository -s /cache/.m2/settings.xml -DskipTests=true clean install
           [INFO] Scanning for projects...
           ...
           [INFO] -------------------------------------------
           [INFO] BUILD SUCCESS
           [INFO] -------------------------------------------
           [INFO] Total time: 11.771s
           [INFO] Finished at: Tue Dec 03 00:00:03 UTC 2013
           [INFO] Final Memory: 12M/142M
           [INFO] -------------------------------------------
    -----> Discovering process types
           Procfile declares types -> web
    
    -----> Compiled slug size: 63.5 MB
           Launching... done, v2

    -----> peachy-waxworks deployed to Deis
           http://peachy-waxworks.example.com ...

    $ curl -s http://peachy-waxworks.example.com
    Powered by Deis!

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

.. _`Linux container engine`: http://docker.io/
.. _`twelve-factor methodology`: http://12factor.net/
.. _`Heroku Buildpacks`: https://devcenter.heroku.com/articles/buildpacks
.. _`Dockerfiles`: http://docs.docker.io/en/latest/use/builder/