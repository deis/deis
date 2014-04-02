:title: Manage the Deis PaaS Controller
:description: Learn how to manage your Deis controller.
:keywords: guide, howto, deis, controller, sysadmins, operations, syslog

.. _manage-controller:

Manage the Controller
=====================

The following documentation will help you get set up with managing your Deis PaaS
controller to your liking.

Sending Logs to a Remote Syslog Server
--------------------------------------

You can forward messages from Deis's syslog server to another remote server hosted outside
of the cluster. This is used in a number of cases:

- there is a legal requirement to consolidate all logs on a single system
- the remote server needs to have a full picture of the cluster's activity to perform
  correctly

In our case, you can forward all messages from your Deis cluster to your remote syslog
server by setting two keys in :ref:`Discovery`. SSH into your controller and run the
following commands:

.. code-block:: console

    $ etcdctl set /deis/logger/remoteHost <your_remote_syslog_hostname>
    $ etcdctl set /deis/logger/remotePort <your_remote_syslog_port>

Setting up new Routes
---------------------

By default, Deis's :ref:`Router` sets up only one route, which is just a route from your
publicly viewable IP address at port 80 to port 8000 (the port that your :ref:`Controller`
typically resides). No domain-specfic routes are defined. However, You can set up new
upstream routes for the Deis :ref:`Router` to handle. In the background, Deis's router talks
to Redis to retrieve these upstream routes. For example, if you want to set up a new route
to the :ref:`Registry`, you can set a new value in Redis like so:

.. code-block:: console

    $ dig registry.deisapp.com
    [...]
    ;; ANSWER SECTION:
    registry.deisapp.com. 1649 IN  A   <ip-address>
    $ sudo apt-get install -yq redis-server
    $ redis-cli
    redis 127.0.0.1:6379> SET registry.deisapp.com <ip-address>:5000
    redis 127.0.0.1:6379> exit
    $ curl -L registry.deisapp.com
    "docker-registry server (production)"
