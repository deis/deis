:title: Customizing database
:description: Learn how to tune custom Deis settings.

.. _database_settings:

Customizing database
=========================
The following settings are tunable for the :ref:`database` component.

Dependencies
------------
Requires: none

Required by: :ref:`controller <controller_settings>`

Considerations: none

Settings set by database
------------------------
The following etcd keys are set by the database component, typically in its /bin/boot script.

===========================              =================================================================================
setting                                  description
===========================              =================================================================================
/deis/database/host                      IP address of the host running database
/deis/database/port                      port used by the database service (default: 5432)
/deis/database/engine                    database engine (default: postgresql_psycopg2)
/deis/database/adminUser                 database admin user (default: postgres)
/deis/database/adminPass                 database admin password (default: changeme123)
/deis/database/user                      database user (default: deis)
/deis/database/password                  database password (default: changeme123)
/deis/database/name                      database name (default: deis)
===========================              =================================================================================

Settings used by database
-------------------------
The database component uses no keys from etcd other than the ones it sets.

Using a custom database image
-----------------------------
You can use a custom Docker image for the database component instead of the image
supplied with Deis:

.. code-block:: console

    $ etcdctl set /deis/database/image myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ etcdctl set /deis/database/image registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock database image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock database image`: https://github.com/deis/deis/tree/master/database
