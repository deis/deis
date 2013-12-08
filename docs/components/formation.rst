:title: Formation
:description: A Deis formation is a set of infrastructure used to host applications.
:keywords: formation, deis

.. _formation:

Formation
=========
A :ref:`formation` is a set of infrastructure used to host :ref:`Applications <application>`.
Each formation includes :ref:`Layers <layer>` of :ref:`Nodes <node>` 
that provide different services to the formation.

Creating a Formation
--------------------
Creating a formation is easy...

.. code-block:: console

    $ deis formations:create test
    Creating formation... done, created test  

Viewing a Formation
-------------------
We can take a peek at our new formation using "deis formations:info":

.. code-block:: console

    $ deis formations:info test
    === test Formation
    {
      "updated": "2013-11-26T22:43:37.854Z", 
      "uuid": "6955c2a8-8505-413b-90b2-305daebdbbf9", 
      "created": "2013-11-26T22:43:37.854Z", 
      "domain": null, 
      "owner": "gabrtv", 
      "nodes": "{}", 
      "id": "test"
    }
    
    === test Layers
    
    === test Nodes

This formation has no :ref:`Layers <layer>` and no :ref:`Nodes <node>` 
-- so it can't do much yet.  
It's also important to note that "domain" is null.
Without a domain a formation can only host one application at a time.

Updating a Formation
--------------------
We can update the formation with "deis formations:update".  Let's give our
formation a "domain" so we can have it host multiple applications later
once we scale some proxy nodes and setup wildcard DNS.

.. code-block:: console

    $ deis formations:update test --domain=deisapp.com
    {
      "updated": "2013-11-26T22:54:46.708Z", 
      "uuid": "6955c2a8-8505-413b-90b2-305daebdbbf9", 
      "created": "2013-11-26T22:43:37.854Z", 
      "domain": "deisapp.com", 
      "owner": "gabrtv", 
      "nodes": "{}", 
      "id": "test"
    }

Converging a Formation
----------------------

