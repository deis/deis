:title: Customizing cache
:description: Learn how to tune custom Deis settings.

.. _cache_settings:

Customizing cache
=========================
The following settings are tunable for the :ref:`cache` component. Values are stored in etcd.

Dependencies
------------
Requires: none

Required by: :ref:`controller <controller_settings>`, :ref:`registry <registry_settings>`

Considerations: none

Settings set by cache
---------------------
The following etcd keys are set by the cache component, typically in its /bin/boot script.

================              ==============================================
setting                       description
================              ==============================================
/deis/cache/host              IP address of the host running cache
/deis/cache/port              port used by the cache service (default: 6379)
================              ==============================================

Settings used by cache
----------------------

.. important::

   All changes require a restart.


The following etcd keys are used by the cache component.

====================================      =============================================================================================================================================================================================
setting                                   description
====================================      =============================================================================================================================================================================================
/deis/cache/maxmemory                     maximum memory used for caching (default: 50mb)
====================================      =============================================================================================================================================================================================

Using a custom cache image
--------------------------
You can use a custom Docker image for the cache component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config cache set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config cache set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock cache image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock cache image`: https://github.com/deis/deis/tree/master/cache
