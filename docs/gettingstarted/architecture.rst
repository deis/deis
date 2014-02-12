:title: Deis Architecture
:description: Architecture of the Deis application platform (PaaS)
:keywords: deis, paas, application platform, architecture

.. _architecture:

Architecture
============

Deis consists of 8 modules that combine to create a distributed PaaS.
Each Deis module is deployed as one or more `Docker`_ containers.

.. _controller:

Controller
----------
The controller module is the "brains" of the Deis platform, in charge of:

* Processing client API calls
* Managing nodes that host containers and provide services
* Managing containers that perform work
* Managing proxies that route traffic to containers
* Managing users, providers, flavors, keys and other base configuration

The controller module includes:

* `Django`_ for processing API calls
* `Celery`_ for managing task queues

.. _database:

Database
--------
The database module uses `PostgreSQL`_  to store durable platform state.

.. _cache:

Cache
-----
The cache module uses `Redis`_ to:

 * Store work queue data for Celery
 * Cache sessions and synchronize locks for Django
 * Store recent log data for the :ref:`Controller`

.. _builder:

Builder
-------
The builder module uses a `Git`_ server to process :ref:`Application` builds.
The builder:

 #. Receives incoming ``git push`` requests over SSH
 #. Authenticates the user via SSH key fingerprint
 #. Authorizes the user's access to write to the Git repository
 #. Builds a new `Docker` image from the updated git repository
 #. Adds the latest :ref:`Config` to the resulting Docker image
 #. Pushes the new Docker image to the platform's :ref:`Registry`
 #. Creates a new :ref:`Release` on the :ref:`Controller`

Once a new :ref:`Release` is generated, a new set of containers 
is deployed across the platform automatically.

.. _registry:

Registry
--------
The registry module hosts `Docker`_ images on behalf of the platform.
Image data is typically stored on a storage service like 
`Amazon S3`_ or `OpenStack Storage`_.

.. _logserver:

Log Server
----------
The log server module uses `rsyslog`_ to aggregate log data from 
across the platform.
This data can then be queried by the :ref:`Controller`.

.. _runtime:

Runtime
-------
The runtime module uses `Docker`_ to run containers for deployed applications.

.. _proxy:

Proxy
-----
The proxy module uses `Nginx`_ to route traffic to application containers.
 
.. _`Django`: https://www.djangoproject.com/
.. _`Celery`: http://www.celeryproject.org/
.. _`PostgreSQL`: http://www.postgresql.org/
.. _`Redis`: http://redis.io/
.. _`Git`: http://git-scm.com/
.. _`Docker`: http://docker.io/
.. _`Amazon S3`: http://aws.amazon.com/s3/
.. _`OpenStack Storage`: http://www.openstack.org/software/openstack-storage/
.. _`rsyslog`: http://www.rsyslog.com/
.. _`Nginx`: http://nginx.org/
