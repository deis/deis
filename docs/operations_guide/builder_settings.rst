:title: Customizing builder
:description: Learn how to tune custom Deis settings.

.. _builder_settings:

Customizing builder
=========================
The following settings are tunable for the :ref:`builder` component. Values are stored in etcd.

Dependencies
------------
Requires: :ref:`controller <controller_settings>`, :ref:`registry <registry_settings>`

Required by: :ref:`router <router_settings>`

Considerations: must live on the same host as controller (see `#985`_)

Settings set by builder
-----------------------
The following etcd keys are set by the builder component, typically in its /bin/boot script.

==================              ================================================
setting                         description
==================              ================================================
/deis/builder/host              IP address of the host running builder
/deis/builder/port              port used by the builder service (default: 2223)
==================              ================================================

Settings used by builder
---------------------------
The following etcd keys are used by the builder component.

====================================      ===========================================================
setting                                   description
====================================      ===========================================================
/deis/builder/users/*                     user SSH keys to provision (set by controller)
/deis/controller/builderKey               used to communicate with the controller (set by controller)
/deis/controller/host                     host of the controller component (set by controller)
/deis/controller/port                     port of the controller component (set by controller)
/deis/controller/protocol                 protocol of the controller component (set by controller)
/deis/registry/host                       host of the controller component (set by registry)
/deis/registry/port                       port of the controller component (set by registry)
====================================      ===========================================================

Using a custom builder image
----------------------------
You can use a custom Docker image for the builder component instead of the image
supplied with Deis:

.. code-block:: console

    $ etcdctl set /deis/builder/image myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ etcdctl set /deis/builder/image registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock builder image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock builder image`: https://github.com/deis/deis/tree/master/builder
.. _`#985`: https://github.com/deis/deis/issues/985
