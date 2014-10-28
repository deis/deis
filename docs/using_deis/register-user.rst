:title: Register a new Deis user using the client
:description: First steps for developers using Deis to deploy and scale applications.


Register a User
===============
To use Deis, you must first register a user on the :ref:`Controller`.
To ``git push`` you must provide your SSH public key for authentication.

Register with a Controller
--------------------------
Use ``deis register`` with the :ref:`Controller` URL (supplied by your Deis administrator)
to create a new account.  You will be logged in automatically.

.. code-block:: console

    $ deis register http://deis.example.com
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

.. important::

    The first user to register with Deis receives "superuser" privileges.

Upload Your SSH Public Key
--------------------------
If you plan on using ``git push`` to deploy applications to Deis, you must provide your SSH public key.  Use the ``deis keys:add`` command to upload your default SSH public key, usually one of:

 * ~/.ssh/id_rsa.pub
 * ~/.ssh/id_dsa.pub

.. code-block:: console

    $ deis keys:add
    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

Logout from a Controller
------------------------
Logout of an existing controller session using ``deis logout``.

.. code-block:: console

    $ deis logout
    Logged out as deis

Login to a Controller
---------------------
If you already have an account, use ``deis login`` to authenticate against the Deis :ref:`Controller`.

.. code-block:: console

    $ deis login http://deis.example.com
    username: deis
    password:
    Logged in as deis

.. note::

    Deis session information is stored in your user's ~/.deis directory.
