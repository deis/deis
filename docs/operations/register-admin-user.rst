:title: Register an Admin User with Deis
:description: Learn how to operate a Deis formation using the Deis command line interface.

.. _register-admin-user:

Register an Admin User
======================
Once your :ref:`Controller` is running you need to register an admin user
using the Deis command-line client.

Install the Deis Client
-----------------------
Install the latest Deis client using Python's `pip`_:

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

Register a User
---------------
Now that the client is installed, create a user account on the Deis :ref:`Controller`.

.. important:: First User Gets Admin
   The first user to register with Deis receives "superuser" priviledges.

.. code-block:: console

    $ deis register http://deis.example.com:8000
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

.. note::

    As of v0.5.1, the proxy was removed for Deis platform services. It has yet to be added
    back in. See `issue 535`_ for more details.

    As a workaround, use the following:

    :code:`deis register http://deis.example.com:8000`


.. _`pip`: http://www.pip-installer.org/en/latest/installing.html
.. _`issue 535`: https://github.com/opdemand/deis/issues/535
