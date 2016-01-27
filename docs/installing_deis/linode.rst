:title: Installing Deis on Linode
:description: How to provision a multi-node Deis cluster on Linode

.. _deis_on_linode:

Linode
======

In this tutorial, we will show you how to set up your own 3-node cluster on Linode.

Please :ref:`get the source <get_the_source>` and refer to the scripts in `contrib/linode`_
while following this documentation.

.. important::

    Linode support is untested by the Deis team, so we rely on the community to
    improve this documentation and fix bugs. We greatly appreciate the help!


Prerequisites
-------------

Before we can begin to provision a cluster on Linode, let's get a few things squared away.


Enable KVM Hypervisor
^^^^^^^^^^^^^^^^^^^^^

Navigate to the `Linode Account Settings`_ page and change the Hypervisor Preference to ``KVM``.

Although it is possible to provision CoreOS under Xen on Linode it is much more difficult and
the tools included only work with the KVM Hypervisor.


Obtain Linode API Key
^^^^^^^^^^^^^^^^^^^^^

Next, navigate to the `Linode API Keys`_ page and generate an API key. Take note of the key,
as you will need it later when you provision your cluster.


Install Python and Dependencies
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The scripts used to provision the cluster are written for Python 2.7 and require a few
dependencies be installed with the `pip`_ package manager. If you are on OS X or Linux you
likely have these already available.

Lets install our dependencies:

.. code-block:: console

    $ pip install -r contrib/linode/requirements.txt


Generate SSH Key
^^^^^^^^^^^^^^^^

If you don't already have a SSH key, the following command will generate
a new keypair named "deis":

.. code-block:: console

    $ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis


Check System Requirements
-------------------------

Please refer to :ref:`system-requirements` for resource considerations when choosing a Linode
plan to run Deis. A Deis cluster must have 3 or more nodes. See :ref:`cluster-size` for more details.


Create Cloud Init
-----------------

Create your cloud init file using Deis' ``contrib/linode/create-linode-user-data.py`` script.

First navigate to the ``contrib/linode`` directory:

.. code-block:: console

    $ cd contrib/linode

Then, create the ``linode-user-data.yaml`` file:

.. code-block:: console

    $ ./create-linode-user-data.py --public-key /path/to/key/deis.pub

It is possible to specify multiple authorized keys and/or specify an etcd token to use for the cluster.
See the full command usage below:

.. code-block:: console

    usage: create-linode-user-data.py [-h] --public-key PUBLIC_KEY_FILES
                                  [--etcd-token ETCD_TOKEN]

    Create Linode User Data

    optional arguments:
      -h, --help            show this help message and exit
      --public-key PUBLIC_KEY_FILES
                            Authorized SSH Keys
      --etcd-token ETCD_TOKEN
                            Etcd Token

Provision Cluster
-----------------

The time has finally come to provision our cluster and all it takes is a single command!

.. code-block:: console

    $ ./provision-linode-cluster.py --api-key=YOUR_LINODE_API_KEY provision

This command will create a 3 node cluster of Linode 4096s in Dallas TX, however by passing additional
arguments you can specify the data center, size of nodes, number of nodes, and a bunch more:

.. code-block:: console

    usage: provision-linode-cluster.py provision [-h] [--num NUM_NODES]
                                                 [--name-prefix NODE_NAME_PREFIX]
                                                 [--display-group NODE_DISPLAY_GROUP]
                                                 [--plan NODE_PLAN]
                                                 [--datacenter NODE_DATA_CENTER]
                                                 [--cloud-config CLOUD_CONFIG]
                                                 [--coreos-version COREOS_VERSION]
                                                 [--coreos-channel COREOS_CHANNEL]

    optional arguments:
      -h, --help            show this help message and exit
      --num NUM_NODES       Number of nodes to provision
      --name-prefix NODE_NAME_PREFIX
                            Node name prefix
      --display-group NODE_DISPLAY_GROUP
                            Node display group
      --plan NODE_PLAN      Node plan id. Use list-plans to find the id.
      --datacenter NODE_DATA_CENTER
                            Node data center id. Use list-data-centers to find the
                            id.
      --cloud-config CLOUD_CONFIG
                            CoreOS cloud config user-data file
      --coreos-version COREOS_VERSION
                            CoreOS version number to install
      --coreos-channel COREOS_CHANNEL
                            CoreOS channel to install from


Additionally, the provision tool contains two utilities to list available data centers and plans that
can help find the command argument values.

.. code-block:: console

    $ ./provision-linode-cluster.py --api-key=YOUR_LINODE_API_KEY list-data-centers

.. code-block:: console

    $ ./provision-linode-cluster.py --api-key=YOUR_LINODE_API_KEY list-plans


Apply Security Group Settings
-----------------------------

Because Linode does not have a security group feature, we'll need to add some custom
``iptables`` rules so our components are not accessible to the outside world.


If you are on the Linode private network, run:

.. code-block:: console

    $ ./apply-firewall.py --private-key /path/to/key/deis


If you are outside the private network, you will have to manually specify the public ip address of
each host. To do so, run:

.. code-block:: console

    $ ./apply-firewall.py --private-key /path/to/key/deis --hosts 1.2.3.4 11.22.33.44 111.222.33.44

    
Or, you can provide the display group (NOTE: the default display group is ``deis``) to search for the
nodes using the Linode API, by running:

.. code-block:: console

    $ ./apply-firewall.py --private-key /path/to/key/deis --api-key YOUR_LINODE_API_KEY --display-group YOUR_DISPLAY_GROUP


The script will use either the Linode API or the etcd discovery url to find all of the nodes in your
cluster and create iptables rules to allow connections between nodes while blocking outside connections
automatically. Note that when discovering node ips, the ``--display-group`` parameter has highest priority,
then manual specification via ``--nodes`` and ``--hosts`` (i.e. public and private ips), then the etcd
discovery url via parameter ``--display-url`` or the ``linode-user-data.yaml`` file. Full command usage:

.. code-block:: console

    usage: apply-firewall.py [-h] --private-key PRIVATE_KEY [--private]
                             [--adding-new-nodes]
                             [--discovery-url DISCOVERY_URL]
                             [--display-group DISPLAY_GROUP]
                             [--hosts HOSTS [HOSTS ...]]
                             [--nodes HOSTS [HOSTS ...]]

    Apply a "Security Group" to a Deis cluster

    optional arguments:
      -h, --help            show this help message and exit
      --private-key PRIVATE_KEY
                            Cluster SSH Private Key
      --private             Only allow access to the cluster from the private
                            network
      --adding-new-nodes    When adding new nodes to existing cluster, allows access to etcd
      --display-group DISPLAY_GROUP
                            Linode display group for nodes 
      --discovery-url DISCOVERY_URL
                            Etcd discovery url
      --hosts HOSTS [HOSTS ...]
                            The public IP addresses of the hosts
      --nodes HOSTS [HOSTS ...]
                            The private IP addresses of the hosts 


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.


Adding Nodes to an Existing Cluster
-----------------------------------

When adding one or more nodes to an existing CoreOS setup, ``etcd`` will be `added as a proxy to
the existing cluster`_. The setup of a proxy requires access to ports 2379 and 2380 of the existing
nodes in the cluster.

In order to open up these ports, before cluster provisioning, run:

.. code-block:: console

    $ ./apply-firewall.py --private-key /path/to/key/deis --hosts 1.2.3.4 11.22.33.44 111.222.33.44
                          --adding-new-nodes

    
Then provision the cluster as described above and afterwards reapply the firewall using
``./apply-firewall.py`` without the ``--adding-new-nodes`` parameter.


.. _`added as a proxy to the existing cluster`: https://coreos.com/etcd/docs/latest/clustering.html#public-etcd-discovery-service
.. _`contrib/linode`: https://github.com/deis/deis/tree/master/contrib/linode
.. _`Linode Account Settings`: https://manager.linode.com/account/settings
.. _`Linode API Keys`: https://manager.linode.com/profile/api
.. _`pip`: https://pip.pypa.io/en/stable/

