:title: Configure DNS
:description: Configure name resolution for your Deis Cluster

.. _configure-dns:

Configure DNS
-------------

For local clusters, we've created the DNS record ``local.deisapp.com`` which resolves to the IP of the first VM, 172.17.8.100.
You can use ``local.deisapp.com`` to both log into the controller and to access applications that you've deployed (they will be subdomains of ``local.deisapp.com``, like ``happy-unicorn.local.deisapp.com``). Similarly, you can use ``local3.deisapp.com`` or ``local5.deisapp.com`` for 3- and 5-node clusters, respectively. No DNS configuration is necessary for local clusters.

For Deis clusters hosted elsewhere (EC2, Rackspace, bare metal, etc.), DNS records will need to be created to point to the cluster. For a one-node cluster, we schedule and launch one router, and deis-router and deis-controller will run on the same host. So, both DNS records specified below can be configured to point to this one machine.

On a multi-node cluster, however, there are probably multiple routers, and the controller will likely be scheduled on a separate machine. As mentioned in :ref:`configure-load-balancers`, a load balancer is recommended in this scenario.

Note that the controller will eventually live behind the routers so that all external traffic will flow through the load balancer - configuring a DNS record which points to a service whose IP could change is less than ideal.

Necessary DNS records
---------------------

The DNS records for Deis should be configured as such:
* ``deis.example.org`` should resolve to the IP of the machine that runs ``deis-controller``
* ``*.deis.example.org`` (a wildcard DNS entry) should point to the load balancer (or the same machine for 1-node Vagrant, or any single instance of ``deis-router`` if one likes to live life on the edge)

These records are necessary for all deployments of Deis (EC2, Rackspace, bare metal, etc.). Local clusters can use the hostnames ``local.deisapp.com``, ``local3.deisapp.com``, or ``local5.deiaspp.com``.
