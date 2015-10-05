:title: Customizing logger
:description: Learn how to tune custom Deis settings.

.. _logger_settings:

Customizing logger
=========================
The following settings are tunable for the :ref:`logger` component.

Dependencies
------------
Requires: none

Required by: :ref:`controller <controller_settings>`

Considerations: none

Settings set by logger
------------------------
The following etcd keys are set by the logger component, typically in its /bin/boot script.

===========================              =================================================================================
setting                                  description
===========================              =================================================================================
/deis/logs/host                          IP address of the host running logger
/deis/logs/port                          port used by the logger service (default: 514)
===========================              =================================================================================

Settings used by logger
-------------------------
The following etcd keys are used by the logger component.

====================================      ================================================================================
setting                                   description
====================================      ================================================================================
/deis/logs/storageAdapterType             Type of storage adapter to use: ``file`` or ``memory``; if not set, ``file`` is assumed.  It is also possible so specify the size of the in-memory adapter's internal ring buffer (in lines; a line is a max of 65k) using a value like: ``memory:<size>``.  1000 is the default size.
/deis/logs/drain                          URL for an external service that logs can be forwarded to for long-term archival. If not set, no drain is used.  URLs beginning with ``udp://``, ``syslog://`` use UDP for transport.  URLs beginning with ``tcp://`` use TCP.
====================================      ================================================================================

.. note::

  Those running the stateless (Ceph-less) platform should prefer the in-memory storage adapter.

Using a custom logger image
---------------------------

.. note::

  Instead of using a custom logger image, it is possible to redirect Deis logs to an external location.
  For more details, see :ref:`platform_logging`.

You can use a custom Docker image for the logger component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config logger set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config logger set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock logger image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock logger image`: https://github.com/deis/deis/tree/master/logger
