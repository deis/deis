:title: Customizing database
:description: Learn how to tune custom Deis settings.

.. _database_settings:

Customizing database
=========================
The following settings are tunable for the :ref:`database` component.

Dependencies
------------
Requires: :ref:`store-gateway <store_gateway_settings>`

Required by: :ref:`controller <controller_settings>`

Considerations: none

Settings set by database
------------------------
The following etcd keys are set by the database component, typically in its /bin/boot script.

===========================              =================================================================================
setting                                  description
===========================              =================================================================================
/deis/database/adminPass                 database admin password (default: changeme123)
/deis/database/adminUser                 database admin user (default: postgres)
/deis/database/bucketName                store component bucket used for database WAL logs and backups (default: db_wal)
/deis/database/engine                    database engine (default: postgresql_psycopg2)
/deis/database/host                      IP address of the host running database
/deis/database/name                      database name (default: deis)
/deis/database/password                  database password (default: changeme123)
/deis/database/port                      port used by the database service (default: 5432)
/deis/database/user                      database user (default: deis)
===========================              =================================================================================

Settings used by database
-------------------------
The following etcd keys are used by the database component.

====================================      ====================================================================================
setting                                   description
====================================      ====================================================================================
/deis/store/gateway/accessKey             S3 API access used to access the deis store gateway (set by store-gateway)
/deis/store/gateway/host                  host of the store gateway component (set by store-gateway)
/deis/store/gateway/port                  port of the store gateway component (set by store-gateway)
/deis/store/gateway/secretKey             S3 API secret key used to access the deis store gateway (set by store-gateway)
====================================      ====================================================================================

Using a custom database image
-----------------------------
You can use a custom Docker image for the database component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config database set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config database set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock database image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock database image`: https://github.com/deis/deis/tree/master/database
