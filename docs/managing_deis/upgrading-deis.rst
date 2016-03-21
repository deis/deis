:title: Upgrading Deis
:description: Guide to upgrading Deis to a new release.


.. _upgrading-deis:

Upgrading Deis
==============

There are currently two strategies for upgrading a Deis cluster:

* In-place Upgrade (recommended)
* Migration Upgrade

Before attempting an upgrade, it is strongly recommended to :ref:`backup your data <backing_up_data>`.

In-place Upgrade
----------------

An in-place upgrade swaps out platform containers for newer versions on the same set of hosts,
leaving your applications and platform data intact.  This is the easiest and least disruptive upgrade strategy.
The general approach is to use ``deisctl`` to uninstall all platform components, update the platform version
and then reinstall platform components.

.. important::

    Always use a version of ``deisctl`` that matches the Deis release.
    Verify this with ``deisctl --version``.

Use the following steps to perform an in-place upgrade of your Deis cluster.

First, use the current ``deisctl`` to stop and uninstall the Deis platform.

.. code-block:: console

    $ deisctl --version  # should match the installed platform
    1.0.2
    $ deisctl stop platform && deisctl uninstall platform

Finally, update ``deisctl`` to the new version and reinstall:

.. code-block:: console

    $ curl -sSL http://deis.io/deisctl/install.sh | sh -s 1.13.0
    $ deisctl --version  # should match the desired platform
    1.13.0
    $ deisctl config platform set version=v1.13.0
    $ deisctl install platform
    $ deisctl start platform

.. attention::

    In-place upgrades incur approximately 10-30 minutes of downtime for deployed applications, the router mesh
    and the platform control plane.  Please plan your maintenance windows accordingly.

.. note::

    When upgrading an AWS cluster older than Deis v1.6, a :ref:`migration_upgrade` is
    preferable.

    On AWS, Deis enables the :ref:`PROXY protocol <proxy_protocol>` by default.
    If an in-place upgrade is required, run ``deisctl config router set proxyProtocol=1``,
    enable PROXY protocol for ports 80 and 443 on the ELB, add a ``TCP 443:443`` listener, and
    change existing targets and health checks from HTTP to TCP.

Upgrade Deis clients
^^^^^^^^^^^^^^^^^^^^
As well as upgrading ``deisctl``, make sure to upgrade the :ref:`deis client <install-client>` to
match the new version of Deis.

Graceful Upgrade
----------------

Alternatively, an experimental feature exists to provide the ability to perform a graceful upgrade. This process is
available for version 1.9.0 moving forward and is intended to facilitate upgrades within a major version (for example,
from 1.9.0 to 1.9.1 or 1.11.2). Upgrading between major versions is not supported (for example, from 1.9.0 to a
future 2.0.0). Unlike the in-place process above, this process keeps the platform's routers and publishers up during
the upgrade process. This means that there should only be a maximum of around 1-2 seconds of downtime while the
routers boot up. Many times, there will be no downtime at all.

.. note::

    Your loadbalancer configuration is the determining factor for how much downtime will occur during a successful upgrade.
    If your loadbalancer is configured to quickly reactivate failed hosts to its pool of active hosts, its quite possible to
    achieve zero downtime upgrades. If your loadbalancer is configured to be more pessimistic, such as requiring multiple
    successful healthchecks before reactivating a node, then the chance for downtime increases. You should review your
    loadbalancers configuration to determine what to expect during the upgrade process.

The process involves two ``deisctl`` subcommands, ``upgrade-prep`` and ``upgrade-takeover``, in coordination with a few other important commands.

.. note::

    If you are using Deis in :ref:`stateless mode <running-deis-without-ceph>`, you should add the option `--stateless`
    to `upgrade-prep` and `upgrade-takeover` subcommands to start only the necessary components.

First, a new ``deisctl`` version should be installed to a temporary location, reflecting the desired version to upgrade
to. Care should be taken not to overwrite the existing ``deisctl`` version.

.. code-block:: console

    $ mkdir /tmp/upgrade
    $ curl -sSL http://deis.io/deisctl/install.sh | sh -s 1.13.0 /tmp/upgrade
    $ /tmp/upgrade/deisctl --version  # should match the desired platform
    1.13.0
    $ /tmp/upgrade/deisctl refresh-units
    $ /tmp/upgrade/deisctl config platform set version=v1.13.0

Now it is possible to prepare the cluster for the upgrade using the old ``deisctl`` binary. This command will shutdown
and uninstall all components of the cluster except the router and publisher. This means your services should still be
serving traffic afterwards, but nothing else in the cluster will be functional.

.. code-block:: console

    $ /opt/bin/deisctl upgrade-prep

Finally, the rest of the components are brought up by the new binary. First, a rolling restart is done on the routers,
replacing them one by one. Then the rest of the components are brought up. The end result should be an upgraded cluster.

.. code-block:: console

    $ /tmp/upgrade/deisctl upgrade-takeover

It is recommended to move the newer ``deisctl`` into ``/opt/bin`` once the procedure is complete.

If the process were to fail, the old version can be restored manually by reinstalling and starting the old components.

.. code-block:: console

    $ /tmp/upgrade/deisctl stop platform
    $ /tmp/upgrade/deisctl uninstall platform
    $ /tmp/upgrade/deisctl config platform set version=v1.13.0
    $ /opt/bin/deisctl refresh-units
    $ /opt/bin/deisctl install platform
    $ /opt/bin/deisctl start platform

Upgrade Deis clients
^^^^^^^^^^^^^^^^^^^^
As well as upgrading ``deisctl``, make sure to upgrade the :ref:`deis client <install-client>` to
match the new version of Deis.


.. _migration_upgrade:

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


.. _upgrading-coreos:

Upgrading CoreOS
----------------

By default, Deis disables CoreOS automatic updates. This is partially because in the case of a
machine reboot, Deis components will be scheduled to a new host and will need a few minutes to start
and restore to a running state. This results in a short downtime of the Deis control plane,
which can be disruptive if unplanned.

Additionally, because Deis customizes the CoreOS cloud-config file, upgrading the CoreOS host to
a new version without accounting for changes in the cloud-config file could cause Deis to stop
functioning properly.

.. important::

  Enabling updates for CoreOS will result in the machine upgrading to the latest CoreOS release
  available in a particular channel. Sometimes, new CoreOS releases make changes that will break
  Deis. It is always recommended to provision a Deis release with the CoreOS version specified
  in that release's provision scripts or documentation.

.. important::

  Upgrading a cluster can result in simultaneously running different etcd versions,
  which may introduce incompatibilities that result in a broken etcd cluster. It is
  always recommended to first test upgrades in a non-production cluster whenever possible.

While typically not recommended, it is possible to trigger an update of a CoreOS machine. Some
Deis releases may recommend a CoreOS upgrade - in these cases, the release notes for a Deis release
will point to this documentation.

Checking the CoreOS version
^^^^^^^^^^^^^^^^^^^^^^^^^^^

You can check the CoreOS version by running the following command on the CoreOS machine:

.. code-block:: console

    $ cat /etc/os-release

Or from your local machine:

.. code-block:: console

    $ ssh core@<server ip> 'cat /etc/os-release'


Triggering an upgrade
^^^^^^^^^^^^^^^^^^^^^

To upgrade CoreOS, run the following commands:

.. code-block:: console

    $ ssh core@<server ip>
    $ sudo su
    $ echo GROUP=stable > /etc/coreos/update.conf
    $ systemctl unmask update-engine.service
    $ systemctl start update-engine.service
    $ update_engine_client -update
    $ systemctl stop update-engine.service
    $ systemctl mask update-engine.service
    $ reboot

.. warning::

  You should only upgrade one host at a time. Removing multiple hosts from the cluster
  simultaneously can result in failure of the etcd cluster. Ensure the recently-rebooted host
  has returned to the cluster with ``fleetctl list-machines`` before moving on to the next host.

After the host reboots, ``update-engine.service`` should be unmasked and started once again:

.. code-block:: console

    $ systemctl unmask update-engine.service
    $ systemctl start update-engine.service

It may take a few minutes for CoreOS to recognize that the update has been applied successfully, and
only then will it update the boot flags to use the new image on subsequent reboots. This can be confirmed
by watching the ``update-engine`` journal:

.. code-block:: console

    $ journalctl -fu update-engine

Seeing a message like ``Updating boot flags...`` means that the update has finished, and the service
should be stopped and masked once again:

.. code-block:: console

    $ systemctl stop update-engine.service
    $ systemctl mask update-engine.service

The update is now complete.

.. note::

    Users have reported that some cloud providers do not allow the boot partition to be updated,
    resulting in CoreOS reverting to the originally installed version on a reboot.
