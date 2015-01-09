:title: Customizing registry
:description: Learn how to tune custom Deis settings.

.. _registry_settings:

Customizing registry
=========================
The following settings are tunable for the :ref:`registry` component.

Dependencies
------------
Requires: :ref:`cache <cache_settings>`, :ref:`store-gateway <store_gateway_settings>`

Required by: :ref:`builder <builder_settings>`, :ref:`controller <controller_settings>`

Considerations: none

Settings set by registry
--------------------------
The following etcd keys are set by the registry component, typically in its /bin/boot script.

===================================      =================================================================================
setting                                  description
===================================      =================================================================================
/deis/registry/auth/docker/username      username in docker public registry
/deis/registry/auth/docker/password      password in docker public registry
/deis/registry/auth/url/username         username in private registry. This will force https. (url: hostname and port)
/deis/registry/auth/url/password         password in private registry. This will force https. (url: hostname and port)
/deis/registry/bucketName                store component bucket used for registry image layers (default: registry)
/deis/registry/host                      IP address of the host running registry
/deis/registry/port                      port used by the registry service (default: 5000)
/deis/registry/protocol                  protocol for registry (default: http)
/deis/registry/secretKey                 used for secrets (default: randomly generated)
===================================      =================================================================================

Settings used by registry
---------------------------
The following etcd keys are used by the registry component.

====================================      =================================================================================
setting                                   description
====================================      =================================================================================
/deis/cache/host                          host of the cache component (set by cache)
/deis/cache/port                          port of the cache component (set by cache)
/deis/store/gateway/accessKey             S3 API access used to access store-gateway (set by store-gateway)
/deis/store/gateway/host                  host of the store-gateway component (set by store-gateway)
/deis/store/gateway/port                  port of the store-gateway component (set by store-gateway)
/deis/store/gateway/secretKey             S3 API secret key used to access store-gateway (set by store-gateway)
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

Private registry authentication
-------------------------------

Docker provides the command `docker login` to authenticate users. It is possible to use the generated configuration file
in deis to gain access to registries with authentication

Procedure
---------

In any node of the cluster follow the docker procedure to use an account (`docker login`) with any private registry required

.. note::
    To use a private registry use: docker login https://example.com:5000


Compile and run the utility read-dockercfg located in `docker-authentication`

.. code-block:: console

    $ read-dockercfg [-dockercfg optional path to the .dockercfg directory]

.. note::
    read-dockercfg will read the content of the file .dockercfg and set the required keys in etcd


Copy the file $HOME/.dockercfg to all the nodes in the cluster.
Also copy the file to the /root folder (docker run as root in coreos)

.. _`stock registry image`: https://github.com/deis/deis/tree/master/registr
.. _`README`: https://github.com/dotcloud/docker-registry/blob/master/README.md
.. _`docker login`: https://docs.docker.com/reference/commandline/cli/#login
.. _`docker-authentication`: https://github.com/deis/deis/tree/master/contrib/docker-authentication
