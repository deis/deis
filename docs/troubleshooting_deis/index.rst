:title: Troubleshooting Deis
:description: Resolutions for common issues encountered when running Deis.

.. _troubleshooting_deis:

Troubleshooting Deis
====================

:Release: |version|
:Date: |today|

.. toctree::

    troubleshooting-store

Common issues that users have run into when provisioning Deis are detailed below.

Logging in to the cluster
-------------------------

To open a interactive shell on a machine in your cluster:

.. code-block:: console

    $ deisctl ssh <unit>

For example, to open a shell session on the machine that is running Controller,
you can run this:

.. code-block:: console

    $ deisctl ssh controller

You can execute just a single command instead of opening a shell:

.. code-block:: console

    $ deisctl ssh <unit> <command>

You can also connect directly to the Docker instance of that unit:

.. code-block:: console

    $ deisctl dock <unit> <command>

For example, to start a Bash session on the Builder Docker container, you can
run the following command:

.. code-block:: console

    $ deisctl dock builder bash`


Troubleshooting etcd
--------------------

Sometimes issues with Deis are caused by latency between CoreOS hosts. A telltale sign of this is
if all of the Deis components on a single machine crash. To aid in debugging etcd, we've created
a system service that is installed but not started when you deploy CoreOS using our provision scripts.

To start this service, run ``sudo systemctl start debug-etcd`` on a CoreOS machine in your cluster.
This starts a service which queries etcd's state once per second. Watching this output with
``journalctl -fu debug-etcd`` makes it easy to spot heartbeat timeouts or other abnormalities
which will lead to issues running Deis successfully.

A deis-store component fails to start
-------------------------------------

For information on troubleshooting a ``deis-store`` component, see :ref:`troubleshooting-store`.

Any component fails to start
----------------------------

Use ``deisctl status <component>`` to view the status of the component.
You can also use ``deisctl journal <component>`` to tail logs for a component, or ``deisctl list``
to list all components.

Failed initializing SSH client
------------------------------

A ``deisctl`` command fails with: 'Failed initializing SSH client: ssh: handshake failed: ssh: unable to authenticate'.
Did you remember to add your SSH key to the ssh-agent? ``ssh-add -L`` should list the key you used
to provision the servers. If it's not there, ``ssh-add -K /path/to/your/key``.

All the given peers are not reachable
-------------------------------------

A ``deisctl`` command fails with: 'All the given peers are not reachable (Tried to connect to each peer twice and failed)'.
The most common cause of this issue is that a new discovery URL wasn't generated and updated in
``contrib/coreos/user-data`` before the cluster was launched. Each Deis cluster must have a unique
discovery URL, or else ``etcd`` will try and fail to connect to old hosts. Try destroying the cluster
and relaunching the cluster with a fresh discovery URL.

You can use ``make discovery-url`` to automatically fetch a new discovery URL.

Could not find unit template...
-------------------------------

If you built ``deisctl`` locally or didn't use its installer, you may see an error like this:

    .. code-block:: console

        $ deisctl install platform

        Storage subsystem...
        Could not find unit template for store-daemon

This is because ``deisctl`` could not find unit files for Deis locally. Run
``deisctl help refresh-units`` to see where ``deisctl`` searches, and then run a command such as
``deisctl refresh-units --tag=v1.13.0``, or set the ``$DEISCTL_UNITS`` environment variable to a directory
containing the unit files.

Other issues
------------

Running into something not detailed here? Please `open an issue`_ or hop into #deis on Freenode IRC and we'll help!

.. _`open an issue`: https://github.com/deis/deis/issues/new
