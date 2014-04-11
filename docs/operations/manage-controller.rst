:title: Manage the Deis PaaS Controller
:description: Learn how to manage your Deis controller.

.. _manage-controller:

Manage the Controller
=====================
The following documentation will help you get set up with managing your
Deis PaaS controller to your liking.

Sending Logs to a Remote Syslog Server
--------------------------------------
You can forward messages from Deis's syslog server to another remote
server hosted outside of the cluster. This is used in a number of cases:

- there is a legal requirement to consolidate all logs on a
  single system
- the remote server needs to have a full picture of the cluster's
  activity to perform correctly

In our case, you can forward all messages from your Deis cluster to your
remote syslog server by setting two keys in *etcd*. SSH into your
controller and run the following commands:

.. code-block:: console

    $ etcdctl set /deis/logger/remoteHost <your_remote_syslog_hostname>
    $ etcdctl set /deis/logger/remotePort <your_remote_syslog_port>
