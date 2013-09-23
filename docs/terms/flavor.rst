:title: Flavor
:description: A Deis flavor defines the configuration for nodes in a layer, including their provider type and launch parameters.
:keywords: flavor, deis, nodes, configuration

.. _flavor:

Flavor
======
A flavor defines the configuration for :ref:`Nodes <node>` in a 
:ref:`Layer`, including their:

* Provider Type (e.g. EC2)
* Launch Parameters (region, zone, etc)

The :ref:`Controller` comes pre-seeded with default flavors for supported providers.
The default flavors typically use Deis-optimized images for faster provisioning.
