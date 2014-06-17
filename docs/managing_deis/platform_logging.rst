:title: Platform logging
:description: Configuring platform logging.

.. _platform_logging:

Platform logging
=========================

Comprehensive platform logging is a goal for Deis 1.0. We are currently investigating solutions for
this, and progress can be tracked in GitHub issue `#980`_.

In the meantime, however, ``journalctl`` can be used to pipe log output to a remote host. For example,
to send service log output to Papertrail:

.. code-block:: console

    journalctl -o short -f | ncat --udp logs.papertrailapp.com 34000

.. _`#980`: https://github.com/deis/deis/issues/980
