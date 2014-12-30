:title: Installing Deis on Microsoft Azure
:description: How to provision a multi-node Deis cluster on Microsoft Azure

.. _deis_on_azure:

Microsoft Azure
===============

This section will show you how to create a 3-node Deis cluster on Microsoft Azure.

Before you start, :ref:`get the Deis source <get_the_source>` and change directory into `contrib/azure`_ while following this documentation.


Install Python and Azure SDK for Python
---------------------------------------

The cluster creation tool uses Python and the Python Azure library to create a CoreOS cluster. If you haven't already, install these on your development machine:

.. code-block:: console

    $ brew install python
    ...

    $ sudo pip install azure
    ...


And check to make sure they are configured correctly:

.. code-block:: console

    $ python -c "import azure; print(azure.__version__)"
    0.9.0  <-- everything working ok!

Generate Certificates
---------------------

The azure-coreos-cluster creation tool uses the Azure management REST API to create the CoreOS cluster which uses a management certificate to authenticate.

If you don't have a management certificate already configured, the script generate-mgmt-cert.sh can create this certificate for you. Otherwise, you can skip to the next section.

If you need to create a certificate, edit cert.conf in contrib/azure with your company's details and then run:

.. code-block:: console

    $ ./generate-mgmt-cert.sh

Upload Management Cert
----------------------

If you haven't uploaded your management certificate to Azure (azure-cert.cer if you used the script in the previous section), do that now using the `management certificates tab`_ of the Azure portal's Settings.

Also copy the Azure subscription id from this table and save it for the cluster creation script below.

Create Cluster Cloud Config
---------------------------

Before we can create a cluster, we need to create a cloud config for it. The script create-azure-user-data does this for you. This script takes the stock cluster instance config in ../coreos/user-data.example and customizes it for Azure and inserts a unique cluster discovery url:

.. code-block:: console

    $ ./create-azure-user-data $(curl -s https://discovery.etcd.io/new)

This will create a azure-user-data cloud config file. We'll use this with the script in the next section during cluster creation.

Create CoreOS Cluster
---------------------

With the management certificate and cloud config in place, we are ready to create our cluster.

* Create a container called 'vhds' within a storage account in the same region as your cluster using the Azure portal. Note the URL of the container for the cluster creation script below.
* Choose a cloud service name for your Deis cluster for the script below. The script will automatically create this cloud service for you.

With that, let's run the azure-coreos-cluster script which will create the CoreOS cluster. Fill in the bracketed values with the values for your deployment you created above.

.. code-block:: console

    $ ./azure-coreos-cluster [cloud service name]
         --subscription [subscription id]
         --azure-cert azure-cert.pem 
         --num-nodes 3
         --location "West US"    
         --vm-size Large  
         --pip
         --deis
         --blob-container-url https://[blob container].blob.core.windows.net/vhds/
         --data-disk
         --custom-data azure-user-data

This script will by default provision a 3 node cluster but you can increase this with the --num-nodes parameter. Likewise, you can increase the vm size using the --vm-size. It is not recommended that you use smaller than Large (A3) sized instances.

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always have an odd number of members. For more information, see `optimal etcd cluster size`_.


Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

IMPORTANT NOTE: Once you have installed deisctl, you will need to use a customized deis-builder component for Azure since Azure uses routable IP addresses for each instance. Configure this using the following command before you run 'deisctl install platform': 

.. code-block:: console

    $ deisctl config builder set image=deis/builder:v1.1.1-azure

.. _`management certificates tab`: https://manage.windowsazure.com/#Workspaces/AdminTasks/ListManagementCertificates
.. _`contrib/azure`: https://github.com/deis/deis/tree/master/contrib/azure
.. _`etcd`: https://github.com/coreos/etcd
.. _`optimal etcd cluster size`: https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md
