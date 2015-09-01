:title: Backing Up and Restoring Data
:description: Backing up stateful data on Deis.

.. _backing_up_data:

Backing Up and Restoring Data
=============================

While applications deployed on Deis follow the Twelve-Factor methodology and are thus stateless,
Deis maintains platform state in the :ref:`Store` component.

The store component runs `Ceph`_, and is used by the :ref:`Database`, :ref:`Registry`,
:ref:`Controller`, and :ref:`Logger` components as a data store. Database and registry
use store-gateway and controller and logger use store-volume. Being backed by the store component
enables these components to move freely around the cluster while their state is backed by store.

The store component is configured to still operate in a degraded state, and will automatically
recover should a host fail and then rejoin the cluster. Total data loss of Ceph is only possible
if all of the store containers are removed. However, backup of Ceph is fairly straightforward, and
is recommended before :ref:`Upgrading Deis <upgrading-deis>`.

Data stored in Ceph is accessible in two places: on the CoreOS filesystem at ``/var/lib/deis/store``
and in the store-gateway component. Backing up this data is straightforward - we can simply tarball
the filesystem data, and use any S3-compatible blob store tool to download all files in the
store-gateway component.

Setup
-----

The ``deis-store-gateway`` component exposes an S3-compatible API, so we can use a tool like `s3cmd`_
to work with the object store. First, `download s3cmd`_ and install it (you'll need at least version
1.5.0 for Ceph support).

We'll need the generated access key and secret key for use with the gateway. We can get these using
``deisctl``, either on one of the cluster machines or on a remote machine with ``DEISCTL_TUNNEL`` set:

.. code-block:: console

    $ deisctl config store get gateway/accessKey
    $ deisctl config store get gateway/secretKey

Back on the local machine, run ``s3cmd --configure`` and enter your access key and secret key.

When prompted with the ``Use HTTPS protocol`` option, answer ``No``. Other settings can be left at
the defaults. If the configure script prompts to test the credentials, skip that step - it will
try to authenticate against Amazon S3 and fail.

You'll need to change two configuration settings - edit ``~/.s3cfg`` and change
``host_base`` and ``host_bucket`` to match ``deis-store.<your domain>``. For example, for my local
Vagrant setup, I've changed the lines to:

.. code-block:: console

    host_base = deis-store.local3.deisapp.com
    host_bucket = deis-store.local3.deisapp.com

We can now use ``s3cmd`` to back up and restore data from the store-gateway.

.. note::

    Some users have reported that the data transferred in this process can overwhelm the gateway
    component, and that scaling up to multiple gateways with ``deisctl scale`` before both the backup
    and restore alleviates this issue.

Backing up
----------

Database backups and registry data
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The store-gateway component stores database backups and is used to store data for the registry.
On our local machine, we can use ``s3cmd sync`` to copy the objects locally:

.. code-block:: console

    $ s3cmd sync s3://db_wal .
    $ s3cmd sync s3://registry .

Log data
~~~~~~~~

The store-volume service mounts a filesystem which is used by the controller and logger components
to store and retrieve application and component logs.

Since this is just a POSIX filesystem, you can simply tarball the contents of this directory
and rsync it to a local machine:

.. code-block:: console

    $ ssh core@<hostname> 'cd /var/lib/deis/store && sudo tar cpzf ~/store_file_backup.tar.gz .'
    tar: /var/lib/deis/store/logs/deis-registry.log: file changed as we read it
    $ rsync -avhe ssh core@<hostname>:~/store_file_backup.tar.gz .

Note that you'll need to specify the SSH port when using Vagrant:

.. code-block:: console

    $ rsync -avhe 'ssh -p 2222' core@127.0.0.1:~/store_file_backup.tar.gz .

Note the warning - in a running cluster the log files are constantly being written to, so we are
preserving a specific moment in time.

Database data
~~~~~~~~~~~~~

While backing up the Ceph data is sufficient (as database ships backups and WAL logs to store),
we can also back up the PostgreSQL data using ``pg_dumpall`` so we have a text dump of the database.

We can identify the machine running database with ``deisctl list``, and from that machine:

.. code-block:: console

    core@deis-1 ~ $ docker exec deis-database sudo -u postgres pg_dumpall > dump_all.sql
    core@deis-1 ~ $ docker cp deis-database:/app/dump_all.sql .

Restoring
---------

.. note::

    Restoring data is only necessary when deploying a new cluster. Most users will use the normal
    in-place upgrade workflow which does not require a restore.

We want to restore the data on a new cluster before the rest of the Deis components come up and
initialize. So, we will install the whole platform, but only start the store components:

.. code-block:: console

    $ deisctl install platform
    $ deisctl start store-monitor
    $ deisctl start store-daemon
    $ deisctl start store-metadata
    $ deisctl start store-gateway@1
    $ deisctl start store-volume

We'll also need to start a router so we can access the gateway:

.. code-block:: console

    $ deisctl start router@1

The default maximum body size on the router is too small to support large uploads to the gateway,
so we need to increase it:

.. code-block:: console

    $ deisctl config router set bodySize=100m

The new cluster will have generated a new access key and secret key, so we'll need to get those again:

.. code-block:: console

    $ deisctl config store get gateway/accessKey
    $ deisctl config store get gateway/secretKey

Edit ``~/.s3cfg`` and update the keys.

Now we can restore the data!

Database backups and registry data
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Because neither the database nor registry have started, the bucket we need to restore to will not
yet exist. So, we'll need to create those buckets:

.. code-block:: console

    $ s3cmd mb s3://db_wal
    $ s3cmd mb s3://registry

Now we can restore the data:

.. code-block:: console

    $ s3cmd sync basebackups_005 s3://db_wal
    $ s3cmd sync wal_005 s3://db_wal
    $ s3cmd sync registry s3://registry

Log data
~~~~~~~~

Once we copy the tarball back to one of the CoreOS machines, we can extract it:

.. code-block:: console

    $ rsync -avhe ssh store_file_backup.tar.gz core@<hostname>:~/store_file_backup.tar.gz
    $ ssh core@<hostname> 'cd /var/lib/deis/store && sudo tar -xzpf ~/store_file_backup.tar.gz --same-owner'

Note that you'll need to specify the SSH port when using Vagrant:

.. code-block:: console

    $ rsync -avhe 'ssh -p 2222' store_file_backup.tar.gz core@127.0.0.1:~/store_file_backup.tar.gz

Finishing up
~~~~~~~~~~~~

Now that the data is restored, the rest of the cluster should come up normally with a ``deisctl start platform``.

The controller will automatically re-write user keys, application data, and domains from the
restored database to etcd.

That's it! The cluster should be fully restored.

Tools
-----

Various community members have developed tools to assist in automating the backup and restore process outlined above.
Information on the tools can be found on the `Community Contributions`_ page.

.. _`Ceph`: http://ceph.com
.. _`download s3cmd`: http://s3tools.org/download
.. _`Community Contributions`: https://github.com/deis/deis/blob/master/contrib/README.md#backup-tools
.. _`s3cmd`: http://s3tools.org/
