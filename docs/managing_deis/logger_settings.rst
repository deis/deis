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

Considerations: must live on the same host as controller (see `#985`_)

Settings set by logger
------------------------
The following etcd keys are set by the database component, typically in its /bin/boot script.

===========================              =================================================================================
setting                                  description
===========================              =================================================================================
/deis/logger/host                        IP address of the host running logger
/deis/logger/port                        port used by the logger service (default: 514)
===========================              =================================================================================

Settings used by logger
-------------------------
The logger component uses no keys from etcd.

Using a custom logger image
---------------------------
You can use a custom Docker image for the logger component instead of the image
supplied with Deis:

.. code-block:: console

    $ etcdctl set /deis/logger/image myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ etcdctl set /deis/logger/image registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock logger image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock logger image`: https://github.com/deis/deis/tree/master/logger
.. _`#985`: https://github.com/deis/deis/issues/985
