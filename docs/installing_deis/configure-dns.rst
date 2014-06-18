:title: Configure DNS
:description: Configure name resolution for your Deis Cluster

.. _configure-dns:

Configure DNS
=============

For local clusters, we've created the DNS record ``local.deisapp.com`` which resolves to the IP of the first VM, 172.17.8.100.
You can use ``local.deisapp.com`` to both log into the controller and to access applications that you've deployed (they will be subdomains of ``local.deisapp.com``, like ``happy-unicorn.local.deisapp.com``). Similarly, you can use ``local3.deisapp.com`` or ``local5.deisapp.com`` for 3- and 5-node clusters, respectively. No DNS configuration is necessary for local clusters.

For Deis clusters hosted elsewhere (EC2, Rackspace, bare metal, etc.), DNS records will need to be created to point to the cluster. For a one-node cluster, we schedule and launch one router, and deis-router and deis-controller will run on the same host. So, the DNS record specified below can be configured to point to this one machine.

If you do not modify your DNS settings on EC2 but still want to register a user, try using xip.io:

.. code-block:: console

    $ deis register deis.ec2-12-345-678-90.us-west-1.compute.amazonaws.com.xip.io

On a multi-node cluster, however, there are probably multiple routers, and the controller will likely be scheduled on a separate machine. As mentioned in :ref:`configure-load-balancers`, a load balancer is recommended in this scenario.

Note that the controller will eventually live behind the routers so that all external traffic will flow through the load balancer - configuring a DNS record which points to a service whose IP could change is less than ideal.

Necessary DNS records
---------------------

Deis requires one wildcard DNS record. Assuming ``myapps.com`` is the top-level domain apps will live under:

* ``*.myapps.com`` should have A-record entries for each of the load balancer IP addresses

Apps can then be accessed via ``appname.myapps.com``, and the Deis controller can be accessed at ``deis.myapps.com``.

This record is necessary for all deployments of Deis (EC2, Rackspace, bare metal, etc.). Local clusters can use the domain ``local.deisapp.com``, ``local3.deisapp.com``, or ``local5.deiaspp.com``.
