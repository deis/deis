:title: Install the Deis Client on your Workstation
:description: First steps for developers using Deis to deploy and scale applications.

.. _install-client:

Install the Client
==================
The Deis command-line interface (CLI), or client, allows you to interact
with a Deis :ref:`Controller`. You must install the client to use Deis.

Install the Deis Client
-----------------------
Your Deis client should match your server's version. For development, an
easy way to ensure this is to run `client/deis.py` in the code repository
you used to provision the server. You can make a symlink or shell alias for
`deis` to that file:

.. code-block:: console

    $ pip install docopt==0.6.2 python-dateutil==2.2 requests==2.3.0 termcolor==1.1.0
    $ sudo ln -fs $(pwd)/client/deis.py /usr/local/bin/deis
    $ deis
    Usage: deis <command> [<args>...]

If you don't have Python_, install the latest `deis` binary executable for
Linux or Mac OS X with this command:

.. code-block:: console

    $ curl -sSL http://deis.io/deis-cli/install.sh | sh

The installer puts `deis` in your current directory, but you should move it
somewhere in your $PATH.

Proxy Support
-------------
Set the ```http_proxy``` or ```https_proxy``` environment variable to enable proxy support:

.. code-block:: console

    $ export http_proxy="http://proxyip:port"
    $ export https_proxy="http://proxyip:port"

Integrated Help
---------------
The Deis client comes with comprehensive documentation for every command.
Use ``deis help`` to explore the commands available to you:

.. code-block:: console

    $ deis help
    The Deis command-line client issues API calls to a Deis controller.

    Usage: deis <command> [<args>...]

    Auth commands::

      register      register a new user with a controller
      login         login to a controller
      logout        logout from the current controller

    Subcommands, use ``deis help [subcommand]`` to learn more::
    ...

To get help on subcommands, use ``deis help [subcommand]``:

.. code-block:: console

    $ deis help apps
    Valid commands for apps:

    apps:create        create a new application
    apps:list          list accessible applications
    apps:info          view info about an application
    apps:open          open the application in a browser
    apps:logs          view aggregated application logs
    apps:run           run a command in an ephemeral app container
    apps:destroy       destroy an application

    Use `deis help [command]` to learn more

.. _pip: http://www.pip-installer.org/en/latest/installing.html
.. _Python: https://www.python.org/
