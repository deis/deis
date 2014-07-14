:title: Configure load balancers
:description: Configure load balancers for your Deis Cluster

.. _configure-load-balancers:

Configure load balancers
------------------------

.. image:: DeisLoadBalancerDiagram.png
    :alt: Deis Load Balancer Diagram

For a one-node Deis cluster, there is one router and one controller, so load balancing is unnecessary.
You can proceed with the next section: :ref:`configure-dns`.

On a multi-node cluster, however, there are probably multiple routers scheduled to the cluster, and
these can potentially move hosts. Therefore, it is recommended that you configure a load balancer
to operate in front of the Deis cluster to serve application traffic.

These ports need to be open on the load balancers:

* 80 (for application traffic and for API calls to the controller)
* 2222 (for traffic to the builder)

Optionally, you can also open port 443 and configure SSL termination on the load balancers, but
requests should still be forwarded to port 80 on the routers. Communication between Deis components
is currently unencrypted.

A health check should be configured on the load balancer to send an HTTP request to /health-check at
port 80 on all nodes in the Deis cluster. The health check endpoint returns an HTTP 200. This enables
the load balancer to serve trafic to whichever hosts happen to be running the deis-router component
at any moment.

.. note::

  Elastic load balancers on EC2 appear to have a default timeout of 60 seconds, which will disrupt
  a ``git push`` when using Deis. Users can request an increased timeout from Amazon. More details
  are in this AWS `support thread`_.

.. _`support thread`: https://forums.aws.amazon.com/thread.jspa?messageID=423862
