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
Runtime layers service requests and run background tasks for the formation.
Nodes in a runtime layer use a configuration management system (e.g. Chef Server)
to deploy :ref:`Containers <container>` running a specific :ref:`Release`.

Proxy Layers
^^^^^^^^^^^^
Proxy layers expose the formation to the outside world.
Nodes in a proxy layer use configuration management to configure routing of 
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
