:title: Configure DNS
:description: Configure name resolution for your Deis Cluster

.. _configure-dns:

Configure DNS
=============

For local clusters, we've created the DNS record ``local.deisapp.com`` which resolves to the IP of the first VM, 172.17.8.100.
You can use ``local.deisapp.com`` to both log into the controller and to access applications that you've deployed (they will be subdomains of ``local.deisapp.com``, like ``happy-unicorn.local.deisapp.com``). Similarly, you can use ``local3.deisapp.com`` or ``local5.deisapp.com`` for 3- and 5-node clusters, respectively. No DNS configuration is necessary for local clusters.

For Deis clusters hosted elsewhere (EC2, Rackspace, DigitalOcean, Google Compute Engine, bare metal, etc.), DNS records will need to be created to point to the cluster. For a one-node cluster, we schedule and launch one router, and deis-router and deis-controller will run on the same host. So, the DNS record specified below can be configured to point to this one machine.

On a multi-node cluster, however, there are probably multiple routers, and the controller will likely be scheduled on a separate machine. As mentioned in :ref:`configure-load-balancers`, a load balancer is recommended in this scenario.

Note that the controller will eventually live behind the routers so that all external traffic will flow through the load balancer - configuring a DNS record which points to a service whose IP could change is less than ideal.

.. _dns_records:

Necessary DNS records
---------------------

Deis requires one wildcard DNS record. Assuming ``myapps.com`` is the top-level domain apps will live under:

* ``*.myapps.com`` should have A-record entries for each of the load balancer IP addresses

Apps can then be accessed by browsers at ``appname.myapps.com``, and the controller will be available to the Deis client at ``deis.myapps.com``.

This record is necessary for all deployments of Deis (EC2, Rackspace, DigitalOcean, Google Compute Engine, bare metal, etc.). Local clusters can use the domain ``local.deisapp.com``, ``local3.deisapp.com``, or ``local5.deiaspp.com``.

.. _xip_io:

Using xip.io
------------
An alternative to configuring your own DNS records is to use `xip`_ to reference the IP of your load balancer. For example:

.. code-block:: console

    $ deis register http://deis.10.21.12.2.xip.io

You would then create the cluster with ``10.21.12.2.xip.io`` as the cluster domain.

Note that xip does not seem to work for EC2 ELBs - you will have to use an actual DNS record.

.. _`xip`: http://xip.io/
