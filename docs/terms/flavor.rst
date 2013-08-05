:title: Flavor
:description: What is a Deis Flavor?
:keywords: flavor

.. _flavor:

Flavor
======
A flavor defines the configuration for :ref:`Nodes <node>` in a 
:ref:`Layer`, including their:

* Provider Type (e.g. EC2)
* Launch Parameters (region, zone, etc)
* Initial Configuration using `cloud-config`_

The :ref:`Controller` comes pre-seeded with default flavors for EC2
that use 64-bit Deis-optimized AMIs with an m1.medium instance size.

.. _`cloud-config`: http://cloudinit.readthedocs.org/en/latest/topics/examples.html#install-and-run-chef-recipes
