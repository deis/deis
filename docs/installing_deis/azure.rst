:title: Installing Deis on Microsoft Azure
:description: How to provision a multi-node Deis cluster on Microsoft Azure

.. _deis_on_azure:

Microsoft Azure
===============

This section will show you how to create a 3-node Deis cluster on Microsoft Azure.

Before you start, :ref:`get the Deis source <get_the_source>` and change directory into `contrib/azure`_
while following this documentation.


Install Python and Azure SDK for Python
---------------------------------------

The cluster creation tool uses Python and the Python Azure library to create a CoreOS cluster.
If you haven't already, install these on your development machine:

.. code-block:: console

    $ brew install python
    $ sudo pip install azure pyyaml

Generate Certificates
---------------------

The azure-coreos-cluster creation tool uses the Azure management REST API to create the CoreOS
cluster which uses a management certificate to authenticate.

If you don't have a management certificate already configured, the script generate-mgmt-cert.sh can
create this certificate for you. Otherwise, you can skip to the next section.

If you need to create a certificate, edit cert.conf in contrib/azure with your company's details and then run:

.. code-block:: console

    $ ./generate-mgmt-cert.sh

Upload Management Cert
----------------------

If you haven't uploaded your management certificate to Azure (azure-cert.cer if you used the script
in the previous section), do that now using the `management certificates tab`_ of the
Azure portal's Settings.

Also copy the Azure subscription id from this table and save it for the cluster creation script below.


Create CoreOS Cluster
---------------------

With the management certificate and cloud config in place, we are ready to create our cluster.

* Create a container called ``vhds`` within a storage account in the same region as your cluster using the Azure portal. Note the URL of the container for the cluster creation script below.
* Choose a cloud service name for your Deis cluster for the script below. The script will automatically create this cloud service for you.
* Create an `affinity group`_ if you already don't have one. Supply it in quotes with the ``--affinity-group`` parameter. Although *using an affinity group is not mandatory*, it is **highly recommended** since it tells the Azure fabric to place all VMs in the cluster physically close to each other, reducing inter-node latency by a great deal. If you don't want ot use affinity groups, specify a `region`_ for Azure to use with a ``--location`` parameter. The default is ``"West US"``. If you specify both parameters, ``location`` will be ignored. Please note that the script *will not* create an affinity group by itself; it expects the affinity group exists.

This script calls the ``./create-azure-user-data`` script which takes the stock cluster instance config in ``../coreos/user-data.example``, customizes it for Azure, and inserts a unique cluster discovery
endpoint. It will then use the newly created CoreOS config on the newly provisioned cluster.

With that, let's run the azure-coreos-cluster script which will create the CoreOS cluster. Fill in the bracketed values with the values for your deployment you created above.

.. code-block:: console

    $ ./azure-coreos-cluster [cloud service name]
         --subscription [subscription id]
         --azure-cert azure-cert.pem
         --num-nodes 3
         --affinity-group [affinity group name]
         --vm-size Large
         --pip
         --deis
         --blob-container-url https://[blob container].blob.core.windows.net/vhds/
         --data-disk
         --custom-data azure-user-data

This script will by default provision a 3 node cluster but you can increase this with the
``--num-nodes`` parameter. Likewise, you can increase the VM size using ``--vm-size``.
It is not recommended that you use smaller than Large (A3) sized instances.

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always
have an odd number of members. For more information, see `etcd disaster recovery`_.


Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

.. _`management certificates tab`: https://manage.windowsazure.com/#Workspaces/AdminTasks/ListManagementCertificates
.. _`contrib/azure`: https://github.com/deis/deis/tree/master/contrib/azure
.. _`etcd`: https://github.com/coreos/etcd
.. _`etcd disaster recovery`: https://github.com/coreos/etcd/blob/master/Documentation/admin_guide.md#disaster-recovery
.. _`region`: http://azure.microsoft.com/en-us/regions/
.. _`affinity group`: https://msdn.microsoft.com/en-gb/library/azure/jj156085.aspx
