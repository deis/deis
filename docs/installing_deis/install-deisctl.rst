:title: Installing the Deis Control Utility
:description: Learn how to install the Deis Control Utility

.. _install_deisctl:

Install deisctl
===============

The Deis Control Utility, or ``deisctl`` for short, is a command-line client used to configure and
manage the Deis Platform.

Building from Installer
-----------------------

To install the latest version of deisctl, change to the directory where you would like to install
the binary. Then, install the Deis Control Utility by downloading and running the install script
with the following command:

.. code-block:: console

    $ cd ~/bin
    $ curl -sSL http://deis.io/deisctl/install.sh | sh -s 1.0.1

This installs deisctl to the current directory, and refreshes the Deis systemd unit files used to
schedule the components. Link it to /usr/local/bin, so it will be in your PATH:

.. code-block:: console

    $ sudo ln -fs $PWD/deisctl /usr/local/bin/deisctl

To change installation options, save the installer directly:

.. image:: download-linux-brightgreen.svg
   :target: https://s3-us-west-2.amazonaws.com/opdemand/deisctl-1.0.1-linux-amd64.run

.. image:: download-osx-brightgreen.svg
   :target: https://s3-us-west-2.amazonaws.com/opdemand/deisctl-1.0.1-darwin-amd64.run

Then run the downloaded file as a shell script. Append ``--help`` to see what options
are available.

.. important::

    Always use a version of ``deisctl`` that matches the Deis release.
    Verify this with ``deisctl --version``.

Builds are hosted on an S3 bucket at this URL format:

``https://s3-us-west-2.amazonaws.com/opdemand/deisctl-<VERSION>-<darwin|linux>-amd64.run``

For example, the deisctl release for Deis version 1.0.1 can be downloaded here:

.. image:: download-linux-brightgreen.svg
   :target: https://s3-us-west-2.amazonaws.com/opdemand/deisctl-1.0.1-linux-amd64.run

.. image:: download-osx-brightgreen.svg
   :target: https://s3-us-west-2.amazonaws.com/opdemand/deisctl-1.0.1-darwin-amd64.run

Building from Source
--------------------

If you want to install from source, ensure you have `godep`_ installed and run:

.. code-block:: console

	$ make -C deisctl build

You can then move or link the client so it will be in your path:

.. code-block:: console

	$ sudo ln -fs $PWD/deisctl/deisctl /usr/local/bin/deisctl


.. _`godep`: https://github.com/tools/godep
