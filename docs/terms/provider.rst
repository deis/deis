:title: Provider
:description: A provider is a pluggable connector to a third-party cloud API. Supported providers come pre-installed on the Deis controller.
:keywords: provider, deis

.. _provider:

Provider
========
A provider is a pluggable connector to a third-party cloud API, such as `Amazon EC2`_.
Deis's supported providers come pre-installed on the :ref:`Controller`.

Building a custom provider is simple.  It must publish 5 methods:

* build_layer - to create any shared infrastructure needed by the layer's nodes
* destroy_layer - to destroy any shared infrastructure
* build_node - to provision a node and prepare it for bootstrapping by Chef
* destroy_node - to destroy a node after it has been purged from the Chef Server
* seed_flavors - to seed the controller database with default flavors

Provider developers can review the `EC2 Reference Implementation`_.

.. _`Amazon EC2`: http://aws.amazon.com/ec2/
.. _`EC2 Reference Implementation`: https://github.com/opdemand/deis/blob/master/provider/ec2.py
