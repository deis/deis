:title: Running Deis without Ceph
:description: Configuring the cluster to remove Ceph from the control plane.

.. _running-deis-without-ceph:

Running Deis without Ceph
=========================

.. include:: ../_includes/_ceph-dependency-description.rst

This guide is intended to assist users who are interested in removing the Ceph
dependency of the Deis control plane.

.. note::

  This guide was adapted from content graciously provided by Deis community member
  `Arne-Christian Blystad`_.

Requirements
------------

External services are required to replace the internal store components:

* S3-compatible blob store (like `Amazon S3`_)
* PostgreSQL database (like `Amazon RDS`_)
* Log drain service with syslog log format compatibility (like `Papertrail`_)

Understanding component changes
-------------------------------

Either directly or indirectly, all components in the :ref:`control-plane`
require Ceph (:ref:`store`). Some components require changes to accommodate
the removal of Ceph. The necessary changes are described below.

Logger
^^^^^^

The :ref:`logspout` component attaches to Docker containers on each host and
listens for log events from platform components and running applications.  It
ships these to the :ref:`logger` component.  By default, the logger writes the
logs to a distributed Ceph filesystem. These logs can then be fetched by the
:ref:`controller` component via HTTP.

In a Ceph-less cluster, the Logger component should be configured, instead, to
use in-memory log storage.  Optionally, a drain may also be configured to forward
logs to an external log service (such as Papertrail) for longer-term archival.

Database
^^^^^^^^

The :ref:`database` runs PostgreSQL and uses the Ceph S3 API (provided by
``deis-store-gateway``) to store PostgreSQL backups and WAL logs.
Should the host running database fail, the database component will fail over to
a new host, start up, and replay backups and WAL logs to recover to its
previous state.

We will not be using the database component in the Ceph-less cluster, and will
instead rely on an external database.

When provisioning the database, it is strongly recommended to use an `m3.medium`
instance or greater.

Registry
^^^^^^^^

The :ref:`registry` component is an instance of the offical Docker registry, and
is used to store application releases. The registry supports any S3 store, so
a Ceph-less cluster will simply reconfigure registry to use another store (typically
Amazon S3 itself).

Builder
^^^^^^^

The :ref:`builder` component is responsible for building applications deployed
to Deis via the ``git push`` workflow. It pushes to registry to store releases,
so it will require no changes.

Store
^^^^^

The :ref:`store` components implement Ceph itself. In a Ceph-less cluster, we
will skip the installation and starting of these components.

Deploying the cluster
---------------------

This guide assumes a typical deployment on AWS by following the :ref:`deis_on_aws`
guide.

Deploy an AWS cluster
^^^^^^^^^^^^^^^^^^^^^

Follow the :ref:`deis_on_aws` installation documentation through the "Configure
DNS" portion.

Configure logger
^^^^^^^^^^^^^^^^

The :ref:`logger` component should be configured to use in-memory storage. Optionally
it may also be configured to drain logs to an external service for longer-term
archival.

.. code-block:: console

    $ STORAGE_ADAPTER=memory
    $ DRAIN=udp://logs.somewhere.com:12345 # Supported protocols are udp and tcp; for backwards compatibility, "syslog" is an alias for udp
    $ deisctl config logs set storageAdapterType=${STORAGE_ADAPTER} drain=${DRAIN}

Configure registry
^^^^^^^^^^^^^^^^^^

The :ref:`registry` component won't start until it's configured with a store.

S3 store configuration sample:

.. code-block:: console

    $ BUCKET=MYS3BUCKET
    $ AWS_S3_REGION=some-aws-region #(e.g., us-west-1)
    $ deisctl config registry set s3bucket=${BUCKET} \
                                  s3region=${AWS_S3_REGION} \
                                  s3path=/ \
                                  s3encrypt=false \
                                  s3secure=false

Due to `issue 4568`_, you'll also need to run the following to ensure confd will template out the
registry's configuration:

.. code-block:: console

    $ deisctl config store set gateway=' '

By default, the registry will try to authenticate to S3 using the instance role.
If your cluster is not running on EC2, you can supply hard coded API access and
secret key:

.. code-block:: console

    $ deisctl config registry set s3accessKey=your-access-key \
                                  s3secretKey=your-secret-key

For reference, here's example of a policy you could attach to the role/user used by
the registry:

.. code-block:: javascript

    {
      "Statement": [
        {
          "Resource": [
            "arn:aws:s3:::MYBUCKET"
          ],
          "Action": [
            "s3:ListBucket",
            "s3:GetBucketLocation"
          ],
          "Effect": "Allow"
        },
        {
          "Resource": [
            "arn:aws:s3:::MYBUCKET/*"
          ],
          "Action": [
            "s3:GetObject",
            "s3:PutObject",
            "s3:DeleteObject"
          ],
          "Effect": "Allow"
        }
      ],
      "Version": "2012-10-17"
    }

Openstack-swift support requires `Swift3`_ middleware to be installed. Here is a sample configuration:

.. code-block:: console

    $ SWIFT_CONTAINER=mycontainer
    $ SWIFT_USER=system:root
    $ SWIFT_SECRET_KEY=testpass
    $ deisctl config registry set bucketName=${SWIFT_CONTAINER}
    $ deisctl config store set gateway/accessKey=${SWIFT_USER} \
                               gateway/secretKey=${SWIFT_SECRET_KEY} \
                               gateway/host=10.1.50.1 \
                               gateway/port=8080

Configure database settings
^^^^^^^^^^^^^^^^^^^^^^^^^^^

Since we won't be running the :ref:`database`, we need to configure these settings
so the controller knows where to connect.

.. code-block:: console

    $ HOST=something.rds.amazonaws.com
    $ DB_USER=deis
    $ DB_PASS=somethingsomething
    $ DATABASE=deis
    $ deisctl config database set engine=postgresql_psycopg2 \
                                  host=${HOST} \
                                  port=5432 \
                                  name=${DATABASE} \
                                  user=${DB_USER} \
                                  password=${DB_PASS}

Deploy the platform
^^^^^^^^^^^^^^^^^^^

The typical :ref:`install_deis_platform` documentation can be followed, with
one caveat: since we won't be deploying many of the typical Deis components, we cannot
use ``deisctl install platform`` or ``deisctl start platform`` -- instead, we
use ``deisctl install stateless-platform`` and ``deisctl start stateless-platform``.

These commands tell ``deisctl`` to skip the components that we don't need to use.

Confirm installation
^^^^^^^^^^^^^^^^^^^^

That's it! Deis is now running without Ceph. Issue a ``deisctl list`` to confirm
that the services are started, and see :ref:`using_deis` to start using the cluster.

Upgrading Deis
--------------

When following the :ref:`upgrading-deis` documentation, be sure to use
``stateless-platform`` instead of ``platform``.

.. _`Amazon RDS`: http://aws.amazon.com/rds/
.. _`Amazon S3`: http://aws.amazon.com/s3/
.. _`Arne-Christian Blystad`: https://github.com/blystad
.. _`issue 4568`: https://github.com/deis/deis/issues/4568
.. _`Papertrail`: https://papertrailapp.com/
.. _`Swift3`: https://github.com/openstack/swift3
