:title: Operational tasks
:description: Common operational tasks for your Deis cluster.

.. _operational_tasks:

Operational tasks
~~~~~~~~~~~~~~~~~

Below are some common operational tasks for managing the Deis platform.


Managing users
==============

There are two classes of Deis users: normal users and administrators.

* Users can use most of the features of Deis - creating and deploying applications, adding/removing domains, etc.
* Administrators can perform all the actions that users can, but they also have owner access to all applications.

The first user created on a Deis installation is automatically an administrator.


Promoting users to administrators
---------------------------------

You can use the ``deis perms`` command to promote a user to an administrator:

.. code-block:: console

    $ deis perms:create john --admin

.. _disable_user_registration:

Disabling user registration
---------------------------

You can disable user registration for everybody except admins:

.. code-block:: console

    $ deisctl config controller set registrationMode="admin_only"

If you want to entirely disable user registration:

.. code-block:: console

    $ deisctl config controller set registrationMode="disabled"

Re-issuing User Authentication Tokens
-------------------------------------

The controller API uses a simple token-based HTTP Authentication scheme. Token authentication is
appropriate for client-server setups, such as native desktop and mobile clients. Each user of the
platform is issued a token the first time that they sign up on the platform. If this token is
compromised, it will need to be regenerated.

A user can regenerate their own token like this:

.. code-block:: console

    $ deis auth:regenerate

An administrator can also regenerate the token of another user like this:

.. code-block:: console

    $ deis auth:regenerate -u test-user


At this point, the user will no longer be able to authenticate against the controller with his auth
token:

.. code-block:: console

    $ deis apps
    401 UNAUTHORIZED
    Detail:
    Invalid token

They will need to log back in to use their new auth token.

If there is a cluster wide security breach, an administrator can regenerate everybody's auth token like this:

.. code-block:: console

    $ deis auth:regenerate --all=true
