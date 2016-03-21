:title: Installing the Deis Control Utility
:description: Learn how to install the Deis Control Utility

.. _install_deisctl:

Install deisctl
===============

The Deis Control Utility, or ``deisctl`` for short, is a command-line client used to configure and
manage the Deis Platform.

Run the Installer
-----------------

Change to the directory where you would like the ``deisctl`` binary to be installed, then download
and run the latest installer:

.. code-block:: console

    $ cd ~/bin
    $ curl -sSL http://deis.io/deisctl/install.sh | sh -s 1.13.0
    $ # on CoreOS, add "sudo" to install to /opt/bin/deisctl
    $ curl -sSL http://deis.io/deisctl/install.sh | sudo sh -s 1.13.0

This installs ``deisctl`` version 1.13.0 to the current directory, and downloads the matching
Deis systemd unit files used to schedule the components. Link ``deisctl`` into /usr/local/bin, so
it will be in your ``$PATH``:

.. code-block:: console

    $ sudo ln -fs $PWD/deisctl /usr/local/bin/deisctl

To change installation options, save the installer directly:

.. image:: download-linux-brightgreen.svg
   :target: https://s3-us-west-2.amazonaws.com/get-deis/deisctl-1.13.0-linux-amd64.run

.. image:: download-osx-brightgreen.svg
   :target: https://s3-us-west-2.amazonaws.com/get-deis/deisctl-1.13.0-darwin-amd64.run

Then run the downloaded file as a shell script. Append ``--help`` to see what options
are available.

.. important::

    Always use a version of ``deisctl`` that matches the Deis release.
    Verify this with ``deisctl --version``.


Building from Source
--------------------

To build ``deisctl`` locally, first :ref:`get the source <get_the_source>`, ensure
you have `godep`_ installed, and run:

.. code-block:: console

	$ make -C deisctl build

You can then move or link the client so it will be in your path:

.. code-block:: console

	$ sudo ln -fs $PWD/deisctl/deisctl /usr/local/bin/deisctl

.. note::

    Remember to run ``deisctl refresh-units`` or to set ``$DEISCTL_UNITS`` to an appropriate
    directory if you do not use the ``deisctl`` installer. Run the command
    ``deisctl help refresh-units`` for more information.


.. _`godep`: https://github.com/tools/godep
