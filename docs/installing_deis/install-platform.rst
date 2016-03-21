:title: Installing the Deis Platform
:description: Learn how to install the Deis Platform

.. _install_deis_platform:

Install the Deis Platform
=========================

We will use the ``deisctl`` utility to provision the Deis platform
from a CoreOS host or a workstation that has SSH access to CoreOS.

First check that you have ``deisctl`` installed and the version is correct.

.. code-block:: console

    $ deisctl --version
    1.13.0

If not, follow instructions to :ref:`install_deisctl`.

Ensure your SSH agent is running and select the private key that corresponds to the SSH key added
to your CoreOS nodes:

.. code-block:: console

    $ eval `ssh-agent -s`
    $ ssh-add ~/.ssh/deis

.. note::

    For Vagrant clusters: ``ssh-add ~/.vagrant.d/insecure_private_key``

Find the public IP address of one of your nodes, and export it to the DEISCTL_TUNNEL environment
variable (substituting your own IP address):

.. code-block:: console

    $ export DEISCTL_TUNNEL=104.131.93.162

If you set up the "convenience" DNS records, you can just refer to them via

.. code-block:: console

    $ export DEISCTL_TUNNEL="deis-1.example.com"

.. note::

    For Vagrant clusters: ``export DEISCTL_TUNNEL=172.17.8.100``

This is the IP address where deisctl will attempt to communicate with the cluster. You can test
that it is working properly by running ``deisctl list``. If you see a single line of output, the
control utility is communicating with the nodes.

Before provisioning the platform, we'll need to add the SSH key to Deis so it can connect to remote
hosts during ``deis run``:

.. code-block:: console

    $ deisctl config platform set sshPrivateKey=~/.ssh/deis

.. note::

    For Vagrant clusters: ``deisctl config platform set sshPrivateKey=${HOME}/.vagrant.d/insecure_private_key``

We'll also need to tell the controller which domain name we are deploying applications under:

.. code-block:: console

    $ deisctl config platform set domain=example.com

.. note::

    For Vagrant clusters: ``deisctl config platform set domain=local3.deisapp.com``

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

Now that you've finished provisioning a cluster, please refer to :ref:`using_deis` to get started
using the platform.
