:title: Create a cluster
:description: Create your default Deis cluster.

.. _create_cluster:

Create a Cluster
======================

Applications on Deis are deployed to a :ref:`cluster`. Before you can deploy
applications, you need to create a cluster.

The parameters to `deis clusters:create` are:

* cluster name - the name used by Deis to reference the cluster
* cluster hostname - the hostname for the cluster -- applications are accessible under this domain
* cluster members (`--hosts`) - a comma-separated list of IP addresses of cluster members -- not necessarily all members, but at least one
* auth SSH key (`--auth`) - the SSH private key used to provision servers (for EC2 and Rackspace, this key is likely `~/.ssh/deis`)

For example, to create a cluster on a local Deis installation:

.. code-block:: console

    $ deis clusters:create dev local.deisapp.com --hosts=local.deisapp.com --auth=~/.vagrant.d/insecure_private_key
