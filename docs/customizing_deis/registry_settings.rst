:title: Customizing registry
:description: Learn how to tune custom Deis settings.

.. _registry_settings:

Customizing registry
=========================
The following settings are tunable for the :ref:`registry` component.

Dependencies
------------
Requires: :ref:`store-gateway <store_gateway_settings>`

Required by: :ref:`builder <builder_settings>`, :ref:`controller <controller_settings>`

Considerations: none

Settings set by registry
--------------------------
The following etcd keys are set by the registry component, typically in its /bin/boot script.

===========================              =================================================================================
setting                                  description
===========================              =================================================================================
/deis/registry/bucketName                store component bucket used for registry image layers. Periods are not allowed in the name (default: registry)
/deis/registry/host                      IP address of the host running registry
/deis/registry/port                      port used by the registry service (default: 5000)
/deis/registry/protocol                  protocol for registry (default: http)
===========================              =================================================================================

Settings used by registry
---------------------------
The following etcd keys are used by the registry component.

====================================      =================================================================================
setting                                   description
====================================      =================================================================================
/deis/cache/host                          host of a Redis cache (optional)
/deis/cache/port                          port of a Redis cache (optional)
/deis/store/gateway/accessKey             S3 API access used to access store-gateway (set by store-gateway)
/deis/store/gateway/host                  host of the store-gateway component (set by store-gateway)
/deis/store/gateway/port                  port of the store-gateway component (set by store-gateway)
/deis/store/gateway/secretKey             S3 API secret key used to access store-gateway (set by store-gateway)
====================================      =================================================================================

If the ``/deis/registry/s3bucket`` key is supplied, the registry
will use Amazon S3 as its storage backend and use the following values.

====================================      =================================================================================
setting                                   description
====================================      =================================================================================
/deis/registry/s3accessKey                S3 API access key. If not specified, the registry will get it from the instance role
/deis/registry/s3secretKey                S3 API secret key, required if s3accessKey is specified
/deis/registry/s3region                   S3 region to connect to, will use boto default if not specified
/deis/registry/s3bucket                   S3 bucket to store images. Periods are not allowed in the name.
/deis/registry/s3path                     path in the bucket (default: "/registry")
/deis/registry/s3encrypt                  whether the object is encrypted while at rest on the server (default: true)
/deis/registry/s3secure                   use secure protocol to establish connection with S3 (default: true)
====================================      =================================================================================

The Deis registry component inherits from the Docker registry container, so additional configuration
options can be supplied. For a full explanation of these settings, see the Docker registry `README`_.

Using a custom registry image
-----------------------------
You can use a custom Docker image for the registry component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config registry set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config registry set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock registry image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock registry image`: https://github.com/deis/deis/tree/master/registry
.. _`README`: https://github.com/dotcloud/docker-registry/blob/master/README.md
