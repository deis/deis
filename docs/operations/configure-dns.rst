:title: Configure DNS
:description: Configure name resolution for your Deis Cluster

.. _configure-dns:

Configure DNS
-------------

For a one-node cluster, both deis-router and deis-controller will run on the same host. For convenience, we've created the DNS record ``local.deisapp.com`` which resolves to the IP of the first VM, 172.17.8.100.
You can use ``local.deisapp.com`` to both log into the controller and to access applications that you've deployed (they will be subdomains of ``local.deisapp.com``, like ``happy-unicorn.local.deisapp.com``).

On a multi-node cluster, however, the router and controller will likely be scheduled on separate machines. Since we cannot know the IP addresses ahead of time, you'll need to setup resolution yourself using your own domain (unfortunately, wildcard hostnames are not permitted in ``/etc/hosts``). The records should be as follows:

* ``deis.example.org`` should resolve to the IP of the machine that runs ``deis-controller``
* ``*.deis.example.org`` (a wildcard DNS entry) should resolve to the IP of the machine that runs ``deis-router``

These records are necessary for multi-node Vagrant as well as any other multi-node deployments of Deis (EC2, Rackspace, etc.).
