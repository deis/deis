:title: Customizing controller
:description: Learn how to tune custom Deis settings.

.. _controller_settings:

Customizing controller
=========================
The following settings are tunable for the :ref:`controller` component.

Dependencies
------------
Requires: :ref:`database <database_settings>`, :ref:`registry <registry_settings>`

Required by: :ref:`router <router_settings>`

Considerations: none

Settings set by controller
--------------------------
The following etcd keys are set by the controller component, typically in its /bin/boot script.

=============================            =================================================================================
setting                                  description
=============================            =================================================================================
/deis/controller/host                    IP address of the host running controller
/deis/controller/port                    port used by the controller service (default: 8000)
/deis/controller/protocol                protocol for controller (default: http)
/deis/controller/secretKey               used for secrets (default: randomly generated)
/deis/controller/builderKey              used by builder to authenticate with the controller (default: randomly generated)
/deis/controller/unitHostname            See `Unit hostname`_. (default: "default")
/deis/builder/users/*                    stores user SSH keys (used by builder)
/deis/domains/*                          domain configuration for applications (used by router)
/deis/logs/host                          IP address of the host running logger
=============================            =================================================================================

Settings used by controller
---------------------------
The following etcd keys are used by the controller component.

====================================      ======================================================
setting                                   description
====================================      ======================================================
/deis/controller/registrationMode         set registration to "enabled", "disabled", or "admin_only" (default: "enabled")
/deis/controller/subdomain                subdomain used by the router for API requests (default: "deis")
/deis/controller/webEnabled               enable controller web UI (default: 0)
/deis/controller/workers                  number of web worker processes (default: CPU cores * 2 + 1)
/deis/database/host                       host of the database component (set by database)
/deis/database/port                       port of the database component (set by database)
/deis/database/engine                     database engine (set by database)
/deis/database/name                       database name (set by database)
/deis/database/user                       database user (set by database)
/deis/database/password                   database password (set by database)
/deis/registry/host                       host of the registry component (set by registry)
/deis/registry/port                       port of the registry component (set by registry)
/deis/registry/protocol                   protocol of the registry component (set by registry)
====================================      ======================================================

Using a custom controller image
-------------------------------
You can use a custom Docker image for the controller component instead of the image
supplied with Deis:

.. code-block:: console

    $ deisctl config controller set image=myaccount/myimage:latest

This will pull the image from the public Docker registry. You can also pull from a private
registry:

.. code-block:: console

    $ deisctl config controller set image=registry.mydomain.org:5000/myaccount/myimage:latest

Be sure that your custom image functions in the same way as the `stock controller image`_ shipped with
Deis. Specifically, ensure that it sets and reads appropriate etcd keys.

.. _`stock controller image`: https://github.com/deis/deis/tree/master/controller

Unit hostname
-------------
Per default, Docker automatically generates a hostname for your application unit, such as:
``5c149b397cd6``. Auto generated hostnames is not always preferred. For instance,
New Relic would classify each Docker container as an unique server since they use hostname
for grouping applications running on the same server together.

Deis supports configuring hostname assignment through the ``unitHostname`` setting.
You can change the assignment solution using the following command:

.. code-block:: console

    $ deisctl config controller set unitHostname=application

The valid ``unitHostname`` values are:

default
    Docker will generate the hostname. Example: ``5c149b397cd6``

application
    The hostname is assigned based on the unit name. Example: ``dancing-cat.v2.web.1``

server
    The hostname is assigned based on the CoreOS hostname. Example:
    ``ip-10-21-2-168.eu-west-1.compute.internal``

.. note::

    Changes to ``/deis/controller/unitHostname`` requires either pushing a new build to
    every application or scaling them down and up.
    The change is only detected when a container unit is deployed.

Changing the Registration Mode
------------------------------

By default, anybody can register a user with the Deis controller.
However, this is often undesirable from a security point of view.

Deis supports configuring the registration mode through the ``registrationMode`` setting.

Registration Modes
^^^^^^^^^^^^^^^^^^
========== =========================================================
mode       description
========== =========================================================
enabled    Default. Anybody can register a user with the controller.
disabled   Nobody can register a user with the controller.
admin_only Only admins can register a user with the controller.
========== =========================================================

This will set the registration mode to admin_only.

.. code-block:: console

    $ deisctl config controller set registrationMode="admin_only"

Using a LDAP Auth
-----------------
The Deis controller supports Single Sign On access control, for now Deis is able to authenticate using LDAP or Active Directory.

Settings used by LDAP
^^^^^^^^^^^^^^^^^^^^^
=========================================           =================================================================================
setting                                             description
=========================================           =================================================================================
/deis/controller/auth/ldap/endpoint                 The full LDAP endpoint. (Ex.: ldap://ldap.company.com)
/deis/controller/auth/ldap/bind/dn                  Full user for bind. (Ex.: user@company.com. For Anonymous bind leave blank)
/deis/controller/auth/ldap/bind/password            Password of the user for bind. (For anonymous bind leave blank)
/deis/controller/auth/ldap/user/basedn              The BASE DN where your LDAP Users are placed. (Ex.: OU=TeamX,DC=Company,DC=com)
/deis/controller/auth/ldap/user/filter              The field that we will match with username of Deis. (In most cases is uuid, AD uses sAMAccountName)
/deis/controller/auth/ldap/group/basedn             The BASE DN where the groups of your LDAP are are located. (Ex.: OU=Groups,OU=TeamX,DC=Company,DC=com)
/deis/controller/auth/ldap/group/filter             The field that we will locate your groups with LDAPSearch. (In most cases is objectClass)
/deis/controller/auth/ldap/group/type               The Groups type of LDAP. (Use groupOfNames if you don't know)
=========================================           =================================================================================

Configuring LDAP on Controller
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. important::

    It's important that you register the first user of the default auth in order to have an admin ( see :ref:`Register a User <register-user>` ) without this you don't have any deis admin because LDAP users haven't this permission, you will need to set this later.
    After this you need to disable the registration ( see :ref:`disable_user_registration` ) avoiding that "ghost" users register and access your Deis. The auth model of controller by default allows multiple source auths so LDAP and non-LDAP users will be able to login.


.. code-block:: console

    $ deisctl config controller set auth/ldap/endpoint=<ldap-endpoint>
    $ deisctl config controller set auth/ldap/bind/dn=<bind-dn-full-user>
    $ deisctl config controller set auth/ldap/bind/password=<bind-dn-user-password>
    $ deisctl config controller set auth/ldap/user/basedn=<user-base-dn>
    $ deisctl config controller set auth/ldap/user/filter=<user-filter>
    $ deisctl config controller set auth/ldap/group/basedn=<group-base-dn>
    $ deisctl config controller set auth/ldap/group/filter=<group-filter>
    $ deisctl config controller set auth/ldap/group/type=<group-type>

.. note::

    You can set a LDAP user as admin by using ``deis perms:create <LDAP User> --admin`` with the admin created before.

.. note::

    LDAP support was contributed by community member Pedro Spagiari (`@phspagiari <http://github.com/phspagiari/>`_) and is unsupported by the Deis core team.
