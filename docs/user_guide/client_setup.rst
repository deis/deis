:title: Introduction to the Deis client
:description: An introduction to the Deis command-line interface.

.. _client:

The Client
==========

The Deis client is the primary user interface to Deis. It accepts commands from the user
and communicates with a Deis cluster. The client communicates to the server via sockets or
though a RESTful API. You can use it in a similar fashion to Heroku's client: create
applications, add SSH keys, share applications between users, add custom domains to your
applications, as well as numerous other management operations.

The :ref:`Client Reference <client_ref>` has full descriptions of all client commands and
options. The documentation for these commands are also available in the client with
:code:`deis help`.

.. _install_client:

Install the Client
------------------

The Deis client allows you to interact with a Deis :ref:`Controller`. You must install the
client to use Deis.

Download Binaries
^^^^^^^^^^^^^^^^^

You can download a binary executable version of the Deis client for Mac OS X, Windows, or Debian Linux:

    - https://s3-us-west-2.amazonaws.com/opdemand/deis-osx-0.9.0.tgz
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-win32-0.9.0.zip
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-deb-wheezy-0.9.0.tgz

Extract the Deis client and place it in your workstation path.

Install with Pip
^^^^^^^^^^^^^^^^

You can also install the latest client using pip_, a third-party package manager
for Python applications:

.. code-block:: console

    $ sudo pip install --upgrade deis
    Downloading/unpacking deis
      Downloading deis-0.9.0.tar.gz
      Running setup.py egg_info for package deis
      ...
    Successfully installed deis
    Cleaning up...

    $ deis
    Usage: deis <command> [<args>...]

Getting Help
------------

The Deis client comes with comprehensive documentation for every command.
Use :code:`deis help` to explore the commands available to you:

.. code-block:: console

    $ deis help
    The Deis command-line client issues API calls to a Deis controller.

    Usage: deis <command> [<args>...]

    Auth commands::

      register      register a new user with a controller
      login         login to a controller
      logout        logout from the current controller

    [...]

To get help on subcommands, use :code:`deis help [subcommand]`:

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

See also the :ref:`Client Reference <client_ref>` for a full list of commands.

Register a User
---------------

To use Deis, you must first register a user on the :ref:`Controller`. Use
:ref:`deis register <deis_auth>` with the :ref:`Controller` URL (supplied by your Deis
administrator) to create a new account.  You will be logged in automatically.

.. code-block:: console

    $ deis register http://deis.example.com
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

Upload Your SSH Public Key
--------------------------

If you plan on using :code:`git push` to deploy applications to Deis, you must provide
your SSH public key.  Use the :ref:`deis keys:add <deis_keys>` command to upload your
default SSH public key, usually one of:

- ~/.ssh/id_rsa.pub
- ~/.ssh/id_dsa.pub

.. code-block:: console

    $ deis keys:add
    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

Login to a Controller
---------------------

If you already have an account, use :ref:`deis login <deis_auth>` to authenticate against
the Deis :ref:`Controller`.

.. code-block:: console

    $ deis login http://deis.example.com
    username: deis
    password:
    Logged in as deis

.. note::

    Deis session information is stored in your user's ~/.deis directory.

Logout from a Controller
------------------------

Logout of an existing controller session using :ref:`deis logout <deis_auth>`.

.. code-block:: console

    $ deis logout
    Logged out as deis


.. _pip: http://www.pip-installer.org/en/latest/installing.html
.. _Python: https://www.python.org/
