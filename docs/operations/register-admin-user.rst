:title: Register an Admin User with Deis
:description: Learn how to operate a Deis formation using the Deis command line interface.
:keywords: tutorial, guide, walkthrough, howto, deis, formations

.. _register-admin-user:

Register an Admin User
======================
Once your :ref:`Controller` is running you need to register an admin user
using the Deis command-line client.

Install the Deis Client
-----------------------
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
    $ deis
    Usage: deis <command> [<args>...]

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

Discover Provider Credentials
-----------------------------
.. important:: Provider API
   If you don't want the ability to scale servers automatically
   using the Deis :ref:`Provider` API you can skip this section.

If you want to use automated provisioning, you'll need to provide Deis 
with cloud provider credentials used to bootstrap :ref:`Nodes <node>`.

The ``deis providers:discover`` command
will look at standard environment variables on your workstation to discover
credentials for supported cloud providers.
The table below shows how the environment variables on your workstation map to
provider types and fields stored in your Deis user account.

======================= =============== ==============
Variable Name           Provider Type   Provider Field
======================= =============== ==============
AWS_ACCESS_KEY_ID       ec2             access_key
AWS_SECRET_ACCESS_KEY   ec2             secret_key
AWS_ACCESS_KEY          ec2             access_key
AWS_SECRET_KEY          ec2             secret_key
RACKSPACE_USERNAME      rackspace       username
RACKSPACE_API_KEY       rackspace       api_key
DIGITALOCEAN_CLIENT_ID  digitalocean    client_id
DIGITALOCEAN_API_KEY    digitalocean    api_key
======================= =============== ==============

----

To discover providers using the Deis client:

.. code-block:: console

    $ deis providers:discover
    Discovered EC2 credentials: AAAAAAAAAAAAAAAAAAAA
    Import EC2 credentials? (y/n) : y
    Uploading EC2 credentials... done
    No Rackspace credentials discovered.
    No DigitalOcean credentials discovered.
    No Vagrant VMs discovered.


.. _`pip`: http://www.pip-installer.org/en/latest/installing.html