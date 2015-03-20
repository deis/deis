:title: Using WebSockets
:description: Learn how to enable WebSockets on your Deis installation

.. _configure-websockets:

Configure WebSockets
=========================
The Deis platform supports WebSockets out of the box, however, most Deis deployments have a load
balancer in front of the Deis :ref:`router`, as seen in the diagram:

.. image:: DeisLoadBalancerDiagram.png
    :alt: Deis Load Balancer Diagram

The load balancer in front of Deis must also support WebSockets in order to have a fully working
WebSockets environment on  your Deis deployment. `Amazon Elastic Load Balancer`_,
`Rackspace Cloud Load Balancer`_, and possibly other providers' load balancers does not support
WebSockets out of the box without switching the load balancers transport protocol from **HTTP** to
**TCP**.

.. _`Amazon Elastic Load Balancer`: https://forums.aws.amazon.com/thread.jspa?threadID=84606
.. _`Rackspace Cloud Load Balancer`: https://community.rackspace.com/products/f/25/t/3362

AWS (EC2)
----------
Since AWS ELB does not support WebSockets out of the box in HTTP mode, and setting the ELB in TCP
mode results in the loss of origin IP address and scheme information (``X-Forwarded-For``,
``X-Forwarded-Scheme``), the recommended workaround until AWS ELB supports WebSockets is creating a
extra load balancer for WebSocket traffic in TCP Mode.

To enable the extra load balancer, see :ref:`AWS - WebSockets <deis_on_aws_websockets>`.

Example application setup:
~~~~~~~~~~~~~~~~~~~~~~~~~~
Application at ``example.com`` with WebSocket endpoint at ``example.com/ws``

Configure DNS for ``www.example.com`` to access ``DeisWebELB`` (CNAME)
Configure DNS for  ``ws.example.com`` to access ``DeisTCPELB`` (CNAME)

Add both ``www.example.com`` and ``ws.example.com`` as domains to the
application.

Access websocket endpoint at ``ws.example.com/ws``.

Rackspace
---------
Rackspaces Cloud Load Balancer does not support WebSockets out of the box. By changing the protocol
from HTTP(S) to TCP_CLIENT_FIRST, WebSockets should work.

``X-Forwarded-For`` will be incorrect, and Rackspace does not currently support the PROXY protocol.

An alternative solution would be to have one load balancer use HTTP, and then another load balancer
for the websocket specific part of the application set to ``TCP_CLIENT_FIRST``.

DigitalOcean
------------
The current DigitalOcean setup uses DNS based load balancing, meaning that there is not an
intermediate load balancer between the end user/customer and Deis. WebSockets should work
out of the box with proper ``X-Forwarded-For`` headers being sent to your applications.
