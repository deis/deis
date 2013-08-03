:title: Technical Overview
:description: -- Technical Overview of the Deis platform
:keywords: deis, documentation, technical, overview

.. _overview:

Overview
========

Deis is an application platform that deploys and scales `Twelve Factor`_ apps 
using a formation of `Chef`_ Nodes, `Docker`_ containers and 
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
with ease.  Formations have two types of layers.

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
The Deis :ref:`Controller` includes a `Gitosis Server`_ that receives 
incoming git push requests over SSH and builds application
inside an ephemeral Docker container. A tarball of the /app directory is 
extracted into a :ref:`slug` and exposed on an Nginx static file server. 
The slug is later downloaded by the the runtime layer and bind-mounted
into a Docker container.

Release Stage
^^^^^^^^^^^^^
A :ref:`release` is a :ref:`build` combined with :ref:`config`.  
When a new build is created or config is changed,
a new release is rolled automatically.  Releases make it easy to
rollback of code and configuration.

Run Stage
^^^^^^^^^
The run stage updates Chef databags and converges all nodes in the formation, 
deploying the latest release on containers and reconfiguring proxies.   
SSH is used to converge all of the nodes in the runtime layer followed 
by all of the nodes in the proxy layer.

Backing Services
----------------
In keeping with `Twelve Factor`_ app methodology `backing services`_ like
databases, queues and storage are decoupled and attached via `environment
variables`_.  This allows formations to use backing services provided via
different formations (via their proxy layer), or external/third-party 
services accessible over the network.  The use of environment variables
also allows formations to easily swap backing services when necessary.

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
