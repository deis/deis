:title: Register an Admin User with Deis
:description: Learn how to operate a Deis formation using the Deis command line interface.

.. _register-admin-user:

Register an Admin User
======================
Once your :ref:`Controller` is running, you must register an admin user.
using the Deis client. The Deis command-line interface (CLI), or client,
allows you to interact with a Deis :ref:`Controller`. You must install
the client to use Deis.

Install with Pip
----------------
Install the latest Deis client using Python's pip_ package manager:

.. code-block:: console

    $ pip install deis
    Downloading/unpacking deis
      Downloading deis-0.8.0.tar.gz
      Running setup.py egg_info for package deis
      ...
    Successfully installed deis
    Cleaning up...
    $ deis
    Usage: deis <command> [<args>...]

If you don't have Python_ installed, you can download a binary executable
version of the Deis client for Mac OS X, Windows, or Debian Linux:

    - https://s3-us-west-2.amazonaws.com/opdemand/deis-osx-0.8.0.tgz
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-win32-0.8.0.zip
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-deb-wheezy-0.8.0.tgz

Register a User
---------------
Now that the client is installed, create a user account on the Deis :ref:`Controller`.

.. important:: First User Gets Admin
   The first user to register with Deis receives "superuser" priviledges.

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
