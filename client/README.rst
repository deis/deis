Deis Client
===========
The Deis client is a Python CLI that issues API calls to a private
Deis controller, providing a Heroku-inspired PaaS workflow.

.. image:: https://badge.fury.io/py/deis.png
    :target: http://badge.fury.io/py/deis

.. image:: https://travis-ci.org/opdemand/deis.png?branch=master
    :target: https://travis-ci.org/opdemand/deis

.. image:: https://pypip.in/d/deis/badge.png
    :target: https://crate.io/packages/deis/

What is Deis?
-------------

Deis is an open source PaaS that makes it easy to deploy and scale LXC
containers and Chef nodes used to host applications, databases, middleware
and other services. Deis leverages Chef, Docker and Heroku Buildpacks to
provide a private PaaS that is lightweight and flexible.

Deis comes with out-of-the-box support for Ruby, Python, Node.js, Java,
Clojure, Scala, Play, PHP, Perl, Dart and Go. However, Deis can deploy
anything using Heroku Buildpacks, Docker images or Chef recipes. Deis is
designed to work with any cloud provider, although only EC2 is currently
supported.


Why Deis?
=========

Deploy anything
---------------

Deploy a wide range of languages and frameworks with a simple git push 
using Heroku Buildpacks or Dockerfiles. Use custom Chef layers to deploy
databases, middleware and other add-on services.


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


Getting Started
---------------

Installing the Deis client from the Python Package Index is simple:

::

    $ pip install deis

The client will automatically install the following dependencies:

-  `pyyaml <https://bitbucket.org/xi/pyyaml>`__
-  `requests <http://python-requests.org>`__


License
-------

Copyright 2013, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.
