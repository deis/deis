:title: Multi Host Topology for a Deis PaaS Formation
:description: Learn how to build with Deis formations using the Deis command line interface.
:keywords: tutorial, guide, walkthrough, howto, deis, formations

.. _multi-host:

Multi Host
==========
In a multi host topology, the proxy and runtime services of a :ref:`formation`
are split into their own :ref:`Layers <layer>` and managed separately.  
A multi host formation can range from 2 :ref:`Nodes <node>` to 
(theoretically) unlimited nodes.
To achieve high availability a multi host formation
should have a minimum of 4 nodes, 2 proxy and 2 runtime.

.. include:: create-formation.txt

Create the Layers
-----------------

.. include:: choose-flavor.txt

Create a Proxy Layer
````````````````````
.. code-block:: console

    $ deis layers:create dev proxy ec2-us-west-2 --proxy=y --runtime=n
    Creating proxy layer... done in 3s

We create a new layer in the "dev" formation called "proxy".
The layer has proxy set to "yes" and runtime set to "no", which means
the layer will route traffic to containers, but will not host containers.

Create a Runtime Layer
``````````````````````
.. code-block:: console

    $ deis layers:create dev runtime ec2-us-west-2 --proxy=n --runtime=y
    Creating runtime layer... done in 4s

We create a new layer in the "dev" formation called "runtime" which
has proxy set to "no" and runtime set to "yes".
This layer is purely for hosting containers.

.. include:: scale-nodes.txt

Scale Automatically
```````````````````
If you've configured automated provisioning using the Deis :ref:`Provider` API,
you can use the ``deis nodes:scale`` command.

.. code-block:: console

    $ deis nodes:scale dev proxy=2 runtime=2
    Scaling nodes... but first, coffee!
    done in 263s

.. include:: scale-manually.txt

.. include:: wildcard-dns.txt

.. code-block:: console

    $ dig testing123.deisapp.com +noall +answer
    ...
    testing123.deisapp.com.	45	IN	A	54.245.11.172
    testing123.deisapp.com.	45	IN	A	54.202.163.190
