:title: Customizing router
:description: Learn how to tune custom Deis settings.

.. _router_settings:

Customizing router
=========================
The following settings are tunable for the :ref:`router` component.

Dependencies
------------
Requires: :ref:`builder <builder_settings>`, :ref:`controller <controller_settings>`, :ref:`store-gateway <store_gateway_settings>`

Required by: none

Considerations: none

Settings set by router
--------------------------
The following etcd keys are set by the router component, typically in its /bin/boot script.

=============================            ===================================================================================
setting                                  description
=============================            ===================================================================================
/deis/router/hosts/$HOST                 IP address and port of the host running this router (there can be multiple routers)
=============================            ===================================================================================

Settings used by router
---------------------------
The following etcd keys are used by the router component.

=======================================      ==================================================================================================================================================================================================================================================================================================================================
setting                                      description
=======================================      ==================================================================================================================================================================================================================================================================================================================================
/deis/builder/host                           host of the builder component (set by builder)
/deis/builder/port                           port of the builder component (set by builder)
/deis/config/\*/deis_whitelist               comma separated list of IPs (or CIDR) allowed to connect to the application containers (set by controller) Example: "0.0.0.0:some_optional_label,10.0.0.0/8"
/deis/controller/host                        host of the controller component (set by controller)
/deis/controller/port                        port of the controller component (set by controller)
/deis/domains/\*                             domain configuration for applications (set by controller)
/deis/router/affinityArg                     for requests with the indicated query string variable, hash its contents to perform session affinity (default: undefined)
/deis/router/bodySize                        nginx body size setting (default: 1m)
/deis/router/defaultTimeout                  default timeout value in seconds. Should be greater then the frontfacing load balancers timeout value (default: 1300)
/deis/router/builder/timeout/connect         proxy_connect_timeout for deis-builder (default: 10000). Unit in milliseconds
/deis/router/builder/timeout/tcp             proxy_timeout for deis-builder (default: 1200000). Unit in milliseconds
/deis/router/controller/timeout/connect      proxy_connect_timeout for deis-controller (default: 10m)
/deis/router/controller/timeout/read         proxy_read_timeout for deis-controller (default: 20m)
/deis/router/controller/timeout/send         proxy_send_timeout for deis-controller (default: 20m)
/deis/router/controller/whitelist            comma separated list of IPs (or CIDR) allowed to connect to the controller (default: not set) Example: "0.0.0.0:some_optional_label,10.0.0.0/8"
/deis/router/enableNginxStatus               enable vhost traffic status page
/deis/router/enforceHTTPS                    redirect all HTTP traffic to HTTPS (default: false)
/deis/router/enforceWhitelist                deny all connections unless specifically whitelisted (default: false)
/deis/router/firewall/enabled                nginx naxsi firewall enabled (default: false)
/deis/router/firewall/errorCode              nginx default firewall error code (default: 400)
/deis/router/errorLogLevel                   nginx error_log level (default: error) Valid options: debug, info, notice, warn, error, crit, alert, emerg
/deis/router/gzip                            nginx gzip setting (default: on)
/deis/router/gzipCompLevel                   nginx gzipCompLevel setting (default: 5)
/deis/router/gzipDisable                     nginx gzipDisable setting (default: "msie6")
/deis/router/gzipHttpVersion                 nginx gzipHttpVersion setting (default: 1.1)
/deis/router/gzipMinLength                   nginx gzipMinLength setting (default: 256)
/deis/router/gzipProxied                     nginx gzipProxied setting (default: any)
/deis/router/gzipTypes                       nginx gzipTypes setting (default: "application/atom+xml application/javascript application/json application/rss+xml application/vnd.ms-fontobject application/x-font-ttf application/x-web-app-manifest+json application/xhtml+xml application/xml font/opentype image/svg+xml image/x-icon text/css text/plain text/x-component")
/deis/router/gzipVary                        nginx gzipVary setting (default: on)
/deis/router/gzipDisable                     nginx gzipDisable setting (default: "msie6")
/deis/router/gzipTypes                       nginx gzipTypes setting (default: "application/x-javascript application/xhtml+xml application/xml application/xml+rss application/json text/css text/javascript text/plain text/xml")
/deis/router/hsts/enabled                    enable HTTP Strict Transport Security headers for HTTPS requests (default: false)
/deis/router/hsts/maxAge                     maximum number of seconds user agents should observe HSTS rewrites (default: 10886400)
/deis/router/hsts/includeSubDomains          enforce HSTS for requests on all subdomains (default: false)
/deis/router/hsts/preload                    allow the domain to be included in the HSTS preload list (default: false)
/deis/router/maxWorkerConnections            maximum number of simultaneous connections that can be opened by a worker process (default: 768)
/deis/router/serverNameHashMaxSize           nginx server_names_hash_max_size setting (default: 512)
/deis/router/serverNameHashBucketSize        nginx server_names_hash_bucket_size (default: 64)
/deis/router/sslCert                         cluster-wide SSL certificate
/deis/router/sslCiphers                      cluster-wide enabled SSL ciphers
/deis/router/sslKey                          cluster-wide SSL private key
/deis/router/sslDhparam                      cluster-wide SSL dhparam
/deis/router/sslProtocols                    nginx ssl_protocols setting (default: TLSv1 TLSv1.1 TLSv1.2)
/deis/router/sslSessionCache                 nginx ssl_session_cache setting (default: not set)
/deis/router/sslSessionTickets               nginx ssl_session_tickets setting (default: on)
/deis/router/sslSessionTimeout               nginx ssl_session_timeout setting (default: 10m)
/deis/router/sslBufferSize                   nginx ssl_buffer_size setting (default: 4k)
/deis/router/trafficStatusZoneSize           nginx vhost_traffic_status_zone size setting (default: 1m)
/deis/router/workerProcesses                 nginx number of worker processes to start (default: auto i.e. available CPU cores)
/deis/router/proxyProtocol                   nginx PROXY protocol enabled
/deis/router/proxyRealIpCidr                 nginx IP with CIDR used by the load balancer in front of deis-router (default: 10.0.0.0/8)
/deis/services/*                             healthy application containers reported by deis/publisher
/deis/store/gateway/host                     host of the store gateway component (set by store-gateway)
/deis/store/gateway/port                     port of the store gateway component (set by store-gateway)
=======================================      ==================================================================================================================================================================================================================================================================================================================================

Using a custom router image
---------------------------
You can use a custom Docker image for the router component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config router set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config router set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock router image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock router image`: https://github.com/deis/deis/tree/master/router


.. _proxy_protocol:

PROXY Protocol
--------------

PROXY is a simple protocol supported by nginx, HAProxy, Amazon ELB, and others. It provides a method
to obtain information about the original requests IP address sent to a load
balancer in front of Deis :ref:`router`.

The Protocol works by prepending, for example, the following to the request:

.. code-block:: text

	PROXY TCP4 129.164.129.164\r\n

The :ref:`router` will pick up the IP information and forward it to the application in the
``X-Forwarded-For`` header.

Load Balancers supporting the HTTP protocol may not need this, except in cases where one would run
WebSockets on a Load Balancer without support for WebSockets (for example AWS ELB) and one also
wants to know the IP address of the original request.
