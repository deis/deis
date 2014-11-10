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
/deis/controller/host                        host of the controller component (set by controller)
/deis/controller/port                        port of the controller component (set by controller)
/deis/domains/*                              domain configuration for applications (set by controller)
/deis/router/bodySize                        nginx body size setting (default: 1m)
/deis/router/builder/timeout/connect         proxy_connect_timeout for deis-builder (default: 10000). Unit in miliseconds
/deis/router/builder/timeout/read            proxy_read_timeout for deis-builder (default: 1200000). Unit in miliseconds
/deis/router/builder/timeout/send            proxy_send_timeout for deis-builder (default: 1200000). Unit in miliseconds
/deis/router/builder/timeout/tcp             timeout for deis-builder (default: 1200000). Unit in miliseconds
/deis/router/controller/timeout/connect      proxy_connect_timeout for deis-controller (default: 10m)
/deis/router/controller/timeout/read         proxy_read_timeout for deis-controller (default: 20m)
/deis/router/controller/timeout/send         proxy_send_timeout for deis-controller (default: 20m)
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
/deis/router/sslCert                         cluster-wide SSL certificate
/deis/router/sslKey                          cluster-wide SSL private key
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
