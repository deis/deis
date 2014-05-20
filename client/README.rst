Deis Client
===========
The Deis client is a Python CLI that issues API calls to a private
Deis controller, providing a Heroku-inspired PaaS workflow.

.. image:: https://badge.fury.io/py/deis.png
    :target: http://badge.fury.io/py/deis

.. image:: https://travis-ci.org/deis/deis.png?branch=master
    :target: https://travis-ci.org/deis/deis

.. image:: https://pypip.in/d/deis/badge.png
    :target: https://pypi.python.org/pypi/deis/
    :alt: Downloads

.. image:: https://pypip.in/license/deis/badge.png
    :target: https://pypi.python.org/pypi/deis/
    :alt: License

What is Deis?
-------------

Deis is an open source PaaS that makes it easy to deploy and scale containers
to host applications, databases, middleware and other services. Deis leverages
Docker, CoreOS and Heroku Buildpacks to provide a private PaaS that is
lightweight and flexible.

Deis comes with out-of-the-box support for Ruby, Python, Node.js, Java,
Clojure, Scala, Play, PHP, Perl, Dart and Go. However, Deis can deploy
anything using Docker images or Heroku Buildpacks. Deis is designed to work
with any cloud provider. Currently Amazon EC2 and Rackspace are supported.


Why Deis?
=========

Deploy anything
---------------

Deploy a wide range of languages and frameworks with a simple git push
using Heroku Buildpacks or Dockerfiles.


Control everything
------------------

Choose your hosting provider configuration. Define a cluster to meet your own
needs. Retain full root access to every node. Manage your platform with a
private Deis controller.


Scale effortlessly
------------------

Add nodes automatically and scale containers with a single command. Smart
scheduling, container balancing and proxy reconfiguration are completely
automated.


100% Open Source
----------------

Free, transparent and easily customized. Join the open-source PaaS
and DevOps community by using Deis and complimentary projects like
Docker, CoreOS and Heroku Buildpacks.


Get Started
===========

1. `Install the Client`_:

.. code-block:: console

    $ sudo pip install deis

If you don't have `Python`_ installed, you can download a binary executable
version of the Deis client for Mac OS X, Windows, or Debian Linux:

    - https://s3-us-west-2.amazonaws.com/opdemand/deis-osx-0.8.0.tgz
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-win32-0.8.0.zip
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-deb-wheezy-0.8.0.tgz

2. `Register a User`_:

.. code-block:: console

    $ deis register http://deis.local.deisapp.com
    $ deis keys:add


3. `Deploy an Application`_:

.. code-block:: console

    $ deis create
    Creating application... done, created peachy-waxworks
    Git remote deis added
    $ git push deis master
           Java app detected
    -----> Installing OpenJDK 1.6... done
    ...
    -----> Compiled slug size: 63.5 MB
           Launching... done, v2

    -----> peachy-waxworks deployed to Deis
           http://peachy-waxworks.example.com ...

    $ curl -s http://peachy-waxworks.example.com
    Powered by Deis!


4. `Manage an Application`_:

.. code-block:: console

    $ deis config:set DATABASE_URL=postgres://user:pass@example.com:5432/db
    $ deis scale web=8
    $ deis run ls -l  # the view from inside a container
    total 28
    -rw-r--r-- 1 root root  553 Dec  2 23:59 LICENSE
    -rw-r--r-- 1 root root   60 Dec  2 23:59 Procfile
    -rw-r--r-- 1 root root   33 Dec  2 23:59 README.md
    -rw-r--r-- 1 root root 1622 Dec  2 23:59 pom.xml
    drwxr-xr-x 3 root root 4096 Dec  2 23:59 src
    -rw-r--r-- 1 root root   25 Dec  2 23:59 system.properties
    drwxr-xr-x 6 root root 4096 Dec  3 00:00 target


To learn more, use ``deis help`` or browse `the documentation`_.

.. _`Install the Client`: http://docs.deis.io/en/latest/developer/install-client/
.. _`Python`: http://www.python.org/
.. _`Register a User`: http://docs.deis.io/en/latest/developer/register-user/
.. _`Deploy an Application`: http://docs.deis.io/en/latest/developer/deploy-application/
.. _`Manage an Application`: http://docs.deis.io/en/latest/developer/manage-application/
.. _`the documentation`: http://docs.deis.io/


License
-------

Copyright 2013, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at `<http://www.apache.org/licenses/LICENSE-2.0>`__.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.
