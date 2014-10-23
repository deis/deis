:title: Register an Admin User with Deis
:description: Learn how to operate a Deis formation using the Deis command line interface.

.. _register-admin-user:

Register an Admin User
======================
Once your :ref:`Controller` is running, you must register an admin user
using the Deis client. The Deis command-line interface (CLI), or client,
allows you to interact with a Deis :ref:`Controller`. You must install
the client to use Deis.

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

Register a User
---------------
Now that the client is installed, create a user account on the Deis :ref:`Controller`.

.. important:: First User Gets Admin
   The first user to register with Deis receives "superuser" privileges.

.. code-block:: console

    $ deis register http://deis.example.com
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

.. _pip: http://www.pip-installer.org/en/latest/installing.html
.. _Python: https://www.python.org/
