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

Deploying a new Cluster with External Components
------------------------------------------------

If you're upgrading from a cluster where you have outsourced your components outside of
Deis (such as migrating deis-database onto Amazon Relational Database Services), you have
the benefit of preserving existing data, but you still need to update DNS records and the
like.

Provision Servers
^^^^^^^^^^^^^^^^^

Provision the CoreOS cluster as you normally would with any release of Deis. However, do
not install any components onto this cluster. We need to point etcd to the components
which are running outside of the cluster.

Export Etcd Keys
^^^^^^^^^^^^^^^^

To migrate over, start by pointing the new cluster at the old cluster's endpoints:

.. code-block:: console

    $ deisctl config database set host pqsl.example.org
    $ deisctl config database set port 1234
    ...

Next, you'll also want to migrate over the application directories:

    $ etcdctl mkdir /deis/services/appname

Start new Components
^^^^^^^^^^^^^^^^^^^^

The Makefile takes care of this logic for us:

.. code-block:: console

    dev $ make run

Re-deploy Apps to the new Cluster
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

With this process, re-deploying apps couldn't be easier. Just scale the processes down to
0 for each application, then scale back up.

.. code-block:: console

    $ deis scale --app example web=0
    $ deis scale --app example web=3

.. note::

    Support for ``deis ps:restart`` is being tracked in `#467`_.

Test applications
^^^^^^^^^^^^^^^^^

Test to make sure applications work as expected on the new Deis cluster.

Update DNS records
^^^^^^^^^^^^^^^^^^

Once you've finished migrating over to the new cluster, just update your wildcard DNS to
point at your new load balancer. The application names are all the same, so no CNAME
modification needs to occur.

.. _`#710`: https://github.com/deis/deis/issues/710
.. _`#467`: https://github.com/deis/deis/issues/467
