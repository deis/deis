:title: Managing users
:description: Managing users for your Deis cluster.

.. _managing_users:

Managing users
=========================

There are two classes of Deis users: normal users and administrators.

* Users can use most of the features of Deis - creating and deploying applications, adding/removing domains, etc.
* Administrators can perform all the actions that users can, but they can also create, edit, and destroy clusters.

The first user created on a Deis installation is automatically an administrator.

Promoting users to administrators
---------------------------------

You can use the ``deis perms`` command to promote a user to an administrator:

.. code-block:: console

    $ deis perms:create john --admin
