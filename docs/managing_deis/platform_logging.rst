:title: Platform logging
:description: Configuring platform logging.

.. _platform_logging:

Platform logging
================

Logging for Deis components and deployed applications is handled by two components:
:ref:`logger` and :ref:`logspout`.

``deis-logspout`` runs on all CoreOS hosts, collects logs from running containers
and sends their logs to ``/deis/logs/host`` and ``/deis/logs/port``.

``deis-logger`` collects the logs sent by logspout and archives them for use by :ref:`Controller`
when a client runs ``deis logs``. This component publishes its host and port to ``/deis/logs/host``
and ``/deis/logs/port``, and is typically the service which consumes logs from ``deis-logspout``.

Application log drain
---------------------

Application logs can be drained to an external syslog server (or compatible service such as Logstash, Papertrail, Splunk etc).

.. code-block:: console

    $ deisctl config logs set drain=syslog://logs2.papertrailapp.com:23654

This will send all application logs - there is currently no way to drain logs per application.

Routing host logs to a custom location
--------------------------------------

Logging to an external location can be achieved without modifying the log flow within Deis -
we can simply send the master journal on a CoreOS host using ``ncat``. For example, if I'm using the
`Papertrail`_ hosted log service, I can forward all logs on a host to Papertrail using the host
and port provided to me by Papertrail:

.. code-block:: console

    $ journalctl -o short -f | ncat --ssl logs2.papertrailapp.com 23654

This is really only useful when shipped as a service that we don't have to run ourselves in
a shell. We can use a fleet service for this:

.. code-block:: console

    [Unit]
    Description=Log forwarder

    [Service]
    ExecStart=/bin/sh -c "journalctl -o short -f | ncat --ssl logs2.papertrailapp.com 23654"

    [Install]
    WantedBy=multi-user.target

    [X-Fleet]
    Global=true

Save the file as ``log-forwarder.service``. Load and start the service with
``fleetctl load log-forwarder.service && fleetctl start log-forwarder.service``.

Shortly thereafter, you should start to see logs from every host in your cluster appear in the
Papertrail dashboard.

.. _`logspout`: https://github.com/progrium/logspout
.. _`papertrail`: https://papertrailapp.com/
