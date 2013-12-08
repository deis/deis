:title: Install the Deis Client on your Workstation
:description: First steps for developers using Deis to deploy and scale applications.
:keywords: tutorial, guide, walkthrough, howto, deis, developer, dev

Install the Client
==================
The Deis client allows you to interact with a Deis :ref:`Controller`.
You'll need to install the client before you can use Deis.

Install with Pip
----------------
Install the latest stable client using Python's `pip`_:

.. code-block:: console

    $ sudo pip install deis
    Password:
    Downloading/unpacking deis
      Downloading deis-0.3.0.tar.gz
      Running setup.py egg_info for package deis
    ...
    Successfully installed deis
    Cleaning up...

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

.. _`pip`: http://www.pip-installer.org/en/latest/installing.html
