:title: Concepts
:description: Concepts of the Deis application platform. Deploy and scale
 twelve-factor apps. Learn about clusters, builds, releases, run and
 backing services.

.. _concepts:

Concepts
========

Deis is an application platform that deploys and scales `Twelve-Factor`_ apps
using clusters of `CoreOS`_ nodes, `Docker`_ containers and `Nginx`_ proxies.

Clusters
--------
A :ref:`cluster` is a set of nodes used to host applications
or services.

Applications
------------
An :ref:`application` lives on a :ref:`cluster`, where it uses
:ref:`Containers <container>` to process requests and to run background tasks
for a deployed git repository.
Developers use :ref:`Applications <application>` to push code, change config,
scale containers, view logs, or run admin commands --
regardless of the cluster's underlying infrastructure.

Build, Release, Run
-------------------
Deis enforces strict separation between Build and Run stages, following
the `Twelve-Factor model`_.

Build Stage
^^^^^^^^^^^
The :ref:`Controller` includes a *gitreceive* hook that receives
incoming git push requests over SSH and builds applications
inside ephemeral Docker containers.
Tarballs of the /app directory are extracted into a slug and exposed
on the Controller using an Nginx static file server.
The slug is later downloaded by the runtime layer and bind-mounted
into a Docker container for execution.

Release Stage
^^^^^^^^^^^^^
During the release stage, a :ref:`build` is combined with :ref:`config`
to create a new numbered :ref:`release`.
The release stage is triggered any time a new build is created or
config is changed, making it easy to rollback code and configuration.

Run Stage
^^^^^^^^^
The run stage updates Chef databags and `converges`_ all nodes in the formation.
The databag specifies the current application releases,
the placement of containers across the runtime layer,
and the configuration of the proxy layer.
SSH is used to converge nodes in runtime layers followed
by nodes in proxy layers, making zero downtime deployment possible.

Backing Services
----------------
In keeping with `Twelve-Factor`_ methodology, `backing services`_ like
databases, queues and storage are decoupled and attached using `environment
variables`_.  This allows applications to use backing services provided by
other applications, or external/third-party services accessible over the network.
The use of environment variables makes it easy to swap backing services
when necessary.

See Also
--------
* :ref:`Developer Guide <developer>`
* :ref:`Operations Guide <developer>`
* `The Twelve-Factor App <http://12factor.net/>`_

.. _`Twelve-Factor`: http://12factor.net/
.. _`CoreOS`: https://coreos.com/
.. _`Docker`: http://docker.io/
.. _`Nginx`: http://wiki.nginx.org/Main
.. _`Chef Databag`: http://docs.opscode.com/essentials_data_bags.html
.. _`Twelve-Factor model`: http://12factor.net/build-release-run
.. _`backing services`: http://12factor.net/backing-services
.. _`environment variables`: http://12factor.net/config
.. _`converges`: http://docs.opscode.com/essentials_nodes_chef_run.html
