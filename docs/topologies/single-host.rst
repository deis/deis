:title: Single Host Topology for a Deis PaaS
:description: Learn how to build with Deis formations using the Deis command line interface.
:keywords: tutorial, guide, walkthrough, howto, deis, formations

.. _single-host:

Single Host
===========
In a single host topology, the proxy and runtime services of a 
:ref:`formation` run on a single :ref:`Node`.
Though an obvious single point of failure, a single host
formation can work well for small use-cases.

.. include:: create-formation.txt

Create the Layers
-----------------

.. include:: choose-flavor.txt

Create a Dual Proxy/Runtime Layer
`````````````````````````````````
.. code-block:: console

    $ deis layers:create dev nodes ec2-us-west-2 --proxy=y --runtime=y
    Creating nodes layer... done in 3s

We create a new layer in the "dev" formation called "nodes".
The layer has proxy set to "yes" and runtime set to "yes", which means
the layer will host containers as well as route inbound traffic to them.

.. include:: scale-nodes.txt

Scale Automatically
```````````````````
If you've configured automated provisioning using the Deis :ref:`Provider` API,
you can use the ``deis nodes:scale`` command.

.. code-block:: console

    $ deis nodes:scale dev nodes=1
    Scaling nodes... but first, coffee!
    done in 263s

This will automatically provision a new node (separate from the controller) 
which will host the entire formation.

.. include:: scale-manually.txt
   
.. include:: wildcard-dns.txt

.. code-block:: console

    $ dig testing123.deisapp.com +noall +answer
    ...
    testing123.deisapp.com.	45	IN	A	54.245.11.172
