:title: Quick Start
:description: How to start provisioning a multi-node Deis cluster

Quick Start
===========

These steps will help you provision a Deis cluster.

.. _get_the_source:

Get the Source
--------------

.. include:: ../_includes/_get-the-source.rst


.. _generate_ssh_key:

Generate SSH Key
----------------

.. include:: ../_includes/_generate-ssh-key.rst


.. _generate_discovery_url:

Generate a New Discovery URL
----------------------------

.. include:: ../_includes/_generate-discovery-url.rst


Check System Requirements
-------------------------

The Deis provision scripts default to a machine size which should be adequate to run Deis, but this
can be customized. Please refer to :ref:`system-requirements` for resource considerations when
choosing a machine size to run Deis.

Choose a Provider
-----------------

Choose one of the following providers and deploy a new cluster:

- :ref:`deis_on_aws`
- :ref:`deis_on_bare_metal`
- :ref:`deis_on_digitalocean`
- :ref:`deis_on_gce`
- :ref:`deis_on_azure`
- :ref:`deis_on_linode`
- :ref:`deis_on_openstack`
- :ref:`deis_on_vagrant`


Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a CoreOS cluster,
please :ref:`install_deis_platform`.


.. _`CoreOS`: https://coreos.com/
