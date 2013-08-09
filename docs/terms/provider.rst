:title: Provider
:description: A provider is a pluggable connector to a third-party cloud API. Supported providers come pre-installed on the Deis controller.
:keywords: provider, deis

.. _provider:

Provider
========
A provider is a pluggable connector to a third-party cloud API, such as `Amazon EC2`_.
Deis's supported providers come pre-installed on the :ref:`Controller`.

Building a custom provider is simple.  It must publish 5 methods as Celery tasks:

* build_layer - to create any shared infrastructure needed by the layer's nodes
* destroy_layer - to destroy any shared infrastructure
* launch_node - to provision a node and register it with the Chef Server
* terminate_node - to destroy a node a remove its records from the Chef Server
* converge_node - to force converge a node (using SSH or other means)

Provider developers can review the `EC2 Reference Implementation`_.

.. _`Amazon EC2`: http://aws.amazon.com/ec2/
.. _`EC2 Reference Implementation`: https://github.com/opdemand/deis/blob/master/celerytasks/ec2.py
