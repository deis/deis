Deis Client
===========
The Deis client is a Python CLI that issues API calls to a private
Deis controller, providing a Heroku-inspired PaaS workflow.

.. image:: https://badge.fury.io/py/deis.png
    :target: http://badge.fury.io/py/deis

.. image:: https://travis-ci.org/opdemand/deis.png?branch=master
    :target: https://travis-ci.org/opdemand/deis

What is Deis?
-------------

Deis is an open source PaaS that makes it easy to deploy and scale LXC
containers and Chef nodes used to host applications, databases, middleware
and other services. Deis leverages Chef, Docker and Heroku Buildpacks to
provide a private PaaS that is lightweight and flexible.

Deis comes with out-of-the-box support for Ruby, Python, Node.js, Java,
Clojure, Scala, Play, PHP, Perl, Dart and Go. However, Deis can deploy
anything using Heroku Buildpacks, Docker images or Chef recipes. Deis is
designed to work with any cloud provider. Currently Amazon EC2, Rackspace,
and DigitalOcean are supported.


Why Deis?
=========

Deploy anything
---------------

Deploy a wide range of languages and frameworks with a simple git push
using Heroku Buildpacks or (coming soon) Dockerfiles. Use custom Chef layers
to deploy databases, middleware and other add-on services.


Control everything
------------------

Choose your hosting provider configuration. Define a formation with your
own proxy and runtime layers. Retain full root access to every node.
Manage your platform with a private Deis controller.


Scale effortlessly
------------------

Scale nodes and containers with a single command. Node provisioning,
container balancing and proxy reconfiguration are completely automated.


100% Open Source
----------------

Free, transparent and easily customized. Join the open-source PaaS
and DevOps community by using Deis and complimentary projects like
Docker, Chef and Heroku Buildpacks.


Get Started
===========

1. `Install the Client`_:

.. code-block:: console

    $ sudo pip install deis


2. `Register a User`_:

.. code-block:: console

    $ deis register http://deis.example.com
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
