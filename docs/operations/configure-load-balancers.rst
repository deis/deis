:title: Configure load balancers
:description: Configure load balancers for your Deis Cluster

.. _configure-load-balancers:

Configure load balancers
------------------------

For a one-node Deis cluster, there is one router and one controller, so load balancing is unnecessary.
You can proceed with the next section: :ref:`configure-dns`.

On a multi-node cluster, however, there are probably multiple routers scheduled to the cluster, and
these can potentially move hosts. Therefore, it is recommended that you configure a load balancer
to operate in front of the Deis cluster to serve application traffic. A simple configuration is one
that has all Deis machines listed in its configuration file, but a host is only considered 'healthy'
when it is responding to ports 80 and 2222. This enables the load balancer to serve trafic to whichever
hosts happen to be running the deis-router component at any one time.

These ports need to be open on the load balancers:

* 80 (for application traffic and for API calls to the controller)
* 2222 (for traffic to the builder)

Optionally, you can also open port 443 and configure SSL termination on the load balancers, but
requests should still be forwarded to port 80 on the routers. Communication between Deis components
is currently unencrypted.
