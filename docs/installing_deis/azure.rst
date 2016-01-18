:title: Installing Deis on Microsoft Azure
:description: How to provision a multi-node Deis cluster on Microsoft Azure

.. _deis_on_azure:

Microsoft Azure
===============

This section will show you how to create a 3-node Deis cluster on Microsoft Azure.

Before you start, :ref:`get the Deis source <get_the_source>` and change directory into `contrib/azure`_
while following this documentation.


Install Tools
-------------

The cluster creation tool uses a Python script to generate a configuration file.
This script uses PyYAML, a Python library, to do its work.

If you haven't already, install these on your development machine:

For OSX users:

.. code-block:: console

    $ brew install python
    $ sudo pip install pyyaml

For Ubuntu users:

.. code-block:: console

    $ sudo apt-get install -y python-yaml

Additionally, we'll also need to install the `Azure CLI`_ from Microsoft.

Create CoreOS Cluster
---------------------

First, login to the Azure CLI:

.. code-block:: console

    $ azure login

.. note::

    Deis makes use of `Azure Resource Manager`_ to submit a template
    describing the infrastructure that we'd like to create. You'll need an
    `organizational account`_ (not a typical Microsoft or Live account) in order to
    use this template.

Instruct the client to switch to ARM mode:

.. code-block:: console

    $ azure config mode arm

Switch to the ``contrib/azure`` directory:

.. code-block:: console

    $ cd contrib/azure

Generate a new discovery URL for the deployment so the hosts can find each other:

.. code-block:: console

    $ ./create-azure-user-data $(curl -s https://discovery.etcd.io/new)

Next, edit ``parameters.json`` to configure the parameters required for the
cluster. For ``publicDomainName``, specify the prefix of domain name (like ``deisNode``).
For ``sshKeyData``, use the public key material for the SSH key you'd like
to use to log into the hosts. For ``customData``, you'll need to supply the
base64-encoded version of ``azure-user-data``. This can be generated using ``base64``:

.. code-block:: console

    $ base64 azure-user-data

Paste the result into ``parameters.json``. Any of the values in ``parameters.json``
are defaults and can be customized if desired.

.. note::

  For best performance, Deis clusters on Azure default to using `premium storage`_.
  This incurs an additional cost. Using standard storage is possible, but is unsupported
  as it resulted in cluster issues during testing. Premium storage is only available
  in `some regions`_.

Finally, we can deploy. Choose a valid location to deploy -- you can list all locations
with ``azure location list``.

As an example, to create a deployment named "deis" in the "West US" region:

.. code-block:: console

    $ azure group create --name deis --location "West US" --deployment-name deis --template-file arm-template.json --parameters-file parameters.json

Each instance will have a public IP address which can be used to log in via SSH
or as a tunnel endpoint for ``deisctl``. You can get these IPs from the Azure Portal
or via the CLI with ``azure vm show``:

.. code-block:: console

    $ azure vm show deisNode0 --resource-group deis | grep 'Public IP address'

Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

.. _`Azure CLI`: https://azure.microsoft.com/en-us/documentation/articles/xplat-cli-install/
.. _`Azure Resource Manager`: https://azure.microsoft.com/en-us/documentation/articles/resource-manager-deployment-model/
.. _`contrib/azure`: https://github.com/deis/deis/tree/master/contrib/azure
.. _`organizational account`: http://www.brucebnews.com/2013/04/the-difference-between-a-microsoft-account-and-an-office-365-account/
.. _`premium storage`: https://azure.microsoft.com/en-us/services/storage/premium-storage/
.. _`some regions`: https://azure.microsoft.com/en-us/regions/#services
