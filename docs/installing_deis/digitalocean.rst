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

In order to provision the cluster, we will need to install a couple of administrative tools.
`docl`_ is a convenience tool to help provision DigitalOcean Droplets. We will also require the
`Deis Control Utility`_, which will assist us with installing, configuring and managing the Deis
platform.

Check System Requirements
-------------------------

Please refer to :ref:`system-requirements` for resource considerations when choosing a droplet
size to run Deis.


Generate SSH Key
----------------

.. include:: ../_includes/_generate-ssh-key.rst


Generate a New Discovery URL
----------------------------

.. include:: ../_includes/_generate-discovery-url.rst


Create CoreOS Droplets
----------------------

Now that we have the user-data file, we can provision some Droplets. We've made this process simple
by supplying a script that does all the heavy lifting for you. If you want to provision manually,
however, start by uploading the SSH key you wish to use to log into each of these servers. After
that, create at least three Droplets with the following specifications:

 - All Droplets deployed in the same region
 - Region must have private networking enabled
 - Region must have User Data enabled. Supply the user-data file here
 - Select CoreOS Stable channel
 - Select your SSH key from the list

If private networking is not available in your region, swap out ``$private_ipv4`` with
``$public_ipv4`` in the user-data file.

If you want to use the script:

.. code-block:: console

    $ gem install docl
    $ docl authorize
    $ docl upload_key deis ~/.ssh/deis.pub
    $ # retrieve your SSH key's ID
    $ docl keys
    deis (id: 12345)
    $ # retrieve the region name
    $ docl regions --metadata --private-networking
    Amsterdam 2 (ams2)
    Amsterdam 3 (ams3)
    London 1 (lon1)
    New York 3 (nyc3)
    Singapore 1 (sgp1)
    $ ./contrib/digitalocean/provision-do-cluster.sh nyc3 12345 4GB

Which will provision 3 CoreOS nodes for use.

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
