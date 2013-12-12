:title: Provider
:description: A provider is a pluggable connector to a third-party cloud API. Supported providers come pre-installed on the Deis controller.
:keywords: provider, deis

.. _provider:

Provider
========
A provider is a pluggable connector to a third-party cloud API such as `Amazon EC2`_.
Deis's supported providers come pre-installed on the :ref:`Controller`.

If you are using Deis on a public or private cloud
there's a good chance Deis can scale your cloud servers with a single command.
To make this possible, Deis integrates with various cloud provider APIs,
referred to as :ref:`Providers <provider>`.

Storing Provider Credentials
----------------------------
Each Deis user account contains a set of records for each supported :ref:`provider`.
These records hold secure keys and credentials used to make cloud API calls.
Credentials are private to each user account.

Static Provider
---------------
If you are using Deis on bare metal or with an unsupported cloud provider, you
won't be able to scale nodes automatically.
The "static" provider allows you to manually add and remove servers
without using any automated provisioning.

Discovering Providers
---------------------
The Deis client will look at the environment variables on your workstation to discover
credentials for supported cloud providers. 
The table below shows how the environment variables on your workstation map to
provider types and fields stored in your Deis user account.

======================= =============== ==============
Variable Name           Provider Type   Provider Field
======================= =============== ==============
AWS_ACCESS_KEY_ID       ec2             access_key
AWS_SECRET_ACCESS_KEY   ec2             secret_key
AWS_ACCESS_KEY          ec2             access_key
AWS_SECRET_KEY          ec2             secret_key
RACKSPACE_USERNAME      rackspace       username
RACKSPACE_API_KEY       rackspace       api_key
DIGITALOCEAN_CLIENT_ID  digitalocean    client_id
DIGITALOCEAN_API_KEY    digitalocean    api_key
======================= =============== ==============

----

To discover providers using the Deis client:

.. code-block:: console

    $ deis providers:discover
    Discovered EC2 credentials: AAAAAAAAAAAAAAAAAAAA
    Import EC2 credentials? (y/n) : y
    Uploading EC2 credentials... done
    No Rackspace credentials discovered.
    No DigitalOcean credentials discovered.
    No Vagrant VMs discovered.

Advanced Provider Management
----------------------------

.. code-block:: console

    $ deis help providers
    Valid commands for providers:

    providers:list        list available providers for the logged in user
    providers:discover    discover provider credentials using envvars
    providers:create      create a new provider for use by deis
    providers:info        print information about a specific provider

    Use `deis help [command]` to learn more

The client allows you to list providers, get detailed information about a provider, or
create a new provider altogether (useful for using a different set of credentials).

Building Custom Providers
-------------------------
Building a custom provider is simple.  It must publish 5 methods:

* build_layer - to create any shared infrastructure needed by the layer's nodes
* destroy_layer - to destroy any shared infrastructure
* build_node - to provision a node and prepare it for bootstrapping by Chef
* destroy_node - to destroy a node after it has been purged from the Chef Server
* seed_flavors - to seed the controller database with default flavors

Provider developers can review the `EC2 Reference Implementation`_.

.. _`Amazon EC2`: http://aws.amazon.com/ec2/
.. _`EC2 Reference Implementation`: https://github.com/opdemand/deis/blob/master/provider/ec2.py
