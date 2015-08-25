:title: Installing Deis on DigitalOcean
:description: How to provision a multi-node Deis cluster on DigitalOcean

.. _deis_on_digitalocean:

DigitalOcean
============

In this tutorial, we will show you how to set up your own 3-node cluster on DigitalOcean.

Please :ref:`get the source <get_the_source>` and refer to the scripts in `contrib/digitalocean`_
while following this documentation.


Prerequisites
-------------

To complete this guide, you must have the following:

 - A domain to point to the cluster
 - The ability to provision at least 3 DigitalOcean Droplets that are 4GB or greater

Additionally, we'll need to install `Terraform`_ to do the heavy lifting for us.


Check System Requirements
-------------------------

Please refer to :ref:`system-requirements` for resource considerations when choosing a droplet
size to run Deis.


Generate SSH Key
----------------

.. include:: ../_includes/_generate-ssh-key.rst

Upload this key to DigitalOcean so we can use it for the rest of the provisioning
process.

Generate a New Discovery URL
----------------------------

.. include:: ../_includes/_generate-discovery-url.rst


Create CoreOS Droplets
----------------------

The only other pieces of information we'll need are your DigitalOcean API token
and the fingerprint of your SSH key, both of which can be obtained from the
DigitalOcean interface.

From the source code root directory, invoke Terraform:

.. code-block:: console

    $ terraform apply -var 'token=a1b2c3d3e4f5' \
                      -var 'ssh_keys=c1:d3:a2:b4:e4:f5' \
                      -var 'region=nyc3' \
                      -var 'prefix=deis' \
                      -var 'instances=3' \
                      -var 'size=8GB' \
                      contrib/digitalocean


Note that only ``token`` and ``ssh_keys`` are required - if unset, the other variables
will default to 3 hosts in the ``sfo1`` region with a size of 8GB and a prefix
of ``deis``. Additionally, ``ssh_keys`` can be just one key, or a comma-separated
list of keys to be added to the hosts for the ``core`` user.

The ``region`` option must specify a region with private networking.

Configure DNS
-------------

.. note::

    If you're using your own third-party DNS registrar, please refer to their documentation on this
    setup, along with the :ref:`dns_records` required.

.. note::

    If you don't have an available domain for testing, you can refer to the :ref:`xip_io`
    documentation on setting up a wildcard DNS for Deis.

Deis requires a wildcard DNS record to function properly. If the top-level domain (TLD) that you
are using is ``example.com``, your applications will exist at the ``*.example.com`` level. For example, an
application called ``app`` would be accessible via ``app.example.com``.

One way to configure this on DigitalOcean is to setup round-robin DNS via the `DNS control panel`_.
To do this, add the following records to your domain:

 - A wildcard CNAME record at your top-level domain, i.e. a CNAME record with * as the name, and @
   as the canonical hostname
 - For each CoreOS machine created, an A-record that points to the TLD, i.e. an A-record named @,
   with the droplet's public IP address

The zone file will now have the following entries in it: (your IP addresses will be different)

.. code-block:: console

    *   CNAME   @
    @   IN A    104.131.93.162
    @   IN A    104.131.47.125
    @   IN A    104.131.113.138

For convenience, you can also set up DNS records for each node:

.. code-block:: console

    deis-1   IN A    104.131.93.162
    deis-2   IN A    104.131.47.125
    deis-3   IN A    104.131.113.138

If you need help using the DNS control panel, check out `this tutorial`_ on DigitalOcean's
community site.

Apply Security Group Settings
-----------------------------

Because DigitalOcean does not have a security group feature, we'll need to add some custom
``iptables`` rules so our components are not accessible from the outside world. To do this, there
is a script in ``contrib/`` which will help us with that. To run it, use:

.. code-block:: console

    $ for i in 1 2 3; do ssh core@deis-$i.example.com 'bash -s' < contrib/util/custom-firewall.sh; done

Our components should now be locked down from external sources.

Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.


.. _`contrib/digitalocean`: https://github.com/deis/deis/tree/master/contrib/digitalocean
.. _`docl`: https://github.com/nathansamson/docl#readme
.. _`Deis Control Utility`: https://github.com/deis/deis/tree/master/deisctl#readme
.. _`DNS control panel`: https://cloud.digitalocean.com/domains
.. _`this tutorial`: https://www.digitalocean.com/community/tutorials/how-to-set-up-a-host-name-with-digitalocean
.. _`Terraform`: https://terraform.io/
