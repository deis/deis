:title: CLI Plugins
:description: How to manage plugins for the Deis CLI.

.. _cli_plugins:

CLI Plugins
===========

Plugins allow developers to extend the functionality of the :ref:`Deis Client <install-client>`,
adding new commands or features.

If an unknown command is specified, the Client will attempt to execute the command as a
dash-separated command. In this case, ``deis resource:command`` will execute ``deis-resource`` with
the argument list ``command``. In full form:

.. code-block:: console

    $ # these two are identical
    $ deis accounts:list
    $ deis-accounts list

Any flags after the command will also be sent to the plugin as an argument:

.. code-block:: console

    $ # these two are identical
    $ deis accounts:list --debug
    $ deis-accounts list --debug

But flags preceding the command will not:

.. code-block:: console

    $ # these two are identical
    $ deis --debug accounts:list
    $ deis-accounts list
