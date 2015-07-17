:title: Configure DNS
:description: Configure name resolution for your Deis Cluster

.. _configure-dns:

Configure DNS
=============

For Deis clusters on Amazon Web Services, Azure, DigitalOcean, Google Compute Engine,
Linode, OpenStack, or bare metal, :ref:`DNS records <dns_records>` must be created.
The cluster runs multiple routers infront of the Deis controller and apps
you deploy, so a :ref:`load balancer <configure-load-balancers>` is recommended.

Vagrant
-------

For local Vagrant clusters, no DNS configuration is required. The domain
``local3.deisapp.com`` already resolves to the IPs of the first 3 VMs provisioned
by Deis' Vagrantfile: 172.17.8.100, 172.17.8.101, 172.17.8.102.

Use ``deis.local3.deisapp.com`` to log in to the controller on a 3-node Vagrant
cluster. Apps that you deploy will have their name prefixed to the domain, such
as "golden-chinbone.local3.deisapp.com".

Similarly, use ``local5.deisapp.com`` for a 5-node Vagrant cluster.

.. _dns_records:

Necessary DNS records
---------------------

Deis requires a wildcard DNS record. Assuming ``myapps.com`` is the top-level domain
apps will live under:

* ``*.myapps.com`` should have "A" record entries for each of the load balancer's IP addresses

Apps can then be accessed by browsers at ``appname.myapps.com``, and the controller will be available to the Deis client at ``deis.myapps.com``.

`AWS recommends`_ against creating "A" record entries; instead, create a wildcard "CNAME" record entry for the load balancer's DNS name, or use Amazon `Route 53`_.

These records are necessary for all deployments of Deis other than Vagrant clusters.

.. _xip_io:

Using xip.io
------------
An alternative to configuring your own DNS records is to use `xip`_ to reference the IP of your load balancer. For example:

.. code-block:: console

    $ deis register http://deis.10.21.12.2.xip.io

You would then create the cluster with ``10.21.12.2.xip.io`` as the cluster domain.

Note that xip does not seem to work for AWS ELBs - you will have to use an actual DNS record.

.. _`AWS recommends`: https://docs.aws.amazon.com/ElasticLoadBalancing/latest/DeveloperGuide/using-domain-names-with-elb.html
.. _`Route 53`: http://aws.amazon.com/route53/
.. _`xip`: http://xip.io/
