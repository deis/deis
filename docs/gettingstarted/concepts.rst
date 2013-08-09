:title: Concepts
:description: Concepts of the Deis application platform, which deploys and scales twelve factor apps. Learn about formations, layers, build, release, run & backing services.
:keywords: deis, formations, layers, build, release, run, backing services

.. _concepts:

Concepts
========

Deis is an application platform that deploys and scales `Twelve Factor`_ apps 
using a formation of `Chef`_ nodes, `Docker`_ containers and 
`Nginx`_ proxies.

Formations
----------
A :ref:`formation` is a set of infrastructure used to host a single application
or service backed by a single git repository. Each formation includes
:ref:`Layers <layer>` of :ref:`Nodes <node>` used to host services, a set of 
:ref:`Containers <container>` used to run isolated processes, and a 
:ref:`Release` that defines the current :ref:`Build` and :ref:`config` 
deployed by containers.

Layers
------
:ref:`Layers <layer>` are homogeneous groups of :ref:`Nodes <node>` that 
perform work on behalf of a formation.  Each node in a layer has 
the same :ref:`Flavor` and Chef configuration, allowing them to be scaled
easily.  Formations have two types of layers.

Runtime Layers
^^^^^^^^^^^^^^
Runtime layers service requests and run background tasks for the formation.
Nodes in a runtime layer use a `Chef Databag`_  to deploy
:ref:`Containers <container>` running a specific :ref:`Release`.  

Proxy Layers
^^^^^^^^^^^^
Proxy layers expose the formation to the outside world.
Nodes in a proxy layer use a `Chef Databag`_ to configure routing of 
inbound requests to :ref:`Containers <container>` hosted on runtime layers.

Build, Release, Run
------------------- 
Deis enforces strict separation between Build, Release and Run stages
following the `Twelve Factor model`_.

Build Stage
^^^^^^^^^^^
The :ref:`Controller` includes a `Gitosis Server`_ that receives
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
The databag specifies the current release, the placement of containers across 
the runtime layer, and the configuration of the proxy layer.
SSH is used to converge all of the nodes in the runtime layer followed 
by all of the nodes in the proxy layer, making zero downtime deployment possible.

Backing Services
----------------
In keeping with `Twelve Factor`_ methodology, `backing services`_ like
databases, queues and storage are decoupled and attached using `environment
variables`_.  This allows formations to use backing services provided via
different formations (through their proxy layer), or external/third-party 
services accessible over the network.  The use of environment variables
also allows formations to easily swap backing services when necessary.

See Also
--------
* :ref:`Installation`
* :ref:`Usage`
* :ref:`Tutorial`
* `The Twelve Factor App <http://12factor.net/>`_


.. _`Twelve Factor`: http://12factor.net/
.. _`Chef`: http://www.opscode.com/chef/
.. _`Docker`: http://docker.io/
.. _`Nginx`: http://wiki.nginx.org/Main
.. _`Chef Databag`: http://docs.opscode.com/essentials_data_bags.html
.. _`Twelve Factor model`: http://12factor.net/build-release-run
.. _`backing services`: http://12factor.net/backing-services
.. _`environment variables`: http://12factor.net/config
.. _`Gitosis Server`: https://github.com/opdemand/gitosis
.. _`Buildstep`: https://github.com/opdemand/buildstep
.. _`converges`: http://docs.opscode.com/essentials_nodes_chef_run.html
