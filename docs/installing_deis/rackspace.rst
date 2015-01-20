:title: Installing Deis on Rackspace
:description: How to provision a multi-node Deis cluster on Rackspace

.. _deis_on_rackspace:

Rackspace
=========

We'll mostly be following the `CoreOS on Rackspace`_ guide. You'll need to have a sane python
environment with ``pip`` already installed (``sudo easy_install pip``).

Please :ref:`get the source <get_the_source>` and refer to the scripts in `contrib/rackspace`_
while following this documentation.


Install supernova
-----------------

.. code-block:: console

    $ sudo pip install keyring
    $ sudo pip install rackspace-novaclient
    $ sudo pip install supernova


Configure supernova
-------------------

Edit ``~/.supernova`` to match the following:

.. code-block:: console

    [production]
    OS_AUTH_URL = https://identity.api.rackspacecloud.com/v2.0/
    OS_USERNAME = {rackspace_username}
    OS_PASSWORD = {rackspace_api_key}
    OS_TENANT_NAME = {rackspace_account_id}
    OS_REGION_NAME = DFW (or ORD or another region)
    OS_AUTH_SYSTEM = rackspace

Your account ID is displayed in the upper right-hand corner of the cloud control panel UI, and your
API key can be found on the Account Settings page.


Generate SSH Key
----------------

.. include:: ../_includes/_generate-ssh-key.rst


Set up your keys
----------------

Tell supernova about the key pair and give it an identifiable name:

.. code-block:: console

    $ supernova production keypair-add --pub-key ~/.ssh/deis.pub deis-key


Generate a New Discovery URL
----------------------------

.. include:: ../_includes/_generate-discovery-url.rst


Choose number of instances
--------------------------

By default, the provision script will provision 3 servers. You can override this by setting `DEIS_NUM_INSTANCES`:

.. code-block:: console

    $ DEIS_NUM_INSTANCES=5 ./provision-rackspace-cluster.sh deis-key

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always have an odd number of members.
For more information, see `etcd disaster recovery`_.

Deis clusters of less than 3 nodes are unsupported.

Run the Provision Script
------------------------

Run the Rackspace provision script to spawn a new CoreOS cluster. You'll need to provide the name
of the key pair you just added. Optionally, you can also specify a flavor name, as well as the supernova environment to use. The environment will default to `production`, and the flavor to `performance1-2`.

.. code-block:: console

    $ cd contrib/rackspace
    $ ./provision-rackspace-cluster.sh
    Usage: provision-rackspace-cluster.sh <key pair name> [flavor] [environment]
    $ ./provision-rackspace-cluster.sh deis-key

.. note::

    The provision script will by default provision ``performance1-2`` servers. Choosing a smaller
    server size is not recommended. Please refer to :ref:`system-requirements` for resource
    considerations when choosing a server size to run Deis.

Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Configure Load Balancer
-----------------------

You'll need to create two load balancers on Rackspace to handle your cluster:

.. code-block:: console

    Load Balancer 1
    Port 80
    Protocol HTTP
    Health Monitoring -
      Monitor Type HTTP
      HTTP Path /health-check

    Load Balancer 2
    Virtual IP Shared VIP on Another Load Balancer (select Load Balancer 1)
    Port 2222
    Protocol TCP

Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.


.. _`contrib/rackspace`: https://github.com/deis/deis/tree/master/contrib/rackspace
.. _`CoreOS on Rackspace`: https://coreos.com/docs/running-coreos/cloud-providers/rackspace/
.. _`etcd disaster recovery`: https://github.com/coreos/etcd/blob/master/Documentation/admin_guide.md#disaster-recovery
.. _Rackspace: https://github.com/deis/deis/tree/master/contrib/rackspace#readme
.. _`contrib/rackspace`: https://github.com/deis/deis/tree/master/contrib/rackspace
