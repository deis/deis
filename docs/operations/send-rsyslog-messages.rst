:title: Sending Logs to a Remote Syslog Server
:description: Guide to dumping the deis cluster's logs to a remote Syslog server
:keywords: tutorial, guide, walkthrough, howto, deis, sysadmins, operations, rsyslog

.. _send-rsyslog-messages:

Exporting Syslog Messages
=========================

In this document, we will demonstrate forwarding messages from Deis's syslog server to
another one. This is used in a number of cases:

- there is a legal requirement to consolidate all logs on a single system
- the remote server needs to have a full picture of the cluster's activity to perform
  correctly

In our case, you can forward all messages from your Deis cluster to your remote syslog
server by setting two keys in :ref:`Discovery`. SSH into your controller and run the
following commands:

.. code-block:: console

    $ etcdctl set /deis/logger/remoteHost <your_remote_syslog_hostname>
    $ etcdctl set /deis/logger/remotePort <your_remote_syslog_port>
