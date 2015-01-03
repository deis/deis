:title: Customizing controller
:description: Learn how to tune custom Deis settings.

.. _controller_settings:

Customizing controller
=========================
The following settings are tunable for the :ref:`controller` component.

Dependencies
------------
Requires: :ref:`cache <cache_settings>`, :ref:`database <database_settings>`, :ref:`registry <registry_settings>`

Required by: :ref:`router <router_settings>`

Considerations: none

Settings set by controller
--------------------------
The following etcd keys are set by the controller component, typically in its /bin/boot script.

=============================            =================================================================================
setting                                  description
=============================            =================================================================================
/deis/controller/host                    IP address of the host running controller
/deis/controller/port                    port used by the controller service (default: 8000)
/deis/controller/protocol                protocol for controller (default: http)
/deis/controller/secretKey               used for secrets (default: randomly generated)
/deis/controller/builderKey              used by builder to authenticate with the controller (default: randomly generated)
/deis/controller/unitHostname            See `Unit hostname`_. (default: "default")
/deis/builder/users/*                    stores user SSH keys (used by builder)
/deis/domains/*                          domain configuration for applications (used by router)
=============================            =================================================================================

Settings used by controller
---------------------------
The following etcd keys are used by the controller component.

====================================      ======================================================
setting                                   description
====================================      ======================================================
/deis/controller/registrationEnabled      enable registration for new Deis users (default: true)
/deis/controller/webEnabled               enable controller web UI (default: false)
/deis/cache/host                          host of the cache component (set by cache)
/deis/cache/port                          port of the cache component (set by cache)
/deis/database/host                       host of the database component (set by database)
/deis/database/port                       port of the database component (set by database)
/deis/database/engine                     database engine (set by database)
/deis/database/name                       database name (set by database)
/deis/database/user                       database user (set by database)
/deis/database/password                   database password (set by database)
/deis/registry/host                       host of the registry component (set by registry)
/deis/registry/port                       port of the registry component (set by registry)
/deis/registry/protocol                   protocol of the registry component (set by registry)
====================================      ======================================================

Using a custom controller image
-------------------------------
You can use a custom Docker image for the controller component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config controller set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config controller set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock controller image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock controller image`: https://github.com/deis/deis/tree/master/controller

Unit hostname
-------------
Per default, Docker automatically generates a hostname for your application unit, such as: 
``5c149b397cd6``. Auto generated hostnames is not always preferred. For instance, 
New Relic would classify each Docker container as an unique server since they use hostname 
for grouping applications running on the same server together.

Deis supports configuring hostname assignment through the ``unitHostname`` setting.
You can change the assignment solution using the following command:

.. code-block:: console

    $ deisctl config controller set unitHostname=application

The valid ``unitHostname`` values are:

default
    Docker will generate the hostname. Example: ``5c149b397cd6``

application
    The hostname is assigned based on the unit name. Example: ``dancing-cat.v2.web.1``

server
    The hostname is assigned based on the CoreOS hostname. Example: 
    ``ip-10-21-2-168.eu-west-1.compute.internal``

.. note::

    Changes to ``/deis/controller/unitHostname`` requires either pushing a new build to 
    every application or scaling them down and up.
    The change is only detected when a container unit is deployed.