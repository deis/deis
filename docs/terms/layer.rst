:title: Layer
:description: Deis layers are homogeneous groups of nodes that perform work on behalf of a formation. Details on runtime layers, proxy layers and custom layers.
:keywords: layer, layers, deis

.. _layer:

Layer
=====
:ref:`Layers <layer>` are homogeneous groups of :ref:`Nodes <node>` that 
perform work on behalf of a formation.  Each node in a layer has 
the same :ref:`Flavor` and Chef configuration, allowing them to be scaled
easily.

Runtime Layers
^^^^^^^^^^^^^^
Runtime layers service requests and run background tasks for the formation.
Nodes in a runtime layer use a Chef databag  to deploy
:ref:`Containers <container>` running a specific :ref:`Release`.  

Proxy Layers
^^^^^^^^^^^^
Proxy layers expose the formation to the outside world.
Nodes in a proxy layer use a Chef databag to configure routing of 
inbound requests to :ref:`Containers <container>` hosted on runtime layers.

Custom Layers
^^^^^^^^^^^^^
It is also possible to create custom layers that contain custom run-lists.