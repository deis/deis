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

Installing the Deis client using `pip`_ is simple::

    $ pip install deis

`pip`_ will automatically install the following dependencies:

-  `docopt <http://docopt.org>`__
-  `pyyaml <https://bitbucket.org/xi/pyyaml>`__
-  `requests <http://python-requests.org>`__

You should know the fully-qualified domain name of an existing
Deis controller. To set up a Deis controller, see the
`Installation`_ documentation.

Registration will discover SSH keys automatically and use the
`standard environment variables`_ to configure the EC2 provider.

.. code-block:: console

    $ deis register http://my-deis-controller.fqdn
    username: myuser
    password:
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

    Found EC2 credentials: AKIAJTVXXXXXXXXXXXXX
    Import these credentials? (y/n) : y
    Uploading EC2 credentials... done

Find an application you’d like to deploy, or clone `an example app`_.

Change into the application directory and use ``deis create`` to initialize
a new formation in a specific EC2 region.

Use the ``deis layers:scale`` command to provision nodes that will be
dedicated to this formation.

.. code-block:: console

    $ cd <my-application-repo>
    $ deis create --flavor=ec2-us-west-2
    Creating formation... done, created peachy-waxworks
    Git remote deis added

    Creating runtime layer... done
    Creating proxy layer... done

    Use deis layers:scale proxy=1 runtime=1 to scale a basic formation

    $ deis layers:scale proxy=1 runtime=1
    Scaling layers... but first, coffee!
    ...done in 232s

    Use `git push deis master` to deploy to your formation


Use ``git push deis master`` to deploy your application.

Deis will automatically deploy Docker containers and configure Nginx proxies
to route requests to your application.

.. code-block:: console

    $ git push deis master
    Counting objects: 146, done.
    Delta compression using up to 8 threads.
    Compressing objects: 100% (122/122), done.
    Writing objects: 100% (146/146), 21.54 KiB, done.
    Total 146 (delta 84), reused 47 (delta 22)
           Node.js app detected
    -----> Resolving engine versions
           Using Node.js version: 0.10.15
           Using npm version: 1.2.30
    ...
    -----> Building runtime environment
    -----> Discovering process types
           Procfile declares types -> web

    -----> Compiled slug size: 4.7 MB
           Launching... done, v2

    -----> peachy-waxworks deployed to Deis
           http://ec2-198.51.100.36.us-west-2.compute.amazonaws.com ...

    $ curl -s http://ec2-198.51.100.36.us-west-2.compute.amazonaws.com
    Powered by Deis!

To learn more, use ``deis help`` or browse `the documentation`_.

.. _`pip`: http://www.pip-installer.org/en/latest/installing.html
.. _`Installation`: http://docs.deis.io/en/latest/gettingstarted/installation/
.. _`standard environment variables`: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux
.. _`an example app`: https://github.com/opdemand/example-nodejs-express
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
