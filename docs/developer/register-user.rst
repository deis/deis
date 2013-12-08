:title: Register a new Deis user using the client
:description: First steps for developers using Deis to deploy and scale applications.
:keywords: tutorial, guide, walkthrough, howto, deis, developer, dev

Register a User
===============
To deploy an :ref:`Application`, you must be logged into a Deis :ref:`Controller`.
To ``git push`` you must provide your SSH public key for authentication.

Create a User Account
---------------------
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

Upload Your SSH Public Key
--------------------------
Use the ``deis keys:add`` command to upload your default SSH public key, usually one of:

 * ~/.ssh/id_rsa.pub
 * ~/.ssh/id_dsa.pub

.. code-block:: console

    $ deis keys:add
    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

