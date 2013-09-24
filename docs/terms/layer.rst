:title: Layer
:description: Deis layers are homogeneous groups of nodes that perform work on behalf of a formation. Details on runtime layers, proxy layers and custom layers.
:keywords: layer, nodes, deis

.. _layer:

Layer
=====
Layers are homogeneous groups of :ref:`Nodes <node>` that 
perform work on behalf of a formation.  Each node in a layer has 
the same :ref:`Flavor` and configuration, allowing them to be scaled
easily.

Runtime Layers
^^^^^^^^^^^^^^
Runtime layers host :ref:`Containers <container>` for a formation.
Nodes in a runtime layer use a `Chef Databag`_ to deploy containers for 
each :ref:`application` in the formation.

Proxy Layers
^^^^^^^^^^^^
Proxy layers expose :ref:`Applications <application>` to the outside world.
Nodes in a proxy layer use a `Chef Databag`_ to configure routing of 
inbound requests to :ref:`Containers <container>` hosted on runtime layers.

Dual Layers
^^^^^^^^^^^
A layer can be serve as both a proxy and a runtime.  By default, most 
formations are created with an initial "runtime" layer that is dual proxy/runtime,
allowing an entire formation to live on one :ref:`node`. 

Custom Layers
^^^^^^^^^^^^^
It is also possible to create custom layers that don't provide runtime or proxy
services to the formation.  This is useful for scaling and managing backing
services built with Chef that need to be managed alongside a formation.

.. _`Chef Databag`: http://docs.opscode.com/essentials_data_bags.html
