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


Re-issuing User Authentication Tokens
-------------------------------------

The controller API uses a simple token-based HTTP Authentication scheme. Token authentication is
appropriate for client-server setups, such as native desktop and mobile clients. Each user of the
platform is issued a token the first time that they sign up on the platform. If this token is
compromised, you'll need to manually intervene to re-issue a new authentication token for the user.
To do this, SSH into the node running the controller and drop into a Django shell:

.. code-block:: console

    $ fleetctl ssh deis-controller
    $ docker exec -it deis-controller python manage.py shell
    >>>

At this point, let's re-issue an auth token for this user. Let's assume that the name for the user
is Bob (poor Bob):

.. code-block:: console

    >>> from django.contrib.auth.models import User
    >>> from rest_framework.authtoken.models import Token
    >>> bob = User.objects.get(username='bob')
    >>> token = Token.objects.get(user=bob)
    >>> token.delete()
    >>> exit()

At this point, Bob will no longer be able to authenticate against the controller with his auth
token:

.. code-block:: console

    $ deis apps
    401 UNAUTHORIZED
    Detail:
    Invalid token

For Bob to be able to use the API again, he will have to authenticate against the controller to be
re-issued a new token:

.. code-block:: console

    $ deis login http://deis.example.com
    username: bob
    password:
    Logged in as bob
    $ deis apps
    === Apps
