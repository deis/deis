:title: Upgrading Deis
:description: Guide to upgrading Deis to a new release.


.. _upgrading-deis:

Upgrading Deis
==============

There are currently two strategies for upgrading a Deis cluster:

* In-place Upgrade (recommended)
* Migration Upgrade

In-place Upgrade
----------------

An in-place upgrade swaps out platform containers for newer versions on the same set of hosts,
leaving your applications and platform data intact.  This is the easiest and least disruptive upgrade strategy.
The general approach is to use ``deisctl`` to uninstall all platform components, update the platform version
and then reinstall platform components.

.. note::

    In-place upgrades are supported starting from Deis version 0.14.0

Use the following steps to perform an in-place upgrade of your Deis cluster.

.. code-block:: console

    $ deisctl uninstall platform
    $ deisctl config platform set version=v0.15.0
    $ deisctl install platform
    $ deisctl start platform

.. attention::

    In-place upgrades incur approximately 10-30 minutes of downtime for deployed applications, the router mesh
    and the platform control plane.  Please plan your maintenance windows accordingly.


Migration Upgrade
-----------------

This upgrade method provisions a new cluster running in parallel to the old one. Applications are
migrated to this new cluster one-by-one, and DNS records are updated to cut over traffic on a
per-application basis. This results in a no-downtime controlled upgrade, but has the caveat that no
data from the old cluster (users, releases, etc.) is retained. Future ``deisctl`` tooling will have
facilities to export and import this platform data.

.. note::

    Migration upgrades are useful for moving Deis to a new set of hosts,
    but should otherwise be avoided due to the amount of manual work involved.

.. important::

    In order to migrate applications, your new cluster must have network access
    to the registry component on the old cluster

Enumerate Existing Applications
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Each application will need to be deployed to the new cluster manually.
Log in to the existing cluster as an admin user and use the ``deis`` client to
gather information about your deployed applications.

List all applications with:

.. code-block:: console

    $ deis apps:list

Gather each application's version with:

.. code-block:: console

    $ deis apps:info -a <app-name>

Provision servers
^^^^^^^^^^^^^^^^^
Follow the Deis documentation to provision a new cluster using your desired target release.
Be sure to use a new etcd discovery URL so that the new cluster doesn't interfere with the running one.

Upgrade Deis clients
^^^^^^^^^^^^^^^^^^^^
If changing versions, make sure you upgrade your ``deis`` and ``deisctl`` clients
to match the cluster's release.

Register and login to the new controller
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Register an account on the new controller and login.

.. code-block:: console

    $ deis register http://deis.newcluster.example.org
    $ deis login http://deis.newcluster.example.org

Migrate applications
^^^^^^^^^^^^^^^^^^^^
The ``deis pull`` command makes it easy to migrate existing applications from
one cluster to another.  However, you must have network access to the existing
cluster's registry component.

Migrate a single application with:

.. code-block:: console

    $ deis create <app-name>
    $ deis pull registry.oldcluster.example.org:5000/<app-name>:<version>

This will move the application's Docker image across clusters, ensuring the application
is migrated bit-for-bit with an identical build and configuration.

Now each application is running on the new cluster, but they are still running (and serving traffic)
on the old cluster.  Use ``deis domains:add`` to tell Deis that this application can be accessed
by its old name:

.. code-block:: console

    $ deis domains:add oldappname.oldcluster.example.org

Repeat for each application.

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
