:title: SSL Endpoints
:description: Configure SSL termination for your Deis cluster


.. _platform_ssl:

Installing SSL for the Platform
===============================

SSL/TLS is the standard security technology for establishing an encrypted link
between a web server and a browser. This link ensures that all data passed between the web server
and browsers remain private and integral.

To enable SSL for your cluster and all apps running upon it, you can add an SSL key to your load
balancer. You must either provide an SSL certificate that was registered with a CA or provide
:ref:`your own self-signed SSL certificate <creating_self_signed_ssl>`.


.. _load_balancer_ssl:

Installing SSL on a Load Balancer
---------------------------------

On most cloud-based load balancers, you can install a SSL certificate onto the load balancer
itself. Any communication inbound to the cluster will be encrypted while the internal components
of Deis will still communicate over HTTP.

.. note::

    On AWS, Deis enables the :ref:`PROXY protocol <proxy_protocol>` by default, requiring
    installation of :ref:`SSL on the Deis routers <router_ssl>`, as described below.
    Disable the PROXY protocol with ``deisctl config router rm proxyProtocol`` and change
    existing targets and health checks from TCP to HTTP to terminate SSL connections at an
    Amazon ELB instead.

To enable SSL, you will need to open port 443 on the load balancer and forward it to port 80 on the
routers. For AWS, you'll also need to add port 443 in the security group settings for your load
balancer.

See your vendor's specific instructions on installing SSL on your load balancer. For AWS, see their
documentation on `installing an SSL cert for load balancing`_.

.. _router_ssl:

Installing SSL on the Deis Routers
----------------------------------

You can also use the Deis routers to terminate SSL connections.
Use ``deisctl`` to install the certificate and private keys:

.. code-block:: console

    $ deisctl config router set sslKey=<path-to-key> sslCert=<path-to-cert>

If your certificate has intermediate certs that need to be presented as part of a
certificate chain, append the intermediate certs to the bottom of the sslCert value.

.. note::

    To secure all endpoints on the platform domain, you must use a wildcard certificate.


Redirecting traffic to HTTPS
----------------------------

Once your cluster is serving traffic over HTTPS, you can optionally instruct the router component
to forward all traffic on HTTP to HTTPS (application traffic and requests to the controller component).

This is achieved with ``deisctl``:

.. code-block:: console

    $ deisctl config router set enforceHTTPS=true


.. _`installing an SSL cert for load balancing`: http://docs.aws.amazon.com/ElasticLoadBalancing/latest/DeveloperGuide/ssl-server-cert.html
