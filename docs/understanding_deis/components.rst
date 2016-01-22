:title: Components
:description: Components of the Deis application platform (PaaS)

.. _components:

Components
==========

Deis consists of a number of components that combine to create a distributed PaaS.
Each Deis component is deployed as a container or set of containers.

.. _comp_controller:

Controller
----------
The controller component is an HTTP API server. Among other functions, the
controller contains :ref:`the scheduler <scheduler>`, which decides
where to run app containers.
The ``deis`` command-line client interacts with this component.

.. _database:

Database
--------
The database component is a `PostgreSQL`_ server used to store durable
platform state. Backups and WAL logs are pushed to :ref:`Store`.

.. _builder:

Builder
-------
The builder component uses a `Git`_ server to process
:ref:`Application` builds. The builder:

 #. Receives incoming ``git push`` requests over SSH
 #. Authenticates the user via SSH key fingerprint
 #. Authorizes the user's access to write to the Git repository
 #. Builds a new `Docker` image from the updated git repository
 #. Adds the latest :ref:`Config` to the resulting Docker image
 #. Pushes the new Docker image to the platform's :ref:`Registry`
 #. Triggers a new :ref:`Release` through the :ref:`Controller`

.. note::

    The builder component does not incorporate :ref:`Config` directly into the
    images it produces.   A :ref:`Release` is a pairing of an application image
    with application configuration maintained separately in the Deis
    :ref:`Database`.

    Once a new :ref:`Release` is generated, a new set of containers
    is deployed across the platform automatically.

.. _registry:

Registry
--------
The registry component hosts `Docker`_ images on behalf of the platform.
Image data is stored by :ref:`Store`.

.. _logspout:

Logspout
--------
The logspout component is a customized version of `progrium's logspout`_ that runs
on all CoreOS hosts in the cluster and collects logs from running containers.
It sends the logs to the :ref:`logger` component.

.. _logger:

Logger
------
The logger component is a syslog server that collects logs from :ref:`logspout`
components spread across the platform.
This data can then be queried by the :ref:`Controller`.

.. _publisher:

Publisher
---------
The publisher component is a microservice written in Go that publishes
containers to etcd so they can be exposed by the platform :ref:`router`.

.. _router:

Router
------
The router component uses `Nginx`_ to route traffic to application containers.

.. _store:

Store
------
The store component uses `Ceph`_ to store data for Deis components
which need to store state, including :ref:`Registry`, :ref:`Database`
and :ref:`Logger`.

.. _`Amazon S3`: http://aws.amazon.com/s3/
.. _`Celery`: http://www.celeryproject.org/
.. _`Ceph`: http://ceph.com
.. _`Docker`: http://docker.io/
.. _`etcd`: https://github.com/coreos/etcd
.. _`Git`: http://git-scm.com/
.. _`Nginx`: http://nginx.org/
.. _`OpenStack Storage`: http://www.openstack.org/software/openstack-storage/
.. _`PostgreSQL`: http://www.postgresql.org/
.. _`progrium's logspout`: https://github.com/progrium/logspout
.. _`Redis`: http://redis.io/
