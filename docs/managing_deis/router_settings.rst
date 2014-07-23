:title: Customizing router
:description: Learn how to tune custom Deis settings.

.. _router_settings:

Customizing router
=========================
The following settings are tunable for the :ref:`router` component.

Dependencies
------------
Requires: :ref:`builder <builder_settings>`, :ref:`controller <controller_settings>`

Required by: none

Considerations: none

Settings set by router
--------------------------
The following etcd keys are set by the router component, typically in its /bin/boot script.

===========================              =================================================================================
setting                                  description
===========================              =================================================================================
/deis/router/$HOST/host                  IP address of the host running this router (there can be multiple routers)
/deis/router/$HOST/port                  port used by this router service (there can be multiple routers) (default: 80)
===========================              =================================================================================

Settings used by router
---------------------------
The following etcd keys are used by the router component.

====================================      =============================================================================================================================================================================================
setting                                   description
====================================      =============================================================================================================================================================================================
/deis/domains/*                           domain configuration for applications (set by controller)
/deis/services/*                          application configuration (set by application unit files)
/deis/builder/host                        host of the builder component (set by builder)
/deis/builder/port                        port of the builder component (set by builder)
/deis/controller/host                     host of the controller component (set by controller)
/deis/controller/port                     port of the controller component (set by controller)
/deis/router/gzip                         nginx gzip setting (default: on)
/deis/router/gzipHttpVersion              nginx gzipHttpVersion setting (default: 1.0)
/deis/router/gzipCompLevel                nginx gzipCompLevel setting (default: 2)
/deis/router/gzipProxied                  nginx gzipProxied setting (default: any)
/deis/router/gzipVary                     nginx gzipVary setting (default: on)
/deis/router/gzipDisable                  nginx gzipDisable setting (default: "msie6")
/deis/router/gzipTypes                    nginx gzipTypes setting (default: "application/x-javascript, application/xhtml+xml, application/xml, application/xml+rss, application/json, text/css, text/javascript, text/plain, text/xml")
====================================      =============================================================================================================================================================================================

Using a custom router image
---------------------------
You can use a custom Docker image for the router component instead of the image
supplied with Deis:

.. code-block:: console

    $ etcdctl set /deis/router/image myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ etcdctl set /deis/router/image registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock router image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock router image`: https://github.com/deis/deis/tree/master/router
