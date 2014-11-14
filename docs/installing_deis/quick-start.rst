:title: Quick Start
:description: How to start provisioning a multi-node Deis cluster

Quick Start
===========

These steps will help you provision a Deis cluster.


.. _get_the_source:

.. include:: ../_includes/_get-the-source.rst


.. _generate_ssh_key:

Generate an SSH key
-------------------

The ``deisctl`` utility communicates with remote machines over an SSH tunnel.
If you don't already have an SSH key, the following command will generate
a new keypair named "deis":

.. code-block:: console

    $ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis


.. _generate_discovery_url:

Generate a New Discovery URL
----------------------------

Discovery URLs help connect `etcd`_ instances together by storing a list of peer addresses and metadata under a
unique address. You can generate a new discovery URL for use in your platform by
running the following from the root of the repository:

.. code-block:: console

    $ make discovery-url

This will write a new discovery URL to the user-data file. Some essential scripts are supplied in
this user-data file, so it is mandatory for provisioning Deis.

Check System Requirements
-------------------------

The Deis provision scripts default to a machine size which should be adequate to run Deis, but this
can be customized. Please refer to :ref:`system-requirements` for resource considerations when
choosing a machine size to run Deis.

Choose a Provider
-----------------

Choose one of the following providers and deploy a new cluster:

- :ref:`deis_on_aws`
- :ref:`deis_on_digitalocean`
- :ref:`deis_on_gce`
- :ref:`deis_on_rackspace`
- :ref:`deis_on_vagrant`
- :ref:`deis_on_bare_metal`


Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a CoreOS cluster,
please refer to :ref:`install_deisctl` and :ref:`install_deis_platform`.


.. _`CoreOS`: https://coreos.com/
.. _`etcd`: https://github.com/coreos/etcd
