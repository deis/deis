:title: Permissions
:description: Usage and descriptions of permissions in deis.

.. _user_permissions:

User Permissions
================

Deis has many options for granting and restricting a user's access to parts of the platform.

Administrators
--------------

Administrators have complete access to every part of the platform.

They have every permission listed granted and cannot be limited.

App Owners
----------

App owners have complete access to their app.

They have every app permission listed granted and their app permissions cannot be limited.

Types of Permissions
--------------------

There are two types of permissions, cluster level permissions and app level permissions.

Cluster level permissions apply to the entire cluster, such as the permission to create apps.
They can only be set by administrators.

App level permissions apply to an app, such as the permission to set config.
They can only be set by the app owner.

Default Permissions
-------------------

By default, users have the following permissions: ``app``, ``certs``.
When invited to an app, users have the ``push``, ``config``, ``domains``, ``scale`` permissions by default.

Changing the Default Permissions
--------------------------------

Default apps permissions are stored in etcd and can be updated with ``deisctl``.
Updating default permissions is retroactive and affects current and users.

This does not affect users that have had a permission explicitly set.
For example, if you banned a user from creating apps but you then updated
the default permission to allow users to create apps, that user still couldn't create
apps.

See :ref:`configure_default_permissions` to learn how to configure default permissions.

Setting Permissions
-------------------

To grant a permission to a user, you can run ``deis perms:create <username> --<permission>``
For example, to grant user foo permission to create an app ``deis perms:create foo --app``

To revoke a permission, you can run ``deis perms:delete <username> --<permission>``
For example, to prevent user foo from creating apps ``deis perms:delete foo --app``

Viewing Permissions
-------------------

To view what permissions you have been granted, you can run ``deis perms:view``.

Administrators or app owners can view other user's permissions with ``deis perms:view --username=<username>``.


List of Permissions
-------------------

Cluster Level Permissions
^^^^^^^^^^^^^^^^^^^^^^^^^
============== ==========================
permission     description
============== ==========================
app            can create apps
app-management full access to every app
certs          can view and set ssl certs
============== ==========================

App Level Permissions
^^^^^^^^^^^^^^^^^^^^^
========== ==========================
permission     description
========== ==========================
config     can view and set config
domains    can view and set domains
push       can push code to app
scale      can scale an app
========== ==========================
