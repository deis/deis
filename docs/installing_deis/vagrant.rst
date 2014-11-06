:title: Installing Deis on Vagrant
:description: How to provision a multi-node Deis cluster on Vagrant

.. _deis_on_vagrant:

Vagrant
=======

`Vagrant`_ is a tool for building complete development environments with a focus on automation.
This guide demonstrates how you can stand up a Deis cluster for development purposes using Vagrant.


Install Prerequisites
---------------------

Please install `Vagrant`_ v1.6.5+ and `VirtualBox`_.

.. note::

    For Ubuntu users: the VirtualBox package in Ubuntu has some issues when running in
    RAM-constrained environments. Please install the latest version of VirtualBox from Oracle's
    website.


Generate a New Discovery URL
----------------------------

To get started with provisioning Deis, we will need to generate a new Discovery URL. Discovery URLs
help connect `etcd`_ instances together by storing a list of peer addresses and metadata under a
unique address. You can generate a new discovery URL for use in your platform by
running the following from the root of the repository:

.. code-block:: console

    $ make discovery-url

This will write a new discovery URL to the user-data file. Some convenience scripts are supplied in
this user-data file, so it is mandatory for provisioning Deis.


Boot CoreOS
-----------

Start the CoreOS cluster on VirtualBox. From a command prompt, switch directories to the root of
the Deis project and type:

.. code-block:: console

    $ vagrant up

This instructs Vagrant to spin up 3 VMs. To be able to connect to the VMs, you must add your
Vagrant-generated SSH key to the ssh-agent (``deisctl`` requires the agent to have this key):

.. code-block:: console

    $ ssh-add ~/.vagrant.d/insecure_private_key


Configure DNS
-------------

For convenience, we have set up a few DNS records for users running on Vagrant.
``local3.deisapp.com`` is set up for 3-node clusters and ``local5.deisapp.com`` is set up for
5-node clusters. If you want to set up your own DNS records, see :ref:`configure-dns` for more
information.


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.


.. _Vagrant: http://www.vagrantup.com/
.. _VirtualBox: https://www.virtualbox.org/wiki/Downloads
.. _etcd: https://github.com/coreos/etcd
