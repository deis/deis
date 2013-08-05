:title: Formation
:description: What is a Deis Formation?
:keywords: formation, deis

.. _formation:

Formation
=========
A :ref:`formation` is a set of infrastructure used to host a single application
or service backed by a single git repository. Each formation includes
:ref:`Layers <layer>` of :ref:`Nodes <node>` used to host services, a set of 
:ref:`Containers <container>` used to run isolated processes, and a 
:ref:`Release` that defines the current :ref:`Build` and :ref:`config` 
deployed by containers.