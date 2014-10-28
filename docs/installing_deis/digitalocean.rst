:title: Installing Deis on DigitalOcean
:description: How to provision a multi-node Deis cluster on DigitalOcean

.. _deis_on_digitalocean:

DigitalOcean
============

In this tutorial, we will show you how to set up your own 3-node cluster on DigitalOcean. This
guide is also available in DigitalOcean's `Community site`_, so check out their guide as well!

Prerequisites
-------------

To complete this guide, you must have the following:

 - An SSH key for running operator's commands against the cluster using ``deisctl``
 - An SSH key for authorizing yourself against Deis' builder
 - A domain to point to the cluster
 - The ability to provision at least 3 DigitalOcean Droplets that are 2GB or greater

In order to provision the cluster, we will need to install a couple of administrative tools.
`docl`_ is a convenience tool to help provision DigitalOcean Droplets. We will also require the
`Deis Control Utility`_, which will assist us with installing, configuring and managing the Deis
platform.

Generate a New Discovery URL
----------------------------

To get started with provisioning the Droplets, we will need to generate a new Discovery URL.
Discovery URLs help connect `etcd`_ instances together by storing a list of peer addresses and
metadata under a unique address. You can generate a new discovery URL for use in your platform by
running the following from the root of the repository:

.. code-block:: console

    $ make discovery-url

This will write a new discovery URL to the user-data file. This file is used by DigitalOcean's v2
metadata API to create and customize each machine in our cluster to our liking. Some convenience
scripts are supplied in this user-data file, so it is mandatory for provisioning Deis.

Create CoreOS Droplets
----------------------

Now that we have the user-data file, we can provision some Droplets. We've made this process simple
by supplying a script that does all the heavy lifting for you. If you want to provision manually,
however, start by uploading the SSH key you wish to use to log into each of these servers. After
that, create at least three Droplets with the following specifications:

 - At least 2GB -- more is recommended
 - All Droplets deployed in the same region
 - Region must have private networking enabled
 - Region must have User Data enabled. Supply the user-data file here
 - Select CoreOS Alpha channel
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
    $ ./contrib/digitalocean/provision-do-cluster nyc3 12345 4GB

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

Install Deis Control Utility
----------------------------

Now that we have the CoreOS cluster set up, we will install the Deis Control Utility. This client
will help us configure and install the platform on top of our CoreOS cluster. Please see
:ref:`install_deisctl` for instructions.

Install Deis Platform
---------------------

From the computer you installed the Deis tools on, we will provision the Deis platform. Ensure your
SSH agent is running (and select the private key that corresponds to the SSH keys added to your
CoreOS droplets):

.. code-block:: console

    $ eval `ssh-agent -s`
    $ ssh-add ~/.ssh/deis

Find the public IP address of one of your CoreOS droplets, and export it to the DEISCTL_TUNNEL
environment variable (substitute your own IP address):

.. code-block:: console

    $ export DEISCTL_TUNNEL=104.131.93.162

If you set up the "convenience" DNS records, you can just refer to them via

.. code-block:: console

    $ export DEISCTL_TUNNEL="deis-1.example.com"

This is the IP address where deisctl will attempt to communicate with the cluster. You can test
that it is working properly by running deisctl list. If you see a single line of output, the
control utility is communicating with the CoreOS machines.

Before provisioning the platform, we'll need to add the SSH key to deis so it can connect to remote
hosts during ``deis run``:

.. code-block:: console

    $ deisctl config platform set sshPrivateKey=~/.ssh/deis

We'll also need to tell the controller which domain name we are deploying applications under:

.. code-block:: console

    $ deisctl config platform set domain=example.com

Once finished, run this command to provision the Deis platform:

.. code-block:: console

    $ deisctl install platform

You will see output like the following, which indicates that the units required to run Deis have
been loaded on the CoreOS cluster:

.. code-block:: console

    ● ▴ ■
    ■ ● ▴ Installing Deis...
    ▴ ■ ●

    Scheduling data containers...
    ...
    Deis installed.
    Please run `deisctl start platform` to boot up Deis.

Run this command to start the Deis platform:

.. code-block:: console

    $ deisctl start platform

Once you see "Deis started.", your Deis platform is running on a cluster! You may verify that all
of the Deis units are loaded and active by running the following command:

.. code-block:: console

    $ deisctl list

All of the units should be active.

Now that you've finished provisioning a cluster, please refer to :ref:`using_deis` to get
started using the platform.


.. _`Community site`: https://www.digitalocean.com/community/tutorials/how-to-set-up-a-deis-cluster-on-digitalocean
.. _`docl`: https://github.com/nathansamson/docl#readme
.. _`Deis Control Utility`: https://github.com/deis/deis/tree/master/deisctl#readme
.. _`DNS control panel`: https://cloud.digitalocean.com/domains
.. _`etcd`: https://github.com/coreos/etcd
.. _`this tutorial`: https://www.digitalocean.com/community/tutorials/how-to-set-up-a-host-name-with-digitalocean
