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

====================================      ======================================================
setting                                   description
====================================      ======================================================
/deis/logs/handlertype                    Type of handler 'standard' or 'ringbuffer'
====================================      ======================================================

In memory ring buffer
-----------------------
In case of cephless clusters logger can store some logs in memory.
To enable ring buffer logger mode set /deis/logs/handlertype to 'ringbuffer'

By default logger will try to write all logs to file system (ceph mount),
in case of ring buffer logger store by default 1000 lines of log for each appication in own memory.

Also logger will start web service on 8088 port which uses controller to handle
user requests for application logs.

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
