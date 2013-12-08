:title: Deploy a Formation with Deis
:description: Learn how to deploy a Deis formation using the Deis command line interface.
:keywords: tutorial, guide, walkthrough, howto, deis, formations

Deploy a Formation
==================
:ref:`Formations <formation>` are a set of :ref:`Nodes <node>` 
used to host :ref:`Applications <application>`. 
The high-level process for deploying a formation is:

 #. Create the Formation
 #. Create Layers
 #. Create Nodes (manually or automatically using the :ref:`Provider` API)
 #. Publish the Proxy Layer using Wildcard DNS

Formations are central to Deis platform operations.
Before we deploy one, let's take a moment to review 
what they include and how they can be structured.

About Layers
------------
Formations are organized by :ref:`layer`, groups of :ref:`Nodes <node>` 
that provide services to the formation.  These services include:

 * Runtime Services - host Docker containers
 * Proxy Services - route traffic to Docker containers
 * Custom Services - custom services defined via a Chef

Since layers are unique to Deis, let's explore the options for
creating a layer using the Deis client:

.. code-block:: console

    $ deis help layers:create
    Create a layer of nodes
    
    Usage: deis layers:create <formation> <id> <flavor> [options]
    
    Options:
    --proxy=<yn>                    layer can be used for proxy [default: y]
    --runtime=<yn>                  layer can be used for runtime [default: y]
    --ssh_username=USERNAME         username for ssh connections [default: ubuntu]
    --ssh_private_key=PRIVATE_KEY   private key for ssh comm (default: auto-gen)
    --ssh_public_key=PUBLIC_KEY     public key for ssh comm (default: auto-gen)
    --ssh_port=<port>               port number for ssh comm (default: 22)

Layer SSH Configuration
```````````````````````
Note that SSH keys and settings are stored in the layer.
When Deis accesses nodes via SSH, it uses layer's SSH configuration. 
By default, SSH keys are generated and managed automatically.

Proxy/Runtime Booleans
``````````````````````
The layer has two important true/false flags: proxy and runtime.
These flags tell Deis if nodes in this layer should act as a:

 * runtime - host Docker containers for the formation
 * proxy - route traffic to Docker containers in the formation

You can configure "dual" layers that act as proxy and a runtime simultaneously.
Both can also be set to false -- useful for custom Chef layers.

Formation Topologies
--------------------
Deis supports many different host/network topologies for each formation.
Topologies range from single-host deployments to N-node formations with 
dedicated proxy and runtime layers.

Single Host
```````````
In a single host topology, the proxy and runtime services of a 
:ref:`formation` run on a single :ref:`Node`.
Though an obvious single point of failure, a single host
formation can work well for small use-cases.

Dual Host
`````````
In a dual host topology, the proxy and runtime services of a :ref:`formation` 
run on two :ref:`Nodes <node>`.
This provides improved availability, since the failure of a single node still 
leaves another node with both proxy and runtime services intact.

Multi Host
``````````
In a multi host topology, the proxy and runtime services of a :ref:`formation`
are split into their own :ref:`Layers <layer>` and managed separately.  
A multi host formation can range from 2 :ref:`Nodes <node>` to 
(theoretically) unlimited nodes.
To achieve high availability a multi host formation
should have a minimum of 4 nodes, 2 proxy and 2 runtime.

Deploy a Formation
------------------
Now that we understand the moving parts, it's time to select a topology and deploy it.

 * Deploy a :ref:`single-host` Formation
 * Deploy a :ref:`dual-host` Formation
 * Deploy a :ref:`multi-host` Formation
