:title: Upgrading Deis
:description: Guide to upgrading Deis to a new release.


.. _upgrading-deis:

Upgrading Deis
==============

This guide provides some general information and considerations around upgrading a Deis cluster.
Additional tooling around upgrading Deis is planned for a future Deis release
(tracked in `#710`_). The upgrade strategies outlined below
reflect the current state of the Deis platform - they are certainly not ideal, and significant work
is planned to make upgrading Deis much less painful for future releases.

The current recommended upgrade paths for Deis are either to provision a new cluster and cut over
DNS to point to new application endpoints, or to perform an in-place upgrade of Deis components.

Deploying a New Cluster
-----------------------

This upgrade method provisions a new cluster running in parallel to the old one. Applications are
pushed to this new cluster one-by-one, and DNS records are updated to cut over traffic on a
per-application basis. This results in a no-downtime controlled upgrade, but has the caveat that no
data from the old cluster (users, releases, etc.) is retained. Future upgrade tooling will have
facilities to export and import cluster data.

Provision servers
^^^^^^^^^^^^^^^^^
A new cluster can be provisioned by following the release's README or contrib/ directory tooling.
Be sure to use a new etcd discovery URL so that the new cluster doesn't interfere with the running one.

Upgrade Deis client and fleetctl
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
The latest Deis client should be installed. You can either use pip or download a pre-compiled binary.
See the release notes or the release's README for a link to the latest client.

Also, new Deis releases frequently require upgrades to the fleetctl client. Again, see the release
notes for the new release to see if an upgrade is necessary.

Before upgrading the client, we should logout with the old client:

.. code-block:: console

    dev $ deis logout

Register and login to the new controller
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Register an account on the new controller and login.

.. code-block:: console

    dev $ deis register http://deis.newcluster.example.org
    dev $ deis login http://deis.newcluster.example.org
    dev $ deis keys:add

Push apps to the new cluster
^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Each application will need to be deployed to the new cluster. You can use ``deis apps:list`` to help
enumerate all existing applications.

For each existing application, rename the existing deis remote, and create a new application.

.. code-block:: console

    dev $ git remote rename deis deis-old
    dev $ deis create
    dev $ git push deis master

Note that you'll also need to ``deis config:set`` for each environment variable an application
references - use ``deis config:list`` to enumerate these.

Now each application is running on the new cluster, but they are still running (and serving traffic)
on the old cluster.

We need to tell Deis that this application can be accessed by its old name:

.. code-block:: console

    dev $ deis domains:add oldappname.oldcluster.example.org

Test applications
^^^^^^^^^^^^^^^^^
Test to make sure applications work as expected on the new Deis cluster.

Update DNS records
^^^^^^^^^^^^^^^^^^
For each application, create CNAME records to point the old application names to the new. Note that
once these records propagate, the new cluster is serving live traffic. You can perform cutover on a
per-application basis and slowly retire the old cluster.

If an application is named 'happy-bandit' on the old Deis cluster and 'jumping-cuddlefish' on the
new cluster, you would create a DNS record that looks like the following:

.. code-block:: console

    happy-bandit.oldcluster.example.org.        CNAME       jumping-cuddlefish.newcluster.example.org

Retire the old cluster
^^^^^^^^^^^^^^^^^^^^^^
Once all applications have been validated, the old cluster can be retired.

In-place upgrade
----------------

This upgrade method involves shutting down the existing cluster, upgrading Deis components, and then
restarting them. It has the benefit of preserving existing data, but has the caveat that it results
in cluster downtime.

Stop components
^^^^^^^^^^^^^^^
The Makefile has the ability to stop all Deis components. Once this occurs, all apps will
be unresponsive.

.. code-block:: console

    dev $ make stop

You can ``make status`` to ensure that all services are stopped.

Export data containers
^^^^^^^^^^^^^^^^^^^^^^
Four Deis components, builder, database, logger, and registry, run separate containers to store
their stateful data. We export these as tarballs before upgrading the containers.

.. code-block:: console

    dev $ fleetctl ssh deis-builder.service
    coreos $ sudo docker export deis-builder-data > /home/coreos/deis-builder-data-backup.tar
    dev $ fleetctl ssh deis-database.service
    coreos $ sudo docker export deis-database-data > /home/coreos/deis-database-data-backup.tar
    dev $ fleetctl ssh deis-logger.service
    coreos $ sudo docker export deis-logger-data > /home/coreos/deis-logger-data-backup.tar
    dev $ fleetctl ssh deis-registry.service
    coreos $ sudo docker export deis-registry-data > /home/coreos/deis-registry-data-backup.tar

Upgrade Deis client and fleetctl
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
The latest Deis client should be installed. You can either use pip or download a pre-compiled binary.
See the release notes or the release's README for a link to the latest client.

Also, new Deis releases frequently require upgrades to the fleetctl client. Again, see the
release notes for the new release to see if an upgrade is necessary.

Upgrade components
^^^^^^^^^^^^^^^^^^
Now we instruct all servers to download the release's Docker containers. This will take some time...

.. code-block:: console

    dev $ make pull

Import data containers
^^^^^^^^^^^^^^^^^^^^^^
We need to reimport the saved data containers we exported earlier.

.. code-block:: console

    dev $ fleetctl ssh deis-builder.service
    coreos $ cat /home/coreos/deis-builder-data-backup.tar | sudo docker import - deis-builder-data
    dev $ fleetctl ssh deis-database.service
    coreos $ cat /home/coreos/deis-database-data-backup.tar | sudo docker import - deis-database-data
    dev $ fleetctl ssh deis-logger.service
    coreos $ cat /home/coreos/deis-logger-data-backup.tar | sudo docker import - deis-logger-data
    dev $ fleetctl ssh deis-registry.service
    coreos $ cat /home/coreos/deis-registry-data-backup.tar | sudo docker import - deis-registry-data

Start the cluster
^^^^^^^^^^^^^^^^^
Again, the Makefile takes care of this logic for us:

.. code-block:: console

    dev $ make run

Test
^^^^
Ensure all applications function as expected.

.. _`#710`: https://github.com/deis/deis/issues/710
